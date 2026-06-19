package command

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// StopCommand 定義 /stop 指令。
var StopCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "stop",
		Description: "停止播放並清空佇列",
	},
	Handler: stopCommandHandler,
}

// stopCommandHandler 處理 /stop 指令，停止播放並清空佇列。
func stopCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	if musicService == nil {
		respond(event, "音樂服務尚未初始化。")
		return
	}

	guildID := *event.GuildID()
	guildIDStr := guildID.String()

	// 停止播放並離開語音頻道
	err := StopPlayback(event.Client(), guildID)
	if err != nil {
		respond(event, fmt.Sprintf("❌ 停止失敗：%v", err))
		return
	}

	removed := musicService.RemovePlayer(guildIDStr)
	if removed {
		respond(event, "⏹️ 已停止播放並清空佇列。")
	} else {
		respond(event, "⏹️ 已停止播放。")
	}
}
