package command

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// NowPlayingCommand 定義 /nowplaying 指令。
var NowPlayingCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "nowplaying",
		Description: "顯示目前正在播放的歌曲",
	},
	Handler: nowPlayingCommandHandler,
}

// nowPlayingCommandHandler 處理 /nowplaying 指令，顯示目前播放的歌曲資訊。
func nowPlayingCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	if musicService == nil {
		respond(event, "音樂服務尚未初始化。")
		return
	}

	guildID := event.GuildID().String()
	player := musicService.GetOrCreatePlayer(guildID)
	song, ok := player.CurrentSong()

	if !ok {
		respond(event, "目前沒有播放任何歌曲。")
		return
	}

	message := fmt.Sprintf("🎵 正在播放：**%s**\n🔗 %s", song.Title, song.URL)
	RespondWithControlButton(event, message)
}
