package main

import (
	"context"
	"flag"
	"log"

	"discordbot/internal/bot"
	"discordbot/internal/config"

	"github.com/disgoorg/snowflake/v2"
)

func main() {
	// 添加清理指令的 flag
	cleanCommands := flag.Bool("clean", false, "清理所有舊指令後退出")
	cleanGlobal := flag.Bool("global", false, "清理全域指令而非 Guild 指令")
	listCommands := flag.Bool("list", false, "列出所有現有指令")
	flag.Parse()

	// 載入配置
	cfg := config.Load()

	// 創建 Bot
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	ctx := context.Background()
	appID := b.Client.ApplicationID()

	log.Printf("Bot Application ID: %s", appID)

	// 如果是列出指令模式
	if *listCommands {
		var guildID snowflake.ID
		if !*cleanGlobal && cfg.GuildID != "" {
			guildID = snowflake.MustParse(cfg.GuildID)
			log.Printf("Listing Guild commands for Guild ID: %s", cfg.GuildID)
		} else {
			log.Println("Listing Global commands")
		}

		commands, err := bot.ListCommands(ctx, b.Client, appID, guildID)
		if err != nil {
			log.Fatalf("Failed to list commands: %v", err)
		}

		log.Printf("Found %d commands:", len(commands))
		for i, cmd := range commands {
			log.Printf("  %d. %s (ID: %s)", i+1, cmd.Name(), cmd.ID())
		}
		return
	}

	// 如果是清理模式
	if *cleanCommands {
		log.Println("🧹 Cleaning up old commands...")

		var guildID snowflake.ID
		if !*cleanGlobal && cfg.GuildID != "" {
			guildID = snowflake.MustParse(cfg.GuildID)
			log.Printf("Cleaning Guild commands for Guild ID: %s", cfg.GuildID)
		} else {
			log.Println("Cleaning Global commands")
		}

		if err := bot.CleanupCommands(ctx, b.Client, appID, guildID); err != nil {
			log.Fatalf("Failed to cleanup commands: %v", err)
		}

		log.Println("✅ Cleanup completed!")
		return
	}

	// 正常啟動模式
	if err := b.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	log.Println("Bot is running. Press Ctrl+C to exit.")

	// 阻塞主 goroutine
	select {}
}
