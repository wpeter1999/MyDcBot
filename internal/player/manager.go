package player

import "sync"

// Manager 管理多個 GuildPlayer 實例，確保每個 Guild 的播放狀態相互隔離。
type Manager struct {
	mu            sync.Mutex
	players       map[string]*GuildPlayer
	queueCapacity int
}

// NewManager 建立 Manager 並為每個 GuildPlayer 配置指定的佇列容量。
func NewManager(queueCapacity int) *Manager {
	return &Manager{
		players:       make(map[string]*GuildPlayer),
		queueCapacity: queueCapacity,
	}
}

// GetOrCreate 回傳指定 Guild 的播放器；若不存在則建立新的播放器。
func (m *Manager) GetOrCreate(guildID string) *GuildPlayer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if player, ok := m.players[guildID]; ok {
		return player
	}

	player := NewGuildPlayer(guildID, m.queueCapacity)
	m.players[guildID] = player
	return player
}

// Get 回傳指定 Guild 的現有播放器；不存在時 ok 為 false。
func (m *Manager) Get(guildID string) (*GuildPlayer, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	player, ok := m.players[guildID]
	return player, ok
}

// Remove 停止並刪除指定 Guild 的播放器；播放器不存在時回傳 false。
func (m *Manager) Remove(guildID string) bool {
	m.mu.Lock()
	player, ok := m.players[guildID]
	if ok {
		delete(m.players, guildID)
	}
	m.mu.Unlock()

	if !ok {
		return false
	}

	player.Stop()
	return true
}
