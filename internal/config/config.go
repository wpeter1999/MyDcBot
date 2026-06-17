package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 存放應用程式配置
type Config struct {
	BotToken  string
	GuildID   string
	CwaApiKey string
}

// Load 載入 .env 並從環境變數讀取配置
func Load() *Config {
	// 嘗試載入 .env（如果存在）
	_ = godotenv.Load()

	cfg := loadFromEnv()
	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN is not set in .env or environment variables")
	}

	return cfg
}

func loadFromEnv() *Config {
	return &Config{
		BotToken:  os.Getenv("BOT_TOKEN"),
		GuildID:   os.Getenv("GUILD_ID"),
		CwaApiKey: os.Getenv("CWA_API_KEY"),
	}
}
