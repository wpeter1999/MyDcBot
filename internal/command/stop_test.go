package command

import (
	"testing"

	"discordbot/internal/player"
)

func TestStopCommand(t *testing.T) {
	t.Run("StopCommand_定義正確", func(t *testing.T) {
		if StopCommand == nil {
			t.Fatal("StopCommand should not be nil")
		}

		if StopCommand.Command.Name != "stop" {
			t.Errorf("StopCommand.Name = %v, want 'stop'", StopCommand.Command.Name)
		}

		if StopCommand.Handler == nil {
			t.Error("StopCommand.Handler should not be nil")
		}
	})
}

func TestExecuteStop(t *testing.T) {
	t.Run("停止播放應該清空佇列", func(t *testing.T) {
		mockPlayer := &MockPlayerControllerExt{
			queue: []player.Song{
				{Title: "歌曲1"},
				{Title: "歌曲2"},
			},
		}

		// 設定當前歌曲
		mockPlayer.SetCurrentSong(player.Song{Title: "當前歌曲"})

		// 執行停止
		mockPlayer.Stop()

		// 驗證佇列被清空
		if mockPlayer.QueueLen() != 0 {
			t.Error("停止後佇列應該為空")
		}

		// 驗證當前歌曲被清除
		_, hasSong := mockPlayer.CurrentSong()
		if hasSong {
			t.Error("停止後不應該有當前歌曲")
		}
	})
}
