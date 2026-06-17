package bot

import (
	"errors"
	"os"
	"testing"

	"discordbot/internal/command"
	"discordbot/internal/config"

	"github.com/bwmarrin/discordgo"
)

// TestNew 測試 New 會建立 Discord session 並保存設定
func TestNew(t *testing.T) {
	cfg := &config.Config{BotToken: "test-token", GuildID: "guild-id"}

	b, err := New(cfg)
	if err != nil {
		t.Fatalf("預期 New 不應回傳錯誤，但得到: %v", err)
	}
	if b.Session == nil {
		t.Fatal("New 應建立 Discord session")
	}
	if b.cfg != cfg {
		t.Fatal("New 應保存傳入的 config")
	}
}

// TestStart_Success 測試 Start 會開啟 session、註冊指令並保存 handlers
func TestStart_Success(t *testing.T) {
	restore := replaceBotDependencies(t)
	defer restore()

	openCalled := false
	openSession = func(s *discordgo.Session) error {
		openCalled = true
		return nil
	}

	createdCommands := []*discordgo.ApplicationCommand{{ID: "cmd-1", Name: "ping"}}
	createdHandlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {},
	}
	registerBotCommands = func(_ command.CommandRegistrar, appID, guildID string) ([]*discordgo.ApplicationCommand, map[string]func(*discordgo.Session, *discordgo.InteractionCreate), error) {
		if appID != "app-id" {
			t.Errorf("appID 應為 app-id，實際為 %q", appID)
		}
		if guildID != "guild-id" {
			t.Errorf("guildID 應為 guild-id，實際為 %q", guildID)
		}
		return createdCommands, createdHandlers, nil
	}

	b := newTestBot("app-id")
	if err := b.Start(); err != nil {
		t.Fatalf("Start 不應回傳錯誤，但得到: %v", err)
	}
	if !openCalled {
		t.Fatal("Start 應開啟 session")
	}
	if len(b.registeredCommands) != 1 || b.registeredCommands[0].ID != "cmd-1" {
		t.Fatalf("Start 應保存已註冊指令，實際為 %#v", b.registeredCommands)
	}
	if _, ok := b.commandHandlers["ping"]; !ok {
		t.Fatal("Start 應保存 command handlers")
	}
}

// TestStart_ReturnsOpenError 測試開啟 session 失敗時 Start 會回傳錯誤
func TestStart_ReturnsOpenError(t *testing.T) {
	restore := replaceBotDependencies(t)
	defer restore()

	wantErr := errors.New("open failed")
	openSession = func(s *discordgo.Session) error {
		return wantErr
	}

	b := newTestBot("app-id")
	err := b.Start()
	if !errors.Is(err, wantErr) {
		t.Fatalf("Start 應回傳 open 錯誤，實際為 %v", err)
	}
}

// TestStart_ReturnsRegisterError 測試註冊指令失敗時 Start 會回傳錯誤
func TestStart_ReturnsRegisterError(t *testing.T) {
	restore := replaceBotDependencies(t)
	defer restore()

	wantErr := errors.New("register failed")
	registerBotCommands = func(_ command.CommandRegistrar, appID, guildID string) ([]*discordgo.ApplicationCommand, map[string]func(*discordgo.Session, *discordgo.InteractionCreate), error) {
		return nil, nil, wantErr
	}

	b := newTestBot("app-id")
	err := b.Start()
	if !errors.Is(err, wantErr) {
		t.Fatalf("Start 應回傳註冊錯誤，實際為 %v", err)
	}
}

// TestStop_DeletesCommandsAndClosesSession 測試 Stop 會刪除已註冊指令並關閉 session
func TestStop_DeletesCommandsAndClosesSession(t *testing.T) {
	restore := replaceBotDependencies(t)
	defer restore()

	var deleted []string
	deleteApplicationCommand = func(_ *discordgo.Session, appID, guildID, commandID string) error {
		if appID != "app-id" {
			t.Errorf("appID 應為 app-id，實際為 %q", appID)
		}
		if guildID != "guild-id" {
			t.Errorf("guildID 應為 guild-id，實際為 %q", guildID)
		}
		deleted = append(deleted, commandID)
		return nil
	}

	closed := false
	closeSession = func(s *discordgo.Session) error {
		closed = true
		return nil
	}

	b := newTestBot("app-id")
	b.registeredCommands = []*discordgo.ApplicationCommand{
		{ID: "cmd-1", Name: "ping"},
		{ID: "cmd-2", Name: "weather"},
	}

	b.Stop()

	if len(deleted) != 2 || deleted[0] != "cmd-1" || deleted[1] != "cmd-2" {
		t.Fatalf("Stop 應刪除所有已註冊指令，實際為 %#v", deleted)
	}
	if !closed {
		t.Fatal("Stop 應關閉 session")
	}
}

// TestWaitForShutdown_ReturnsAfterSignal 測試收到關閉訊號後 WaitForShutdown 會返回
func TestWaitForShutdown_ReturnsAfterSignal(t *testing.T) {
	restore := replaceBotDependencies(t)
	defer restore()

	notifyShutdownSignal = func(stop chan<- os.Signal) {
		stop <- os.Interrupt
	}

	b := newTestBot("app-id")
	b.WaitForShutdown()
}

// TestInteractionCreate_DispatchesKnownCommand 測試 interactionCreate 會分派已知 slash command
func TestInteractionCreate_DispatchesKnownCommand(t *testing.T) {
	handled := false
	b := &Bot{
		commandHandlers: map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
			"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {
				handled = true
			},
		},
	}

	b.interactionCreate(&discordgo.Session{}, applicationCommandInteraction("ping"))

	if !handled {
		t.Fatal("interactionCreate 應分派已知指令")
	}
}

// TestInteractionCreate_IgnoresUnknownCommand 測試 interactionCreate 會忽略未知指令
func TestInteractionCreate_IgnoresUnknownCommand(t *testing.T) {
	handled := false
	b := &Bot{
		commandHandlers: map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
			"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {
				handled = true
			},
		},
	}

	b.interactionCreate(&discordgo.Session{}, applicationCommandInteraction("unknown"))

	if handled {
		t.Fatal("interactionCreate 不應分派未知指令")
	}
}

// TestInteractionCreate_IgnoresNonCommand 測試 interactionCreate 會忽略非 slash command 互動
func TestInteractionCreate_IgnoresNonCommand(t *testing.T) {
	handled := false
	b := &Bot{
		commandHandlers: map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
			"ping": func(_ *discordgo.Session, _ *discordgo.InteractionCreate) {
				handled = true
			},
		},
	}

	b.interactionCreate(&discordgo.Session{}, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{Type: discordgo.InteractionPing},
	})

	if handled {
		t.Fatal("interactionCreate 不應分派非 ApplicationCommand 互動")
	}
}

func replaceBotDependencies(t *testing.T) func() {
	t.Helper()

	originalOpenSession := openSession
	originalCloseSession := closeSession
	originalRegisterBotCommands := registerBotCommands
	originalDeleteApplicationCommand := deleteApplicationCommand
	originalNotifyShutdownSignal := notifyShutdownSignal

	openSession = func(s *discordgo.Session) error { return nil }
	closeSession = func(s *discordgo.Session) error { return nil }
	registerBotCommands = func(_ command.CommandRegistrar, appID, guildID string) ([]*discordgo.ApplicationCommand, map[string]func(*discordgo.Session, *discordgo.InteractionCreate), error) {
		return nil, nil, nil
	}
	deleteApplicationCommand = func(s *discordgo.Session, appID, guildID, commandID string) error { return nil }
	notifyShutdownSignal = func(stop chan<- os.Signal) { stop <- os.Interrupt }

	return func() {
		openSession = originalOpenSession
		closeSession = originalCloseSession
		registerBotCommands = originalRegisterBotCommands
		deleteApplicationCommand = originalDeleteApplicationCommand
		notifyShutdownSignal = originalNotifyShutdownSignal
	}
}

func newTestBot(appID string) *Bot {
	state := discordgo.NewState()
	state.User = &discordgo.User{ID: appID}

	return &Bot{
		Session: &discordgo.Session{
			State: state,
		},
		cfg: &config.Config{GuildID: "guild-id"},
	}
}

func applicationCommandInteraction(name string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{Name: name},
		},
	}
}
