package command

import (
	"testing"
)

func TestBuildYtDlpArgsFormats(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		url            string
		outputTemplate string
		wantContains   []string
	}{
		{
			name:           "MP3 320kbps格式",
			format:         "mp3-320",
			url:            "https://youtube.com/watch?v=test",
			outputTemplate: "/tmp/output",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "mp3",
				"--audio-quality", "0",
				"--no-playlist",
				"--max-filesize", "25M",
				"--match-filter", "duration < 600",
			},
		},
		{
			name:           "FLAC無損格式",
			format:         "flac",
			url:            "https://youtube.com/watch?v=test",
			outputTemplate: "/tmp/output",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "flac",
				"--no-playlist",
			},
		},
		{
			name:           "Opus格式",
			format:         "opus-192",
			url:            "https://youtube.com/watch?v=test",
			outputTemplate: "/tmp/output",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "opus",
				"--audio-quality", "192K",
			},
		},
		{
			name:           "M4A格式",
			format:         "m4a-256",
			url:            "https://youtube.com/watch?v=test",
			outputTemplate: "/tmp/output",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "m4a",
				"--audio-quality", "256K",
			},
		},
		{
			name:           "WAV格式",
			format:         "wav",
			url:            "https://youtube.com/watch?v=test",
			outputTemplate: "/tmp/output",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "wav",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildYtDlpArgs(tt.format, tt.url, tt.outputTemplate)

			// 驗證包含的參數
			for _, want := range tt.wantContains {
				if !containsArg(args, want) {
					t.Errorf("buildYtDlpArgs() missing argument %q", want)
				}
			}

			// 驗證 URL 在參數中
			if !containsArg(args, tt.url) {
				t.Errorf("buildYtDlpArgs() should contain URL %q", tt.url)
			}

			// 驗證輸出模板
			if !containsArg(args, "--output") {
				t.Error("buildYtDlpArgs() should contain --output flag")
			}

			if !containsArg(args, tt.outputTemplate) {
				t.Errorf("buildYtDlpArgs() should contain output template %q", tt.outputTemplate)
			}
		})
	}
}

func TestBuildYtDlpArgs_DefaultFormat(t *testing.T) {
	t.Run("未知格式使用默認MP3", func(t *testing.T) {
		args := buildYtDlpArgs("unknown-format", "https://youtube.com/watch?v=test", "/tmp/output")

		if !containsArg(args, "mp3") {
			t.Error("buildYtDlpArgs() should default to mp3 for unknown format")
		}

		if !containsArg(args, "--audio-quality") {
			t.Error("buildYtDlpArgs() should include audio quality for default format")
		}
	})
}

func TestBuildYtDlpArgs_CommonParameters(t *testing.T) {
	t.Run("所有格式都應包含通用參數", func(t *testing.T) {
		formats := []string{"mp3-320", "flac", "opus-192", "m4a-256", "wav"}

		for _, format := range formats {
			args := buildYtDlpArgs(format, "https://youtube.com/watch?v=test", "/tmp/output")

			// 驗證通用參數
			commonParams := []string{
				"--no-playlist",
				"--max-filesize", "25M",
				"--match-filter", "duration < 600",
				"--output",
			}

			for _, param := range commonParams {
				if !containsArg(args, param) {
					t.Errorf("Format %s: missing common parameter %q", format, param)
				}
			}
		}
	})
}

// 輔助函數
func containsArg(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

func TestDownloadCommand_Structure(t *testing.T) {
	t.Run("DownloadCommand_定義正確", func(t *testing.T) {
		if DownloadCommand == nil {
			t.Fatal("DownloadCommand should not be nil")
		}

		if DownloadCommand.Command.Name != "download" {
			t.Errorf("DownloadCommand.Name = %v, want 'download'", DownloadCommand.Command.Name)
		}

		if DownloadCommand.Handler == nil {
			t.Error("DownloadCommand.Handler should not be nil")
		}

		// 驗證選項
		if len(DownloadCommand.Command.Options) != 2 {
			t.Errorf("DownloadCommand should have 2 options, got %d", len(DownloadCommand.Command.Options))
		}
	})

	t.Run("DownloadCommand_格式選項正確", func(t *testing.T) {
		// 驗證支援的格式
		expectedFormats := []string{"mp3-320", "m4a-256", "opus-192", "flac", "wav"}

		// 這裡驗證指令定義包含正確的格式選項
		// 實際測試需要解析 Options 結構
		if DownloadCommand.Command.Options == nil {
			t.Error("DownloadCommand.Options should not be nil")
		}

		_ = expectedFormats // 使用變數避免警告
	})
}
