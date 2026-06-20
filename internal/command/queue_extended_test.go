package command

import (
	"testing"

	"discordbot/internal/player"
)

func TestFormatQueueDisplay_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		currentSong  *player.Song
		queueSongs   []player.Song
		wantContains []string
	}{
		{
			name: "佇列有超過10首歌",
			currentSong: &player.Song{
				Title:       "當前歌曲",
				RequestedBy: "user123",
			},
			queueSongs: []player.Song{
				{Title: "歌曲1"}, {Title: "歌曲2"}, {Title: "歌曲3"},
				{Title: "歌曲4"}, {Title: "歌曲5"}, {Title: "歌曲6"},
				{Title: "歌曲7"}, {Title: "歌曲8"}, {Title: "歌曲9"},
				{Title: "歌曲10"}, {Title: "歌曲11"}, {Title: "歌曲12"},
			},
			wantContains: []string{"當前歌曲", "歌曲1", "歌曲10", "還有 2 首歌曲"},
		},
		{
			name: "佇列剛好10首歌",
			currentSong: &player.Song{
				Title:       "當前歌曲",
				RequestedBy: "user123",
			},
			queueSongs: []player.Song{
				{Title: "歌曲1"}, {Title: "歌曲2"}, {Title: "歌曲3"},
				{Title: "歌曲4"}, {Title: "歌曲5"}, {Title: "歌曲6"},
				{Title: "歌曲7"}, {Title: "歌曲8"}, {Title: "歌曲9"},
				{Title: "歌曲10"},
			},
			wantContains: []string{"當前歌曲", "歌曲1", "歌曲10", "11 首歌曲"},
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
					t.Errorf("FormatQueueDisplay() missing '%s' in output", want)
				}
			}
		})
	}
}

func TestFormatQueueDisplay_QueueCount(t *testing.T) {
	t.Run("計算總歌曲數包含當前播放", func(t *testing.T) {
		mockPlayer := &MockPlayerControllerExt{
			currentSong: &player.Song{Title: "當前歌曲"},
			queue: []player.Song{
				{Title: "歌曲1"},
				{Title: "歌曲2"},
			},
		}

		got := FormatQueueDisplay(mockPlayer)

		// 應該顯示 3 首歌曲（1 當前 + 2 佇列）
		if !contains(got, "3 首歌曲") {
			t.Error("FormatQueueDisplay() should count current song in total")
		}
	})

	t.Run("計算總歌曲數不包含當前播放_當無歌曲時", func(t *testing.T) {
		mockPlayer := &MockPlayerControllerExt{
			currentSong: nil,
			queue: []player.Song{
				{Title: "歌曲1"},
				{Title: "歌曲2"},
			},
		}

		got := FormatQueueDisplay(mockPlayer)

		// 應該顯示 2 首歌曲
		if !contains(got, "2 首歌曲") {
			t.Error("FormatQueueDisplay() should count only queue when no current song")
		}
	})
}
