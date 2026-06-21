package command

import (
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// ShuffleCommand 定義 shuffle 指令
var ShuffleCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "shuffle",
		Description: "打亂佇列中的歌曲順序",
	},
	Handler: shuffleCommandHandler,
}

// shuffleCommandHandler 處理 /shuffle 指令
func shuffleCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	// 檢查使用者是否在語音頻道
	guildID, _, ok := getVoiceContext(event)
	if !ok {
		updateResponse(event, "❌ 你必須在語音頻道中才能使用此指令")
		return
	}

	// 取得播放器
	guildPlayer := musicService.GetOrCreatePlayer(guildID.String())

	// 檢查佇列是否為空
	queueLen := guildPlayer.QueueLen()
	if queueLen == 0 {
		updateResponse(event, "⚠️ 佇列中沒有歌曲可以打亂")
		return
	}

	// 打亂佇列
	guildPlayer.Shuffle()

	// 回應使用者
	message := fmt.Sprintf("🔀 **已打亂佇列**\n共 %d 首歌曲已隨機排序", queueLen)

	if err := event.CreateMessage(discord.MessageCreate{
		Content: message,
	}); err != nil {
		log.Printf("[Shuffle] 回應失敗: %v", err)
	}
}
