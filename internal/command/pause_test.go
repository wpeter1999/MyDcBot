package command

import (
	"testing"

	"discordbot/internal/player"

	"github.com/bwmarrin/discordgo"
)

// TestPauseCommand_IsRegistered 測試 PauseCommand 是否已註冊。
func TestPauseCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "pause" {
			found = true
			break
		}
	}
	if !found {
		t.Error("PauseCommand 未註冊到 CommandRegistry")
	}
}

// TestPauseCommand_TogglesPause 測試 /pause 會切換暫停狀態。
func TestPauseCommand_TogglesPause(t *testing.T) {
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
				Name: "pause",
			},
		},
	}

	// 第一次應該暫停
	pauseCommandHandler(&discordgo.Session{}, interaction)
	if !fakePlayer.paused {
		t.Error("第一次呼叫應該切換為 paused=true")
	}

	// 第二次應該繼續
	pauseCommandHandler(&discordgo.Session{}, interaction)
	if fakePlayer.paused {
		t.Error("第二次呼叫應該切換為 paused=false")
	}

	if gotContent == "" {
		t.Fatal("pauseCommandHandler 應該回應訊息")
	}
}

// TestPauseCommand_RespondsWhenNothingPlaying 測試沒有播放時 /pause 的回應。
func TestPauseCommand_RespondsWhenNothingPlaying(t *testing.T) {
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
				Name: "pause",
			},
		},
	}

	pauseCommandHandler(&discordgo.Session{}, interaction)

	if gotContent != "目前沒有播放任何歌曲。" {
		t.Fatalf("沒有播放時應回應「目前沒有播放任何歌曲。」，實際回應：%q", gotContent)
	}
}
