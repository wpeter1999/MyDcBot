package main

import (
	"fmt"
	"log"

	"discordbot/internal/bot"
	"discordbot/internal/config"
)

type runnableBot interface {
	Start() error
	Stop()
	WaitForShutdown()
}

var (
	loadConfig = config.Load
	newBot     = func(cfg *config.Config) (runnableBot, error) {
		return bot.New(cfg)
	}
	printLine = fmt.Println
	logFatalf = log.Fatalf
	runApp    = run
)

func main() {
	if err := runApp(); err != nil {
		logFatalf("%v", err)
	}
}

func run() error {
	cfg := loadConfig()

	b, err := newBot(cfg)
	if err != nil {
		return fmt.Errorf("error creating bot: %w", err)
	}

	if err := b.Start(); err != nil {
		return fmt.Errorf("error starting bot: %w", err)
	}
	defer b.Stop()

	printLine("Bot is running. Press Ctrl+C to exit.")
	b.WaitForShutdown()
	printLine("Shutting down...")

	return nil
}
