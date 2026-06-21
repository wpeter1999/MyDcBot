package command

import (
	"testing"

	"discordbot/internal/player"

	"github.com/disgoorg/snowflake/v2"
)

// TestLoopCommand_Definition 測試 loop 指令定義
func TestLoopCommand_Definition(t *testing.T) {
	if LoopCommand.Command.Name != "loop" {
		t.Errorf("expected command name 'loop', got %q", LoopCommand.Command.Name)
	}

	if LoopCommand.Command.Description == "" {
		t.Error("command description should not be empty")
	}

	if LoopCommand.Handler == nil {
		t.Error("command handler should not be nil")
	}
}

// TestLoopCommandHandler_Integration 測試 loop 指令處理器（整合測試）
func TestLoopCommandHandler_Integration(t *testing.T) {
	// 設定測試環境
	manager := player.NewManager(50)
	SetMusicService(NewDefaultMusicService(manager))

	guildID := snowflake.ID(123456789)
	userID := snowflake.ID(987654321)

	// 建立測試播放器
	guildPlayer := musicService.GetOrCreatePlayer(guildID.String())

	// 測試切換循環模式
	tests := []struct {
		name         string
		expectedMode player.LoopMode
	}{
		{"第一次切換應為單曲循環一次", player.LoopSingleOnce},
		{"第二次切換應為單曲無限循環", player.LoopSingleInfinite},
		{"第三次切換應為關閉", player.LoopOff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newMode := guildPlayer.ToggleLoopMode()
			if newMode != tt.expectedMode {
				t.Errorf("expected mode %v, got %v", tt.expectedMode, newMode)
			}
		})
	}

	// 清理
	_ = userID // 避免未使用變數警告
}
