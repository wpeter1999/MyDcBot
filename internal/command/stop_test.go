package command

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestStopCommand_IsRegistered 測試 StopCommand 是否已註冊。
func TestStopCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "stop" {
			found = true
			break
		}
	}
	if !found {
		t.Error("StopCommand 未註冊到 CommandRegistry")
	}
}

// TestStopCommand_RemovesPlayer 測試 /stop 會移除播放器。
func TestStopCommand_RemovesPlayer(t *testing.T) {
	originalService := GetMusicService()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	// 建立播放器
	_ = fakeService.GetOrCreatePlayer("guild-1")

	var gotContent string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		gotContent = content
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: "guild-1",
			Type:    discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "stop",
			},
		},
	}

	stopCommandHandler(&discordgo.Session{}, interaction)

	if _, ok := fakeService.players["guild-1"]; ok {
		t.Error("stopCommandHandler 應該移除播放器")
	}
	if gotContent == "" {
		t.Fatal("stopCommandHandler 應該回應訊息")
	}
}

// TestStopCommand_RespondsWhenNoPlayer 測試沒有播放器時 /stop 的回應。
func TestStopCommand_RespondsWhenNoPlayer(t *testing.T) {
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
				Name: "stop",
			},
		},
	}

	stopCommandHandler(&discordgo.Session{}, interaction)

	if gotContent != "目前沒有正在播放的內容。" {
		t.Fatalf("沒有播放器時應回應「目前沒有正在播放的內容。」，實際回應：%q", gotContent)
	}
}
