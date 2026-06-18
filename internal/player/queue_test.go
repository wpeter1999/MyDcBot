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
