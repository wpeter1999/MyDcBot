package command

import (
	"context"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// InteractionHandler 定義處理 slash command 互動的函式類型
type InteractionHandler func(*events.ApplicationCommandInteractionCreate)

// BotCommand 封裝一個 Discord slash command 的定義與處理函式
type BotCommand struct {
	Command discord.SlashCommandCreate
	Handler InteractionHandler
}

var respondToInteraction = func(event *events.ApplicationCommandInteractionCreate, content string) {
	if err := event.CreateMessage(discord.MessageCreate{
		Content: content,
	}); err != nil {
		log.Printf("failed to respond: %v", err)
	}
}

// respond 統一回應使用者訊息的輔助函式
func respond(event *events.ApplicationCommandInteractionCreate, content string) {
	respondToInteraction(event, content)
}

// deferAndFollowUp 延遲回應並發送 follow-up 訊息的輔助函式
func deferAndFollowUp(event *events.ApplicationCommandInteractionCreate, content string) {
	ctx := context.Background()

	// Defer response
	if err := event.DeferCreateMessage(false); err != nil {
		log.Printf("failed to defer response: %v", err)
		return
	}

	// Send follow-up
	if _, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
		Content: &content,
	}); err != nil {
		log.Printf("failed to send follow-up: %v", err)
	}

	_ = ctx // 保留供未來使用
}
