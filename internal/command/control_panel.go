package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"discordbot/internal/player"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
	ButtonLoop       = "music_loop"
	ButtonShuffle    = "music_shuffle"
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
			discord.ButtonComponent{
				Style:    discord.ButtonStyle(player.GetLoopMode().ButtonStyle()),
				CustomID: ButtonLoop,
				Emoji:    &discord.ComponentEmoji{Name: player.GetLoopMode().Icon()},
			},
			discord.ButtonComponent{
				Style:    getShuffleButtonStyle(player),
				CustomID: ButtonShuffle,
				Emoji:    &discord.ComponentEmoji{Name: "🔀"},
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

	case ButtonLoop:
		handleLoopButton(event, player)

	case ButtonShuffle:
		handleShuffleButton(event, player)

	default:
		respondToComponentInteraction(event, "未知的按鈕操作。")
	}
}

// 處理暫停按鈕
func handlePauseButton(event *events.ComponentInteractionCreate, player PlayerController) {
	newPaused, isPlaying, err := ExecutePauseToggle(*event.GuildID())

	if !isPlaying {
		respondToComponentInteraction(event, "目前沒有播放任何歌曲。")
		return
	}

	if err != nil {
		respondToComponentInteraction(event, fmt.Sprintf("❌ 操作失敗：%v", err))
		return
	}

	if newPaused {
		respondToComponentInteraction(event, "⏸️ 已暫停播放。")
	} else {
		respondToComponentInteraction(event, "▶️ 已恢復播放。")
	}
}

// 處理跳過按鈕
func handleSkipButton(event *events.ComponentInteractionCreate, player PlayerController) {
	hasNext := ExecuteSkip(*event.GuildID(), player)

	if !hasNext {
		respondToComponentInteraction(event, "⏭️ 已跳過，但佇列中沒有更多歌曲。")
		return
	}

	respondToComponentInteraction(event, "⏭️ 正在跳過...")

	// 獲取用戶的語音頻道
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		log.Printf("無法獲取用戶的語音頻道")
		return
	}

	// 開始播放下一首（會自動重試失敗的歌曲）
	go func() {
		playedSong, err := PlayNextSongFromQueue(event.Client(), *event.GuildID(), *voiceState.ChannelID)
		if err != nil {
			log.Printf("Skip button: Failed to play any song: %v", err)
			sendComponentFollowupMessage(event, "❌ 無法播放佇列中的任何歌曲")
		} else if playedSong != nil {
			log.Printf("Skip button: Now playing: %s", playedSong.Title)
			sendComponentFollowupMessage(event, fmt.Sprintf("✅ 現在播放：**%s**", playedSong.Title))
		}
	}()
}

// sendComponentFollowupMessage 發送 component 的 followup 訊息
func sendComponentFollowupMessage(event *events.ComponentInteractionCreate, content string) {
	_, err := event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
		Content: content,
	})
	if err != nil {
		log.Printf("Failed to send followup message: %v", err)
	}
}

// 處理停止按鈕
func handleStopButton(event *events.ComponentInteractionCreate, player PlayerController) {
	err := ExecuteStop(event.Client(), *event.GuildID(), player)
	if err != nil {
		respondToComponentInteraction(event, fmt.Sprintf("❌ %v", err))
		return
	}

	respondToComponentInteraction(event, "⏹️ 已停止播放並清空佇列。")
}

// 處理當前播放按鈕
func handleNowPlayingButton(event *events.ComponentInteractionCreate, player PlayerController) {
	message, _ := FormatNowPlaying(player)
	respondToComponentInteraction(event, message)
}

// 處理佇列按鈕
func handleQueueButton(event *events.ComponentInteractionCreate, player PlayerController) {
	message := FormatQueueDisplay(player)
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

// 處理循環按鈕
func handleLoopButton(event *events.ComponentInteractionCreate, player PlayerController) {
	// 切換循環模式
	newMode := player.ToggleLoopMode()

	// 建構回應訊息
	icon := newMode.Icon()
	modeName := newMode.String()

	message := fmt.Sprintf("%s **循環模式：%s**", icon, modeName)

	respondToComponentInteraction(event, message)
}

// 處理隨機按鈕
func handleShuffleButton(event *events.ComponentInteractionCreate, player PlayerController) {
	// 檢查佇列是否為空
	queueLen := player.QueueLen()
	if queueLen == 0 {
		respondToComponentInteraction(event, "⚠️ 佇列中沒有歌曲可以打亂")
		return
	}

	// 打亂佇列
	player.Shuffle()

	// 回應訊息
	message := fmt.Sprintf("🔀 **已打亂佇列**\n共 %d 首歌曲已隨機排序", queueLen)
	respondToComponentInteraction(event, message)
}

// getShuffleButtonStyle 取得 shuffle 按鈕樣式
func getShuffleButtonStyle(player PlayerController) discord.ButtonStyle {
	if player.IsShuffled() {
		return discord.ButtonStyleSuccess // 綠色
	}
	return discord.ButtonStyleSecondary // 灰色
}

// HandleModalSubmit 處理 Modal 提交事件
func HandleModalSubmit(event *events.ModalSubmitInteractionCreate) {
	if event.Data.CustomID != "music_search_modal" {
		return
	}

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

	// 驗證語音頻道
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		updateModalResponse(event, "⚠️ 你必須先加入語音頻道才能播放")
		return
	}

	// 搜尋歌曲
	song, err := searchSongFromModal(query, event.User().ID.String())
	if err != nil {
		updateModalResponse(event, fmt.Sprintf("❌ 搜尋失敗：%v", err))
		return
	}

	// 播放或加入佇列
	handleModalSongPlayback(event, song, voiceState)
}

// searchSongFromModal 從 Modal 搜尋歌曲
func searchSongFromModal(query, userID string) (player.Song, error) {
	if youtubeResolver == nil {
		return player.Song{}, fmt.Errorf("YouTube 解析服務尚未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("搜尋音樂: %s", query)
	song, err := youtubeResolver.Resolve(ctx, query)
	if err != nil {
		log.Printf("搜尋失敗: %v", err)
		return player.Song{}, err
	}

	song.RequestedBy = userID
	return song, nil
}

// handleModalSongPlayback 處理 Modal 搜尋後的播放邏輯
func handleModalSongPlayback(event *events.ModalSubmitInteractionCreate, song player.Song, voiceState discord.VoiceState) {
	guildID := event.GuildID().String()
	guildPlayer := musicService.GetOrCreatePlayer(guildID)

	if err := guildPlayer.Enqueue(song); err != nil {
		updateModalResponse(event, fmt.Sprintf("❌ 加入佇列失敗：%v", err))
		return
	}

	_, hasCurrentSong := guildPlayer.CurrentSong()

	if !hasCurrentSong {
		// 開始播放
		firstSong, ok := guildPlayer.Dequeue()
		if ok {
			guildPlayer.SetCurrentSong(firstSong)
			song = firstSong
		}

		err := JoinVoiceAndPlayWithYtDlp(event.Client(), *event.GuildID(), *voiceState.ChannelID, song.URL)
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

// respondToModalInteraction 回應 Modal 互動（使用 Embed）
func respondToModalInteraction(event *events.ModalSubmitInteractionCreate, content string) {
	embed := discord.NewEmbedBuilder().
		SetColor(0x5865F2).
		SetDescription(content).
		Build()

	if err := event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
		Flags:  discord.MessageFlagEphemeral,
	}); err != nil {
		log.Printf("failed to respond to modal: %v", err)
	}
}

// respondToComponentInteraction 回應組件互動（使用 Embed）
func respondToComponentInteraction(event *events.ComponentInteractionCreate, content string) {
	embed := discord.NewEmbedBuilder().
		SetColor(0x5865F2).
		SetDescription(content).
		Build()

	if err := event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
		Flags:  discord.MessageFlagEphemeral, // 僅發送者可見
	}); err != nil {
		log.Printf("failed to respond to component interaction: %v", err)
	}
}

// 處理隨機按鈕
