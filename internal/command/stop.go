package command

import (
	"context"
	"fmt"
	"log"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// StopCommand 定義 /stop 指令。
var StopCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "stop",
		Description: "停止播放並清空佇列",
	},
	Handler: stopCommandHandler,
}

// ExecuteStop 停止播放並清空佇列（核心業務邏輯）
func ExecuteStop(client bot.Client, guildID snowflake.ID, guildPlayer PlayerController) error {
	// 停止播放並離開語音頻道
	err := StopPlayback(client, guildID)
	if err != nil {
		return fmt.Errorf("停止失敗：%w", err)
	}

	// 清空佇列
	guildPlayer.Stop()

	// 離開語音頻道
	if err := client.UpdateVoiceState(context.Background(), guildID, nil, false, false); err != nil {
		log.Printf("離開語音頻道時出錯: %v", err)
	}

	return nil
}

// stopCommandHandler 處理 /stop 指令，停止播放並清空佇列。
func stopCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	if musicService == nil {
		respond(event, "音樂服務尚未初始化。")
		return
	}

	guildID := *event.GuildID()
	guildIDStr := guildID.String()
	guildPlayer := musicService.GetOrCreatePlayer(guildIDStr)

	err := ExecuteStop(event.Client(), guildID, guildPlayer)
	if err != nil {
		respond(event, fmt.Sprintf("❌ %v", err))
		return
	}

	// 移除 player
	musicService.RemovePlayer(guildIDStr)
	RespondWithControlButton(event, "⏹️ 已停止播放並清空佇列。")
}
