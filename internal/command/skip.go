package command

import "github.com/bwmarrin/discordgo"

// SkipCommand 定義 /skip 指令。
var SkipCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "skip",
		Description: "跳過目前播放的歌曲",
	},
	Handler: skipCommandHandler,
}

// skipCommandHandler 處理 /skip 指令，跳過目前播放的歌曲。
func skipCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	player := musicService.GetOrCreatePlayer(i.GuildID)

	if _, ok := player.CurrentSong(); !ok {
		respond(s, i, "目前沒有播放任何歌曲。")
		return
	}

	if player.Skip() {
		respond(s, i, "⏭️ 已跳過目前的歌曲。")
	} else {
		respond(s, i, "無法跳過歌曲。")
	}
}
