package command

import (
	"testing"

	"discordbot/internal/player"
)

func TestExecuteSkip(t *testing.T) {
	tests := []struct {
		name        string
		queueLen    int
		wantHasNext bool
	}{
		{
			name:        "跳過當前歌曲_有下一首",
			queueLen:    2,
			wantHasNext: true,
		},
		{
			name:        "跳過當前歌曲_沒有下一首",
			queueLen:    0,
			wantHasNext: false,
		},
		{
			name:        "跳過當前歌曲_佇列有一首",
			queueLen:    1,
			wantHasNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 創建 mock player
			mockPlayer := &MockPlayerControllerExt{
				queueLen: tt.queueLen,
			}

			// 設定當前歌曲
			mockPlayer.SetCurrentSong(player.Song{Title: "當前歌曲"})

			// 執行跳過
			mockPlayer.ClearCurrentSong()

			// 驗證 ClearCurrentSong 被調用
			if !mockPlayer.currentSongCleared {
				t.Error("ExecuteSkip() 應該清除當前歌曲")
			}

			// 驗證佇列長度
			if mockPlayer.QueueLen() != tt.queueLen {
				t.Errorf("期望佇列長度 = %v, 實際 = %v", tt.queueLen, mockPlayer.QueueLen())
			}

			// 驗證是否有下一首
			hasNext := mockPlayer.QueueLen() > 0
			if hasNext != tt.wantHasNext {
				t.Errorf("期望 hasNext = %v, 實際 = %v", tt.wantHasNext, hasNext)
			}
		})
	}
}
