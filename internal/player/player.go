package player

import "sync"

// GuildPlayer 保存單一 Discord Guild 的播放狀態與控制訊號。
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

// NewGuildPlayer 建立指定 Guild 的播放器，並初始化佇列、skip signal 與 done channel。
func NewGuildPlayer(guildID string, queueCapacity int) *GuildPlayer {
	return &GuildPlayer{
		guildID: guildID,
		queue:   NewQueue(queueCapacity),
		skip:    make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
}

// GuildID 回傳此播放器所屬的 Discord Guild ID。
func (p *GuildPlayer) GuildID() string {
	return p.guildID
}

// Enqueue 將歌曲加入此 Guild 的播放佇列；播放器停止後會回傳 ErrPlayerStopped。
func (p *GuildPlayer) Enqueue(song Song) error {
	p.mu.RLock()
	stopped := p.stopped
	p.mu.RUnlock()
	if stopped {
		return ErrPlayerStopped
	}

	return p.queue.Enqueue(song)
}

// QueueSnapshot 回傳目前佇列歌曲的複本，不會消費佇列內容。
func (p *GuildPlayer) QueueSnapshot() []Song {
	return p.queue.Snapshot()
}

// QueueLen 回傳目前佇列中的歌曲數量。
func (p *GuildPlayer) QueueLen() int {
	return p.queue.Len()
}

// SetCurrentSong 記錄目前正在播放的歌曲。
func (p *GuildPlayer) SetCurrentSong(song Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	current := song
	p.currentSong = &current
}

// CurrentSong 回傳目前播放歌曲的複本；沒有歌曲時 ok 會是 false。
func (p *GuildPlayer) CurrentSong() (Song, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.currentSong == nil {
		return Song{}, false
	}
	return *p.currentSong, true
}

// ClearCurrentSong 清除目前播放歌曲的狀態。
func (p *GuildPlayer) ClearCurrentSong() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.currentSong = nil
}

// TogglePause 切換暫停狀態，並回傳切換後的新狀態。
func (p *GuildPlayer) TogglePause() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.paused = !p.paused
	return p.paused
}

// IsPaused 回傳播放器目前是否處於暫停狀態。
func (p *GuildPlayer) IsPaused() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.paused
}

// Skip 以非阻塞方式送出跳過訊號；若已有 pending 訊號或播放器已停止則回傳 false。
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

// SkipChan 回傳唯讀 skip channel，供後續播放迴圈監聽跳過訊號。
func (p *GuildPlayer) SkipChan() <-chan struct{} {
	return p.skip
}

// Stop 將播放器標記為停止、清空狀態並只關閉一次 Done channel。
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

// IsStopped 回傳 Stop 是否已經被呼叫過。
func (p *GuildPlayer) IsStopped() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.stopped
}

// Done 回傳播放器停止時會被關閉的 channel，供外部等待清理完成。
func (p *GuildPlayer) Done() <-chan struct{} {
	return p.done
}
