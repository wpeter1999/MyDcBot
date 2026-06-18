package command

import (
	"context"
	"fmt"
	"time"

	"discordbot/internal/youtube"

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

var youtubeResolver youtube.Resolver

// SetYouTubeResolver 設定全域 YouTube Resolver（測試時可注入 fake）。
func SetYouTubeResolver(resolver youtube.Resolver) {
	youtubeResolver = resolver
}

// GetYouTubeResolver 取得目前的 YouTube Resolver。
func GetYouTubeResolver() youtube.Resolver {
	return youtubeResolver
}

// playCommandHandler 處理 /play 指令，將歌曲加入播放佇列。
func playCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if musicService == nil {
		respond(s, i, "音樂服務尚未初始化。")
		return
	}

	if youtubeResolver == nil {
		respond(s, i, "YouTube 解析服務尚未初始化。")
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

	// 使用 context 控制 resolver 超時
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	song, err := youtubeResolver.Resolve(ctx, query)
	if err != nil {
		message := fmt.Sprintf("❌ 無法解析查詢：%v", err)
		respond(s, i, message)
		return
	}

	// 設定 RequestedBy
	song.RequestedBy = i.Member.User.ID

	player := musicService.GetOrCreatePlayer(i.GuildID)
	if err := player.Enqueue(song); err != nil {
		message := fmt.Sprintf("❌ 無法加入佇列：%v", err)
		respond(s, i, message)
		return
	}

	// 嘗試加入語音頻道並啟動播放（如果尚未播放）
	if err := JoinVoiceAndPlay(s, i.GuildID, i.Member.User.ID, player); err != nil {
		// 無法加入語音頻道（例如：使用者不在語音頻道）
		message := fmt.Sprintf("✅ 已加入佇列：**%s**\n⚠️ %v", song.Title, err)
		respond(s, i, message)
		return
	}

	message := fmt.Sprintf("✅ 已加入佇列：**%s**", song.Title)
	respond(s, i, message)
}
