package player

import "sync"

// Manager owns GuildPlayer instances and keeps playback state isolated per guild.
type Manager struct {
	mu            sync.Mutex
	players       map[string]*GuildPlayer
	queueCapacity int
}

// NewManager creates a Manager that gives each GuildPlayer the configured queue capacity.
func NewManager(queueCapacity int) *Manager {
	return &Manager{
		players:       make(map[string]*GuildPlayer),
		queueCapacity: queueCapacity,
	}
}

// GetOrCreate returns the existing player for a guild or creates one if absent.
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

// Get returns the existing player for a guild.
func (m *Manager) Get(guildID string) (*GuildPlayer, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	player, ok := m.players[guildID]
	return player, ok
}

// Remove stops and deletes the player for a guild. It returns false when no player exists.
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
