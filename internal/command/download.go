package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// DownloadCommand 定義 /download 指令
var DownloadCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "download",
		Description: "下載 YouTube 音訊檔案",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "format",
				Description: "音訊格式",
				Required:    true,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{Name: "MP3 320kbps (推薦)", Value: "mp3-320"},
					{Name: "M4A 256kbps", Value: "m4a-256"},
					{Name: "Opus 192kbps (最小)", Value: "opus-192"},
					{Name: "FLAC 無損 (較大)", Value: "flac"},
					{Name: "WAV 原始 (最大)", Value: "wav"},
				},
			},
			discord.ApplicationCommandOptionString{
				Name:        "url",
				Description: "YouTube 影片網址",
				Required:    true,
			},
		},
	},
	Handler: downloadCommandHandler,
}

// downloadCommandHandler 處理 /download 指令
func downloadCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	// Defer response 避免超時
	if err := event.DeferCreateMessage(false); err != nil {
		log.Printf("failed to defer response: %v", err)
		return
	}

	data := event.SlashCommandInteractionData()
	format := data.String("format")
	url := data.String("url")

	if url == "" {
		updateResponse(event, "❌ 請提供 YouTube 網址")
		return
	}

	// 開始下載
	updateResponse(event, "⏳ 正在下載音訊檔案，請稍候...")

	// 創建臨時目錄
	tempDir := filepath.Join(os.TempDir(), "discord-downloads")
	os.MkdirAll(tempDir, 0755)

	// 生成唯一檔案名
	timestamp := time.Now().Unix()
	userID := event.User().ID.String()
	outputTemplate := filepath.Join(tempDir, fmt.Sprintf("%s_%d_%%(title)s.%%(ext)s", userID, timestamp))

	// 根據格式設定 yt-dlp 參數
	args := buildYtDlpArgs(format, url, outputTemplate)

	// 執行 yt-dlp
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("yt-dlp error: %v, output: %s", err, string(output))
		updateResponse(event, fmt.Sprintf("❌ 下載失敗：%v", err))
		return
	}

	// 查找下載的檔案
	files, err := filepath.Glob(filepath.Join(tempDir, fmt.Sprintf("%s_%d_*", userID, timestamp)))
	if err != nil || len(files) == 0 {
		updateResponse(event, "❌ 找不到下載的檔案")
		return
	}

	downloadedFile := files[0]
	defer os.Remove(downloadedFile) // 清理檔案

	// 檢查檔案大小
	fileInfo, err := os.Stat(downloadedFile)
	if err != nil {
		updateResponse(event, fmt.Sprintf("❌ 無法讀取檔案資訊：%v", err))
		return
	}

	fileSize := fileInfo.Size()
	maxSize := int64(25 * 1024 * 1024) // 25 MB

	if fileSize > maxSize {
		sizeMB := float64(fileSize) / 1024 / 1024
		updateResponse(event, fmt.Sprintf("❌ 檔案太大 (%.2f MB)\n💡 建議：使用 Opus 格式或選擇較短的歌曲", sizeMB))
		return
	}

	// 上傳檔案到 Discord
	file, err := os.Open(downloadedFile)
	if err != nil {
		updateResponse(event, fmt.Sprintf("❌ 無法開啟檔案：%v", err))
		return
	}
	defer file.Close()

	fileName := filepath.Base(downloadedFile)
	sizeMB := float64(fileSize) / 1024 / 1024

	message := fmt.Sprintf("✅ **下載完成！**\n📦 檔案：`%s`\n📊 大小：%.2f MB", fileName, sizeMB)

	// 更新回應並附加檔案
	_, err = event.Client().Rest().UpdateInteractionResponse(
		event.ApplicationID(),
		event.Token(),
		discord.MessageUpdate{
			Content: &message,
			Files: []*discord.File{
				{
					Name:   fileName,
					Reader: file,
				},
			},
		},
	)

	if err != nil {
		log.Printf("failed to upload file: %v", err)
		updateResponse(event, fmt.Sprintf("❌ 上傳檔案失敗：%v", err))
		return
	}

	log.Printf("Successfully downloaded and uploaded: %s (%.2f MB)", fileName, sizeMB)
}

// buildYtDlpArgs 根據格式構建 yt-dlp 參數
func buildYtDlpArgs(format, url, outputTemplate string) []string {
	args := []string{
		"--no-playlist",
		"--max-filesize", "25M",
		"--match-filter", "duration < 600", // 限制 10 分鐘
		"--output", outputTemplate,
	}

	switch format {
	case "mp3-320":
		args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0")
	case "m4a-256":
		args = append(args, "--extract-audio", "--audio-format", "m4a", "--audio-quality", "256K")
	case "opus-192":
		args = append(args, "--extract-audio", "--audio-format", "opus", "--audio-quality", "192K")
	case "flac":
		args = append(args, "--extract-audio", "--audio-format", "flac")
	case "wav":
		args = append(args, "--extract-audio", "--audio-format", "wav")
	default:
		args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0")
	}

	args = append(args, url)
	return args
}
