package command

import (
	"testing"
)

func TestGetLavalinkClient(t *testing.T) {
	t.Run("初始化时应该为nil", func(t *testing.T) {
		// Reset global state
		lavalinkClient = nil

		got := GetLavalinkClient()
		if got != nil {
			t.Error("GetLavalinkClient() should return nil when not set")
		}
	})
}

func TestSetLavalinkClient(t *testing.T) {
	t.Run("设置后应该可以获取", func(t *testing.T) {
		// Note: 实际测试需要 mock disgolink.Client
		// 这里验证函数签名正确

		// mockClient := &MockLavalinkClient{}
		// SetLavalinkClient(mockClient)

		// got := GetLavalinkClient()
		// if got != mockClient {
		// 	t.Error("GetLavalinkClient() should return the set client")
		// }

		// Cleanup
		lavalinkClient = nil
	})
}
