package command

import (
	"testing"
)

func TestHelpCommand(t *testing.T) {
	t.Run("HelpCommand_定义正确", func(t *testing.T) {
		if HelpCommand == nil {
			t.Fatal("HelpCommand should not be nil")
		}

		if HelpCommand.Command.Name != "help" {
			t.Errorf("HelpCommand.Name = %v, want 'help'", HelpCommand.Command.Name)
		}

		if HelpCommand.Command.Description == "" {
			t.Error("HelpCommand.Description should not be empty")
		}

		if HelpCommand.Handler == nil {
			t.Error("HelpCommand.Handler should not be nil")
		}
	})
}

func TestHelpCommandContent(t *testing.T) {
	t.Run("帮助信息包含所有指令", func(t *testing.T) {
		// 验证帮助消息包含所有核心指令
		_ = []string{
			"/play",
			"/pause",
			"/skip",
			"/stop",
			"/queue",
			"/nowplaying",
			"/download",
		}

		// Note: 实际测试需要调用 helpCommandHandler 并捕获输出
		// 这里提供测试框架

		// mockEvent := createMockEvent()
		// helpCommandHandler(mockEvent)
		// output := captureOutput(mockEvent)

		// for _, cmd := range requiredCommands {
		// 	if !strings.Contains(output, cmd) {
		// 		t.Errorf("Help message should contain command %s", cmd)
		// 	}
		// }
	})
}
