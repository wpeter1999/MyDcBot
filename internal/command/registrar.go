package command

import "github.com/bwmarrin/discordgo"

// RegisterCommands 將 CommandRegistry 中的指令註冊到 Discord
func RegisterCommands(registrar CommandRegistrar, appID, guildID string) ([]*discordgo.ApplicationCommand, map[string]func(*discordgo.Session, *discordgo.InteractionCreate), error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(CommandRegistry))
	handlers := make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate), len(CommandRegistry))

	for _, cmd := range CommandRegistry {
		created, err := registrar.ApplicationCommandCreate(appID, guildID, cmd.Command)
		if err != nil {
			return nil, nil, err
		}
		registeredCommands = append(registeredCommands, created)
		handlers[cmd.Command.Name] = cmd.Handler
	}

	return registeredCommands, handlers, nil
}

// HandleInteraction 根據指令名稱分派 handler
func HandleInteraction(handlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate), s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if i.Type != discordgo.InteractionApplicationCommand {
		return false
	}

	name := i.ApplicationCommandData().Name
	handler, ok := handlers[name]
	if !ok {
		return false
	}

	handler(s, i)
	return true
}
