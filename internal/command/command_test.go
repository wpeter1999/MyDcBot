package command

import (
	"context"
	"testing"

	"discordbot/internal/player"
)

// ==================== Mock PlayerController ====================

type mockPlayerController struct {
	guildID     string
	queue       []player.Song
	currentSong *player.Song
	paused      bool
	stopped     bool
	loopMode    player.LoopMode
}

func newMockPlayerController(guildID string) *mockPlayerController {
	return &mockPlayerController{
		guildID: guildID,
		queue:   []player.Song{},
	}
}

func (m *mockPlayerController) Enqueue(song player.Song) error {
	if len(m.queue) >= 100 {
		return player.ErrQueueFull
	}
	m.queue = append(m.queue, song)
	return nil
}

func (m *mockPlayerController) Dequeue() (player.Song, bool) {
	if len(m.queue) == 0 {
		return player.Song{}, false
	}
	song := m.queue[0]
	m.queue = m.queue[1:]
	return song, true
}

func (m *mockPlayerController) QueueSnapshot() []player.Song {
	return append([]player.Song{}, m.queue...)
}

func (m *mockPlayerController) QueueLen() int {
	return len(m.queue)
}

func (m *mockPlayerController) CurrentSong() (player.Song, bool) {
	if m.currentSong != nil {
		return *m.currentSong, true
	}
	return player.Song{}, false
}

func (m *mockPlayerController) SetCurrentSong(song player.Song) {
	m.currentSong = &song
}

func (m *mockPlayerController) ClearCurrentSong() {
	m.currentSong = nil
}

func (m *mockPlayerController) TogglePause() bool {
	m.paused = !m.paused
	return m.paused
}

func (m *mockPlayerController) IsPaused() bool {
	return m.paused
}

func (m *mockPlayerController) Skip() {
	// Mock implementation
}

func (m *mockPlayerController) Stop() {
	m.stopped = true
	m.queue = nil
	m.currentSong = nil
}

func (m *mockPlayerController) GetLoopMode() player.LoopMode {
	return m.loopMode
}

func (m *mockPlayerController) SetLoopMode(mode player.LoopMode) {
	m.loopMode = mode
}

func (m *mockPlayerController) ToggleLoopMode() player.LoopMode {
	switch m.loopMode {
	case player.LoopOff:
		m.loopMode = player.LoopSingleOnce
	case player.LoopSingleOnce:
		m.loopMode = player.LoopSingleInfinite
	case player.LoopSingleInfinite:
		m.loopMode = player.LoopOff
	}
	return m.loopMode
}

func (m *mockPlayerController) GuildID() string {
	return m.guildID
}

func (m *mockPlayerController) StartPlayback(ctx context.Context, vc player.VoiceConnection, pipeline player.AudioPipeline) error {
	return nil
}

// ==================== Mock MusicService ====================

type mockMusicService struct {
	players map[string]*mockPlayerController
}

func newMockMusicService() *mockMusicService {
	return &mockMusicService{
		players: make(map[string]*mockPlayerController),
	}
}

func (m *mockMusicService) GetOrCreatePlayer(guildID string) PlayerController {
	if p, ok := m.players[guildID]; ok {
		return p
	}
	p := newMockPlayerController(guildID)
	m.players[guildID] = p
	return p
}

func (m *mockMusicService) RemovePlayer(guildID string) bool {
	if _, ok := m.players[guildID]; ok {
		delete(m.players, guildID)
		return true
	}
	return false
}

// ==================== Tests ====================

func TestMockPlayerController(t *testing.T) {
	t.Run("基本操作", func(t *testing.T) {
		pc := newMockPlayerController("test-guild")

		// 測試 GuildID
		if pc.GuildID() != "test-guild" {
			t.Errorf("期望 GuildID 為 'test-guild'，實際為 '%s'", pc.GuildID())
		}

		// 測試 Enqueue
		song := player.Song{Title: "Test Song", URL: "https://test.com"}
		if err := pc.Enqueue(song); err != nil {
			t.Errorf("Enqueue 失敗: %v", err)
		}

		// 測試 QueueLen
		if pc.QueueLen() != 1 {
			t.Errorf("期望佇列長度為 1，實際為 %d", pc.QueueLen())
		}

		// 測試 SetCurrentSong 和 CurrentSong
		pc.SetCurrentSong(song)
		current, ok := pc.CurrentSong()
		if !ok {
			t.Error("應該有當前歌曲")
		}
		if current.Title != song.Title {
			t.Errorf("期望標題為 '%s'，實際為 '%s'", song.Title, current.Title)
		}

		// 測試 ClearCurrentSong
		pc.ClearCurrentSong()
		_, ok = pc.CurrentSong()
		if ok {
			t.Error("清除後不應該有當前歌曲")
		}

		// 測試 Dequeue
		dequeuedSong, ok := pc.Dequeue()
		if !ok {
			t.Error("Dequeue 應該成功")
		}
		if dequeuedSong.Title != song.Title {
			t.Errorf("期望標題為 '%s'，實際為 '%s'", song.Title, dequeuedSong.Title)
		}

		// 測試 TogglePause
		if !pc.TogglePause() {
			t.Error("第一次 TogglePause 應該返回 true")
		}
		if pc.TogglePause() {
			t.Error("第二次 TogglePause 應該返回 false")
		}

		// 測試 Stop
		pc.Enqueue(song)
		pc.Stop()
		if pc.QueueLen() != 0 {
			t.Error("Stop 後佇列應該為空")
		}
	})
}

func TestMockMusicService(t *testing.T) {
	t.Run("GetOrCreatePlayer 創建新播放器", func(t *testing.T) {
		service := newMockMusicService()
		player1 := service.GetOrCreatePlayer("guild-1")

		if player1 == nil {
			t.Fatal("播放器不應為 nil")
		}

		if player1.GuildID() != "guild-1" {
			t.Errorf("期望 GuildID 為 'guild-1'，實際為 '%s'", player1.GuildID())
		}
	})

	t.Run("GetOrCreatePlayer 返回相同播放器", func(t *testing.T) {
		service := newMockMusicService()
		player1 := service.GetOrCreatePlayer("guild-1")
		player2 := service.GetOrCreatePlayer("guild-1")

		if player1 != player2 {
			t.Error("同一個 Guild 應該返回相同的播放器")
		}
	})

	t.Run("GetOrCreatePlayer 為不同 Guild 創建不同播放器", func(t *testing.T) {
		service := newMockMusicService()
		player1 := service.GetOrCreatePlayer("guild-1")
		player2 := service.GetOrCreatePlayer("guild-2")

		if player1 == player2 {
			t.Error("不同 Guild 應該有不同的播放器")
		}
	})

	t.Run("RemovePlayer 移除播放器", func(t *testing.T) {
		service := newMockMusicService()
		service.GetOrCreatePlayer("guild-1")

		if !service.RemovePlayer("guild-1") {
			t.Error("RemovePlayer 應該返回 true")
		}

		if service.RemovePlayer("guild-1") {
			t.Error("移除不存在的播放器應該返回 false")
		}
	})
}

func TestCommandRegistry(t *testing.T) {
	t.Run("所有核心指令都已註冊", func(t *testing.T) {
		expectedCommands := []string{
			"play", "skip", "pause", "queue",
			"stop", "nowplaying", "download", "help",
		}

		for _, expected := range expectedCommands {
			found := false
			for _, cmd := range CommandRegistry {
				if cmd.Command.CommandName() == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("期望指令 '%s' 已註冊，但未找到", expected)
			}
		}
	})

	t.Run("所有指令都有 Handler", func(t *testing.T) {
		for _, cmd := range CommandRegistry {
			if cmd.Handler == nil {
				t.Errorf("指令 '%s' 的 Handler 為 nil", cmd.Command.CommandName())
			}
		}
	})

	t.Run("所有指令名稱都是唯一的", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, cmd := range CommandRegistry {
			name := cmd.Command.CommandName()
			if seen[name] {
				t.Errorf("指令名稱 '%s' 重複", name)
			}
			seen[name] = true
		}
	})
}

func TestIndividualCommands(t *testing.T) {
	tests := []struct {
		name        string
		command     *BotCommand
		commandName string
	}{
		{"PlayCommand", PlayCommand, "play"},
		{"PauseCommand", PauseCommand, "pause"},
		{"SkipCommand", SkipCommand, "skip"},
		{"StopCommand", StopCommand, "stop"},
		{"QueueCommand", QueueCommand, "queue"},
		{"NowPlayingCommand", NowPlayingCommand, "nowplaying"},
		{"DownloadCommand", DownloadCommand, "download"},
		{"HelpCommand", HelpCommand, "help"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.command == nil {
				t.Fatalf("%s 不應為 nil", tt.name)
			}

			if tt.command.Command.CommandName() != tt.commandName {
				t.Errorf("期望指令名稱為 '%s'，實際為 '%s'",
					tt.commandName, tt.command.Command.CommandName())
			}

			if tt.command.Handler == nil {
				t.Errorf("%s 的 Handler 為 nil", tt.name)
			}
		})
	}
}

func TestBuildYtDlpArgs(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		url          string
		wantContains []string
	}{
		{
			name:   "MP3 320kbps",
			format: "mp3-320",
			url:    "https://youtube.com/watch?v=test",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "mp3",
			},
		},
		{
			name:   "FLAC 無損",
			format: "flac",
			url:    "https://youtube.com/watch?v=test",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "flac",
			},
		},
		{
			name:   "Opus",
			format: "opus-192",
			url:    "https://youtube.com/watch?v=test",
			wantContains: []string{
				"--extract-audio",
				"--audio-format", "opus",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildYtDlpArgs(tt.format, tt.url, "/tmp/test")

			for _, want := range tt.wantContains {
				found := false
				for _, arg := range args {
					if arg == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("期望參數包含 '%s'，但未找到", want)
				}
			}

			// 確保 URL 在參數中
			found := false
			for _, arg := range args {
				if arg == tt.url {
					found = true
					break
				}
			}
			if !found {
				t.Error("URL 應該在參數列表中")
			}
		})
	}
}
