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

// FormatNowPlaying 格式化當前播放訊息（核心業務邏輯）
// 返回：(訊息內容, 是否有歌曲正在播放)
func FormatNowPlaying(guildPlayer PlayerController) (string, bool) {
	song, ok := guildPlayer.CurrentSong()
	if !ok {
		return "目前沒有播放任何歌曲。", false
	}

	message := fmt.Sprintf("🎵 正在播放：**%s**\n🔗 %s", song.Title, song.URL)
	return message, true
}

// nowPlayingCommandHandler 處理 /nowplaying 指令，顯示目前播放的歌曲資訊。
func nowPlayingCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	if musicService == nil {
		respond(event, "音樂服務尚未初始化。")
		return
	}

	guildID := event.GuildID().String()
	player := musicService.GetOrCreatePlayer(guildID)

	message, hasSong := FormatNowPlaying(player)
	if !hasSong {
		respond(event, message)
		return
	}

	RespondWithControlButton(event, message)
}
