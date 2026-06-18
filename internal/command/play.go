package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// PlayCommand 定義 /play 指令。
var PlayCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "播放 YouTube 音樂（輸入 URL 或搜尋關鍵字）",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "YouTube URL 或搜尋關鍵字",
				Required:    true,
			},
		},
	},
	Handler: playCommandHandler,
}

// playCommandHandler 處理 /play 指令，將歌曲加入播放佇列。
func playCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respond(s, i, "請提供 YouTube URL 或搜尋關鍵字。")
		return
	}

	query := options[0].StringValue()
	if query == "" {
		respond(s, i, "請提供 YouTube URL 或搜尋關鍵字。")
		return
	}

	// Phase 2: 暫時回應已收到請求
	// Phase 3 會加入實際的 YouTube 解析
	// Phase 5 會加入實際的播放功能
	message := fmt.Sprintf("✅ 已收到播放請求：`%s`\n（YouTube 解析功能將在 Phase 3 實作）", query)
	respond(s, i, message)
}
