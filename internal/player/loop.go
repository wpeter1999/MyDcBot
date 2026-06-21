package player

// LoopMode 定義循環播放模式
type LoopMode int

const (
	// LoopOff 關閉循環（預設）
	LoopOff LoopMode = iota
	// LoopSingle 單曲循環
	LoopSingle
	// LoopQueue 佇列循環
	LoopQueue
)

// String 回傳 LoopMode 的字串表示
func (m LoopMode) String() string {
	switch m {
	case LoopOff:
		return "關閉"
	case LoopSingle:
		return "單曲循環"
	case LoopQueue:
		return "佇列循環"
	default:
		return "未知"
	}
}

// Icon 回傳 LoopMode 的圖示
func (m LoopMode) Icon() string {
	switch m {
	case LoopOff:
		return "➡️"
	case LoopSingle:
		return "🔂"
	case LoopQueue:
		return "🔁"
	default:
		return "❓"
	}
}
