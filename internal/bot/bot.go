package bot

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discordbot/internal/command"
	"discordbot/internal/config"
	"discordbot/internal/player"
	"discordbot/internal/youtube"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

var (
	openGateway = func(ctx context.Context, client bot.Client) error {
		return client.OpenGateway(ctx)
	}
	closeGateway = func(ctx context.Context, client bot.Client) {
		client.Close(ctx)
	}
	registerBotCommands = command.RegisterCommands
	deleteGuildCommand  = func(ctx context.Context, client bot.Client, appID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID) error {
		return client.Rest().DeleteGuildCommand(appID, guildID, commandID)
	}
	notifyShutdownSignal = func(stop chan<- os.Signal) {
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	}
)

// Bot 封裝 Discord Bot 實例
type Bot struct {
	Client            bot.Client
	Lavalink          disgolink.Client
	registeredCommands []snowflake.ID
	commandHandlers    map[string]command.InteractionHandler
	cfg                *config.Config
	playerManager      *player.Manager
}

// New 建立新的 Bot 實例
func New(cfg *config.Config) (*Bot, error) {
	// 初始化 player manager（佇列容量 50）
	playerManager := player.NewManager(50)

	// 初始化 YouTube resolver
	youtubeRunner := youtube.NewExecCommandRunner()
	youtubeResolver := youtube.NewResolver(youtubeRunner)

	// 設定全域服務（供指令使用）
	command.SetMusicService(command.NewDefaultMusicService(playerManager))
	command.SetYouTubeResolver(youtubeResolver)

	b := &Bot{
		cfg:           cfg,
		playerManager: playerManager,
	}

	// 建立 disgo client
	client, err := disgo.New(cfg.BotToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildVoiceStates,
			),
			gateway.WithPresenceOpts(
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
				gateway.WithListeningActivity("音樂 | /help"),
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(
				cache.FlagGuilds,
				cache.FlagVoiceStates,
			),
		),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: b.onApplicationCommandInteraction,
			OnGuildVoiceStateUpdate: func(event *events.GuildVoiceStateUpdate) {
				log.Printf("[Voice Event] Voice state updated for user %s in guild %s", event.VoiceState.UserID, event.VoiceState.GuildID)
				if b.Lavalink != nil && event.VoiceState.UserID == b.Client.ApplicationID() {
					// 只處理 bot 自己的 voice state 變更
					b.Lavalink.OnVoiceStateUpdate(context.Background(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
				}
			},
			OnVoiceServerUpdate: func(event *events.VoiceServerUpdate) {
				endpointStr := "nil"
				if event.Endpoint != nil {
					endpointStr = *event.Endpoint
				}
				log.Printf("[Voice Event] Voice server updated: %s", endpointStr)
				if b.Lavalink != nil && event.Endpoint != nil {
					b.Lavalink.OnVoiceServerUpdate(context.Background(), event.GuildID, event.Token, *event.Endpoint)
				}
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	b.Client = client

	// 初始化 Lavalink client
	log.Printf("[Lavalink] Initializing Lavalink client...")
	b.Lavalink = disgolink.New(client.ApplicationID())

	// 設定全域服務（供指令使用）
	command.SetLavalinkClient(b.Lavalink)

	return b, nil
}

// Start 啟動 Bot
func (b *Bot) Start() error {
	ctx := context.Background()

	err := openGateway(ctx, b.Client)
	if err != nil {
		return err
	}

	// 連線到 Lavalink
	log.Printf("[Lavalink] Connecting to Lavalink server...")
	_, err = b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     "main",
		Address:  "lavalink:2333",
		Password: "youshallnotpass",
		Secure:   false,
	})
	if err != nil {
		log.Printf("[Lavalink] Failed to connect to Lavalink: %v", err)
		return err
	}
	log.Printf("[Lavalink] Successfully connected to Lavalink")

	// 註冊 Lavalink 事件處理器
	b.Lavalink.AddListeners(&BotEventListener{bot: b})
	log.Printf("[Lavalink] Event handlers registered")

	// 註冊指令
	appID := b.Client.ApplicationID()
	var guildID snowflake.ID
	if b.cfg.GuildID != "" {
		guildID = snowflake.MustParse(b.cfg.GuildID)
	}

	commandIDs, handlers, err := registerBotCommands(b.Client, appID, guildID)
	if err != nil {
		return err
	}

	b.registeredCommands = commandIDs
	b.commandHandlers = handlers

	return nil
}

// Stop 停止 Bot 並清理指令
func (b *Bot) Stop() {
	ctx := context.Background()

	appID := b.Client.ApplicationID()
	var guildID snowflake.ID
	if b.cfg.GuildID != "" {
		guildID = snowflake.MustParse(b.cfg.GuildID)
	}

	for _, cmdID := range b.registeredCommands {
		if err := deleteGuildCommand(ctx, b.Client, appID, guildID, cmdID); err != nil {
			log.Printf("failed to delete command %d: %v", cmdID, err)
		}
	}

	closeGateway(ctx, b.Client)
}

// WaitForShutdown 等待中斷訊號
func (b *Bot) WaitForShutdown() {
	stop := make(chan os.Signal, 1)
	notifyShutdownSignal(stop)
	<-stop
}

// onApplicationCommandInteraction 處理應用程式指令互動事件
func (b *Bot) onApplicationCommandInteraction(event *events.ApplicationCommandInteractionCreate) {
	name := event.Data.CommandName()
	if handler, ok := b.commandHandlers[name]; ok {
		handler(event)
	}
}
