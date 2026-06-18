package youtube

import (
	"context"
	"encoding/json"
	"strings"

	"discordbot/internal/player"
)

// Resolver 定義 YouTube 查詢解析介面。
type Resolver interface {
	// Resolve 將 YouTube URL 或搜尋關鍵字解析為 player.Song。
	Resolve(ctx context.Context, query string) (player.Song, error)
}

// CommandRunner 定義外部指令執行介面，方便測試時注入 fake 實作。
type CommandRunner interface {
	// Run 執行外部指令並回傳標準輸出。
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// ytdlpResolver 使用 yt-dlp 解析 YouTube 查詢的 Resolver 實作。
type ytdlpResolver struct {
	runner CommandRunner
}

// NewResolver 建立使用指定 CommandRunner 的 Resolver。
func NewResolver(runner CommandRunner) Resolver {
	return &ytdlpResolver{runner: runner}
}

// ytdlpOutput 對應 yt-dlp -j 的 JSON 輸出格式。
type ytdlpOutput struct {
	Title      string `json:"title"`
	WebpageURL string `json:"webpage_url"`
	URL        string `json:"url"`
}

// Resolve 實作 Resolver 介面，透過 yt-dlp 解析查詢。
func (r *ytdlpResolver) Resolve(ctx context.Context, query string) (player.Song, error) {
	args := r.buildArgs(query)

	output, err := r.runner.Run(ctx, "yt-dlp", args...)
	if err != nil {
		return player.Song{}, err
	}

	var result ytdlpOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return player.Song{}, err
	}

	return player.Song{
		Title:     result.Title,
		URL:       result.WebpageURL,
		StreamURL: result.URL,
	}, nil
}

// buildArgs 根據查詢字串建立 yt-dlp 指令參數。
func (r *ytdlpResolver) buildArgs(query string) []string {
	args := []string{"-j", "--no-warnings"}

	// 判斷是 URL 還是搜尋關鍵字
	if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
		args = append(args, query)
	} else {
		// 搜尋關鍵字使用 ytsearch1: 前綴
		args = append(args, "ytsearch1:"+query)
	}

	return args
}
