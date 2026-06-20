package command

import (
	"testing"

	"discordbot/internal/player"
)

func TestNowPlayingCommand(t *testing.T) {
	t.Run("NowPlayingCommand_定義正確", func(t *testing.T) {
		if NowPlayingCommand == nil {
			t.Fatal("NowPlayingCommand should not be nil")
		}

		if NowPlayingCommand.Command.Name != "nowplaying" {
			t.Errorf("NowPlayingCommand.Name = %v, want 'nowplaying'", NowPlayingCommand.Command.Name)
		}

		if NowPlayingCommand.Handler == nil {
			t.Error("NowPlayingCommand.Handler should not be nil")
		}
	})
}

func TestFormatNowPlayingLogic(t *testing.T) {
	tests := []struct {
		name         string
		currentSong  *player.Song
		wantContains []string
		wantHasSong  bool
	}{
		{
			name: "有歌曲正在播放",
			currentSong: &player.Song{
				Title:       "測試歌曲",
				URL:         "https://youtube.com/watch?v=test123",
				RequestedBy: "user123",
			},
			wantContains: []string{"測試歌曲", "🎵"},
			wantHasSong:  true,
		},
		{
			name:         "沒有歌曲播放",
			currentSong:  nil,
			wantContains: []string{"目前沒有播放任何歌曲"},
			wantHasSong:  false,
		},
		{
			name: "歌曲包含特殊字符",
			currentSong: &player.Song{
				Title:       "Test Song - (Official Video) [HD]",
				URL:         "https://youtube.com/watch?v=abc",
				RequestedBy: "user456",
			},
			wantContains: []string{"Test Song", "Official Video", "[HD]"},
			wantHasSong:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlayer := &MockPlayerControllerExt{
				currentSong: tt.currentSong,
			}

			got, hasSong := FormatNowPlaying(mockPlayer)

			if hasSong != tt.wantHasSong {
				t.Errorf("FormatNowPlaying() hasSong = %v, want %v", hasSong, tt.wantHasSong)
			}

			for _, want := range tt.wantContains {
				if !contains(got, want) {
					t.Errorf("FormatNowPlaying() missing '%s' in output:\n%s", want, got)
				}
			}

			// 驗證基本結構
			if tt.wantHasSong && got == "" {
				t.Error("FormatNowPlaying() should not return empty string when song exists")
			}
		})
	}
}
