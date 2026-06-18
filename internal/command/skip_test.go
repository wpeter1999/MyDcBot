package command

import (
	"testing"

	"discordbot/internal/player"

	"github.com/bwmarrin/discordgo"
)

// TestSkipCommand_IsRegistered 測試 SkipCommand 是否已註冊。
func TestSkipCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "skip" {
			found = true
			break
		}
	}
	if !found {
		t.Error("SkipCommand 未註冊到 CommandRegistry")
	}
}

// TestSkipCommand_SkipsCurrentSong 測試 /skip 會跳過目前播放的歌曲。
func TestSkipCommand_SkipsCurrentSong(t *testing.T) {
	originalService := GetMusicService()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	fakePlayer := fakeService.GetOrCreatePlayer("guild-1").(*fakePlayerController)
	song := player.Song{Title: "Test Song"}
	fakePlayer.currentSong = &song

	var gotContent string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		gotContent = content
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: "guild-1",
			Type:    discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "skip",
			},
		},
	}

	skipCommandHandler(&discordgo.Session{}, interaction)

	if !fakePlayer.skipCalled {
		t.Error("skipCommandHandler 應該呼叫 player.Skip()")
	}
	if gotContent == "" {
		t.Fatal("skipCommandHandler 應該回應訊息")
	}
}

// TestSkipCommand_RespondsWhenNothingPlaying 測試沒有播放時 /skip 的回應。
func TestSkipCommand_RespondsWhenNothingPlaying(t *testing.T) {
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
				Name: "skip",
			},
		},
	}

	skipCommandHandler(&discordgo.Session{}, interaction)

	if gotContent != "目前沒有播放任何歌曲。" {
		t.Fatalf("沒有播放時應回應「目前沒有播放任何歌曲。」，實際回應：%q", gotContent)
	}
}
