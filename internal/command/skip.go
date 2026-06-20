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

	// 從佇列取出下一首
	nextSong, ok := guildPlayer.Dequeue()
	if !ok {
		respond(event, "⏭️ 佇列中沒有下一首歌曲。")
		return
	}

	// 設定為當前播放
	guildPlayer.SetCurrentSong(nextSong)

	// 獲取 voice channel
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok || voiceState.ChannelID == nil {
		respond(event, "⚠️ 無法找到語音頻道")
		return
	}

	channelID := *voiceState.ChannelID

	RespondWithControlButton(event, fmt.Sprintf("⏭️ 已跳過，正在播放：**%s**", nextSong.Title))

	// 異步播放下一首
	go func() {
		err := JoinVoiceAndPlayWithYtDlp(event.Client(), *event.GuildID(), channelID, nextSong.URL)
		if err != nil {
			log.Printf("Skip: Failed to play next song: %v", err)
			// 嘗試 SoundCloud 備用
			searchQuery := "scsearch:" + nextSong.Title
			err = JoinVoiceAndPlay(event.Client(), *event.GuildID(), channelID, searchQuery)
			if err != nil {
				log.Printf("Skip: Failed with SoundCloud backup: %v", err)
			}
		}
	}()
}
