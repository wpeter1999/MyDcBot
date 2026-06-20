package command

import (
	"context"
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

// SkipCommand 定義 /skip 指令。
var SkipCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "skip",
		Description: "跳過目前播放的歌曲",
	},
	Handler: skipCommandHandler,
}

// skipCommandHandler 處理 /skip 指令，跳過目前播放的歌曲並播放下一首。
func skipCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	if musicService == nil {
		respond(event, "音樂服務尚未初始化。")
		return
	}

	guildID := event.GuildID().String()
	guildPlayer := musicService.GetOrCreatePlayer(guildID)

	// 檢查是否有下一首歌
	if guildPlayer.QueueLen() == 0 {
		RespondWithControlButton(event, "⏭️ 已跳過當前歌曲，但佇列中沒有下一首歌曲了。")
		// 停止播放
		lavalinkClient := GetLavalinkClient()
		if lavalinkClient != nil {
			player := lavalinkClient.Player(*event.GuildID())
			player.Update(context.Background(), lavalink.WithNullTrack())
		}
		guildPlayer.ClearCurrentSong()
		return
	}

	// 停止當前播放
	lavalinkClient := GetLavalinkClient()
	if lavalinkClient != nil {
		player := lavalinkClient.Player(*event.GuildID())
		player.Update(context.Background(), lavalink.WithNullTrack())
	}
	guildPlayer.ClearCurrentSong()

	// 獲取 voice channel
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		respond(event, "⚠️ 無法找到語音頻道")
		return
	}

	channelID := *voiceState.ChannelID

	// 先回應跳過
	RespondWithControlButton(event, "⏭️ 正在跳過...")

	// 異步播放下一首（會自動重試失敗的歌曲）
	go func() {
		playedSong, err := PlayNextSongFromQueue(event.Client(), *event.GuildID(), channelID)
		if err != nil {
			log.Printf("Skip: Failed to play any song from queue: %v", err)
			// 發送失敗訊息
			sendFollowupMessage(event, "❌ 無法播放佇列中的任何歌曲")
		} else if playedSong != nil {
			log.Printf("Skip: Now playing: %s", playedSong.Title)
			// 發送成功訊息
			sendFollowupWithControlButton(event, fmt.Sprintf("✅ 現在播放：**%s**", playedSong.Title))
		}
	}()
}

// sendFollowupMessage 發送 followup 訊息
func sendFollowupMessage(event *events.ApplicationCommandInteractionCreate, content string) {
	_, err := event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
		Content: content,
	})
	if err != nil {
		log.Printf("Failed to send followup message: %v", err)
	}
}

// sendFollowupWithControlButton 發送帶控制面板按鈕的 followup 訊息
func sendFollowupWithControlButton(event *events.ApplicationCommandInteractionCreate, content string) {
	_, err := event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
		Content: content,
		Components: []discord.ContainerComponent{
			discord.NewActionRow(
				discord.NewPrimaryButton("🎵 控制面板", "/control_panel"),
			),
		},
	})
	if err != nil {
		log.Printf("Failed to send followup message: %v", err)
	}
}
