package bot

import (
	"context"
	"log"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// ListCommands 列出所有現有指令
func ListCommands(ctx context.Context, client bot.Client, appID snowflake.ID, guildID snowflake.ID) ([]discord.ApplicationCommand, error) {
	var commands []discord.ApplicationCommand
	var err error

	if guildID != 0 {
		// 獲取 Guild 指令
		commands, err = client.Rest().GetGuildCommands(appID, guildID, false)
	} else {
		// 獲取全域指令
		commands, err = client.Rest().GetGlobalCommands(appID, false)
	}

	return commands, err
}

// CleanupCommands 清理所有舊指令
func CleanupCommands(ctx context.Context, client bot.Client, appID snowflake.ID, guildID snowflake.ID) error {
	log.Printf("Fetching existing commands...")

	var commands []discord.ApplicationCommand
	var err error

	if guildID != 0 {
		// 獲取 Guild 指令
		commands, err = client.Rest().GetGuildCommands(appID, guildID, false)
	} else {
		// 獲取全域指令
		commands, err = client.Rest().GetGlobalCommands(appID, false)
	}

	if err != nil {
		return err
	}

	log.Printf("Found %d existing commands", len(commands))

	// 刪除所有指令
	for _, cmd := range commands {
		log.Printf("Deleting command: %s (ID: %d)", cmd.Name(), cmd.ID())

		if guildID != 0 {
			err = client.Rest().DeleteGuildCommand(appID, guildID, cmd.ID())
		} else {
			err = client.Rest().DeleteGlobalCommand(appID, cmd.ID())
		}

		if err != nil {
			log.Printf("Failed to delete command %s: %v", cmd.Name(), err)
		} else {
			log.Printf("✅ Successfully deleted command: %s", cmd.Name())
		}
	}

	return nil
}
