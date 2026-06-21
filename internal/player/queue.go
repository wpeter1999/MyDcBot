package player

import (
	"math/rand"
	"sync"
)

// Queue 以 FIFO 順序儲存歌曲，並提供安全的快照功能供唯讀指令使用。
type Queue struct {
	mu       sync.Mutex
	songs    []Song
	capacity int
}

// NewQueue 建立一個基於 slice 的 FIFO 佇列，具有固定容量限制。
func NewQueue(capacity int) *Queue {
	if capacity < 0 {
		capacity = 0
	}
	return &Queue{
		songs:    make([]Song, 0, capacity),
		capacity: capacity,
	}
}

// Enqueue 將歌曲加入佇列尾端。
func (q *Queue) Enqueue(song Song) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.songs) >= q.capacity {
		return ErrQueueFull
	}

	q.songs = append(q.songs, song)
	return nil
}

// Dequeue 移除並回傳 FIFO 順序中的下一首歌曲。
func (q *Queue) Dequeue() (Song, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.songs) == 0 {
		return Song{}, false
	}

	song := q.songs[0]
	copy(q.songs, q.songs[1:])
	q.songs = q.songs[:len(q.songs)-1]
	return song, true
}

// Snapshot 回傳已加入佇列的歌曲副本，不會消費佇列內容。
func (q *Queue) Snapshot() []Song {
	q.mu.Lock()
	defer q.mu.Unlock()

	snapshot := make([]Song, len(q.songs))
	copy(snapshot, q.songs)
	return snapshot
}

// Len 回傳目前佇列中的歌曲數量。
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.songs)
}

// Clear 移除所有已加入佇列的歌曲。
func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.songs = q.songs[:0]
}

// EnqueueFront 將歌曲加入佇列最前面（用於循環播放）。
func (q *Queue) EnqueueFront(song Song) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.songs) >= q.capacity {
		return ErrQueueFull
	}

	// 在最前面插入歌曲
	q.songs = append([]Song{song}, q.songs...)
	return nil
}

// Shuffle 打亂佇列中的歌曲順序
func (q *Queue) Shuffle() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 使用 Fisher-Yates shuffle 算法
	for i := len(q.songs) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		q.songs[i], q.songs[j] = q.songs[j], q.songs[i]
	}
}
