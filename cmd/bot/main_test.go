package main

import (
	"errors"
	"strings"
	"testing"

	"discordbot/internal/config"
)

// TestRun_Success 測試 run 能成功載入設定、啟動 Bot、等待關閉並停止 Bot
func TestRun_Success(t *testing.T) {
	restore := replaceMainDependencies(t)
	defer restore()

	fake := &fakeRunnableBot{}
	newBot = func(cfg *config.Config) (runnableBot, error) {
		if cfg.BotToken != "bot-token" {
			t.Errorf("BotToken 應為 bot-token，實際為 %q", cfg.BotToken)
		}
		return fake, nil
	}

	var printed []string
	printLine = func(a ...any) (int, error) {
		printed = append(printed, strings.TrimSpace(strings.Join(toStrings(a), " ")))
		return 0, nil
	}

	if err := run(); err != nil {
		t.Fatalf("run 不應回傳錯誤，但得到: %v", err)
	}
	if !fake.started {
		t.Fatal("run 應啟動 bot")
	}
	if !fake.waited {
		t.Fatal("run 應等待 shutdown")
	}
	if !fake.stopped {
		t.Fatal("run 應停止 bot")
	}
	if len(printed) != 2 {
		t.Fatalf("run 應輸出 2 行訊息，實際為 %#v", printed)
	}
}

// TestRun_ReturnsNewBotError 測試建立 Bot 失敗時 run 會回傳錯誤
func TestRun_ReturnsNewBotError(t *testing.T) {
	restore := replaceMainDependencies(t)
	defer restore()

	wantErr := errors.New("new failed")
	newBot = func(cfg *config.Config) (runnableBot, error) {
		return nil, wantErr
	}

	err := run()
	if !errors.Is(err, wantErr) {
		t.Fatalf("run 應回傳 New 錯誤，實際為 %v", err)
	}
}

// TestRun_ReturnsStartError 測試 Bot 啟動失敗時 run 會回傳錯誤且不呼叫 Stop
func TestRun_ReturnsStartError(t *testing.T) {
	restore := replaceMainDependencies(t)
	defer restore()

	wantErr := errors.New("start failed")
	fake := &fakeRunnableBot{startErr: wantErr}
	newBot = func(cfg *config.Config) (runnableBot, error) {
		return fake, nil
	}

	err := run()
	if !errors.Is(err, wantErr) {
		t.Fatalf("run 應回傳 Start 錯誤，實際為 %v", err)
	}
	if fake.stopped {
		t.Fatal("Start 失敗時不應呼叫 Stop")
	}
}

// TestMain_CallsFatalOnRunError 測試 main 在 run 失敗時會呼叫 fatal logger
func TestMain_CallsFatalOnRunError(t *testing.T) {
	restore := replaceMainDependencies(t)
	defer restore()

	runApp = func() error {
		return errors.New("run failed")
	}

	var got string
	logFatalf = func(format string, v ...any) {
		got = format
		if len(v) != 1 || v[0].(error).Error() != "run failed" {
			t.Fatalf("logFatalf 參數不符合預期: %#v", v)
		}
	}

	main()

	if got != "%v" {
		t.Fatalf("main 應用 %%v 呼叫 logFatalf，實際 format 為 %q", got)
	}
}

// fakeRunnableBot 是測試用 Bot，記錄 Start、Stop、WaitForShutdown 是否被呼叫
type fakeRunnableBot struct {
	started  bool
	stopped  bool
	waited   bool
	startErr error
}

// Start 記錄測試 Bot 已啟動，並回傳預設的啟動錯誤
func (f *fakeRunnableBot) Start() error {
	f.started = true
	return f.startErr
}

// Stop 記錄測試 Bot 已停止
func (f *fakeRunnableBot) Stop() {
	f.stopped = true
}

// WaitForShutdown 記錄測試 Bot 已等待關閉訊號
func (f *fakeRunnableBot) WaitForShutdown() {
	f.waited = true
}

// replaceMainDependencies 替換 main package 的外部相依，並回傳還原函式
func replaceMainDependencies(t *testing.T) func() {
	t.Helper()

	originalLoadConfig := loadConfig
	originalNewBot := newBot
	originalPrintLine := printLine
	originalLogFatalf := logFatalf
	originalRunApp := runApp

	loadConfig = func() *config.Config {
		return &config.Config{BotToken: "bot-token", GuildID: "guild-id"}
	}
	newBot = func(cfg *config.Config) (runnableBot, error) {
		return &fakeRunnableBot{}, nil
	}
	printLine = func(a ...any) (int, error) { return 0, nil }
	logFatalf = func(format string, v ...any) {}
	runApp = run

	return func() {
		loadConfig = originalLoadConfig
		newBot = originalNewBot
		printLine = originalPrintLine
		logFatalf = originalLogFatalf
		runApp = originalRunApp
	}
}

// toStrings 將測試輸出的 any slice 轉成 string slice
func toStrings(values []any) []string {
	strings := make([]string, len(values))
	for i, value := range values {
		strings[i] = value.(string)
	}
	return strings
}
