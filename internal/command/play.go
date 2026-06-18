package command

import (
	"context"
	"fmt"
	"log"
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
	// 立即 defer 回應（避免 3 秒超時）
	deferResponse(s, i)

	if musicService == nil {
		followUp(s, i, "音樂服務尚未初始化。")
		return
	}

	if youtubeResolver == nil {
		followUp(s, i, "YouTube 解析服務尚未初始化。")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		followUp(s, i, "請提供 YouTube URL 或搜尋關鍵字。")
		return
	}

	query := options[0].StringValue()
	if query == "" {
		followUp(s, i, "請提供 YouTube URL 或搜尋關鍵字。")
		return
	}

	// 使用 context 控制 resolver 超時
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	song, err := youtubeResolver.Resolve(ctx, query)
	if err != nil {
		message := fmt.Sprintf("❌ 無法解析查詢：%v", err)
		followUp(s, i, message)
		return
	}

	// 設定 RequestedBy
	song.RequestedBy = i.Member.User.ID

	player := musicService.GetOrCreatePlayer(i.GuildID)
	if err := player.Enqueue(song); err != nil {
		message := fmt.Sprintf("❌ 無法加入佇列：%v", err)
		followUp(s, i, message)
		return
	}

	// 嘗試加入語音頻道並啟動播放（如果尚未播放）
	if err := JoinVoiceAndPlay(s, i.GuildID, i.Member.User.ID, player); err != nil {
		// 無法加入語音頻道（例如：使用者不在語音頻道）
		message := fmt.Sprintf("✅ 已加入佇列：**%s**\n⚠️ 語音功能目前不可用（Discord DAVE 協議支援開發中）", song.Title)
		followUp(s, i, message)
		return
	}

	message := fmt.Sprintf("✅ 已加入佇列：**%s**\n⚠️ 語音功能目前不可用（Discord DAVE 協議支援開發中）", song.Title)
	followUp(s, i, message)
}

// deferResponse 可覆寫的函數，用於發送 deferred response（測試時可注入 fake）。
var deferResponse = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("failed to defer response: %v", err)
	}
}

// followUp 可覆寫的函數，用於發送 follow-up 訊息（測試時可注入 fake）。
var followUp = func(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	if err != nil {
		log.Printf("failed to send follow-up: %v", err)
	}
}
