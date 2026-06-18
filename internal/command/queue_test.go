package command

import (
	"testing"

	"discordbot/internal/player"

	"github.com/bwmarrin/discordgo"
)

// TestQueueCommand_IsRegistered 測試 QueueCommand 是否已註冊。
func TestQueueCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "queue" {
			found = true
			break
		}
	}
	if !found {
		t.Error("QueueCommand 未註冊到 CommandRegistry")
	}
}

// TestQueueCommand_ShowsQueuedSongs 測試 /queue 會顯示已加入佇列的歌曲。
func TestQueueCommand_ShowsQueuedSongs(t *testing.T) {
	originalService := GetMusicService()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	fakePlayer := fakeService.GetOrCreatePlayer("guild-1").(*fakePlayerController)
	fakePlayer.queue = []player.Song{
		{Title: "Song A"},
		{Title: "Song B"},
		{Title: "Song C"},
	}

	var gotContent string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		gotContent = content
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: "guild-1",
			Type:    discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "queue",
			},
		},
	}

	queueCommandHandler(&discordgo.Session{}, interaction)

	if gotContent == "" {
		t.Fatal("queueCommandHandler 應該回應訊息")
	}
	if gotContent == "佇列目前是空的。" {
		t.Fatal("應該顯示佇列內容，而不是空的訊息")
	}
}

// TestQueueCommand_RespondsWhenEmpty 測試佇列為空時 /queue 的回應。
func TestQueueCommand_RespondsWhenEmpty(t *testing.T) {
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
				Name: "queue",
			},
		},
	}

	queueCommandHandler(&discordgo.Session{}, interaction)

	if gotContent != "佇列目前是空的。" {
		t.Fatalf("佇列為空時應回應「佇列目前是空的。」，實際回應：%q", gotContent)
	}
}
