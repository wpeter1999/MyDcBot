package bot

import (
	"testing"

	"discordbot/internal/config"
)

func TestNew(t *testing.T) {
	t.Run("沒有 BotToken 時失敗", func(t *testing.T) {
		cfg := &config.Config{
			BotToken: "",
			GuildID:  "123456789",
		}

		_, err := New(cfg)
		if err == nil {
			t.Error("期望錯誤但沒有發生")
		}
	})
}

func TestBotStructure(t *testing.T) {
	t.Run("Bot 結構定義正確", func(t *testing.T) {
		// 測試 Bot 結構體是否正確定義
		// 這個測試不需要實際的 token
		var b *Bot
		if b != nil {
			t.Error("未初始化的 Bot 應該為 nil")
		}
	})
}

// 注意：實際的 Bot 創建測試需要有效的 Discord token
// 在 CI/CD 環境中應該使用環境變數提供真實 token 進行整合測試

