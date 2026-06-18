# Discord Bot (Go)

使用 `discordgo` 開發的 Discord 音樂 Bot 專案。

## 功能特色

✅ **已實作**
- `/ping` - 測試 bot 連線
- `/weather <city>` - 查詢天氣資訊
- `/play <歌曲名稱或 URL>` - 播放 YouTube 音樂
- `/pause` - 暫停播放
- `/skip` - 跳過當前歌曲
- `/stop` - 停止播放並離開語音頻道
- `/queue` - 查看播放佇列
- `/nowplaying` - 查看當前播放歌曲

⚠️ **已知限制**
- **Discord DAVE 協議問題**：2024 年後 Discord 要求語音連線使用 DAVE (Discord Audio/Video E2E) 加密協議
- 目前 Go 語言的 Discord 函式庫尚未有穩定的 DAVE 支援
- 在啟用 DAVE 的伺服器上，bot 無法加入語音頻道（會收到 `websocket: close 4017` 錯誤）
- 等待社群提供穩定的解決方案後會進行更新

## 安裝與設定

### 方法一：Docker（推薦）

1. 建立 `.env` 檔案：

```bash
cp .env.example .env
```

2. 在 `.env` 中填入你的 Bot Token：

```env
BOT_TOKEN=your_discord_bot_token_here
```

（可選）若想在私人伺服器快速測試 Slash Command，可加入 Guild ID：

```env
GUILD_ID=your_test_guild_id
```

3. 建置並啟動 bot：

```bash
docker compose build
docker compose up bot
```

4. 停止容器：

```bash
docker compose down
```

### 方法二：本地開發

1. 安裝依賴：
   - Go 1.22+
   - FFmpeg
   - Python 3 與 yt-dlp

2. 安裝 Go 依賴：

```bash
go mod tidy
```

3. 建立並設定 `.env` 檔案（同上）

4. 啟動 Bot：

```bash
go run ./cmd/bot
```

## Docker 開發工作流程

### 開發環境

進入開發容器進行互動式開發：

```bash
# 啟動開發工作區
docker compose up -d workspace

# 進入容器
docker compose exec workspace bash

# 在容器內執行
go test ./...
go run ./cmd/bot
```

### 測試

```bash
docker compose exec workspace go test ./...
```

## 專案架構

```
cmd/bot/               # 應用程式入口
internal/
  ├── audio/          # 音訊處理管道 (ffmpeg → opus)
  ├── bot/            # Bot 生命週期與事件處理
  ├── command/        # Discord slash commands
  ├── config/         # 配置載入
  ├── player/         # 播放器與佇列管理
  └── youtube/        # YouTube 影片解析
```

## 測試

所有核心功能都有對應的單元測試：

```bash
go test ./...
```

## 技術細節

### DAVE 協議問題

我們嘗試過的解決方案：
1. ✗ 使用 `cartridge-gg/discordgo` fork（需要 libdave C 函式庫，版本不相容）
2. ✗ 手動編譯 libdave（編譯失敗或找不到符號）
3. ✗ 降級到較舊的 DAVE fork 版本（仍然缺少必要函數）

**目前做法**：暫時使用官方 `bwmarrin/discordgo`，等待社群提供穩定的 DAVE 支援。

### 語音播放架構

```
YouTube URL/搜尋
    ↓ yt-dlp
Audio Stream URL
    ↓ ffmpeg
PCM Audio
    ↓ gopus encoder
Opus Packets
    ↓ discordgo
Discord Voice
```

## 貢獻

歡迎提交 Pull Request 或回報問題！

## 授權

MIT License
