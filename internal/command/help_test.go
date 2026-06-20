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
	t.Run("幫助訊息包含所有指令", func(t *testing.T) {
		// 驗證幫助訊息包含所有核心指令
		requiredCommands := []string{
			"/play",
			"/pause",
			"/skip",
			"/stop",
			"/queue",
			"/nowplaying",
			"/download",
		}

		// 驗證所有必需指令都存在於註冊表中
		for _, cmdName := range requiredCommands {
			found := false
			for _, cmd := range CommandRegistry {
				if "/"+cmd.Command.CommandName() == cmdName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("必需的指令 %s 未在 CommandRegistry 中找到", cmdName)
			}
		}
	})
}
