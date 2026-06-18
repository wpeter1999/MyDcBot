package command

import (
	"context"

	"discordbot/internal/player"
)

// MusicService 定義音樂指令與播放器互動的介面，方便測試時注入 fake 實作。
type MusicService interface {
	// GetOrCreatePlayer 取得或建立指定 Guild 的播放器。
	GetOrCreatePlayer(guildID string) PlayerController

	// RemovePlayer 移除指定 Guild 的播放器。
	RemovePlayer(guildID string) bool
}

// PlayerController 定義單一 Guild 播放器的控制介面。
type PlayerController interface {
	// Enqueue 將歌曲加入播放佇列。
	Enqueue(song player.Song) error

	// Skip 跳過目前播放的歌曲。
	Skip() bool

	// TogglePause 切換暫停/繼續狀態，回傳新的暫停狀態。
	TogglePause() bool

	// IsPaused 回傳目前是否處於暫停狀態。
	IsPaused() bool

	// Stop 停止播放並清空佇列。
	Stop()

	// QueueSnapshot 回傳目前佇列的快照。
	QueueSnapshot() []player.Song

	// CurrentSong 回傳目前播放的歌曲；沒有歌曲時 ok 為 false。
	CurrentSong() (player.Song, bool)

	// GuildID 回傳此播放器所屬的 Guild ID。
	GuildID() string

	// StartPlayback 啟動播放迴圈（來自 player.PlaybackController）。
	StartPlayback(ctx context.Context, vc player.VoiceConnection, pipeline player.AudioPipeline) error
}

// defaultMusicService 是預設的 MusicService 實作，直接使用 player.Manager。
type defaultMusicService struct {
	manager *player.Manager
}

// NewDefaultMusicService 建立預設的 MusicService。
func NewDefaultMusicService(manager *player.Manager) MusicService {
	return &defaultMusicService{manager: manager}
}

// GetOrCreatePlayer 實作 MusicService 介面。
func (s *defaultMusicService) GetOrCreatePlayer(guildID string) PlayerController {
	return s.manager.GetOrCreate(guildID)
}

// RemovePlayer 實作 MusicService 介面。
func (s *defaultMusicService) RemovePlayer(guildID string) bool {
	return s.manager.Remove(guildID)
}

var musicService MusicService

// SetMusicService 設定全域 MusicService，供指令 handler 使用（測試時可注入 fake）。
func SetMusicService(service MusicService) {
	musicService = service
}

// GetMusicService 取得目前的 MusicService。
func GetMusicService() MusicService {
	return musicService
}
