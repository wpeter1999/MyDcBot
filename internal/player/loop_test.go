package player

import (
	"testing"
)

// TestLoopMode_Default 測試預設循環模式應該是關閉
func TestLoopMode_Default(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)

	mode := player.GetLoopMode()
	if mode != LoopOff {
		t.Errorf("expected default loop mode to be LoopOff, got %v", mode)
	}
}

// TestLoopMode_SetAndGet 測試設定和取得循環模式
func TestLoopMode_SetAndGet(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)

	tests := []struct {
		name     string
		loopMode LoopMode
	}{
		{"設定單曲循環", LoopSingle},
		{"設定佇列循環", LoopQueue},
		{"關閉循環", LoopOff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.SetLoopMode(tt.loopMode)

			got := player.GetLoopMode()
			if got != tt.loopMode {
				t.Errorf("expected loop mode %v, got %v", tt.loopMode, got)
			}
		})
	}
}

// TestLoopMode_ToggleOff 測試從關閉切換到單曲循環
func TestLoopMode_ToggleOff(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)

	// 從 LoopOff 切換應該變成 LoopSingle
	newMode := player.ToggleLoopMode()

	if newMode != LoopSingle {
		t.Errorf("expected LoopSingle after toggle from LoopOff, got %v", newMode)
	}

	if player.GetLoopMode() != LoopSingle {
		t.Errorf("expected player loop mode to be LoopSingle, got %v", player.GetLoopMode())
	}
}

// TestLoopMode_ToggleSingle 測試從單曲循環切換到佇列循環
func TestLoopMode_ToggleSingle(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)
	player.SetLoopMode(LoopSingle)

	// 從 LoopSingle 切換應該變成 LoopQueue
	newMode := player.ToggleLoopMode()

	if newMode != LoopQueue {
		t.Errorf("expected LoopQueue after toggle from LoopSingle, got %v", newMode)
	}

	if player.GetLoopMode() != LoopQueue {
		t.Errorf("expected player loop mode to be LoopQueue, got %v", player.GetLoopMode())
	}
}

// TestLoopMode_ToggleQueue 測試從佇列循環切換回關閉
func TestLoopMode_ToggleQueue(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)
	player.SetLoopMode(LoopQueue)

	// 從 LoopQueue 切換應該變成 LoopOff
	newMode := player.ToggleLoopMode()

	if newMode != LoopOff {
		t.Errorf("expected LoopOff after toggle from LoopQueue, got %v", newMode)
	}

	if player.GetLoopMode() != LoopOff {
		t.Errorf("expected player loop mode to be LoopOff, got %v", player.GetLoopMode())
	}
}

// TestLoopMode_ConcurrentAccess 測試並行存取的安全性
func TestLoopMode_ConcurrentAccess(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)

	done := make(chan bool)

	// 啟動多個 goroutine 同時讀寫
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				player.SetLoopMode(LoopSingle)
				player.GetLoopMode()
				player.ToggleLoopMode()
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 只要沒有 panic 就算通過
}

// TestLoopMode_String 測試 LoopMode 的字串表示
func TestLoopMode_String(t *testing.T) {
	tests := []struct {
		mode     LoopMode
		expected string
	}{
		{LoopOff, "關閉"},
		{LoopSingle, "單曲循環"},
		{LoopQueue, "佇列循環"},
		{LoopMode(99), "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.mode.String()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// TestLoopMode_AfterStop 測試停止後循環模式應該重置
func TestLoopMode_AfterStop(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)
	player.SetLoopMode(LoopQueue)

	// 停止播放器
	player.Stop()

	// 停止後應該重置為 LoopOff
	mode := player.GetLoopMode()
	if mode != LoopOff {
		t.Errorf("expected loop mode to reset to LoopOff after stop, got %v", mode)
	}
}

// TestLoopMode_Icon 測試取得循環模式圖示
func TestLoopMode_Icon(t *testing.T) {
	tests := []struct {
		mode     LoopMode
		expected string
	}{
		{LoopOff, "➡️"},
		{LoopSingle, "🔂"},
		{LoopQueue, "🔁"},
	}

	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			got := tt.mode.Icon()
			if got != tt.expected {
				t.Errorf("expected icon %q, got %q", tt.expected, got)
			}
		})
	}
}
