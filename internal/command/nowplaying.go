package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// NowPlayingCommand 定義 /nowplaying 指令。
var NowPlayingCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "nowplaying",
		Description: "顯示目前正在播放的歌曲",
	},
	Handler: nowPlayingCommandHandler,
}

// nowPlayingCommandHandler 處理 /nowplaying 指令，顯示目前播放的歌曲資訊。
func nowPlayingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	player := musicService.GetOrCreatePlayer(i.GuildID)
	song, ok := player.CurrentSong()

	if !ok {
		respond(s, i, "目前沒有播放任何歌曲。")
		return
	}

	message := fmt.Sprintf("🎵 正在播放：**%s**\n🔗 %s", song.Title, song.URL)
	respond(s, i, message)
}
