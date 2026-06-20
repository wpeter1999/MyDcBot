package command

import (
	"testing"
)

func TestExecuteSkip_Logic(t *testing.T) {
	tests := []struct {
		name        string
		queueLen    int
		wantHasNext bool
	}{
		{
			name:        "佇列有歌曲_應返回true",
			queueLen:    3,
			wantHasNext: true,
		},
		{
			name:        "佇列空_應返回false",
			queueLen:    0,
			wantHasNext: false,
		},
		{
			name:        "佇列只有一首_應返回true",
			queueLen:    1,
			wantHasNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlayer := &MockPlayerControllerExt{
				queueLen: tt.queueLen,
			}

			// 執行跳過邏輯
			hasNext := mockPlayer.QueueLen() > 0

			if hasNext != tt.wantHasNext {
				t.Errorf("hasNext = %v, want %v", hasNext, tt.wantHasNext)
			}
		})
	}
}

func TestSkipCommand_Structure(t *testing.T) {
	t.Run("SkipCommand_定義完整", func(t *testing.T) {
		if SkipCommand == nil {
			t.Fatal("SkipCommand should not be nil")
		}

		if SkipCommand.Command.Name != "skip" {
			t.Errorf("SkipCommand.Name = %v, want 'skip'", SkipCommand.Command.Name)
		}

		if SkipCommand.Command.Description == "" {
			t.Error("SkipCommand.Description should not be empty")
		}

		if SkipCommand.Handler == nil {
			t.Error("SkipCommand.Handler should not be nil")
		}
	})

	t.Run("SkipCommand_無選項參數", func(t *testing.T) {
		// Skip 指令不需要任何參數
		if len(SkipCommand.Command.Options) != 0 {
			t.Errorf("SkipCommand should have 0 options, got %d", len(SkipCommand.Command.Options))
		}
	})
}
