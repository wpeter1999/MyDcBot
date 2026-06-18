package command

import (
	"context"
	"fmt"
	"log"
	"sync"

	"discordbot/internal/audio"
	"discordbot/internal/player"

	"github.com/bwmarrin/discordgo"
)

// voiceManager 管理 Discord 語音連線和播放迴圈。
type voiceManager struct {
	mu             sync.Mutex
	activePlayback map[string]context.CancelFunc // guildID -> cancel function
	pipeline       audio.Pipeline
}

var globalVoiceManager = &voiceManager{
	activePlayback: make(map[string]context.CancelFunc),
	pipeline:       audio.NewPipeline(audio.NewOpusStreamer()),
}

// JoinVoiceAndPlay 加入使用者的語音頻道並啟動播放迴圈（如果尚未播放）。
func JoinVoiceAndPlay(s *discordgo.Session, guildID, userID string, player player.PlaybackController) error {
	vm := globalVoiceManager
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// 檢查是否已經在播放
	if _, isPlaying := vm.activePlayback[guildID]; isPlaying {
		// 已經在播放，不需要重新加入
		return nil
	}

	// 找到使用者所在的語音頻道
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return fmt.Errorf("無法取得 Guild 資訊: %w", err)
	}

	var channelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			channelID = vs.ChannelID
			break
		}
	}

	if channelID == "" {
		return fmt.Errorf("你必須先加入語音頻道")
	}

	// 加入語音頻道
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return fmt.Errorf("無法加入語音頻道: %w", err)
	}

	// 啟動播放迴圈
	ctx, cancel := context.WithCancel(context.Background())
	vm.activePlayback[guildID] = cancel

	go func() {
		defer func() {
			vm.mu.Lock()
			delete(vm.activePlayback, guildID)
			vm.mu.Unlock()
			vc.Disconnect()
		}()

		// 包裝 discordgo.VoiceConnection 為我們的介面
		vcWrapper := &voiceConnectionWrapper{vc: vc}

		err := player.StartPlayback(ctx, vcWrapper, vm.pipeline)
		if err != nil && err != context.Canceled {
			log.Printf("播放迴圈錯誤 (Guild %s): %v", guildID, err)
		}
	}()

	return nil
}

// StopPlayback 停止指定 Guild 的播放。
func StopPlayback(guildID string) {
	vm := globalVoiceManager
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if cancel, ok := vm.activePlayback[guildID]; ok {
		cancel()
		delete(vm.activePlayback, guildID)
	}
}

// voiceConnectionWrapper 包裝 discordgo.VoiceConnection 實作 player.VoiceConnection 介面。
type voiceConnectionWrapper struct {
	vc *discordgo.VoiceConnection
}

func (w *voiceConnectionWrapper) Speaking(speaking bool) error {
	return w.vc.Speaking(speaking)
}

func (w *voiceConnectionWrapper) OpusSend() chan<- []byte {
	return w.vc.OpusSend
}

func (w *voiceConnectionWrapper) Disconnect() error {
	return w.vc.Disconnect()
}
