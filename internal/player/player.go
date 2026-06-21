package player

import (
	"context"
	"sync"
)

// PlaybackController 定義播放控制介面，用於啟動播放迴圈。
type PlaybackController interface {
	// StartPlayback 啟動播放迴圈，從佇列取歌並播放。
	// 會阻塞直到播放器停止或 context 被取消。
	StartPlayback(ctx context.Context, vc VoiceConnection, pipeline AudioPipeline) error
}

// VoiceConnection 抽象 Discord 語音連線。
type VoiceConnection interface {
	Speaking(speaking bool) error
	OpusSend() chan<- []byte
	Disconnect() error
}

// AudioPipeline 抽象音訊播放管道。
type AudioPipeline interface {
	Play(ctx context.Context, url string, vc VoiceConnection) error
}

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
	loopMode    LoopMode // 循環播放模式
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

// GuildID 回傳播放器所屬的 Guild ID。
func (p *GuildPlayer) GuildID() string {
	return p.guildID
}

// Enqueue 將歌曲加入佇列尾端，若播放器已停止則回傳 ErrPlayerStopped。
func (p *GuildPlayer) Enqueue(song Song) error {
	p.mu.RLock()
	stopped := p.stopped
	p.mu.RUnlock()
	if stopped {
		return ErrPlayerStopped
	}

	return p.queue.Enqueue(song)
}

// QueueSnapshot 回傳佇列的快照（深拷貝），可用於顯示。
func (p *GuildPlayer) QueueSnapshot() []Song {
	return p.queue.Snapshot()
}

// QueueLen 回傳目前佇列中的歌曲數量。
func (p *GuildPlayer) QueueLen() int {
	return p.queue.Len()
}

// Dequeue 從佇列取出下一首歌曲。
func (p *GuildPlayer) Dequeue() (Song, bool) {
	return p.queue.Dequeue()
}

// SetCurrentSong 設定目前播放的歌曲（線程安全）。
func (p *GuildPlayer) SetCurrentSong(song Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.currentSong = &song
}

// CurrentSong 回傳目前播放的歌曲（若有）。
func (p *GuildPlayer) CurrentSong() (Song, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.currentSong == nil {
		return Song{}, false
	}
	return *p.currentSong, true
}

// ClearCurrentSong 清除目前播放的歌曲。
func (p *GuildPlayer) ClearCurrentSong() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.currentSong = nil
}

// TogglePause 切換暫停狀態，回傳新的暫停狀態。
func (p *GuildPlayer) TogglePause() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.paused = !p.paused
	return p.paused
}

// IsPaused 回傳播放器是否暫停。
func (p *GuildPlayer) IsPaused() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.paused
}

// Skip 發送跳過訊號給播放迴圈。若 skip channel 已滿，則跳過（即不阻塞）。
func (p *GuildPlayer) Skip() {
	select {
	case p.skip <- struct{}{}:
		// 成功送出 skip 訊號
	default:
		// skip channel 已滿，表示已有跳過訊號在等待
	}
}

// HasPendingSkip 回傳是否有待處理的 skip 訊號。
func (p *GuildPlayer) HasPendingSkip() bool {
	select {
	case <-p.skip:
		// 消耗掉 skip 訊號
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
		p.loopMode = LoopOff // 重置循環模式
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

// StartPlayback 實作 PlaybackController 介面，啟動播放迴圈。
func (p *GuildPlayer) StartPlayback(ctx context.Context, vc VoiceConnection, pipeline AudioPipeline) error {
	defer vc.Disconnect()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-p.done:
			return nil
		default:
		}

		// 從佇列取出下一首歌曲
		song, ok := p.queue.Dequeue()
		if !ok {
			// 佇列為空，等待新歌曲或停止訊號
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-p.done:
				return nil
			}
		}

		// 設定目前播放的歌曲
		p.mu.Lock()
		p.currentSong = &song
		p.mu.Unlock()

		// 播放歌曲（使用 StreamURL）
		playCtx, cancel := context.WithCancel(ctx)

		// 監聽 skip 訊號
		go func() {
			select {
			case <-p.skip:
				cancel() // 取消播放
			case <-playCtx.Done():
			}
		}()

		// 實際播放
		err := pipeline.Play(playCtx, song.StreamURL, vc)
		cancel() // 清理 context

		// 清除目前播放的歌曲
		p.mu.Lock()
		p.currentSong = nil
		p.mu.Unlock()

		if err != nil && err != context.Canceled {
			// 播放錯誤（非 skip 造成的取消）
			return err
		}

		// 繼續播放下一首
	}
}

// GetLoopMode 取得目前的循環播放模式
func (p *GuildPlayer) GetLoopMode() LoopMode {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.loopMode
}

// SetLoopMode 設定循環播放模式
func (p *GuildPlayer) SetLoopMode(mode LoopMode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.loopMode = mode
}

// ToggleLoopMode 切換循環播放模式（關閉 -> 單曲循環一次 -> 單曲無限循環 -> 關閉）
// 回傳切換後的模式
func (p *GuildPlayer) ToggleLoopMode() LoopMode {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.loopMode {
	case LoopOff:
		p.loopMode = LoopSingleOnce
	case LoopSingleOnce:
		p.loopMode = LoopSingleInfinite
	case LoopSingleInfinite:
		p.loopMode = LoopOff
	}

	return p.loopMode
}
