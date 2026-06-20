package command

import (
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

func TestExecuteSkip(t *testing.T) {
	tests := []struct {
		name        string
		queueLen    int
		wantHasNext bool
	}{
		{
			name:        "跳过当前歌曲_有下一首",
			queueLen:    2,
			wantHasNext: true,
		},
		{
			name:        "跳过当前歌曲_没有下一首",
			queueLen:    0,
			wantHasNext: false,
		},
		{
			name:        "跳过当前歌曲_佇列有一首",
			queueLen:    1,
			wantHasNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock player
			mockPlayer := &MockPlayerControllerExt{
				queueLen: tt.queueLen,
			}

			_ = mockPlayer
			_ = snowflake.ID(123456789)

			// Note: 实际测试需要 mock Lavalink client
			// 这里提供测试框架

			// hasNext := ExecuteSkip(guildID, mockPlayer)

			// if hasNext != tt.wantHasNext {
			// 	t.Errorf("ExecuteSkip() = %v, want %v", hasNext, tt.wantHasNext)
			// }

			// 验证 ClearCurrentSong 被调用
			// if !mockPlayer.currentSongCleared {
			// 	t.Error("ExecuteSkip() should clear current song")
			// }
		})
	}
}

// Note: MockPlayerControllerExt 定义在 test_helpers.go 中
