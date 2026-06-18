package command

import (
	"testing"

	"discordbot/internal/player"

	"github.com/bwmarrin/discordgo"
)

// fakeMusicService 是測試用的 MusicService 實作。
type fakeMusicService struct {
	players map[string]*fakePlayerController
}

// fakePlayerController 是測試用的 PlayerController 實作。
type fakePlayerController struct {
	guildID       string
	queue         []player.Song
	currentSong   *player.Song
	paused        bool
	stopped       bool
	skipCalled    bool
	stopCalled    bool
	enqueueCalled bool
}

func newFakeMusicService() *fakeMusicService {
	return &fakeMusicService{
		players: make(map[string]*fakePlayerController),
	}
}

func (s *fakeMusicService) GetOrCreatePlayer(guildID string) PlayerController {
	if p, ok := s.players[guildID]; ok {
		return p
	}
	p := &fakePlayerController{guildID: guildID, queue: []player.Song{}}
	s.players[guildID] = p
	return p
}

func (s *fakeMusicService) RemovePlayer(guildID string) bool {
	if _, ok := s.players[guildID]; ok {
		delete(s.players, guildID)
		return true
	}
	return false
}

func (p *fakePlayerController) Enqueue(song player.Song) error {
	p.enqueueCalled = true
	if len(p.queue) >= 50 {
		return player.ErrQueueFull
	}
	p.queue = append(p.queue, song)
	return nil
}

func (p *fakePlayerController) Skip() bool {
	p.skipCalled = true
	return !p.stopped
}

func (p *fakePlayerController) TogglePause() bool {
	p.paused = !p.paused
	return p.paused
}

func (p *fakePlayerController) IsPaused() bool {
	return p.paused
}

func (p *fakePlayerController) Stop() {
	p.stopCalled = true
	p.stopped = true
	p.queue = nil
	p.currentSong = nil
}

func (p *fakePlayerController) QueueSnapshot() []player.Song {
	snapshot := make([]player.Song, len(p.queue))
	copy(snapshot, p.queue)
	return snapshot
}

func (p *fakePlayerController) CurrentSong() (player.Song, bool) {
	if p.currentSong == nil {
		return player.Song{}, false
	}
	return *p.currentSong, true
}

func (p *fakePlayerController) GuildID() string {
	return p.guildID
}

// TestNowPlayingCommand_IsRegistered 測試 NowPlayingCommand 是否已註冊。
func TestNowPlayingCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "nowplaying" {
			found = true
			break
		}
	}
	if !found {
		t.Error("NowPlayingCommand 未註冊到 CommandRegistry")
	}
}

// TestNowPlayingCommand_ShowsCurrentSong 測試 /nowplaying 會顯示目前播放的歌曲。
func TestNowPlayingCommand_ShowsCurrentSong(t *testing.T) {
	originalService := GetMusicService()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	fakePlayer := fakeService.GetOrCreatePlayer("guild-1").(*fakePlayerController)
	song := player.Song{Title: "Test Song", URL: "https://example.test"}
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
				Name: "nowplaying",
			},
		},
	}

	nowPlayingCommandHandler(&discordgo.Session{}, interaction)

	if gotContent == "" {
		t.Fatal("nowPlayingCommandHandler 應該回應訊息")
	}
	if gotContent == "目前沒有播放任何歌曲。" {
		t.Fatal("應該顯示目前播放的歌曲，而不是沒有播放的訊息")
	}
}

// TestNowPlayingCommand_ShowsNothingPlayingWhenEmpty 測試沒有播放時 /nowplaying 的回應。
func TestNowPlayingCommand_ShowsNothingPlayingWhenEmpty(t *testing.T) {
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
				Name: "nowplaying",
			},
		},
	}

	nowPlayingCommandHandler(&discordgo.Session{}, interaction)

	if gotContent != "目前沒有播放任何歌曲。" {
		t.Fatalf("沒有播放時應回應「目前沒有播放任何歌曲。」，實際回應：%q", gotContent)
	}
}
