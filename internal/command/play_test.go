package command

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestPlayCommand_IsRegistered 測試 PlayCommand 是否已註冊。
func TestPlayCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "play" {
			found = true
			break
		}
	}
	if !found {
		t.Error("PlayCommand 未註冊到 CommandRegistry")
	}
}

// TestPlayCommand_HasRequiredQueryOption 測試 /play 是否有必填的 query 參數。
func TestPlayCommand_HasRequiredQueryOption(t *testing.T) {
	if len(PlayCommand.Command.Options) == 0 {
		t.Fatal("PlayCommand 應該有至少一個參數")
	}

	opt := PlayCommand.Command.Options[0]
	if opt.Name != "query" {
		t.Errorf("第一個參數應為 'query'，實際為 '%s'", opt.Name)
	}
	if !opt.Required {
		t.Error("query 參數應該是必填的")
	}
	if opt.Type != discordgo.ApplicationCommandOptionString {
		t.Error("query 參數應該是 String 類型")
	}
}

// TestPlayCommand_RespondsWithQuery 測試 /play 會回應收到的查詢。
func TestPlayCommand_RespondsWithQuery(t *testing.T) {
	originalService := GetMusicService()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	var gotContent string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		gotContent = content
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: "guild-1",
			Type:    discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "play",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "query",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "test song",
					},
				},
			},
		},
	}

	playCommandHandler(&discordgo.Session{}, interaction)

	if gotContent == "" {
		t.Fatal("playCommandHandler 應該回應訊息")
	}
}
