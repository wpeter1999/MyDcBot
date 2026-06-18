package audio

import (
	"context"
	"io"

	"discordbot/internal/player"
)

// Streamer 定義音訊串流介面。
type Streamer interface {
	// Stream 開始串流指定 URL 的音訊，回傳 DCA 格式的 reader。
	Stream(ctx context.Context, url string) (io.ReadCloser, error)
}

// Pipeline 定義完整的音訊播放管道介面。
type Pipeline interface {
	// Play 播放指定 URL 的音訊到 Discord 語音連線。
	// 會阻塞直到播放完成、取消或發生錯誤。
	Play(ctx context.Context, url string, voiceConnection player.VoiceConnection) error
}
