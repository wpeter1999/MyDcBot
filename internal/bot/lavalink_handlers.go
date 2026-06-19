package bot

import (
	"log"

	"discordbot/internal/command"

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

	// 如果是正常結束，自動播放下一首
	if event.Reason == lavalink.TrackEndReasonFinished {
		guildIDStr := player.GuildID().String()
		log.Printf("[Lavalink] Track finished, attempting to play next song...")

		if pm := b.playerManager; pm != nil {
			guildPlayer, ok := pm.Get(guildIDStr)
			if !ok || guildPlayer == nil {
				log.Printf("[Lavalink] No player found for guild %s", guildIDStr)
				return
			}

			// 從佇列取出下一首歌
			nextSong, ok := guildPlayer.Dequeue()
			if !ok {
				log.Printf("[Lavalink] No more songs in queue for guild %s", guildIDStr)
				return
			}

			log.Printf("[Lavalink] Playing next song: %s", nextSong.Title)

			// 設定為當前播放歌曲
			guildPlayer.SetCurrentSong(nextSong)

			// 異步播放下一首，避免阻塞事件處理器
			go func() {
				// 獲取當前的 voice channel ID
				voiceState, ok := b.Client.Caches().VoiceState(player.GuildID(), b.Client.ApplicationID())
				if !ok || voiceState.ChannelID == nil {
					log.Printf("[Lavalink] Could not find voice channel for guild %s", guildIDStr)
					return
				}

				channelID := *voiceState.ChannelID

				// 使用 yt-dlp 提取並播放
				err := command.JoinVoiceAndPlayWithYtDlp(b.Client, player.GuildID(), channelID, nextSong.URL)
				if err != nil {
					log.Printf("[Lavalink] Failed to play next song: %v", err)
					// 嘗試 SoundCloud 備用
					searchQuery := "scsearch:" + nextSong.Title
					err = command.JoinVoiceAndPlay(b.Client, player.GuildID(), channelID, searchQuery)
					if err != nil {
						log.Printf("[Lavalink] Failed to play with SoundCloud backup: %v", err)
					}
				}
			}()
		}
	}
}
