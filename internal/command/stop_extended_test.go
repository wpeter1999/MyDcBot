package command

import (
	"testing"
)

func TestStopCommand_Structure(t *testing.T) {
	t.Run("StopCommand_定義完整", func(t *testing.T) {
		if StopCommand == nil {
			t.Fatal("StopCommand should not be nil")
		}

		if StopCommand.Command.Name != "stop" {
			t.Errorf("StopCommand.Name = %v, want 'stop'", StopCommand.Command.Name)
		}

		if StopCommand.Command.Description == "" {
			t.Error("StopCommand.Description should not be empty")
		}

		if StopCommand.Handler == nil {
			t.Error("StopCommand.Handler should not be nil")
		}
	})

	t.Run("StopCommand_無選項參數", func(t *testing.T) {
		// Stop 指令不需要任何參數
		if len(StopCommand.Command.Options) != 0 {
			t.Errorf("StopCommand should have 0 options, got %d", len(StopCommand.Command.Options))
		}
	})
}

func TestExecuteStop_Logic(t *testing.T) {
	t.Run("執行停止應該清空佇列", func(t *testing.T) {
		// ExecuteStop 的核心邏輯:
		// 1. 調用 StopPlayback (停止 Lavalink player)
		// 2. 調用 player.Stop() (清空佇列和狀態)
		// 3. 調用 UpdateVoiceState (離開語音頻道)

		// 實際測試需要 mock bot.Client 和 PlayerController
	})
}
