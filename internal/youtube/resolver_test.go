package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// fakeCommandRunner 是測試用的 CommandRunner 實作。
type fakeCommandRunner struct {
	output []byte
	err    error
}

func (r *fakeCommandRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return r.output, r.err
}

// TestResolver_ParsesValidJSON 測試 resolver 能正確解析 yt-dlp 的 JSON 輸出。
func TestResolver_ParsesValidJSON(t *testing.T) {
	validJSON := `{
		"title": "Test Song",
		"webpage_url": "https://www.youtube.com/watch?v=test123",
		"url": "https://example.test/stream.m4a"
	}`

	runner := &fakeCommandRunner{output: []byte(validJSON)}
	resolver := NewResolver(runner)

	song, err := resolver.Resolve(context.Background(), "test query")
	if err != nil {
		t.Fatalf("Resolve 不應回傳錯誤: %v", err)
	}

	if song.Title != "Test Song" {
		t.Errorf("Title 應為 'Test Song'，實際為 %q", song.Title)
	}
	if song.URL != "https://www.youtube.com/watch?v=test123" {
		t.Errorf("URL 應為 YouTube 網頁 URL，實際為 %q", song.URL)
	}
	if song.StreamURL != "https://example.test/stream.m4a" {
		t.Errorf("StreamURL 應為直接串流 URL，實際為 %q", song.StreamURL)
	}
}

// TestResolver_ReturnsErrorForInvalidJSON 測試 resolver 對無效 JSON 的處理。
func TestResolver_ReturnsErrorForInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json`

	runner := &fakeCommandRunner{output: []byte(invalidJSON)}
	resolver := NewResolver(runner)

	_, err := resolver.Resolve(context.Background(), "test query")
	if err == nil {
		t.Fatal("無效 JSON 應回傳錯誤")
	}

	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("錯誤應為 json.SyntaxError，實際為 %T", err)
	}
}

// TestResolver_ReturnsErrorForCommandFailure 測試 resolver 對指令執行失敗的處理。
func TestResolver_ReturnsErrorForCommandFailure(t *testing.T) {
	runner := &fakeCommandRunner{err: errors.New("command failed")}
	resolver := NewResolver(runner)

	_, err := resolver.Resolve(context.Background(), "test query")
	if err == nil {
		t.Fatal("指令執行失敗應回傳錯誤")
	}
}

// TestResolver_BuildsCorrectArgsForURL 測試 resolver 對 URL 查詢建立正確的參數。
func TestResolver_BuildsCorrectArgsForURL(t *testing.T) {
	validJSON := `{
		"title": "Test",
		"webpage_url": "https://youtube.com/watch?v=test",
		"url": "https://example.test/stream"
	}`

	var capturedArgs []string
	runner := &fakeCommandRunner{output: []byte(validJSON)}
	originalRunner := runner

	// 包裝 runner 來捕獲參數
	wrappedRunner := &fakeCommandRunner{
		output: originalRunner.output,
		err:    originalRunner.err,
	}

	resolver := &ytdlpResolver{
		runner: &captureArgsRunner{
			delegate: wrappedRunner,
			captured: &capturedArgs,
		},
	}

	_, err := resolver.Resolve(context.Background(), "https://youtube.com/watch?v=test123")
	if err != nil {
		t.Fatalf("Resolve 不應回傳錯誤: %v", err)
	}

	// 驗證參數包含 -j 和查詢字串
	foundJ := false
	foundQuery := false
	for _, arg := range capturedArgs {
		if arg == "-j" {
			foundJ = true
		}
		if arg == "https://youtube.com/watch?v=test123" {
			foundQuery = true
		}
	}

	if !foundJ {
		t.Error("參數應包含 -j")
	}
	if !foundQuery {
		t.Error("參數應包含查詢字串")
	}
}

// TestResolver_BuildsCorrectArgsForSearchTerm 測試 resolver 對搜尋關鍵字建立正確的參數。
func TestResolver_BuildsCorrectArgsForSearchTerm(t *testing.T) {
	validJSON := `{
		"title": "Test",
		"webpage_url": "https://youtube.com/watch?v=test",
		"url": "https://example.test/stream"
	}`

	var capturedArgs []string
	runner := &fakeCommandRunner{output: []byte(validJSON)}

	resolver := &ytdlpResolver{
		runner: &captureArgsRunner{
			delegate: runner,
			captured: &capturedArgs,
		},
	}

	_, err := resolver.Resolve(context.Background(), "test search term")
	if err != nil {
		t.Fatalf("Resolve 不應回傳錯誤: %v", err)
	}

	// 驗證參數包含 ytsearch1: 前綴
	foundSearch := false
	for _, arg := range capturedArgs {
		if arg == "ytsearch1:test search term" {
			foundSearch = true
			break
		}
	}

	if !foundSearch {
		t.Errorf("搜尋關鍵字參數應包含 'ytsearch1:' 前綴，實際參數：%v", capturedArgs)
	}
}

// captureArgsRunner 包裝 CommandRunner 來捕獲參數供測試驗證。
type captureArgsRunner struct {
	delegate CommandRunner
	captured *[]string
}

func (r *captureArgsRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	*r.captured = args
	return r.delegate.Run(ctx, name, args...)
}
