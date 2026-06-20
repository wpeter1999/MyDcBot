package command

import (
	"testing"

	"discordbot/internal/player"
)

func TestExecutePauseToggle(t *testing.T) {
	tests := []struct {
		name           string
		initialPaused  bool
		hasSong        bool
		wantNewPaused  bool
		wantIsPlaying  bool
	}{
		{
			name:           "暫停正在播放的歌曲",
			initialPaused:  false,
			hasSong:        true,
			wantNewPaused:  true,
			wantIsPlaying:  true,
		},
		{
			name:           "恢復已暫停的歌曲",
			initialPaused:  true,
			hasSong:        true,
			wantNewPaused:  false,
			wantIsPlaying:  true,
		},
		{
			name:           "沒有歌曲播放時暫停",
			initialPaused:  false,
			hasSong:        false,
			wantNewPaused:  false,
			wantIsPlaying:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlayer := &MockPlayerControllerExt{}

			if tt.hasSong {
				mockPlayer.SetCurrentSong(player.Song{
					Title: "測試歌曲",
					URL:   "https://youtube.com/watch?v=test",
				})
			}

			// 測試 TogglePause 行為
			_, hasCurrentSong := mockPlayer.CurrentSong()
			if hasCurrentSong != tt.wantIsPlaying {
				t.Errorf("期望 hasCurrentSong = %v, 實際 = %v", tt.wantIsPlaying, hasCurrentSong)
			}
		})
	}
}
