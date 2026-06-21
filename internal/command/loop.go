package command

import (
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// LoopCommand 定義 loop 指令
var LoopCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "loop",
		Description: "切換循環播放模式（關閉 → 單曲循環 → 佇列循環）",
	},
	Handler: loopCommandHandler,
}

// loopCommandHandler 處理 /loop 指令
func loopCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	// 檢查使用者是否在語音頻道
	guildID, _, ok := getVoiceContext(event)
	if !ok {
		updateResponse(event, "❌ 你必須在語音頻道中才能使用此指令")
		return
	}

	// 取得播放器
	guildPlayer := musicService.GetOrCreatePlayer(guildID.String())

	// 切換循環模式
	newMode := guildPlayer.ToggleLoopMode()

	// 建構回應訊息
	icon := newMode.Icon()
	modeName := newMode.String()

	message := fmt.Sprintf("%s **循環模式：%s**", icon, modeName)

	// 回應使用者
	if err := event.CreateMessage(discord.MessageCreate{
		Content: message,
	}); err != nil {
		log.Printf("[Loop] 回應失敗: %v", err)
	}
}

