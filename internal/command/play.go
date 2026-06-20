package command

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"discordbot/internal/player"
	"discordbot/internal/youtube"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// PlayCommand 定義 /play 指令。
var PlayCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "play",
		Description: "播放 YouTube 音樂（輸入 URL 或搜尋關鍵字）",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
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
func playCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	// 立即 defer 回應（避免 3 秒超時）
	if err := event.DeferCreateMessage(false); err != nil {
		log.Printf("failed to defer response: %v", err)
		return
	}

	if !validateServices(event) {
		return
	}

	query := getQueryFromEvent(event)
	if query == "" {
		updateResponse(event, "請提供 YouTube URL 或搜尋關鍵字。")
		return
	}

	song, err := resolveSong(query, event.User().ID.String())
	if err != nil {
		updateResponse(event, fmt.Sprintf("❌ 無法解析查詢：%v", err))
		return
	}

	guildID, channelID, ok := getVoiceContext(event)
	if !ok {
		updateResponse(event, "⚠️ 你必須先加入語音頻道才能播放")
		return
	}

	// 檢查是否為播放清單
	if IsPlaylistURL(query) && strings.HasPrefix(query, "http") {
		handlePlaylist(event, query, guildID, channelID)
		return
	}

	// 處理單曲
	handleSingleSong(event, song, guildID, channelID)
}

// validateServices 驗證必要的服務是否已初始化
func validateServices(event *events.ApplicationCommandInteractionCreate) bool {
	if musicService == nil {
		updateResponse(event, "音樂服務尚未初始化。")
		return false
	}
	if youtubeResolver == nil {
		updateResponse(event, "YouTube 解析服務尚未初始化。")
		return false
	}
	return true
}

// getQueryFromEvent 從事件中取得查詢字串
func getQueryFromEvent(event *events.ApplicationCommandInteractionCreate) string {
	data := event.SlashCommandInteractionData()
	return data.String("query")
}

// resolveSong 解析歌曲
func resolveSong(query, userID string) (*player.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("開始解析 YouTube 查詢: %s", query)
	song, err := youtubeResolver.Resolve(ctx, query)
	if err != nil {
		log.Printf("YouTube 解析失敗: %v", err)
		return nil, err
	}

	log.Printf("YouTube 解析成功: %s", song.Title)
	song.RequestedBy = userID
	return &song, nil
}

// getVoiceContext 取得語音頻道上下文
func getVoiceContext(event *events.ApplicationCommandInteractionCreate) (snowflake.ID, snowflake.ID, bool) {
	guildID := *event.GuildID()
	voiceState, ok := event.Client().Caches().VoiceState(guildID, event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		return 0, 0, false
	}
	return guildID, *voiceState.ChannelID, true
}

// handlePlaylist 處理播放清單
func handlePlaylist(event *events.ApplicationCommandInteractionCreate, query string, guildID, channelID snowflake.ID) {
	log.Printf("檢測到播放清單 URL，提取播放清單...")
	entries, err := ExtractPlaylist(query)
	if err != nil {
		log.Printf("提取播放清單失敗: %v", err)
		updateResponse(event, fmt.Sprintf("❌ 提取播放清單失敗：%v", err))
		return
	}

	if len(entries) == 0 {
		updateResponse(event, "❌ 播放清單是空的")
		return
	}

	// 加入佇列並播放第一首
	guildPlayer := musicService.GetOrCreatePlayer(guildID.String())
	enqueuePlaylistEntries(guildPlayer, entries, event.User().ID.String())

	err = JoinVoiceAndPlayWithYtDlp(event.Client(), guildID, channelID, entries[0].URL)
	if err != nil {
		log.Printf("播放失敗: %v", err)
		updateResponse(event, fmt.Sprintf("❌ 播放失敗：%v", err))
		return
	}

	// 設定第一首為當前播放
	setFirstSongAsCurrent(guildPlayer, entries[0], event.User().ID.String())

	log.Printf("成功開始播放")
	// 播放成功，更新語音頻道狀態
	go UpdateVoiceChannelStatus(event.Client(), channelID, entries[0].Title)
	message := buildPlaylistMessage(entries)
	updateResponseWithControlButton(event, message)
}

// enqueuePlaylistEntries 將播放清單項目加入佇列（跳過第一首）
func enqueuePlaylistEntries(guildPlayer PlayerController, entries []PlaylistEntry, userID string) {
	for i, entry := range entries {
		if i == 0 {
			continue // 跳過第一首，因為會直接播放
		}
		songToAdd := player.Song{
			Title:       entry.Title,
			URL:         entry.URL,
			StreamURL:   "",
			RequestedBy: userID,
		}
		if err := guildPlayer.Enqueue(songToAdd); err != nil {
			log.Printf("加入佇列失敗: %s - %v", entry.Title, err)
		} else {
			log.Printf("加入佇列: %s", entry.Title)
		}
	}
}

// setFirstSongAsCurrent 設定第一首歌為當前播放
func setFirstSongAsCurrent(guildPlayer PlayerController, entry PlaylistEntry, userID string) {
	firstSong := player.Song{
		Title:       entry.Title,
		URL:         entry.URL,
		RequestedBy: userID,
	}
	guildPlayer.SetCurrentSong(firstSong)
}

// buildPlaylistMessage 構建播放清單訊息
func buildPlaylistMessage(entries []PlaylistEntry) string {
	message := fmt.Sprintf("📋 **播放清單載入成功！**\n\n🎵 共 **%d** 首歌曲已加入佇列\n▶️ **正在播放：** %s\n\n", len(entries), entries[0].Title)
	message += "**歌曲清單：**\n"

	maxDisplay := 10
	displayCount := len(entries)
	if displayCount > maxDisplay {
		displayCount = maxDisplay
	}

	for i := 0; i < displayCount; i++ {
		prefix := ""
		if i == 0 {
			prefix = "▶️ "
		}
		message += fmt.Sprintf("%s%d. %s\n", prefix, i+1, entries[i].Title)
	}

	if len(entries) > maxDisplay {
		message += fmt.Sprintf("... 還有 %d 首歌曲", len(entries)-maxDisplay)
	}

	return message
}

// handleSingleSong 處理單曲播放
func handleSingleSong(event *events.ApplicationCommandInteractionCreate, song *player.Song, guildID, channelID snowflake.ID) {
	guildPlayer := musicService.GetOrCreatePlayer(guildID.String())

	// 檢查是否真的有歌曲正在播放
	isPlaying, _, _ := GetPlayerState(guildID)

	if err := guildPlayer.Enqueue(*song); err != nil {
		log.Printf("加入佇列失敗: %v", err)
		updateResponse(event, fmt.Sprintf("❌ 無法加入佇列：%v", err))
		return
	}
	log.Printf("歌曲已加入佇列: %s", song.Title)

	// 如果已經在播放，只回應已加入佇列
	if isPlaying {
		message := fmt.Sprintf("✅ 已加入佇列：**%s**", song.Title)
		updateResponseWithControlButton(event, message)
		return
	}

	// 沒有在播放，從佇列取出並開始播放
	firstSong, ok := guildPlayer.Dequeue()
	if ok {
		guildPlayer.SetCurrentSong(firstSong)
		song = &firstSong
	}

	log.Printf("嘗試加入語音頻道並播放...")
	if err := playWithFallback(event.Client(), guildID, channelID, song); err != nil {
		message := fmt.Sprintf("✅ 已加入佇列：**%s**\n⚠️ 播放失敗：%v", song.Title, err)
		updateResponse(event, message)
		return
	}

	log.Printf("成功開始播放")
	message := fmt.Sprintf("✅ 正在播放：**%s**", song.Title)
	updateResponseWithControlButton(event, message)
}

// playWithFallback 嘗試播放，失敗時使用 SoundCloud 備用
func playWithFallback(client bot.Client, guildID, channelID snowflake.ID, song *player.Song) error {
	err := JoinVoiceAndPlayWithYtDlp(client, guildID, channelID, song.URL)
	if err != nil {
		// yt-dlp 失敗，嘗試 SoundCloud 備用
		log.Printf("yt-dlp 失敗 (%v)，嘗試 SoundCloud", err)
		searchQuery := "scsearch:" + song.Title
		err = JoinVoiceAndPlay(client, guildID, channelID, searchQuery)
		if err != nil {
			log.Printf("播放失敗: %v", err)
			return err
		}
	}

	// 播放成功，更新語音頻道狀態
	go UpdateVoiceChannelStatus(client, channelID, song.Title)
	return nil
}

// updateResponse 更新 deferred response 的輔助函式
func updateResponse(event *events.ApplicationCommandInteractionCreate, content string) {
	_, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
		Content: &content,
	})
	if err != nil {
		log.Printf("failed to update response: %v", err)
	}
}

// updateResponseWithControlButton 更新 deferred response 並附加控制面板按鈕（使用 Embed）
func updateResponseWithControlButton(event *events.ApplicationCommandInteractionCreate, content string) {
	// 創建簡單的 Embed
	embed := discord.NewEmbedBuilder().
		SetColor(0x5865F2).
		SetDescription(content).
		Build()

	components := []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.ButtonComponent{
				Style:    discord.ButtonStylePrimary,
				CustomID: ButtonShowPanel,
				Label:    "🎵 音樂控制面板",
				Emoji:    &discord.ComponentEmoji{Name: "🎛️"},
			},
		},
	}

	embeds := []discord.Embed{embed}
	emptyContent := ""

	_, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
		Content:    &emptyContent,
		Embeds:     &embeds,
		Components: &components,
	})
	if err != nil {
		log.Printf("failed to update response: %v", err)
	}
}
