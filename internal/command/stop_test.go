package command

import (
	"testing"
)

func TestStopCommand(t *testing.T) {
	t.Run("StopCommand_定义正确", func(t *testing.T) {
		if StopCommand == nil {
			t.Fatal("StopCommand should not be nil")
		}

		if StopCommand.Command.Name != "stop" {
			t.Errorf("StopCommand.Name = %v, want 'stop'", StopCommand.Command.Name)
		}

		if StopCommand.Handler == nil {
			t.Error("StopCommand.Handler should not be nil")
		}
	})
}

func TestExecuteStop(t *testing.T) {
	t.Run("停止播放应该清空佇列", func(t *testing.T) {
		// Note: 需要 mock bot.Client 和 PlayerController
		// 这里提供测试框架

		// mockClient := &MockBotClient{}
		// mockPlayer := &MockPlayerController{}
		// guildID := snowflake.ID(123456789)

		// err := ExecuteStop(mockClient, guildID, mockPlayer)

		// if err != nil {
		// 	t.Errorf("ExecuteStop() error = %v, want nil", err)
		// }

		// 验证 Stop 被调用
		// if !mockPlayer.stopped {
		// 	t.Error("ExecuteStop() should call Stop()")
		// }
	})
}
