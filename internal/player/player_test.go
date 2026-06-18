package player

import (
	"errors"
	"testing"
	"time"
)

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

func TestManager_GetOrCreateReturnsSamePlayerForSameGuild(t *testing.T) {
	manager := NewManager(50)

	first := manager.GetOrCreate("guild-1")
	second := manager.GetOrCreate("guild-1")

	if first != second {
		t.Fatal("相同 GuildID 應取得同一個 GuildPlayer")
	}
}

func TestManager_GetOrCreateReturnsDifferentPlayersForDifferentGuilds(t *testing.T) {
	manager := NewManager(50)

	guildA := manager.GetOrCreate("guild-a")
	guildB := manager.GetOrCreate("guild-b")

	if guildA == guildB {
		t.Fatal("不同 GuildID 應取得不同 GuildPlayer")
	}
	if guildA.GuildID() != "guild-a" {
		t.Fatalf("guildA GuildID 應為 guild-a，實際為 %q", guildA.GuildID())
	}
	if guildB.GuildID() != "guild-b" {
		t.Fatalf("guildB GuildID 應為 guild-b，實際為 %q", guildB.GuildID())
	}
}

func TestManager_RemoveStopsAndDeletesPlayer(t *testing.T) {
	manager := NewManager(50)
	player := manager.GetOrCreate("guild-1")

	removed := manager.Remove("guild-1")
	if !removed {
		t.Fatal("Remove 已存在的 player 應回傳 true")
	}
	if !player.IsStopped() {
		t.Fatal("Remove 應停止被移除的 player")
	}
	if _, ok := manager.Get("guild-1"); ok {
		t.Fatal("Remove 後不應再取得 player")
	}
	if manager.Remove("guild-1") {
		t.Fatal("Remove 不存在的 player 應回傳 false")
	}
}

func TestGuildPlayer_EnqueueAndQueueSnapshot(t *testing.T) {
	player := NewGuildPlayer("guild-1", 2)

	if err := player.Enqueue(Song{Title: "Song A"}); err != nil {
		t.Fatalf("Enqueue 不應失敗: %v", err)
	}
	if err := player.Enqueue(Song{Title: "Song B"}); err != nil {
		t.Fatalf("Enqueue 不應失敗: %v", err)
	}

	snapshot := player.QueueSnapshot()
	if len(snapshot) != 2 {
		t.Fatalf("QueueSnapshot 應有 2 首歌，實際有 %d 首", len(snapshot))
	}
	if snapshot[0].Title != "Song A" || snapshot[1].Title != "Song B" {
		t.Fatalf("QueueSnapshot 應保留順序，實際為 %#v", snapshot)
	}
}

func TestGuildPlayer_CurrentSongState(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)
	song := Song{Title: "Song A"}

	if _, ok := player.CurrentSong(); ok {
		t.Fatal("尚未設定 current song 時 CurrentSong 應回傳 ok=false")
	}

	player.SetCurrentSong(song)
	got, ok := player.CurrentSong()
	if !ok {
		t.Fatal("設定 current song 後 CurrentSong 應回傳 ok=true")
	}
	if got.Title != "Song A" {
		t.Fatalf("CurrentSong 應回傳設定的歌曲，實際為 %q", got.Title)
	}

	player.ClearCurrentSong()
	if _, ok := player.CurrentSong(); ok {
		t.Fatal("ClearCurrentSong 後 CurrentSong 應回傳 ok=false")
	}
}

func TestGuildPlayer_TogglePause(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)

	if player.IsPaused() {
		t.Fatal("新 player 預設不應為 paused")
	}
	if paused := player.TogglePause(); !paused {
		t.Fatal("第一次 TogglePause 應切換為 paused=true")
	}
	if !player.IsPaused() {
		t.Fatal("TogglePause 後 IsPaused 應為 true")
	}
	if paused := player.TogglePause(); paused {
		t.Fatal("第二次 TogglePause 應切換為 paused=false")
	}
}

func TestGuildPlayer_SkipIsNonBlockingSignal(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)

	if ok := player.Skip(); !ok {
		t.Fatal("第一次 Skip 應成功送出 signal")
	}
	if ok := player.Skip(); ok {
		t.Fatal("已有 pending skip signal 時，第二次 Skip 應回傳 false 且不可阻塞")
	}

	select {
	case <-player.SkipChan():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Skip 應送出可接收的 signal")
	}
}

func TestGuildPlayer_StopIsIdempotentAndClosesDone(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)
	_ = player.Enqueue(Song{Title: "Song A"})
	player.SetCurrentSong(Song{Title: "Current"})
	player.TogglePause()

	player.Stop()
	player.Stop()

	if !player.IsStopped() {
		t.Fatal("Stop 後 IsStopped 應為 true")
	}
	if player.QueueLen() != 0 {
		t.Fatalf("Stop 應清空 queue，實際 Len 為 %d", player.QueueLen())
	}
	if _, ok := player.CurrentSong(); ok {
		t.Fatal("Stop 應清除 current song")
	}
	if player.IsPaused() {
		t.Fatal("Stop 應重設 paused=false")
	}

	select {
	case <-player.Done():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Stop 應關閉 Done channel")
	}
}

func TestGuildPlayer_EnqueueReturnsErrorAfterStop(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)
	player.Stop()

	err := player.Enqueue(Song{Title: "Song A"})
	if !errors.Is(err, ErrPlayerStopped) {
		t.Fatalf("Stop 後 Enqueue 應回傳 ErrPlayerStopped，實際為 %v", err)
	}
}
