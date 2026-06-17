package command

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestPingCommand_Definition 測試 PingCommand 的基本定義是否正確
func TestPingCommand_Definition(t *testing.T) {
	if PingCommand.Command.Name != "ping" {
		t.Errorf("PingCommand 名稱應為 'ping'，實際為 '%s'", PingCommand.Command.Name)
	}
	if PingCommand.Command.Description == "" {
		t.Error("PingCommand 應該有描述")
	}
	if PingCommand.Handler == nil {
		t.Error("PingCommand 必須有 Handler")
	}
}

// TestPingCommand_IsRegistered 測試 PingCommand 是否已註冊到 CommandRegistry
func TestPingCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "ping" {
			found = true
			break
		}
	}
	if !found {
		t.Error("PingCommand 未註冊到 CommandRegistry")
	}
}

// TestPingCommandHandler_RespondsPong 測試 ping handler 會回覆 Pong!
func TestPingCommandHandler_RespondsPong(t *testing.T) {
	originalResponder := respondToInteraction
	t.Cleanup(func() {
		respondToInteraction = originalResponder
	})

	var got string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		got = content
	}

	pingCommandHandler(&discordgo.Session{}, &discordgo.InteractionCreate{})

	if got != "Pong!" {
		t.Fatalf("pingCommandHandler 回應應為 Pong!，實際為 %q", got)
	}
}
