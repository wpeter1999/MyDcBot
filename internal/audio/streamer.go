package audio

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"layeh.com/gopus"
)

const (
	// Discord 要求的音訊參數
	channels   = 2     // stereo
	sampleRate = 48000 // 48kHz
	frameSize  = 960   // 20ms frame at 48kHz
	maxBytes   = 1024  // Opus 最大封包大小
)

// opusStreamer 使用 ffmpeg + gopus 將音訊 URL 轉換為 Opus 串流。
type opusStreamer struct{}

// NewOpusStreamer 建立新的 Opus streamer。
func NewOpusStreamer() Streamer {
	return &opusStreamer{}
}

// Stream 實作 Streamer 介面，使用 ffmpeg 從 URL 串流音訊並編碼為 Opus。
func (s *opusStreamer) Stream(ctx context.Context, url string) (io.ReadCloser, error) {
	// 使用 ffmpeg 從 URL 取得音訊並轉為 PCM
	// -i: 輸入 URL
	// -f s16le: 16-bit signed little-endian PCM
	// -ar 48000: 48kHz sample rate (Discord 要求)
	// -ac 2: stereo (2 channels)
	// -: 輸出到 stdout
	ffmpegCmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-i", url,
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"-loglevel", "warning",
		"pipe:1",
	)

	// 取得 ffmpeg 的 stdout
	ffmpegOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create ffmpeg stdout pipe: %w", err)
	}

	// 啟動 ffmpeg
	if err := ffmpegCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// 建立 Opus 編碼器
	encoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
	if err != nil {
		ffmpegCmd.Process.Kill()
		return nil, fmt.Errorf("failed to create opus encoder: %w", err)
	}

	// 回傳可以讀取 Opus 封包的 reader
	return &opusReader{
		encoder: encoder,
		pcmIn:   ffmpegOut,
		ffmpeg:  ffmpegCmd,
		pcmBuf:  make([]int16, frameSize*channels),
	}, nil
}

// opusReader 從 PCM 輸入讀取並編碼為 Opus 封包。
type opusReader struct {
	encoder *gopus.Encoder
	pcmIn   io.ReadCloser
	ffmpeg  *exec.Cmd
	pcmBuf  []int16
}

// Read 實作 io.Reader 介面，讀取一個 Opus 封包。
func (r *opusReader) Read(p []byte) (int, error) {
	// 讀取一個 PCM frame (frameSize * channels * 2 bytes per sample)
	bytesNeeded := frameSize * channels * 2
	rawBuf := make([]byte, bytesNeeded)

	n, err := io.ReadFull(r.pcmIn, rawBuf)
	if err != nil {
		return 0, err
	}

	if n != bytesNeeded {
		return 0, io.ErrUnexpectedEOF
	}

	// 轉換 bytes 為 int16 PCM samples
	for i := 0; i < len(r.pcmBuf); i++ {
		r.pcmBuf[i] = int16(rawBuf[i*2]) | int16(rawBuf[i*2+1])<<8
	}

	// 編碼為 Opus
	opusData, err := r.encoder.Encode(r.pcmBuf, frameSize, maxBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to encode opus: %w", err)
	}

	// 複製到輸出 buffer
	if len(opusData) > len(p) {
		return 0, fmt.Errorf("opus packet too large: %d > %d", len(opusData), len(p))
	}

	copy(p, opusData)
	return len(opusData), nil
}

// Close 實作 io.Closer 介面，清理資源。
func (r *opusReader) Close() error {
	r.pcmIn.Close()
	if r.ffmpeg.Process != nil {
		r.ffmpeg.Process.Kill()
	}
	return nil
}
