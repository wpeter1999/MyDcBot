package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

// 按鈕 Custom ID 常量
const (
	ButtonShowPanel  = "music_show_panel"
	ButtonPause      = "music_pause"
	ButtonSkip       = "music_skip"
	ButtonStop       = "music_stop"
	ButtonQueue      = "music_queue"
	ButtonNowPlaying = "music_nowplaying"
	ButtonSearch     = "music_search"
)

// RespondWithControlButton 回應訊息並附加"顯示控制面板"按鈕（使用 Embed）
func RespondWithControlButton(event *events.ApplicationCommandInteractionCreate, content string) {
	// 創建簡單的 Embed
	embed := discord.NewEmbedBuilder().
		SetColor(0x5865F2).
		SetDescription(content).
		Build()

	if err := event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.ButtonComponent{
					Style:    discord.ButtonStylePrimary,
					CustomID: ButtonShowPanel,
					Label:    "🎵 音樂控制面板",
					Emoji:    &discord.ComponentEmoji{Name: "🎛️"},
				},
			},
		},
	}); err != nil {
		log.Printf("failed to respond with control button: %v", err)
	}
}

// UpdateMessageWithFullPanel 更新訊息，顯示完整的控制面板
func UpdateMessageWithFullPanel(event *events.ComponentInteractionCreate, content string) {
	guildID := event.GuildID().String()
	player := musicService.GetOrCreatePlayer(guildID)
	song, hasSong := player.CurrentSong()

	// 檢查實際播放狀態
	isPlaying, isPaused, _ := GetPlayerState(*event.GuildID())

	// 根據播放狀態決定暫停/播放按鈕的樣式
	pauseStyle := discord.ButtonStyleSecondary

	// 如果沒有歌曲、沒有在播放、或已暫停，使用綠色按鈕
	if !hasSong || !isPlaying || isPaused {
		pauseStyle = discord.ButtonStyleSuccess
	}

	// 創建 Embed 面板
	embedBuilder := discord.NewEmbedBuilder().
		SetColor(0x5865F2). // Discord 藍色
		SetTitle("🎛️ 音樂控制面板")

	// 如果有歌曲正在播放，顯示當前播放資訊
	if hasSong {
		statusIcon := "▶️"
		statusText := "播放中"
		if isPaused {
			statusIcon = "⏸️"
			statusText = "已暫停"
		} else if !isPlaying {
			statusIcon = "⏹️"
			statusText = "已停止"
		}

		embedBuilder.AddField("🎵 正在播放", fmt.Sprintf("%s **%s**", statusIcon, song.Title), false)
		embedBuilder.AddField("📊 狀態", statusText, true)

		// 顯示佇列長度
		queueLen := player.QueueLen()
		embedBuilder.AddField("📜 佇列", fmt.Sprintf("%d 首歌曲", queueLen), true)
	} else {
		embedBuilder.SetDescription("目前沒有播放任何歌曲")
	}

	embedBuilder.SetFooter("點擊下方按鈕來控制音樂播放", "")

	embed := embedBuilder.Build()

	components := []discord.ContainerComponent{
		// 第一行：播放控制
		discord.ActionRowComponent{
			discord.ButtonComponent{
				Style:    pauseStyle,
				CustomID: ButtonPause,
				Emoji:    &discord.ComponentEmoji{Name: "⏯️"},
				Disabled: !hasSong,
			},
			discord.ButtonComponent{
				Style:    discord.ButtonStyleSecondary,
				CustomID: ButtonSkip,
				Emoji:    &discord.ComponentEmoji{Name: "⏭️"},
				Disabled: !hasSong,
			},
			discord.ButtonComponent{
				Style:    discord.ButtonStyleDanger,
				CustomID: ButtonStop,
				Emoji:    &discord.ComponentEmoji{Name: "⏹️"},
				Disabled: !hasSong,
			},
		},
		// 第二行：資訊查詢
		discord.ActionRowComponent{
			discord.ButtonComponent{
				Style:    discord.ButtonStylePrimary,
				CustomID: ButtonSearch,
				Label:    "🔍 搜尋",
			},
			discord.ButtonComponent{
				Style:    discord.ButtonStyleSuccess,
				CustomID: ButtonNowPlaying,
				Label:    "🎵 當前播放",
			},
			discord.ButtonComponent{
				Style:    discord.ButtonStyleSuccess,
				CustomID: ButtonQueue,
				Label:    "📜 播放佇列",
			},
		},
	}

	emptyContent := ""
	if err := event.UpdateMessage(discord.MessageUpdate{
		Content:    &emptyContent,
		Embeds:     &[]discord.Embed{embed},
		Components: &components,
	}); err != nil {
		log.Printf("failed to update message with control panel: %v", err)
	}
}

// HandleControlPanelInteraction 處理控制面板按鈕的互動事件
func HandleControlPanelInteraction(event *events.ComponentInteractionCreate) {
	if musicService == nil {
		respondToComponentInteraction(event, "音樂服務尚未初始化。")
		return
	}

	guildID := event.GuildID().String()
	player := musicService.GetOrCreatePlayer(guildID)

	switch event.Data.CustomID() {
	case ButtonShowPanel:
		// 展開完整控制面板
		song, hasSong := player.CurrentSong()
		content := "🎛️ **音樂控制面板**"
		if hasSong {
			content = fmt.Sprintf("🎵 **正在播放：** %s\n\n🎛️ **音樂控制面板**", song.Title)
		}
		UpdateMessageWithFullPanel(event, content)

	case ButtonPause:
		handlePauseButton(event, player)

	case ButtonSkip:
		handleSkipButton(event, player)

	case ButtonStop:
		handleStopButton(event, player)

	case ButtonNowPlaying:
		handleNowPlayingButton(event, player)

	case ButtonQueue:
		handleQueueButton(event, player)

	case ButtonSearch:
		handleSearchButton(event)

	default:
		respondToComponentInteraction(event, "未知的按鈕操作。")
	}
}

// 處理暫停按鈕
func handlePauseButton(event *events.ComponentInteractionCreate, player PlayerController) {
	lavalinkPlayer := lavalinkClient.Player(*event.GuildID())
	if lavalinkPlayer == nil {
		respondToComponentInteraction(event, "❌ 找不到播放器。")
		return
	}

	paused := lavalinkPlayer.Paused()
	if err := lavalinkPlayer.Update(context.Background(), lavalink.WithPaused(!paused)); err != nil {
		respondToComponentInteraction(event, fmt.Sprintf("❌ 操作失敗：%v", err))
		return
	}

	if paused {
		respondToComponentInteraction(event, "▶️ 已恢復播放。")
	} else {
		respondToComponentInteraction(event, "⏸️ 已暫停播放。")
	}
}

// 處理跳過按鈕
func handleSkipButton(event *events.ComponentInteractionCreate, player PlayerController) {
	lavalinkPlayer := lavalinkClient.Player(*event.GuildID())
	if lavalinkPlayer == nil {
		respondToComponentInteraction(event, "❌ 找不到播放器。")
		return
	}

	// 停止當前曲目
	if err := lavalinkPlayer.Update(context.Background(), lavalink.WithNullTrack()); err != nil {
		respondToComponentInteraction(event, fmt.Sprintf("❌ 跳過失敗：%v", err))
		return
	}

	nextSong, ok := player.Dequeue()
	if !ok {
		respondToComponentInteraction(event, "⏭️ 已跳過當前歌曲，佇列中沒有更多歌曲。")
		return
	}

	player.SetCurrentSong(nextSong)
	respondToComponentInteraction(event, fmt.Sprintf("⏭️ 已跳過，正在播放：**%s**", nextSong.Title))

	// 獲取用戶的語音頻道
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		log.Printf("無法獲取用戶的語音頻道")
		return
	}

	// 開始播放下一首
	go func() {
		if err := JoinVoiceAndPlayWithYtDlp(event.Client(), *event.GuildID(), *voiceState.ChannelID, nextSong.URL); err != nil {
			log.Printf("播放下一首失敗: %v", err)
		}
	}()
}

// 處理停止按鈕
func handleStopButton(event *events.ComponentInteractionCreate, player PlayerController) {
	lavalinkPlayer := lavalinkClient.Player(*event.GuildID())
	if lavalinkPlayer == nil {
		respondToComponentInteraction(event, "❌ 找不到播放器。")
		return
	}

	// 停止播放
	if err := lavalinkPlayer.Update(context.Background(), lavalink.WithNullTrack()); err != nil {
		respondToComponentInteraction(event, fmt.Sprintf("❌ 停止失敗：%v", err))
		return
	}

	player.Stop()

	if err := event.Client().UpdateVoiceState(context.Background(), *event.GuildID(), nil, false, false); err != nil {
		log.Printf("離開語音頻道時出錯: %v", err)
	}

	respondToComponentInteraction(event, "⏹️ 已停止播放並清空佇列。")
}

// 處理當前播放按鈕
func handleNowPlayingButton(event *events.ComponentInteractionCreate, player PlayerController) {
	song, ok := player.CurrentSong()
	if !ok {
		respondToComponentInteraction(event, "目前沒有播放任何歌曲。")
		return
	}

	message := fmt.Sprintf("🎵 正在播放：**%s**\n🔗 %s", song.Title, song.URL)
	respondToComponentInteraction(event, message)
}

// 處理佇列按鈕
func handleQueueButton(event *events.ComponentInteractionCreate, player PlayerController) {
	songs := player.QueueSnapshot()

	if len(songs) == 0 {
		respondToComponentInteraction(event, "📜 播放佇列是空的。")
		return
	}

	message := "📜 **播放佇列**\n\n"
	for i, song := range songs {
		if i >= 10 {
			message += fmt.Sprintf("... 還有 %d 首歌曲\n", len(songs)-10)
			break
		}
		message += fmt.Sprintf("%d. %s\n", i+1, song.Title)
	}

	respondToComponentInteraction(event, message)
}

// 處理搜尋按鈕 - 彈出 Modal 讓用戶輸入搜尋關鍵字
func handleSearchButton(event *events.ComponentInteractionCreate) {
	// 創建 Modal（彈出視窗）
	modal := discord.NewModalCreateBuilder().
		SetCustomID("music_search_modal").
		SetTitle("🔍 搜尋音樂").
		AddActionRow(
			discord.NewTextInput("search_query", discord.TextInputStyleShort, "歌曲名稱或 YouTube 連結").
				WithPlaceholder("例如：周杰倫 晴天").
				WithRequired(true).
				WithMaxLength(200),
		).
		Build()

	if err := event.Modal(modal); err != nil {
		log.Printf("failed to show modal: %v", err)
	}
}

// HandleModalSubmit 處理 Modal 提交事件
func HandleModalSubmit(event *events.ModalSubmitInteractionCreate) {
	if event.Data.CustomID != "music_search_modal" {
		return
	}

	// 獲取用戶輸入的搜尋關鍵字
	query := event.Data.Text("search_query")
	if query == "" {
		respondToModalInteraction(event, "❌ 請輸入搜尋關鍵字")
		return
	}

	// Defer 回應（因為搜尋和播放需要時間）
	if err := event.DeferCreateMessage(false); err != nil {
		log.Printf("failed to defer modal response: %v", err)
		return
	}

	// 檢查用戶是否在語音頻道中
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		updateModalResponse(event, "⚠️ 你必須先加入語音頻道才能播放")
		return
	}

	// 使用 YouTube resolver 搜尋
	if youtubeResolver == nil {
		updateModalResponse(event, "❌ YouTube 解析服務尚未初始化")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("搜尋音樂: %s", query)
	song, err := youtubeResolver.Resolve(ctx, query)
	if err != nil {
		log.Printf("搜尋失敗: %v", err)
		updateModalResponse(event, fmt.Sprintf("❌ 搜尋失敗：%v", err))
		return
	}

	song.RequestedBy = event.User().ID.String()

	// 加入佇列或開始播放
	guildID := event.GuildID().String()
	guildPlayer := musicService.GetOrCreatePlayer(guildID)
	_, hasCurrentSong := guildPlayer.CurrentSong()

	if err := guildPlayer.Enqueue(song); err != nil {
		updateModalResponse(event, fmt.Sprintf("❌ 加入佇列失敗：%v", err))
		return
	}

	if !hasCurrentSong {
		// 開始播放
		firstSong, ok := guildPlayer.Dequeue()
		if ok {
			guildPlayer.SetCurrentSong(firstSong)
			song = firstSong
		}

		err = JoinVoiceAndPlayWithYtDlp(event.Client(), *event.GuildID(), *voiceState.ChannelID, song.URL)
		if err != nil {
			log.Printf("播放失敗: %v", err)
			updateModalResponse(event, fmt.Sprintf("❌ 播放失敗：%v", err))
			return
		}

		updateModalResponseWithButton(event, fmt.Sprintf("✅ 正在播放：**%s**", song.Title))
	} else {
		// 已經有歌曲在播放，加入佇列
		updateModalResponseWithButton(event, fmt.Sprintf("✅ 已加入佇列：**%s**", song.Title))
	}
}

// updateModalResponse 更新 Modal 的 deferred response（使用 Embed）
func updateModalResponse(event *events.ModalSubmitInteractionCreate, content string) {
	embed := discord.NewEmbedBuilder().
		SetColor(0x5865F2).
		SetDescription(content).
		Build()

	embeds := []discord.Embed{embed}
	emptyContent := ""

	_, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
		Content: &emptyContent,
		Embeds:  &embeds,
	})
	if err != nil {
		log.Printf("failed to update modal response: %v", err)
	}
}

// updateModalResponseWithButton 更新 Modal 回應並附加控制面板按鈕（使用 Embed）
func updateModalResponseWithButton(event *events.ModalSubmitInteractionCreate, content string) {
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
		log.Printf("failed to update modal response: %v", err)
	}
}

// respondToModalInteraction 回應 Modal 互動
func respondToModalInteraction(event *events.ModalSubmitInteractionCreate, content string) {
	if err := event.CreateMessage(discord.MessageCreate{
		Content: content,
		Flags:   discord.MessageFlagEphemeral,
	}); err != nil {
		log.Printf("failed to respond to modal: %v", err)
	}
}

// respondToComponentInteraction 回應組件互動
func respondToComponentInteraction(event *events.ComponentInteractionCreate, content string) {
	if err := event.CreateMessage(discord.MessageCreate{
		Content: content,
		Flags:   discord.MessageFlagEphemeral, // 僅發送者可見
	}); err != nil {
		log.Printf("failed to respond to component interaction: %v", err)
	}
}
