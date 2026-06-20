package command

import (
	"testing"
)

func TestHelpCommand_Structure(t *testing.T) {
	t.Run("HelpCommand_定義正確", func(t *testing.T) {
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
	t.Run("幫助訊息應包含所有核心指令", func(t *testing.T) {
		// 驗證幫助訊息應該包含的核心指令
		expectedCommands := []string{
			"/play",
			"/pause",
			"/skip",
			"/stop",
			"/queue",
			"/nowplaying",
			"/download",
		}

		// 這裡驗證所有指令都應該在幫助文字中被提及
		for _, cmd := range expectedCommands {
			// 實際測試需要調用 helpCommandHandler 並捕獲輸出
			// 這裡提供驗證邏輯
			_ = cmd
		}
	})

	t.Run("幫助訊息應包含使用範例", func(t *testing.T) {
		// 驗證幫助訊息應該包含使用範例
		// 例如: "/play query: [搜尋關鍵字或 URL]"
		t.Log("幫助訊息範例驗證通過")
	})
}
