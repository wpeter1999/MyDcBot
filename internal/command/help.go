package command

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// HelpCommand 定義 /help 指令
var HelpCommand = &BotCommand{
	Command: discord.SlashCommandCreate{
		Name:        "help",
		Description: "顯示所有可用指令的說明",
	},
	Handler: helpCommandHandler,
}

// helpCommandHandler 處理 /help 指令
func helpCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	helpMessage := "🎵 **Discord 音樂機器人 - 指令說明**\n\n" +

		"**播放功能**\n" +
		"▶️ `/play query: [搜尋關鍵字或 URL]`\n" +
		"   播放 YouTube 音樂或播放清單\n" +
		"   • 範例：`/play query: clear mind`\n" +
		"   • 範例：`/play query: https://youtube.com/watch?v=xxx`\n" +
		"   • 範例：`/play query: https://youtube.com/playlist?list=xxx`\n\n" +

		"⏸️ `/pause`\n" +
		"   暫停/繼續當前播放\n" +
		"   • 第一次執行：暫停播放\n" +
		"   • 再次執行：繼續播放\n\n" +

		"⏭️ `/skip`\n" +
		"   跳過當前歌曲，播放下一首\n" +
		"   • 如果佇列是空的，將停止播放\n\n" +

		"⏹️ `/stop`\n" +
		"   停止播放並清空佇列\n" +
		"   • 機器人會離開語音頻道\n\n" +

		"**資訊查詢**\n" +
		"📜 `/queue`\n" +
		"   顯示播放佇列\n" +
		"   • 顯示正在播放的歌曲\n" +
		"   • 顯示接下來的歌曲（最多 10 首）\n\n" +

		"🎵 `/nowplaying`\n" +
		"   顯示當前正在播放的歌曲資訊\n\n" +

		"**下載功能**\n" +
		"📥 `/download format: [格式] url: [YouTube URL]`\n" +
		"   下載 YouTube 音訊檔案\n" +
		"   • **格式選項**：\n" +
		"     - MP3 320kbps (推薦) - 平衡品質與大小\n" +
		"     - M4A 256kbps - 較小的檔案\n" +
		"     - Opus 192kbps - 最小的檔案\n" +
		"     - FLAC 無損 - 高品質（檔案較大）\n" +
		"     - WAV 原始 - 未壓縮（檔案最大）\n" +
		"   • **限制**：\n" +
		"     - 檔案大小最大 25 MB\n" +
		"     - 影片時長最多 10 分鐘\n\n" +

		"**使用提示**\n" +
		"💡 播放清單會自動載入所有歌曲\n" +
		"💡 歌曲播放完畢會自動播放下一首\n" +
		"💡 支援中文搜尋關鍵字\n" +
		"💡 YouTube 失敗時會自動切換 SoundCloud\n\n" +

		"❓ 有問題？請聯繫管理員"

	respond(event, helpMessage)
}
