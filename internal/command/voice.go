package command

import (
	"context"
	"fmt"
	"log"

	"discordbot/internal/player"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

// JoinVoiceAndPlay 使用 Lavalink 加入語音頻道並播放
func JoinVoiceAndPlay(client bot.Client, guildID snowflake.ID, channelID snowflake.ID, trackURL string) error {
	ctx := context.Background()

	log.Printf("[Voice] Joining voice channel %s in guild %s", channelID, guildID)

	// 1. 更新 Discord voice state (加入頻道)
	err := client.UpdateVoiceState(ctx, guildID, &channelID, false, false)
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %w", err)
	}

	// 2. 獲取 Lavalink client
	lavalinkClient := GetLavalinkClient()
	if lavalinkClient == nil {
		return fmt.Errorf("lavalink client not initialized")
	}

	// 3. 載入音軌
	log.Printf("[Lavalink] Loading track: %s", trackURL)

	node := lavalinkClient.BestNode()
	if node == nil {
		return fmt.Errorf("no lavalink nodes available")
	}

	var loadedTrack lavalink.Track
	var loadErr error

	// LoadTracksHandler 不返回 error，只通過 handler 回調
	node.LoadTracksHandler(ctx, trackURL, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			log.Printf("[Lavalink] Loaded track: %s", track.Info.Title)
			loadedTrack = track
		},
		func(playlist lavalink.Playlist) {
			log.Printf("[Lavalink] Loaded playlist: %s", playlist.Info.Name)
			if len(playlist.Tracks) > 0 {
				loadedTrack = playlist.Tracks[0]
			}
		},
		func(tracks []lavalink.Track) {
			log.Printf("[Lavalink] Loaded %d search results", len(tracks))
			if len(tracks) > 0 {
				loadedTrack = tracks[0]
			}
		},
		func() {
			log.Printf("[Lavalink] No matches found")
			loadErr = fmt.Errorf("no matches found for: %s", trackURL)
		},
		func(err error) {
			log.Printf("[Lavalink] Load failed: %v", err)
			loadErr = err
		},
	))

	if loadErr != nil {
		return loadErr
	}

	if loadedTrack.Info.Title == "" {
		return fmt.Errorf("no track loaded")
	}

	// 4. 播放音軌
	player := lavalinkClient.Player(guildID)
	err = player.Update(ctx, lavalink.WithTrack(loadedTrack))
	if err != nil {
		return fmt.Errorf("failed to play track: %w", err)
	}

	log.Printf("[Lavalink] Now playing: %s", loadedTrack.Info.Title)
	return nil
}

// StopPlayback 停止播放並離開語音頻道
func StopPlayback(client bot.Client, guildID snowflake.ID) error {
	ctx := context.Background()

	log.Printf("[Voice] Stopping playback for guild %s", guildID)

	// 停止 Lavalink player
	lavalinkClient := GetLavalinkClient()
	if lavalinkClient != nil {
		player := lavalinkClient.Player(guildID)
		err := player.Update(ctx, lavalink.WithNullTrack())
		if err != nil {
			log.Printf("[Lavalink] Failed to stop player: %v", err)
		}
	}

	// 離開語音頻道
	err := client.UpdateVoiceState(ctx, guildID, nil, false, false)
	if err != nil {
		return fmt.Errorf("failed to leave voice channel: %w", err)
	}

	return nil
}

// PausePlayback 暫停或恢復播放
func PausePlayback(guildID snowflake.ID, pause bool) error {
	ctx := context.Background()

	lavalinkClient := GetLavalinkClient()
	if lavalinkClient == nil {
		return fmt.Errorf("lavalink client not initialized")
	}

	player := lavalinkClient.Player(guildID)
	err := player.Update(ctx, lavalink.WithPaused(pause))
	if err != nil {
		return fmt.Errorf("failed to update pause state: %w", err)
	}

	return nil
}

// SkipTrack 跳過當前音軌
func SkipTrack(guildID snowflake.ID) error {
	ctx := context.Background()

	lavalinkClient := GetLavalinkClient()
	if lavalinkClient == nil {
		return fmt.Errorf("lavalink client not initialized")
	}

	player := lavalinkClient.Player(guildID)
	// 停止當前音軌，觸發 TrackEnd 事件
	err := player.Update(ctx, lavalink.WithNullTrack())
	if err != nil {
		return fmt.Errorf("failed to skip track: %w", err)
	}

	return nil
}

// GetPlayerState 獲取播放器狀態
func GetPlayerState(guildID snowflake.ID) (isPlaying bool, isPaused bool, track *lavalink.Track) {
	lavalinkClient := GetLavalinkClient()
	if lavalinkClient == nil {
		return false, false, nil
	}

	player := lavalinkClient.Player(guildID)
	currentTrack := player.Track()
	return currentTrack != nil, player.Paused(), currentTrack
}

// PlayNextSongFromQueue 從佇列播放下一首，失敗時自動重試下一首
// 返回實際播放的歌曲信息
func PlayNextSongFromQueue(client bot.Client, guildID snowflake.ID, channelID snowflake.ID) (*player.Song, error) {
	if musicService == nil {
		return nil, fmt.Errorf("music service not initialized")
	}

	guildPlayer := musicService.GetOrCreatePlayer(guildID.String())

	// 嘗試播放佇列中的歌曲，最多嘗試 10 首（避免無限循環）
	maxAttempts := 10
	for attempt := 0; attempt < maxAttempts && guildPlayer.QueueLen() > 0; attempt++ {
		nextSong, ok := guildPlayer.Dequeue()
		if !ok {
			return nil, fmt.Errorf("no songs in queue")
		}

		log.Printf("[AutoPlay] Attempting to play: %s (attempt %d)", nextSong.Title, attempt+1)
		guildPlayer.SetCurrentSong(nextSong)

		// 嘗試用 yt-dlp 播放
		err := JoinVoiceAndPlayWithYtDlp(client, guildID, channelID, nextSong.URL)
		if err == nil {
			log.Printf("[AutoPlay] Successfully playing: %s", nextSong.Title)
			return &nextSong, nil
		}

		log.Printf("[AutoPlay] Failed to play %s: %v, trying SoundCloud...", nextSong.Title, err)

		// 嘗試 SoundCloud 備用
		searchQuery := "scsearch:" + nextSong.Title
		err = JoinVoiceAndPlay(client, guildID, channelID, searchQuery)
		if err == nil {
			log.Printf("[AutoPlay] Successfully playing via SoundCloud: %s", nextSong.Title)
			return &nextSong, nil
		}

		log.Printf("[AutoPlay] SoundCloud also failed for %s: %v, trying next song...", nextSong.Title, err)
		guildPlayer.ClearCurrentSong()
	}

	return nil, fmt.Errorf("failed to play any song from queue after %d attempts", maxAttempts)
}

