package command

import (
	"testing"
)

func TestMusicServiceGlobalVariable(t *testing.T) {
	t.Run("GetMusicService存在", func(t *testing.T) {
		// 验证全局变量可以被访问
		// Note: 实际使用中应该通过依赖注入而非全局变量

		service := GetMusicService()
		_ = service // 可能为 nil，这是预期的
	})
}

func TestSetMusicService(t *testing.T) {
	t.Run("设置和获取MusicService", func(t *testing.T) {
		// Note: 实际测试需要 mock MusicService
		// 这里验证函数签名正确

		// mockService := &MockMusicService{}
		// SetMusicService(mockService)

		// got := GetMusicService()
		// if got != mockService {
		// 	t.Error("GetMusicService() should return the set service")
		// }

		// Cleanup
		SetMusicService(nil)
	})
}
