package command

import (
	"testing"
)

func TestMusicServiceGlobalVariable(t *testing.T) {
	t.Run("GetMusicService存在", func(t *testing.T) {
		// 驗證全域變數可以被存取
		// Note: 實際使用中應該透過依賴注入而非全域變數

		service := GetMusicService()
		_ = service // 可能為 nil，這是預期的
	})
}

func TestSetMusicService(t *testing.T) {
	t.Run("設定和獲取MusicService", func(t *testing.T) {
		// 使用真實的 mock service 測試
		originalService := GetMusicService()
		defer SetMusicService(originalService) // 測試後恢復

		mockService := newMockMusicService()
		SetMusicService(mockService)

		got := GetMusicService()
		if got != mockService {
			t.Error("GetMusicService() 應該返回設定的 service")
		}

		// Cleanup
		SetMusicService(nil)
		if GetMusicService() != nil {
			t.Error("設定為 nil 後應該返回 nil")
		}
	})
}
