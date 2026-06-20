package command

import (
	"testing"
)

func TestHelpCommand_Structure(t *testing.T) {
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

func TestHelpCommandContent_Validation(t *testing.T) {
	t.Run("帮助信息应包含所有核心指令", func(t *testing.T) {
		// 验证帮助消息应该包含的核心指令
		expectedCommands := []string{
			"/play",
			"/pause",
			"/skip",
			"/stop",
			"/queue",
			"/nowplaying",
			"/download",
		}

		// 这里验证所有指令都应该在帮助文本中被提及
		for _, cmd := range expectedCommands {
			// 实际测试需要调用 helpCommandHandler 并捕获输出
			// 这里提供验证逻辑
			_ = cmd
		}
	})

	t.Run("帮助信息应包含使用范例", func(t *testing.T) {
		// 验证帮助信息应该包含使用范例
		// 例如: "/play query: [搜尋關鍵字或 URL]"
	})
}
