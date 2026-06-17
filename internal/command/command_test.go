package command

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// ==================== 共用測試工具 ====================

// fakeRegistrar 用於測試的假註冊器，模擬 Discord API 成功註冊
type fakeRegistrar struct {
	created []*discordgo.ApplicationCommand
	calls   int
}

func (f *fakeRegistrar) ApplicationCommandCreate(appID, guildID string, command *discordgo.ApplicationCommand, _ ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error) {
	f.calls++
	created := *command
	created.ID = "created-" + command.Name
	f.created = append(f.created, &created)
	return &created, nil
}

// failingRegistrar 用於測試註冊失敗的情況
type failingRegistrar struct{}

func (f *failingRegistrar) ApplicationCommandCreate(appID, guildID string, command *discordgo.ApplicationCommand, _ ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error) {
	return nil, errors.New("register failed")
}

// TestRespond_CallsResponder 測試 respond 會把回應內容交給 respondToInteraction
func TestRespond_CallsResponder(t *testing.T) {
	originalResponder := respondToInteraction
	t.Cleanup(func() {
		respondToInteraction = originalResponder
	})

	var gotSession *discordgo.Session
	var gotInteraction *discordgo.InteractionCreate
	var gotContent string
	respondToInteraction = func(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
		gotSession = s
		gotInteraction = i
		gotContent = content
	}

	session := &discordgo.Session{}
	interaction := &discordgo.InteractionCreate{}
	respond(session, interaction, "hello")

	if gotSession != session {
		t.Fatal("respond 應傳遞原本的 Discord session")
	}
	if gotInteraction != interaction {
		t.Fatal("respond 應傳遞原本的 interaction")
	}
	if gotContent != "hello" {
		t.Fatalf("respond 應傳遞回應內容 hello，實際為 %q", gotContent)
	}
}
