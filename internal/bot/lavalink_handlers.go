package bot

import (
	"log"

	"discordbot/internal/command"
	playerPkg "discordbot/internal/player"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

// BotEventListener 實現 disgolink.EventListener 介面
type BotEventListener struct {
	bot *Bot
}

// OnEvent 處理所有 Lavalink 事件
func (l *BotEventListener) OnEvent(player disgolink.Player, event lavalink.Message) {
	switch e := event.(type) {
	case lavalink.PlayerUpdateMessage:
		l.bot.onPlayerUpdate(player, e)
	case lavalink.TrackStartEvent:
		l.bot.onTrackStart(player, e)
	case lavalink.TrackEndEvent:
		l.bot.onTrackEnd(player, e)
	}
}

// onPlayerUpdate 處理 Lavalink player 更新事件
func (b *Bot) onPlayerUpdate(player disgolink.Player, event lavalink.PlayerUpdateMessage) {
	log.Printf("[Lavalink] Player update for guild %d: position=%d", player.GuildID(), event.State.Position)
}

// onTrackStart 處理音軌開始播放事件
func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	log.Printf("[Lavalink] Track started for guild %d: %s", player.GuildID(), event.Track.Info.Title)
}

// onTrackEnd 處理音軌結束事件
func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	log.Printf("[Lavalink] Track ended for guild %d: %s (reason: %s)", player.GuildID(), event.Track.Info.Title, event.Reason)

	// 處理正常結束或載入失敗的情況，自動播放下一首
	// Finished: 正常播放完畢
	// LoadFailed: 無法載入/播放（例如影片無法播放、版權限制等）
	shouldPlayNext := event.Reason == lavalink.TrackEndReasonFinished || event.Reason == lavalink.TrackEndReasonLoadFailed

	if !shouldPlayNext {
		return
	}

	if event.Reason == lavalink.TrackEndReasonLoadFailed {
		log.Printf("[Lavalink] Track failed to load, skipping to next song...")
	}

	// 處理循環播放
	b.handleLoopMode(player)

	b.playNextSongInQueue(player)
}

// playNextSongInQueue 從佇列播放下一首歌曲
func (b *Bot) playNextSongInQueue(player disgolink.Player) {
	guildIDStr := player.GuildID().String()
	log.Printf("[Lavalink] Attempting to play next song...")

	if b.playerManager == nil {
		log.Printf("[Lavalink] Player manager is nil")
		return
	}

	guildPlayer, ok := b.playerManager.Get(guildIDStr)
	if !ok || guildPlayer == nil {
		log.Printf("[Lavalink] No player found for guild %s", guildIDStr)
		return
	}

	// 清除當前播放歌曲
	guildPlayer.ClearCurrentSong()

	// 檢查佇列是否為空
	if guildPlayer.QueueLen() == 0 {
		log.Printf("[Lavalink] No more songs in queue for guild %s", guildIDStr)
		return
	}

	// 異步播放下一首（PlayNextSongFromQueue 會自動從佇列取歌）
	go b.playNextSongAsync(player)
}

// playNextSongAsync 異步播放下一首歌曲
func (b *Bot) playNextSongAsync(player disgolink.Player) {
	// 獲取語音頻道 ID
	voiceState, ok := b.Client.Caches().VoiceState(player.GuildID(), b.Client.ApplicationID())
	if !ok || voiceState.ChannelID == nil {
		log.Printf("[Lavalink] Could not find voice channel for guild %s", player.GuildID().String())
		return
	}

	channelID := *voiceState.ChannelID

	// 使用統一的播放函數，會自動從佇列取出下一首並重試失敗的歌曲
	playedSong, err := command.PlayNextSongFromQueue(b.Client, player.GuildID(), channelID)
	if err != nil {
		log.Printf("[Lavalink] Failed to play any song from queue: %v", err)
	} else if playedSong != nil {
		log.Printf("[Lavalink] Auto-playing: %s", playedSong.Title)
	}
}

// handleLoopMode 處理循環播放邏輯
func (b *Bot) handleLoopMode(player disgolink.Player) {
	guildIDStr := player.GuildID().String()

	if b.playerManager == nil {
		return
	}

	guildPlayer, ok := b.playerManager.Get(guildIDStr)
	if !ok || guildPlayer == nil {
		return
	}

	// 取得當前播放的歌曲
	currentSong, hasSong := guildPlayer.CurrentSong()
	if !hasSong {
		return
	}

	// 取得循環模式
	loopMode := guildPlayer.GetLoopMode()

	switch loopMode {
	case playerPkg.LoopSingleOnce:
		// 單曲循環一次：將當前歌曲插入到佇列最前面，然後關閉循環
		log.Printf("[Loop] Single loop once: inserting %s to front of queue and disabling loop", currentSong.Title)
		if err := guildPlayer.EnqueueFront(currentSong); err != nil {
			log.Printf("[Loop] Failed to enqueue front for single loop once: %v", err)
		} else {
			// 循環一次後自動關閉
			guildPlayer.SetLoopMode(playerPkg.LoopOff)
			log.Printf("[Loop] Loop mode automatically disabled after single repeat")
		}

	case playerPkg.LoopSingleInfinite:
		// 單曲無限循環：將當前歌曲插入到佇列最前面
		log.Printf("[Loop] Single infinite loop: inserting %s to front of queue", currentSong.Title)
		if err := guildPlayer.EnqueueFront(currentSong); err != nil {
			log.Printf("[Loop] Failed to enqueue front for single infinite loop: %v", err)
		}

	case playerPkg.LoopOff:
		// 不循環，什麼都不做
		log.Printf("[Loop] Loop off: normal playback")
	}
}

