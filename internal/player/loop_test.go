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
		{"設定單曲循環一次", LoopSingleOnce},
		{"設定單曲無限循環", LoopSingleInfinite},
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

// TestLoopMode_ToggleOff 測試從關閉切換到單曲循環一次
func TestLoopMode_ToggleOff(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)

	// 從 LoopOff 切換應該變成 LoopSingleOnce
	newMode := player.ToggleLoopMode()

	if newMode != LoopSingleOnce {
		t.Errorf("expected LoopSingleOnce after toggle from LoopOff, got %v", newMode)
	}

	if player.GetLoopMode() != LoopSingleOnce {
		t.Errorf("expected player loop mode to be LoopSingleOnce, got %v", player.GetLoopMode())
	}
}

// TestLoopMode_ToggleSingleOnce 測試從單曲循環一次切換到單曲無限循環
func TestLoopMode_ToggleSingleOnce(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)
	player.SetLoopMode(LoopSingleOnce)

	// 從 LoopSingleOnce 切換應該變成 LoopSingleInfinite
	newMode := player.ToggleLoopMode()

	if newMode != LoopSingleInfinite {
		t.Errorf("expected LoopSingleInfinite after toggle from LoopSingleOnce, got %v", newMode)
	}

	if player.GetLoopMode() != LoopSingleInfinite {
		t.Errorf("expected player loop mode to be LoopSingleInfinite, got %v", player.GetLoopMode())
	}
}

// TestLoopMode_ToggleSingleInfinite 測試從單曲無限循環切換回關閉
func TestLoopMode_ToggleSingleInfinite(t *testing.T) {
	player := NewGuildPlayer("test-guild", 50)
	player.SetLoopMode(LoopSingleInfinite)

	// 從 LoopSingleInfinite 切換應該變成 LoopOff
	newMode := player.ToggleLoopMode()

	if newMode != LoopOff {
		t.Errorf("expected LoopOff after toggle from LoopSingleInfinite, got %v", newMode)
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
				player.SetLoopMode(LoopSingleOnce)
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
		{LoopSingleOnce, "單曲循環一次"},
		{LoopSingleInfinite, "單曲無限循環"},
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
	player.SetLoopMode(LoopSingleInfinite)

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
		{LoopOff, "🔁"},
		{LoopSingleOnce, "🔂"},
		{LoopSingleInfinite, "🔁"},
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

// TestLoopMode_ButtonStyle 測試取得循環模式按鈕樣式
func TestLoopMode_ButtonStyle(t *testing.T) {
	tests := []struct {
		mode     LoopMode
		expected int
	}{
		{LoopOff, 2},            // Secondary (灰色)
		{LoopSingleOnce, 1},     // Primary (藍色)
		{LoopSingleInfinite, 3}, // Success (綠色)
	}

	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			got := tt.mode.ButtonStyle()
			if got != tt.expected {
				t.Errorf("expected button style %d, got %d", tt.expected, got)
			}
		})
	}
}
