package command

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// PauseCommand 定義 /pause 指令。
var PauseCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "pause",
		Description: "暫停或繼續播放",
	},
	Handler: pauseCommandHandler,
}

// pauseCommandHandler 處理 /pause 指令，切換暫停/繼續狀態。
func pauseCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	guildID := *event.GuildID()

	// 獲取當前播放狀態
	isPlaying, isPaused, _ := GetPlayerState(guildID)

	if !isPlaying {
		respond(event, "目前沒有播放任何歌曲。")
		return
	}

	// 切換暫停狀態
	newPauseState := !isPaused
	err := PausePlayback(guildID, newPauseState)
	if err != nil {
		respond(event, fmt.Sprintf("❌ 操作失敗：%v", err))
		return
	}

	if newPauseState {
		RespondWithControlButton(event, "⏸️ 已暫停播放。")
	} else {
		RespondWithControlButton(event, "▶️ 已繼續播放。")
	}
}
