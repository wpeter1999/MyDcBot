package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

// YtDlpInfo yt-dlp JSON 輸出的簡化結構
type YtDlpInfo struct {
	URL    string `json:"url"`
	Title  string `json:"title"`
	Formats []struct {
		URL      string `json:"url"`
		AudioExt string `json:"aext"`
		VideoExt string `json:"vext"`
	} `json:"formats"`
}

// JoinVoiceAndPlayWithYtDlp 使用 yt-dlp 提取音訊 URL 並播放
func JoinVoiceAndPlayWithYtDlp(client bot.Client, guildID snowflake.ID, channelID snowflake.ID, youtubeURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("[yt-dlp] 提取音訊 URL: %s", youtubeURL)

	// 使用 yt-dlp 提取音訊 URL
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--dump-json",
		"--no-playlist",
		"--format", "bestaudio/best",
		youtubeURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("yt-dlp 失敗: %w", err)
	}

	var info YtDlpInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return fmt.Errorf("解析 yt-dlp 輸出失敗: %w", err)
	}

	// 找到最佳音訊格式 URL
	audioURL := info.URL
	for _, format := range info.Formats {
		if format.AudioExt != "none" && format.VideoExt == "none" && format.URL != "" {
			audioURL = format.URL
			break
		}
	}

	if audioURL == "" {
		return fmt.Errorf("找不到音訊 URL")
	}

	log.Printf("[yt-dlp] 提取成功，使用音訊 URL: %s", audioURL[:100]+"...")

	// 更新 Discord voice state (加入頻道)
	err = client.UpdateVoiceState(ctx, guildID, &channelID, false, false)
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %w", err)
	}

	// 獲取 Lavalink client
	lavalinkClient := GetLavalinkClient()
	if lavalinkClient == nil {
		return fmt.Errorf("lavalink client not initialized")
	}

	// 直接使用音訊 URL 播放（HTTP stream）
	node := lavalinkClient.BestNode()
	if node == nil {
		return fmt.Errorf("no lavalink nodes available")
	}

	var loadedTrack lavalink.Track
	var loadErr error

	node.LoadTracksHandler(ctx, audioURL, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			log.Printf("[Lavalink] 載入 HTTP 音訊串流成功")
			loadedTrack = track
		},
		func(playlist lavalink.Playlist) {
			log.Printf("[Lavalink] 載入播放清單（不應該發生）")
		},
		func(tracks []lavalink.Track) {
			log.Printf("[Lavalink] 載入多個音軌（不應該發生）")
		},
		func() {
			log.Printf("[Lavalink] 找不到匹配")
			loadErr = fmt.Errorf("no matches found")
		},
		func(err error) {
			log.Printf("[Lavalink] 載入失敗: %v", err)
			loadErr = err
		},
	))

	if loadErr != nil {
		return loadErr
	}

	if loadedTrack.Info.Title == "" {
		return fmt.Errorf("no track loaded")
	}

	// 播放音軌
	player := lavalinkClient.Player(guildID)
	err = player.Update(ctx, lavalink.WithTrack(loadedTrack))
	if err != nil {
		return fmt.Errorf("failed to play track: %w", err)
	}

	log.Printf("[Lavalink] 開始播放: %s", loadedTrack.Info.Title)
	return nil
}
