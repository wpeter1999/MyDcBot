# 🎵 Discord 音樂機器人

使用 Go 語言開發的功能完整的 Discord 音樂機器人，支援 YouTube 播放、播放清單、下載等功能。

## ✨ 功能特色

### 🎵 播放功能
- ✅ **播放 YouTube 音樂** - 支援關鍵字搜尋和直接 URL
- ✅ **播放清單支援** - 自動載入並播放整個 YouTube 播放清單
- ✅ **自動播放下一首** - 歌曲結束後自動繼續播放佇列
- ✅ **播放控制** - 暫停、繼續、跳過、停止
- ✅ **佇列管理** - 查看和管理播放佇列
- ✅ **當前播放** - 顯示正在播放的歌曲資訊

### 📥 下載功能
- ✅ **音訊下載** - 下載 YouTube 音訊檔案
- ✅ **多格式支援** - MP3, M4A, Opus, FLAC, WAV
- ✅ **自動限制** - 檔案大小和時長限制
- ✅ **直接上傳** - 自動上傳到 Discord（< 25MB）

### 🎯 技術特點
- ✅ **Lavalink 整合** - 使用 Lavalink 處理音訊串流
- ✅ **yt-dlp 提取** - 繞過 YouTube 限制
- ✅ **SoundCloud 備用** - YouTube 失敗時自動切換
- ✅ **Docker 部署** - 完整的 Docker Compose 配置

## 📋 可用指令

| 指令 | 說明 | 範例 |
|------|------|------|
| `/play` | 播放音樂或播放清單 | `/play query: clear mind` |
| `/pause` | 暫停/繼續播放 | `/pause` |
| `/skip` | 跳到下一首 | `/skip` |
| `/stop` | 停止播放並清空佇列 | `/stop` |
| `/queue` | 顯示播放佇列 | `/queue` |
| `/nowplaying` | 顯示當前播放 | `/nowplaying` |
| `/download` | 下載音訊檔案 | `/download format: mp3-320 url: [URL]` |
| `/help` | 顯示使用說明 | `/help` |

## 🚀 快速開始

### 前置需求

- Docker 和 Docker Compose
- Discord Bot Token
- Discord Guild ID（用於測試）

### 安裝步驟

1. **Clone 專案**
```bash
git clone <repository-url>
cd DiscordBot
```

2. **設定環境變數**
```bash
cp .env.example .env
```

編輯 `.env` 並填入：
```env
BOT_TOKEN=your_discord_bot_token_here
GUILD_ID=your_test_guild_id_here
```

3. **啟動服務**
```bash
docker compose up -d
```

4. **查看日誌**
```bash
docker compose logs -f bot
```

5. **停止服務**
```bash
docker compose down
```

## 🔧 開發

### 進入開發環境

```bash
# 啟動開發容器
docker compose up -d workspace

# 進入容器
docker compose exec workspace bash

# 編譯並運行
cd /workspace
go build -o bin/bot ./cmd/bot
./bin/bot
```

### 運行測試

```bash
# 在容器內
docker compose exec workspace go test ./...

# 運行特定測試
docker compose exec workspace go test ./internal/command -v
```

### 清理舊指令

```bash
# 清理 Guild 指令
docker compose exec workspace go run ./cmd/cleanup -clean

# 清理全域指令
docker compose exec workspace go run ./cmd/cleanup -clean -global

# 列出現有指令
docker compose exec workspace go run ./cmd/cleanup -list
```

## 📁 專案結構

```
DiscordBot/
├── cmd/
│   ├── bot/              # 主程式入口
│   └── cleanup/          # 指令清理工具
├── internal/
│   ├── bot/              # Bot 核心和事件處理
│   ├── command/          # Slash Commands 實現
│   │   ├── play.go       # 播放指令（支援播放清單）
│   │   ├── download.go   # 下載指令
│   │   ├── queue.go      # 佇列管理
│   │   └── ...
│   ├── config/           # 配置管理
│   ├── player/           # 播放器和佇列邏輯
│   └── youtube/          # YouTube 解析
├── docker-compose.yml    # Docker 服務配置
├── Dockerfile           # Bot 容器映像
└── README.md            # 本文件
```

## ⚙️ 架構設計

### 系統架構

```
Discord Client
     ↓
Discord Bot (Go)
     ↓
  ┌─────────────┐
  │  Commands   │
  └─────────────┘
     ↓
  ┌─────────────┐
  │  Lavalink   │ ←→ yt-dlp
  └─────────────┘
     ↓
  Audio Stream
```

### 播放流程

```
1. 用戶執行 /play
2. 檢測是否為播放清單
3. 使用 yt-dlp 提取音訊 URL
4. 透過 Lavalink 播放音訊
5. 歌曲結束後自動播放下一首
```

### 下載流程

```
1. 用戶執行 /download
2. yt-dlp 下載並轉換格式
3. 檢查檔案大小（< 25MB）
4. 上傳到 Discord
5. 自動清理暫存檔案
```

## 🔐 所需權限

Bot 需要以下 Discord 權限：

### 文字權限
- View Channels (查看頻道)
- Send Messages (發送訊息)
- Embed Links (嵌入連結)
- Attach Files (附加檔案)
- Use Slash Commands (使用斜線指令)

### 語音權限
- Connect (連接)
- Speak (說話)
- Use Voice Activity (使用語音活動)

### OAuth2 URL
```
https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=36702752&scope=bot%20applications.commands
```

## 🎯 功能限制

- **檔案大小**：下載功能限制 25 MB（Discord 限制）
- **影片時長**：下載功能限制 10 分鐘
- **音訊格式**：支援 MP3, M4A, Opus, FLAC, WAV

## 🐛 已知問題

- 私人播放清單無法訪問（需要 YouTube 登入）
- 部分 YouTube 影片可能因地區限制無法播放

## 🛠️ 疑難排解

### Bot 無法連接

檢查：
1. Bot Token 是否正確
2. Bot 是否已加入伺服器
3. Lavalink 服務是否正常運行

```bash
docker compose ps
docker compose logs lavalink
```

### 無法播放音樂

檢查：
1. 是否在語音頻道內
2. Bot 是否有正確的權限
3. Lavalink 日誌是否有錯誤

```bash
docker compose logs bot
```

### 指令不顯示

執行清理並重啟：
```bash
docker compose exec workspace go run ./cmd/cleanup -clean
docker compose restart bot
```

## 📝 更新日誌

### v2.0.0 (2026-06-20)
- ✅ 完整重寫使用 Lavalink
- ✅ 新增播放清單支援
- ✅ 新增下載功能
- ✅ 新增 /help 指令
- ✅ 自動播放下一首
- ✅ 改善佇列顯示

### v1.0.0 (2024-06-18)
- ✅ 初始版本
- ✅ 基本播放功能

## 🤝 貢獻

歡迎提交 Pull Request 或回報問題！

## 📄 授權

MIT License

## 🙏 致謝

- [disgo](https://github.com/disgoorg/disgo) - Discord API 庫
- [disgolink](https://github.com/disgoorg/disgolink) - Lavalink 客戶端
- [Lavalink](https://github.com/lavalink-devs/Lavalink) - 音訊播放服務
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - YouTube 下載工具
