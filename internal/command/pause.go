package command

import "github.com/bwmarrin/discordgo"

// PauseCommand 定義 /pause 指令。
var PauseCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "pause",
		Description: "暫停或繼續播放",
	},
	Handler: pauseCommandHandler,
}

// pauseCommandHandler 處理 /pause 指令，切換暫停/繼續狀態。
func pauseCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	player := musicService.GetOrCreatePlayer(i.GuildID)

	if _, ok := player.CurrentSong(); !ok {
		respond(s, i, "目前沒有播放任何歌曲。")
		return
	}

	paused := player.TogglePause()
	if paused {
		respond(s, i, "⏸️ 已暫停播放。")
	} else {
		respond(s, i, "▶️ 已繼續播放。")
	}
}
