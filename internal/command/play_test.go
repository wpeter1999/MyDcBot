package command

import (
	"context"
	"errors"
	"testing"

	"discordbot/internal/player"

	"github.com/bwmarrin/discordgo"
)

// fakeYouTubeResolver 是測試用的 YouTube Resolver。
type fakeYouTubeResolver struct {
	song player.Song
	err  error
}

func (r *fakeYouTubeResolver) Resolve(ctx context.Context, query string) (player.Song, error) {
	if r.err != nil {
		return player.Song{}, r.err
	}
	return r.song, nil
}

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

// TestPlayCommand_ResolvesAndEnqueues 測試 /play 會解析查詢並加入佇列。
func TestPlayCommand_ResolvesAndEnqueues(t *testing.T) {
	originalService := GetMusicService()
	originalResolver := GetYouTubeResolver()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		SetYouTubeResolver(originalResolver)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	fakeResolver := &fakeYouTubeResolver{
		song: player.Song{
			Title:     "Test Song",
			URL:       "https://youtube.com/watch?v=test",
			StreamURL: "https://example.test/stream",
		},
	}
	SetYouTubeResolver(fakeResolver)

	var gotContent string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		gotContent = content
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: "guild-1",
			Type:    discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{ID: "user-1"},
			},
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

	fakePlayer := fakeService.GetOrCreatePlayer("guild-1").(*fakePlayerController)
	if !fakePlayer.enqueueCalled {
		t.Error("playCommandHandler 應該呼叫 Enqueue")
	}

	if len(fakePlayer.queue) != 1 {
		t.Fatalf("佇列應有 1 首歌，實際有 %d 首", len(fakePlayer.queue))
	}

	if fakePlayer.queue[0].Title != "Test Song" {
		t.Errorf("加入佇列的歌曲標題應為 'Test Song'，實際為 %q", fakePlayer.queue[0].Title)
	}

	if gotContent == "" {
		t.Fatal("playCommandHandler 應該回應訊息")
	}
}

// TestPlayCommand_ReturnsErrorWhenResolverFails 測試 resolver 失敗時的錯誤處理。
func TestPlayCommand_ReturnsErrorWhenResolverFails(t *testing.T) {
	originalService := GetMusicService()
	originalResolver := GetYouTubeResolver()
	originalResponder := respondToInteraction
	defer func() {
		SetMusicService(originalService)
		SetYouTubeResolver(originalResolver)
		respondToInteraction = originalResponder
	}()

	fakeService := newFakeMusicService()
	SetMusicService(fakeService)

	fakeResolver := &fakeYouTubeResolver{
		err: errors.New("resolver failed"),
	}
	SetYouTubeResolver(fakeResolver)

	var gotContent string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		gotContent = content
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: "guild-1",
			Type:    discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{ID: "user-1"},
			},
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
		t.Fatal("resolver 失敗時應該回應錯誤訊息")
	}
}
