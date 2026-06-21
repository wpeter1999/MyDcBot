package player

// LoopMode 定義循環播放模式
type LoopMode int

const (
	// LoopOff 關閉循環（預設）
	LoopOff LoopMode = iota
	// LoopSingleOnce 單曲循環一次
	LoopSingleOnce
	// LoopSingleInfinite 單曲無限循環
	LoopSingleInfinite
)

// String 回傳 LoopMode 的字串表示
func (m LoopMode) String() string {
	switch m {
	case LoopOff:
		return "關閉"
	case LoopSingleOnce:
		return "單曲循環一次"
	case LoopSingleInfinite:
		return "單曲無限循環"
	default:
		return "未知"
	}
}

// Icon 回傳 LoopMode 的圖示
func (m LoopMode) Icon() string {
	switch m {
	case LoopOff:
		return "🔁"
	case LoopSingleOnce:
		return "🔂"
	case LoopSingleInfinite:
		return "🔁"
	default:
		return "❓"
	}
}

// ButtonStyle 回傳 LoopMode 對應的按鈕樣式
func (m LoopMode) ButtonStyle() int {
	switch m {
	case LoopOff:
		return 2 // Secondary (灰色)
	case LoopSingleOnce:
		return 1 // Primary (藍色)
	case LoopSingleInfinite:
		return 3 // Success (綠色)
	default:
		return 2 // Secondary (灰色)
	}
}

