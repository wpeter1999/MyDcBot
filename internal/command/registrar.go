package command

import (
	"context"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// RegisterCommands 將 CommandRegistry 中的指令註冊到 Discord
func RegisterCommands(client bot.Client, appID snowflake.ID, guildID snowflake.ID) ([]snowflake.ID, map[string]InteractionHandler, error) {
	ctx := context.Background()

	// 準備指令定義
	commands := make([]discord.ApplicationCommandCreate, 0, len(CommandRegistry))
	handlers := make(map[string]InteractionHandler, len(CommandRegistry))

	for _, cmd := range CommandRegistry {
		commands = append(commands, cmd.Command)
		handlers[cmd.Command.CommandName()] = cmd.Handler
	}

	// 註冊指令
	var registeredCommands []discord.ApplicationCommand
	var err error

	if guildID != 0 {
		// Guild commands (faster for testing)
		registeredCommands, err = client.Rest().SetGuildCommands(appID, guildID, commands)
	} else {
		// Global commands (takes up to 1 hour to propagate)
		registeredCommands, err = client.Rest().SetGlobalCommands(appID, commands)
	}

	if err != nil {
		return nil, nil, err
	}

	// 提取 command IDs
	commandIDs := make([]snowflake.ID, len(registeredCommands))
	for i, cmd := range registeredCommands {
		commandIDs[i] = cmd.ID()
	}

	_ = ctx // 保留供未來使用

	return commandIDs, handlers, nil
}
