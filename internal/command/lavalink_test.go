package command

import (
	"testing"
)

func TestGetLavalinkClient(t *testing.T) {
	t.Run("初始化時應該為nil", func(t *testing.T) {
		// Reset global state
		lavalinkClient = nil

		got := GetLavalinkClient()
		if got != nil {
			t.Error("GetLavalinkClient() should return nil when not set")
		}
	})
}

func TestSetLavalinkClient(t *testing.T) {
	t.Run("設定後應該可以獲取", func(t *testing.T) {
		// 保存原始狀態
		originalClient := GetLavalinkClient()
		defer func() {
			lavalinkClient = originalClient
		}()

		// 測試設定為 nil
		SetLavalinkClient(nil)
		got := GetLavalinkClient()
		if got != nil {
			t.Error("設定為 nil 後，GetLavalinkClient() 應該返回 nil")
		}

		// Cleanup
		lavalinkClient = nil
	})
}
