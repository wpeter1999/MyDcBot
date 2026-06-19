package audio

import (
	"context"
	"fmt"
	"io"

	"discordbot/internal/player"
)

// dcaPipeline 實作完整的音訊播放管道。
type dcaPipeline struct {
	streamer Streamer
}

// NewPipeline 建立新的音訊播放管道。
func NewPipeline(streamer Streamer) Pipeline {
	return &dcaPipeline{
		streamer: streamer,
	}
}

// Play 實作 Pipeline 介面，播放音訊到 Discord 語音連線。
func (p *dcaPipeline) Play(ctx context.Context, url string, vc player.VoiceConnection) error {
	// 開始串流
	stream, err := p.streamer.Stream(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}
	defer stream.Close()

	// 通知 Discord 開始說話
	if err := vc.Speaking(true); err != nil {
		return fmt.Errorf("failed to set speaking: %w", err)
	}
	defer vc.Speaking(false)

	// 取得 Opus 發送 channel
	send := vc.OpusSend()

	// 讀取並發送 Opus 封包
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := stream.Read(buf)
		if err == io.EOF {
			return nil // 正常結束
		}
		if err != nil {
			return fmt.Errorf("failed to read stream: %w", err)
		}

		// 複製封包並發送到 Discord
		packet := make([]byte, n)
		copy(packet, buf[:n])

		select {
		case send <- packet:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
