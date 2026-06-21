package player

import (
	"errors"
	"testing"
)

// TestQueue_EnqueueDequeueAndSnapshot 測試佇列的 enqueue、dequeue 與 snapshot 功能。
func TestQueue_EnqueueDequeueAndSnapshot(t *testing.T) {
	queue := NewQueue(3)
	songA := Song{Title: "Song A", URL: "https://example.test/a", RequestedBy: "user-1"}
	songB := Song{Title: "Song B", URL: "https://example.test/b", RequestedBy: "user-2"}

	if err := queue.Enqueue(songA); err != nil {
		t.Fatalf("Enqueue 第一首歌不應回傳錯誤: %v", err)
	}
	if err := queue.Enqueue(songB); err != nil {
		t.Fatalf("Enqueue 第二首歌不應回傳錯誤: %v", err)
	}

	snapshot := queue.Snapshot()
	if len(snapshot) != 2 {
		t.Fatalf("Snapshot 應有 2 首歌，實際有 %d 首", len(snapshot))
	}
	if snapshot[0].Title != "Song A" || snapshot[1].Title != "Song B" {
		t.Fatalf("Snapshot 應保留 enqueue 順序，實際為 %#v", snapshot)
	}
	if queue.Len() != 2 {
		t.Fatalf("Snapshot 不應消費 queue，Len 應為 2，實際為 %d", queue.Len())
	}

	got, ok := queue.Dequeue()
	if !ok {
		t.Fatal("Dequeue 應取得第一首歌")
	}
	if got.Title != "Song A" {
		t.Fatalf("Dequeue 應遵守 FIFO，實際取得 %q", got.Title)
	}
	if queue.Len() != 1 {
		t.Fatalf("Dequeue 後 Len 應為 1，實際為 %d", queue.Len())
	}
}

// TestQueue_ReturnsErrorWhenFull 測試佇列已滿時會回傳 ErrQueueFull 錯誤。
func TestQueue_ReturnsErrorWhenFull(t *testing.T) {
	queue := NewQueue(1)
	if err := queue.Enqueue(Song{Title: "Song A"}); err != nil {
		t.Fatalf("第一次 Enqueue 不應失敗: %v", err)
	}

	err := queue.Enqueue(Song{Title: "Song B"})
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("queue 滿時應回傳 ErrQueueFull，實際為 %v", err)
	}
}

// TestQueue_ClearRemovesAllSongs 測試 Clear 會移除所有已加入佇列的歌曲。
func TestQueue_ClearRemovesAllSongs(t *testing.T) {
	queue := NewQueue(3)
	_ = queue.Enqueue(Song{Title: "Song A"})
	_ = queue.Enqueue(Song{Title: "Song B"})

	queue.Clear()

	if queue.Len() != 0 {
		t.Fatalf("Clear 後 Len 應為 0，實際為 %d", queue.Len())
	}
	if _, ok := queue.Dequeue(); ok {
		t.Fatal("Clear 後 Dequeue 不應取得歌曲")
	}
}

// TestQueue_SnapshotReturnsCopy 測試 Snapshot 回傳的是獨立副本，修改不會影響原佇列。
func TestQueue_SnapshotReturnsCopy(t *testing.T) {
	queue := NewQueue(2)
	_ = queue.Enqueue(Song{Title: "Song A"})

	snapshot := queue.Snapshot()
	snapshot[0].Title = "mutated"

	freshSnapshot := queue.Snapshot()
	if freshSnapshot[0].Title != "Song A" {
		t.Fatalf("Snapshot 應回傳 copy，實際 queue 被改成 %q", freshSnapshot[0].Title)
	}
}

// TestQueue_Shuffle 測試打亂佇列
func TestQueue_Shuffle(t *testing.T) {
	q := NewQueue(10)

	// 加入多首歌曲
	songs := []Song{
		{Title: "Song A"},
		{Title: "Song B"},
		{Title: "Song C"},
		{Title: "Song D"},
		{Title: "Song E"},
	}

	for _, song := range songs {
		if err := q.Enqueue(song); err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
	}

	// 打亂佇列
	q.Shuffle()

	// 檢查佇列長度不變
	if q.Len() != 5 {
		t.Errorf("expected queue length 5, got %d", q.Len())
	}

	// 檢查所有歌曲都還在（雖然順序改變）
	snapshot := q.Snapshot()
	songMap := make(map[string]bool)
	for _, s := range snapshot {
		songMap[s.Title] = true
	}

	for _, expected := range songs {
		if !songMap[expected.Title] {
			t.Errorf("song %s not found after shuffle", expected.Title)
		}
	}
}

// TestQueue_ShuffleEmpty 測試打亂空佇列
func TestQueue_ShuffleEmpty(t *testing.T) {
	q := NewQueue(10)

	// 打亂空佇列不應該 panic
	q.Shuffle()

	if q.Len() != 0 {
		t.Errorf("expected empty queue, got length %d", q.Len())
	}
}

// TestQueue_ShuffleSingleSong 測試打亂只有一首歌的佇列
func TestQueue_ShuffleSingleSong(t *testing.T) {
	q := NewQueue(10)

	song := Song{Title: "Only Song"}
	if err := q.Enqueue(song); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// 打亂
	q.Shuffle()

	// 檢查歌曲還在
	if q.Len() != 1 {
		t.Errorf("expected queue length 1, got %d", q.Len())
	}

	snapshot := q.Snapshot()
	if snapshot[0].Title != "Only Song" {
		t.Errorf("expected 'Only Song', got %s", snapshot[0].Title)
	}
}

// TestQueue_ShuffleConcurrent 測試並行打亂的安全性
func TestQueue_ShuffleConcurrent(t *testing.T) {
	q := NewQueue(100)

	// 加入歌曲
	for i := 0; i < 50; i++ {
		q.Enqueue(Song{Title: "Song"})
	}

	done := make(chan bool)

	// 啟動多個 goroutine 同時打亂
	for i := 0; i < 10; i++ {
		go func() {
			q.Shuffle()
			q.Len()
			q.Snapshot()
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 只要沒有 panic 就算通過
	if q.Len() != 50 {
		t.Errorf("expected queue length 50 after concurrent shuffles, got %d", q.Len())
	}
}
