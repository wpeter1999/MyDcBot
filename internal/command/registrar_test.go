package command

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

// ==================== RegisterCommands 測試 ====================

// TestRegisterCommands_Success 測試 RegisterCommands 成功註冊所有指令
func TestRegisterCommands_Success(t *testing.T) {
	registrar := &fakeRegistrar{}

	created, handlers, err := RegisterCommands(registrar, "app-id", "guild-id")
	if err != nil {
		t.Fatalf("預期不應發生錯誤，但得到: %v", err)
	}

	if len(created) != len(CommandRegistry) {
		t.Fatalf("預期註冊 %d 個指令，實際註冊 %d 個", len(CommandRegistry), len(created))
	}

	if len(handlers) != len(CommandRegistry) {
		t.Fatalf("預期有 %d 個 handler，實際有 %d 個", len(CommandRegistry), len(handlers))
	}

	for _, cmd := range CommandRegistry {
		if _, ok := handlers[cmd.Command.Name]; !ok {
			t.Fatalf("找不到指令 %q 對應的 handler", cmd.Command.Name)
		}
	}
}

// TestRegisterCommands_Error 測試 RegisterCommands 在註冊失敗時正確回傳錯誤
func TestRegisterCommands_Error(t *testing.T) {
	registrar := &failingRegistrar{}

	_, _, err := RegisterCommands(registrar, "app-id", "guild-id")
	if err == nil {
		t.Fatal("預期應該回傳錯誤，但得到 nil")
	}
}

// ==================== HandleInteraction 測試 ====================

// TestHandleInteraction_DispatchesCorrectly 測試 HandleInteraction 能正確分派 handler
func TestHandleInteraction_DispatchesCorrectly(t *testing.T) {
	handled := false

	handlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {
			handled = true
		},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "ping",
			},
		},
	}

	ok := HandleInteraction(handlers, &discordgo.Session{}, interaction)
	if !ok {
		t.Fatal("預期互動應該被處理")
	}
	if !handled {
		t.Fatal("預期 handler 應該被執行")
	}
}

// TestHandleInteraction_IgnoresUnknownCommand 測試遇到未知指令時不會執行 handler
func TestHandleInteraction_IgnoresUnknownCommand(t *testing.T) {
	handlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "unknown-command",
			},
		},
	}

	ok := HandleInteraction(handlers, &discordgo.Session{}, interaction)
	if ok {
		t.Fatal("預期未知指令不應被處理")
	}
}

// TestHandleInteraction_IgnoresNonCommand 測試非 ApplicationCommand 類型的互動會被忽略
func TestHandleInteraction_IgnoresNonCommand(t *testing.T) {
	handlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionPing,
		},
	}

	ok := HandleInteraction(handlers, &discordgo.Session{}, interaction)
	if ok {
		t.Fatal("預期非 ApplicationCommand 互動不應被處理")
	}
}
