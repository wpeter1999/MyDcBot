package player

import (
	"testing"
)

func TestPlayerErrors(t *testing.T) {
	t.Run("ErrQueueFull_定义正确", func(t *testing.T) {
		if ErrQueueFull == nil {
			t.Fatal("ErrQueueFull should not be nil")
		}

		if ErrQueueFull.Error() == "" {
			t.Error("ErrQueueFull should have error message")
		}
	})

	t.Run("ErrPlayerStopped_定义正确", func(t *testing.T) {
		if ErrPlayerStopped == nil {
			t.Fatal("ErrPlayerStopped should not be nil")
		}

		if ErrPlayerStopped.Error() == "" {
			t.Error("ErrPlayerStopped should have error message")
		}
	})

	t.Run("错误消息唯一性", func(t *testing.T) {
		if ErrQueueFull.Error() == ErrPlayerStopped.Error() {
			t.Error("Error messages should be unique")
		}
	})
}
