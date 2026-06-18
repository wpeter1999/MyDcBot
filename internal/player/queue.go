package player

import "sync"

// Queue stores songs in FIFO order and exposes safe snapshots for read-only commands.
type Queue struct {
	mu       sync.Mutex
	songs    []Song
	capacity int
}

// NewQueue creates a slice-backed FIFO queue with a fixed capacity.
func NewQueue(capacity int) *Queue {
	if capacity < 0 {
		capacity = 0
	}
	return &Queue{
		songs:    make([]Song, 0, capacity),
		capacity: capacity,
	}
}

// Enqueue appends a song to the end of the queue.
func (q *Queue) Enqueue(song Song) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.songs) >= q.capacity {
		return ErrQueueFull
	}

	q.songs = append(q.songs, song)
	return nil
}

// Dequeue removes and returns the next song in FIFO order.
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

// Snapshot returns a copy of queued songs without consuming them.
func (q *Queue) Snapshot() []Song {
	q.mu.Lock()
	defer q.mu.Unlock()

	snapshot := make([]Song, len(q.songs))
	copy(snapshot, q.songs)
	return snapshot
}

// Len returns the current number of queued songs.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.songs)
}

// Clear removes all queued songs.
func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.songs = q.songs[:0]
}
