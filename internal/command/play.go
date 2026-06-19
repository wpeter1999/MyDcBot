package command

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"discordbot/internal/player"
	"discordbot/internal/youtube"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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

	if musicService == nil {
		updateResponse(event, "音樂服務尚未初始化。")
		return
	}

	if youtubeResolver == nil {
		updateResponse(event, "YouTube 解析服務尚未初始化。")
		return
	}

	data := event.SlashCommandInteractionData()
	query := data.String("query")

	if query == "" {
		updateResponse(event, "請提供 YouTube URL 或搜尋關鍵字。")
		return
	}

	// 使用 context 控制 resolver 超時
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("開始解析 YouTube 查詢: %s", query)
	song, err := youtubeResolver.Resolve(ctx, query)
	if err != nil {
		log.Printf("YouTube 解析失敗: %v", err)
		message := fmt.Sprintf("❌ 無法解析查詢：%v", err)
		updateResponse(event, message)
		return
	}
	log.Printf("YouTube 解析成功: %s", song.Title)

	// 設定 RequestedBy
	song.RequestedBy = event.User().ID.String()

	guildID := *event.GuildID()
	guildIDStr := guildID.String()

	// 取得使用者所在的語音頻道
	voiceState, ok := event.Client().Caches().VoiceState(guildID, event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		message := fmt.Sprintf("⚠️ 你必須先加入語音頻道才能播放")
		updateResponse(event, message)
		return
	}

	channelID := *voiceState.ChannelID

	// 檢查是否為播放清單 URL
	if IsPlaylistURL(query) && strings.HasPrefix(query, "http") {
		log.Printf("檢測到播放清單 URL，提取播放清單...")
		entries, err := ExtractPlaylist(query)
		if err != nil {
			log.Printf("提取播放清單失敗: %v", err)
			message := fmt.Sprintf("❌ 提取播放清單失敗：%v", err)
			updateResponse(event, message)
			return
		}

		if len(entries) == 0 {
			updateResponse(event, "❌ 播放清單是空的")
			return
		}

		// 將所有歌曲加入佇列（播放清單模式）
		guildPlayer := musicService.GetOrCreatePlayer(guildIDStr)
		for i, entry := range entries {
			songToAdd := player.Song{
				Title:       entry.Title,
				URL:         entry.URL,
				StreamURL:   "",
				RequestedBy: event.User().ID.String(),
			}
			// 跳過第一首，因為會直接播放
			if i == 0 {
				continue
			}
			if err := guildPlayer.Enqueue(songToAdd); err != nil {
				log.Printf("加入佇列失敗: %s - %v", entry.Title, err)
			} else {
				log.Printf("加入佇列: %s", entry.Title)
			}
		}

		// 播放第一首
		err = JoinVoiceAndPlayWithYtDlp(event.Client(), guildID, channelID, entries[0].URL)
		if err != nil {
			log.Printf("播放失敗: %v", err)
			message := fmt.Sprintf("❌ 播放失敗：%v", err)
			updateResponse(event, message)
			return
		}

		// 設定第一首為當前播放歌曲
		if len(entries) > 0 {
			firstSong := player.Song{
				Title:       entries[0].Title,
				URL:         entries[0].URL,
				RequestedBy: event.User().ID.String(),
			}
			guildPlayer.SetCurrentSong(firstSong)
		}

		log.Printf("成功開始播放")

		// 構建播放清單訊息
		playlistMessage := fmt.Sprintf("📋 **播放清單載入成功！**\n\n🎵 共 **%d** 首歌曲已加入佇列\n▶️ **正在播放：** %s\n\n", len(entries), entries[0].Title)

		// 顯示歌曲清單
		playlistMessage += "**歌曲清單：**\n"
		maxDisplay := 10
		if len(entries) > maxDisplay {
			for i := 0; i < maxDisplay; i++ {
				if i == 0 {
					playlistMessage += fmt.Sprintf("▶️ %d. %s\n", i+1, entries[i].Title)
				} else {
					playlistMessage += fmt.Sprintf("%d. %s\n", i+1, entries[i].Title)
				}
			}
			playlistMessage += fmt.Sprintf("... 還有 %d 首歌曲", len(entries)-maxDisplay)
		} else {
			for i, entry := range entries {
				if i == 0 {
					playlistMessage += fmt.Sprintf("▶️ %d. %s\n", i+1, entry.Title)
				} else {
					playlistMessage += fmt.Sprintf("%d. %s\n", i+1, entry.Title)
				}
			}
		}

		updateResponse(event, playlistMessage)
		return
	}

	// 單曲模式：加入佇列
	guildPlayer := musicService.GetOrCreatePlayer(guildIDStr)

	// 檢查是否已經有歌曲在播放
	_, hasCurrentSong := guildPlayer.CurrentSong()

	if err := guildPlayer.Enqueue(song); err != nil {
		log.Printf("加入佇列失敗: %v", err)
		message := fmt.Sprintf("❌ 無法加入佇列：%v", err)
		updateResponse(event, message)
		return
	}
	log.Printf("歌曲已加入佇列: %s", song.Title)

	// 如果沒有正在播放的歌曲，設定為當前播放並開始播放
	if !hasCurrentSong {
		// 從佇列取出第一首（剛剛加入的）
		firstSong, ok := guildPlayer.Dequeue()
		if ok {
			guildPlayer.SetCurrentSong(firstSong)
			song = firstSong // 使用從佇列取出的歌曲
		}
	} else {
		// 已經有歌曲在播放，只回應已加入佇列
		message := fmt.Sprintf("✅ 已加入佇列：**%s**", song.Title)
		updateResponse(event, message)
		return
	}

	// 嘗試加入語音頻道並啟動播放（如果尚未播放）
	log.Printf("嘗試加入語音頻道並播放...")

	err = JoinVoiceAndPlayWithYtDlp(event.Client(), guildID, channelID, song.URL)
	if err != nil {
		// yt-dlp 失敗，嘗試 SoundCloud 備用
		log.Printf("yt-dlp 失敗 (%v)，嘗試 SoundCloud", err)
		searchQuery := "scsearch:" + song.Title
		err = JoinVoiceAndPlay(event.Client(), guildID, channelID, searchQuery)
		if err != nil {
			// 兩者都失敗
			log.Printf("播放失敗: %v", err)
			message := fmt.Sprintf("✅ 已加入佇列：**%s**\n⚠️ 播放失敗：%v", song.Title, err)
			updateResponse(event, message)
			return
		}
	}

	log.Printf("成功開始播放")
	message := fmt.Sprintf("✅ 正在播放：**%s**", song.Title)
	updateResponse(event, message)
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
