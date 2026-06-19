package player

import "errors"

var (
	// ErrQueueFull 當佇列已達容量上限而無法加入新歌曲時回傳此錯誤。
	ErrQueueFull = errors.New("player queue is full")

	// ErrPlayerStopped 當播放器已停止而無法執行操作時回傳此錯誤。
	ErrPlayerStopped = errors.New("player is stopped")
)
