package command

import "github.com/bwmarrin/discordgo"

// StopCommand 定義 /stop 指令。
var StopCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "stop",
		Description: "停止播放並清空佇列",
	},
	Handler: stopCommandHandler,
}

// stopCommandHandler 處理 /stop 指令，停止播放並清空佇列。
func stopCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	// 停止播放迴圈
	StopPlayback(i.GuildID)

	removed := musicService.RemovePlayer(i.GuildID)
	if removed {
		respond(s, i, "⏹️ 已停止播放並清空佇列。")
	} else {
		respond(s, i, "目前沒有正在播放的內容。")
	}
}
