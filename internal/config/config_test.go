package config

import "testing"

// TestLoad 測試 Load 會從環境變數載入必要與選用設定
func TestLoad(t *testing.T) {
	t.Setenv("BOT_TOKEN", "bot-token")
	t.Setenv("GUILD_ID", "guild-id")
	t.Setenv("CWA_API_KEY", "cwa-key")

	cfg := Load()

	if cfg.BotToken != "bot-token" {
		t.Errorf("BotToken 應為 bot-token，實際為 %q", cfg.BotToken)
	}
	if cfg.GuildID != "guild-id" {
		t.Errorf("GuildID 應為 guild-id，實際為 %q", cfg.GuildID)
	}
	if cfg.CwaApiKey != "cwa-key" {
		t.Errorf("CwaApiKey 應為 cwa-key，實際為 %q", cfg.CwaApiKey)
	}
}

// TestLoadFromEnv 測試 loadFromEnv 會讀取 Bot、Guild 與 CWA API 設定
func TestLoadFromEnv(t *testing.T) {
	t.Setenv("BOT_TOKEN", "bot-token")
	t.Setenv("GUILD_ID", "guild-id")
	t.Setenv("CWA_API_KEY", "cwa-key")

	cfg := loadFromEnv()

	if cfg.BotToken != "bot-token" {
		t.Errorf("BotToken 應為 bot-token，實際為 %q", cfg.BotToken)
	}
	if cfg.GuildID != "guild-id" {
		t.Errorf("GuildID 應為 guild-id，實際為 %q", cfg.GuildID)
	}
	if cfg.CwaApiKey != "cwa-key" {
		t.Errorf("CwaApiKey 應為 cwa-key，實際為 %q", cfg.CwaApiKey)
	}
}

// TestLoadFromEnv_AllowsEmptyOptionalValues 測試選用設定允許空值
func TestLoadFromEnv_AllowsEmptyOptionalValues(t *testing.T) {
	t.Setenv("BOT_TOKEN", "bot-token")
	t.Setenv("GUILD_ID", "")
	t.Setenv("CWA_API_KEY", "")

	cfg := loadFromEnv()

	if cfg.BotToken != "bot-token" {
		t.Errorf("BotToken 應為 bot-token，實際為 %q", cfg.BotToken)
	}
	if cfg.GuildID != "" {
		t.Errorf("GuildID 應可為空，實際為 %q", cfg.GuildID)
	}
	if cfg.CwaApiKey != "" {
		t.Errorf("CwaApiKey 應可為空，實際為 %q", cfg.CwaApiKey)
	}
}
