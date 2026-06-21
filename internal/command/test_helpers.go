package command

import (
	"context"

	"discordbot/internal/player"
)

// MockPlayerControllerExt 扩展的 mock player，用于测试
type MockPlayerControllerExt struct {
	queueLen           int
	currentSongCleared bool
	currentSong        *player.Song
	queue              []player.Song
	loopMode           player.LoopMode
}

func (m *MockPlayerControllerExt) QueueLen() int {
	return m.queueLen
}

func (m *MockPlayerControllerExt) ClearCurrentSong() {
	m.currentSongCleared = true
	m.currentSong = nil
}

func (m *MockPlayerControllerExt) Enqueue(song player.Song) error {
	m.queue = append(m.queue, song)
	return nil
}

func (m *MockPlayerControllerExt) Dequeue() (player.Song, bool) {
	if len(m.queue) == 0 {
		return player.Song{}, false
	}
	song := m.queue[0]
	m.queue = m.queue[1:]
	return song, true
}

func (m *MockPlayerControllerExt) SetCurrentSong(song player.Song) {
	m.currentSong = &song
}

func (m *MockPlayerControllerExt) CurrentSong() (player.Song, bool) {
	if m.currentSong == nil {
		return player.Song{}, false
	}
	return *m.currentSong, true
}

func (m *MockPlayerControllerExt) QueueSnapshot() []player.Song {
	return m.queue
}

func (m *MockPlayerControllerExt) Stop() {
	m.queue = nil
	m.currentSong = nil
}

func (m *MockPlayerControllerExt) IsStopped() bool {
	return false
}

func (m *MockPlayerControllerExt) GuildID() string {
	return "test-guild-123"
}

func (m *MockPlayerControllerExt) IsPaused() bool {
	return false
}

func (m *MockPlayerControllerExt) TogglePause() bool {
	return false
}

func (m *MockPlayerControllerExt) Skip() {
	// Mock implementation
}

func (m *MockPlayerControllerExt) SkipChan() <-chan struct{} {
	// Mock implementation - returns nil channel for testing
	return nil
}

func (m *MockPlayerControllerExt) Done() <-chan struct{} {
	// Mock implementation - returns nil channel for testing
	return nil
}

func (m *MockPlayerControllerExt) GetLoopMode() player.LoopMode {
	return m.loopMode
}

func (m *MockPlayerControllerExt) SetLoopMode(mode player.LoopMode) {
	m.loopMode = mode
}

func (m *MockPlayerControllerExt) ToggleLoopMode() player.LoopMode {
	switch m.loopMode {
	case player.LoopOff:
		m.loopMode = player.LoopSingle
	case player.LoopSingle:
		m.loopMode = player.LoopQueue
	case player.LoopQueue:
		m.loopMode = player.LoopOff
	}
	return m.loopMode
}

func (m *MockPlayerControllerExt) StartPlayback(ctx context.Context, vc player.VoiceConnection, pipeline player.AudioPipeline) error {
	// Mock implementation - returns nil for testing
	return nil
}
