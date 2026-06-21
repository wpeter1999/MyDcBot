package player

import (
	"errors"
	"testing"
	"time"
)

// TestGuildPlayer_EnqueueAndQueueSnapshot 測試播放器的 enqueue 與 queue snapshot 功能。
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

// TestGuildPlayer_CurrentSongState 測試播放器的目前播放歌曲狀態管理。
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

// TestGuildPlayer_TogglePause 測試播放器的暫停切換功能。
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

// TestGuildPlayer_SkipIsNonBlockingSignal 測試 Skip 是非阻塞訊號且不會重複送出。
func TestGuildPlayer_SkipIsNonBlockingSignal(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)

	player.Skip()
	if !player.HasPendingSkip() {
		t.Fatal("第一次 Skip 應成功送出 signal")
	}
	player.Skip()
	// 第二次 Skip 不應阻塞

	select {
	case <-player.SkipChan():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Skip 應送出可接收的 signal")
	}
}

// TestGuildPlayer_StopIsIdempotentAndClosesDone 測試 Stop 是冪等的且會關閉 Done channel。
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

// TestGuildPlayer_EnqueueReturnsErrorAfterStop 測試播放器停止後 Enqueue 會回傳錯誤。
func TestGuildPlayer_EnqueueReturnsErrorAfterStop(t *testing.T) {
	player := NewGuildPlayer("guild-1", 50)
	player.Stop()

	err := player.Enqueue(Song{Title: "Song A"})
	if !errors.Is(err, ErrPlayerStopped) {
		t.Fatalf("Stop 後 Enqueue 應回傳 ErrPlayerStopped，實際為 %v", err)
	}
}
