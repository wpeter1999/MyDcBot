package command

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// PauseCommand 定義 /pause 指令。
var PauseCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "pause",
		Description: "暫停或繼續播放",
	},
	Handler: pauseCommandHandler,
}

// ExecutePauseToggle 切換暫停/繼續狀態（核心業務邏輯）
// 返回：(新的暫停狀態, 是否正在播放, 錯誤)
func ExecutePauseToggle(guildID snowflake.ID) (newPaused bool, isPlaying bool, err error) {
	// 獲取當前播放狀態
	isPlaying, isPaused, _ := GetPlayerState(guildID)

	if !isPlaying {
		return false, false, nil
	}

	// 切換暫停狀態
	newPauseState := !isPaused
	err = PausePlayback(guildID, newPauseState)
	if err != nil {
		return false, true, err
	}

	return newPauseState, true, nil
}

// pauseCommandHandler 處理 /pause 指令，切換暫停/繼續狀態。
func pauseCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	guildID := *event.GuildID()

	newPaused, isPlaying, err := ExecutePauseToggle(guildID)

	if !isPlaying {
		respond(event, "目前沒有播放任何歌曲。")
		return
	}

	if err != nil {
		respond(event, fmt.Sprintf("❌ 操作失敗：%v", err))
		return
	}

	if newPaused {
		RespondWithControlButton(event, "⏸️ 已暫停播放。")
	} else {
		RespondWithControlButton(event, "▶️ 已繼續播放。")
	}
}
