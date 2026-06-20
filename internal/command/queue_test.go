package command

import (
	"testing"

	"discordbot/internal/player"
)

func TestFormatQueueDisplay(t *testing.T) {
	tests := []struct {
		name         string
		currentSong  *player.Song
		queueSongs   []player.Song
		wantContains []string
	}{
		{
			name: "显示当前播放和佇列",
			currentSong: &player.Song{
				Title:       "当前歌曲",
				RequestedBy: "user123",
			},
			queueSongs: []player.Song{
				{Title: "歌曲1", RequestedBy: "user123"},
				{Title: "歌曲2", RequestedBy: "user456"},
			},
			wantContains: []string{"当前歌曲", "歌曲1", "歌曲2", "▶️"},
		},
		{
			name:         "没有歌曲播放",
			currentSong:  nil,
			queueSongs:   []player.Song{},
			wantContains: []string{"播放佇列是空的"},
		},
		{
			name: "只有当前歌曲_无佇列",
			currentSong: &player.Song{
				Title:       "唯一歌曲",
				RequestedBy: "user123",
			},
			queueSongs:   []player.Song{},
			wantContains: []string{"唯一歌曲", "▶️"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlayer := &MockPlayerControllerExt{
				currentSong: tt.currentSong,
				queue:       tt.queueSongs,
			}

			got := FormatQueueDisplay(mockPlayer)

			for _, want := range tt.wantContains {
				if !contains(got, want) {
					t.Errorf("FormatQueueDisplay() missing '%s' in output:\n%s", want, got)
				}
			}
		})
	}
}

func TestFormatNowPlaying(t *testing.T) {
	tests := []struct {
		name         string
		currentSong  *player.Song
		wantContains []string
	}{
		{
			name: "有歌曲播放",
			currentSong: &player.Song{
				Title:       "测试歌曲",
				URL:         "https://youtube.com/watch?v=test",
				RequestedBy: "user123",
			},
			wantContains: []string{"测试歌曲", "🎵"},
		},
		{
			name:         "没有歌曲播放",
			currentSong:  nil,
			wantContains: []string{"目前沒有播放任何歌曲"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlayer := &MockPlayerControllerExt{
				currentSong: tt.currentSong,
			}

			got, _ := FormatNowPlaying(mockPlayer)

			for _, want := range tt.wantContains {
				if !contains(got, want) {
					t.Errorf("FormatNowPlaying() missing '%s' in output:\n%s", want, got)
				}
			}
		})
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
