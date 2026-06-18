package player

import "sync"

// GuildPlayer owns playback state and control signals for a single Discord guild.
type GuildPlayer struct {
	guildID string
	queue   *Queue
	skip    chan struct{}
	done    chan struct{}

	mu          sync.RWMutex
	currentSong *Song
	paused      bool
	stopped     bool
	stopOnce    sync.Once
}

// NewGuildPlayer creates a stopped-state-aware player for one guild.
func NewGuildPlayer(guildID string, queueCapacity int) *GuildPlayer {
	return &GuildPlayer{
		guildID: guildID,
		queue:   NewQueue(queueCapacity),
		skip:    make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
}

// GuildID returns the Discord guild ID this player belongs to.
func (p *GuildPlayer) GuildID() string {
	return p.guildID
}

// Enqueue adds a song to this player's queue.
func (p *GuildPlayer) Enqueue(song Song) error {
	p.mu.RLock()
	stopped := p.stopped
	p.mu.RUnlock()
	if stopped {
		return ErrPlayerStopped
	}

	return p.queue.Enqueue(song)
}

// QueueSnapshot returns queued songs without consuming them.
func (p *GuildPlayer) QueueSnapshot() []Song {
	return p.queue.Snapshot()
}

// QueueLen returns the number of queued songs.
func (p *GuildPlayer) QueueLen() int {
	return p.queue.Len()
}

// SetCurrentSong records the song currently being played.
func (p *GuildPlayer) SetCurrentSong(song Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	current := song
	p.currentSong = &current
}

// CurrentSong returns a copy of the current song, if one exists.
func (p *GuildPlayer) CurrentSong() (Song, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.currentSong == nil {
		return Song{}, false
	}
	return *p.currentSong, true
}

// ClearCurrentSong clears the current song state.
func (p *GuildPlayer) ClearCurrentSong() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.currentSong = nil
}

// TogglePause flips the paused state and returns the new value.
func (p *GuildPlayer) TogglePause() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.paused = !p.paused
	return p.paused
}

// IsPaused reports whether the player is currently paused.
func (p *GuildPlayer) IsPaused() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.paused
}

// Skip sends a non-blocking skip signal. It returns false if a signal is already pending or the player is stopped.
func (p *GuildPlayer) Skip() bool {
	p.mu.RLock()
	stopped := p.stopped
	p.mu.RUnlock()
	if stopped {
		return false
	}

	select {
	case p.skip <- struct{}{}:
		return true
	default:
		return false
	}
}

// SkipChan exposes skip signals to the playback loop.
func (p *GuildPlayer) SkipChan() <-chan struct{} {
	return p.skip
}

// Stop marks the player stopped, clears state, and closes Done exactly once.
func (p *GuildPlayer) Stop() {
	p.stopOnce.Do(func() {
		p.mu.Lock()
		p.stopped = true
		p.paused = false
		p.currentSong = nil
		p.mu.Unlock()

		p.queue.Clear()
		close(p.done)
	})
}

// IsStopped reports whether Stop has been called.
func (p *GuildPlayer) IsStopped() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.stopped
}

// Done is closed when the player stops.
func (p *GuildPlayer) Done() <-chan struct{} {
	return p.done
}
