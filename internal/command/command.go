package command

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// BotCommand 封裝一個 Discord slash command 的定義與處理函式
type BotCommand struct {
	Command *discordgo.ApplicationCommand
	Handler func(*discordgo.Session, *discordgo.InteractionCreate)
}

// CommandRegistrar 抽象 Discord 指令註冊介面
type CommandRegistrar interface {
	ApplicationCommandCreate(appID, guildID string, command *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error)
}

var respondToInteraction = func(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("failed to respond: %v", err)
	}
}

// respond 統一回應使用者訊息的輔助函式
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	respondToInteraction(s, i, content)
}
