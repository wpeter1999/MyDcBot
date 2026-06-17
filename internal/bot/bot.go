package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"discordbot/internal/command"
	"discordbot/internal/config"

	"github.com/bwmarrin/discordgo"
)

var (
	openSession = func(s *discordgo.Session) error {
		return s.Open()
	}
	closeSession = func(s *discordgo.Session) error {
		return s.Close()
	}
	registerBotCommands      = command.RegisterCommands
	deleteApplicationCommand = func(s *discordgo.Session, appID, guildID, commandID string) error {
		return s.ApplicationCommandDelete(appID, guildID, commandID)
	}
	notifyShutdownSignal = func(stop chan<- os.Signal) {
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	}
)

// Bot 封裝 Discord Bot 實例
type Bot struct {
	Session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
	commandHandlers    map[string]func(*discordgo.Session, *discordgo.InteractionCreate)
	cfg                *config.Config
}

// New 建立新的 Bot 實例
func New(cfg *config.Config) (*Bot, error) {
	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		Session: dg,
		cfg:     cfg,
	}

	dg.AddHandler(b.interactionCreate)

	return b, nil
}

// Start 啟動 Bot
func (b *Bot) Start() error {
	err := openSession(b.Session)
	if err != nil {
		return err
	}

	// 註冊指令
	commands, handlers, err := registerBotCommands(b.Session, b.Session.State.User.ID, b.cfg.GuildID)
	if err != nil {
		return err
	}

	b.registeredCommands = commands
	b.commandHandlers = handlers

	return nil
}

// Stop 停止 Bot 並清理指令
func (b *Bot) Stop() {
	guildID := b.cfg.GuildID
	for _, cmd := range b.registeredCommands {
		if err := deleteApplicationCommand(b.Session, b.Session.State.User.ID, guildID, cmd.ID); err != nil {
			log.Printf("failed to delete command %q: %v", cmd.Name, err)
		}
	}
	closeSession(b.Session)
}

// WaitForShutdown 等待中斷訊號
func (b *Bot) WaitForShutdown() {
	stop := make(chan os.Signal, 1)
	notifyShutdownSignal(stop)
	<-stop
}

// interactionCreate 處理互動事件
func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name
	if handler, ok := b.commandHandlers[name]; ok {
		handler(s, i)
	}
}
