package command

import "github.com/bwmarrin/discordgo"

// PingCommand 定義 /ping 指令
var PingCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	},
	Handler: pingCommandHandler,
}

func pingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond(s, i, "Pong!")
}
