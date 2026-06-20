package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
)

// PlaylistEntry 播放清單項目
type PlaylistEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// ExtractPlaylist 使用 yt-dlp 提取播放清單
func ExtractPlaylist(youtubeURL string) ([]PlaylistEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Printf("[yt-dlp] 提取播放清單: %s", youtubeURL)

	// 使用 flat-playlist 快速提取播放清單資訊
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--dump-json",
		"--flat-playlist",
		"--skip-download",
		youtubeURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp 提取播放清單失敗: %w", err)
	}

	// 解析每一行 JSON
	lines := strings.Split(string(output), "\n")
	var entries []PlaylistEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry PlaylistEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			log.Printf("[yt-dlp] 解析項目失敗: %v", err)
			continue
		}

		// 構建完整 URL
		if entry.ID != "" {
			entry.URL = "https://www.youtube.com/watch?v=" + entry.ID
			entries = append(entries, entry)
		}
	}

	log.Printf("[yt-dlp] 成功提取 %d 首歌曲", len(entries))
	return entries, nil
}

// PlayPlaylist 播放整個播放清單
func PlayPlaylist(client bot.Client, guildID snowflake.ID, channelID snowflake.ID, entries []PlaylistEntry) error {
	if len(entries) == 0 {
		return fmt.Errorf("播放清單是空的")
	}

	// 播放第一首歌
	firstURL := entries[0].URL
	log.Printf("[Playlist] 開始播放第一首: %s", entries[0].Title)

	err := JoinVoiceAndPlayWithYtDlp(client, guildID, channelID, firstURL)
	if err != nil {
		return fmt.Errorf("播放第一首失敗: %w", err)
	}

	// 注意：佇列功能已在 play.go 的 handlePlaylist 中實現
	// 此函數保留用於未來可能的獨立播放清單處理

	return nil
}

// IsPlaylistURL 檢查 URL 是否包含播放清單
func IsPlaylistURL(url string) bool {
	return strings.Contains(url, "list=") || strings.Contains(url, "/playlist?")
}
