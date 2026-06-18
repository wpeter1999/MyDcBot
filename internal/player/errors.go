package player

import "errors"

var (
	// ErrQueueFull is returned when a song cannot be enqueued because the queue reached capacity.
	ErrQueueFull = errors.New("player queue is full")

	// ErrPlayerStopped is returned when an operation cannot run because the player is stopped.
	ErrPlayerStopped = errors.New("player is stopped")
)
