package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// QueueCommand 定義 /queue 指令。
var QueueCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "queue",
		Description: "顯示目前的播放佇列",
	},
	Handler: queueCommandHandler,
}

// queueCommandHandler 處理 /queue 指令，顯示目前佇列中的歌曲。
func queueCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	player := musicService.GetOrCreatePlayer(i.GuildID)
	queue := player.QueueSnapshot()

	if len(queue) == 0 {
		respond(s, i, "佇列目前是空的。")
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("📜 **播放佇列** (%d 首歌曲)\n\n", len(queue)))

	for idx, song := range queue {
		if idx >= 10 {
			builder.WriteString(fmt.Sprintf("... 還有 %d 首歌曲\n", len(queue)-10))
			break
		}
		builder.WriteString(fmt.Sprintf("%d. **%s**\n", idx+1, song.Title))
	}

	respond(s, i, builder.String())
}
