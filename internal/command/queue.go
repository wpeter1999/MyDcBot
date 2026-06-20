package command

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// QueueCommand 定義 /queue 指令。
var QueueCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "queue",
		Description: "顯示目前的播放佇列",
	},
	Handler: queueCommandHandler,
}

// queueCommandHandler 處理 /queue 指令，顯示目前佇列中的歌曲。
func queueCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	if musicService == nil {
		respond(event, "音樂服務尚未初始化。")
		return
	}

	guildID := event.GuildID().String()
	guildPlayer := musicService.GetOrCreatePlayer(guildID)

	// 取得當前播放的歌曲
	currentSong, hasCurrentSong := guildPlayer.CurrentSong()

	// 取得佇列
	songs := guildPlayer.QueueSnapshot()
	totalSongs := len(songs)
	if hasCurrentSong {
		totalSongs++ // 加上正在播放的歌曲
	}

	if totalSongs == 0 {
		respond(event, "📜 播放佇列是空的")
		return
	}

	var message string
	if hasCurrentSong {
		message = fmt.Sprintf("📜 播放佇列 (%d 首歌曲)\n\n▶️ **正在播放：**\n%s\n", totalSongs, currentSong.Title)

		if len(songs) > 0 {
			message += "\n**接下來：**\n"
		}
	} else {
		message = fmt.Sprintf("📜 播放佇列 (%d 首歌曲)\n\n", totalSongs)
	}

	// 顯示接下來的歌曲
	if len(songs) > 0 {
		maxDisplay := 10
		if len(songs) <= maxDisplay {
			for i, song := range songs {
				message += fmt.Sprintf("%d. %s\n", i+1, song.Title)
			}
		} else {
			for i := 0; i < maxDisplay; i++ {
				message += fmt.Sprintf("%d. %s\n", i+1, songs[i].Title)
			}
			message += fmt.Sprintf("... 還有 %d 首歌曲", len(songs)-maxDisplay)
		}
	}

	RespondWithControlButton(event, message)
}
