# 專案結構文檔

本文檔說明 Discord 音樂機器人的專案結構、檔案組織和架構決策。

## 目錄結構

```
DiscordBot/
├── cmd/                        # 可執行程式入口
│   ├── bot/                    # 主 Bot 程式
│   │   ├── main.go            # 程式入口點
│   │   └── main_test.go       # 整合測試
│   └── cleanup/                # 指令清理工具
│       └── main.go            # 清理舊指令
│
├── internal/                   # 內部程式碼（不可被外部導入）
│   ├── bot/                    # Bot 核心模組
│   │   ├── bot.go             # Bot 初始化與生命週期
│   │   ├── bot_test.go        # Bot 單元測試
│   │   ├── lavalink_handlers.go  # Lavalink 事件處理器
│   │   └── cleanup.go         # 指令清理邏輯
│   │
│   ├── command/                # Discord Slash Commands
│   │   ├── command.go         # 指令介面定義
│   │   ├── command_test.go    # 指令測試
│   │   ├── registry.go        # 指令註冊表
│   │   ├── registrar.go       # 指令註冊器
│   │   ├── registrar_test.go  # 註冊器測試
│   │   ├── music.go           # 音樂服務介面
│   │   ├── voice.go           # 語音連線處理
│   │   │
│   │   ├── play.go            # /play 指令（含播放清單）
│   │   ├── play_test.go       # 播放指令測試
│   │   ├── pause.go           # /pause 指令
│   │   ├── pause_test.go      # 暫停指令測試
│   │   ├── skip.go            # /skip 指令
│   │   ├── skip_test.go       # 跳過指令測試
│   │   ├── stop.go            # /stop 指令
│   │   ├── stop_test.go       # 停止指令測試
│   │   ├── queue.go           # /queue 指令
│   │   ├── queue_test.go      # 佇列指令測試
│   │   ├── nowplaying.go      # /nowplaying 指令
│   │   ├── nowplaying_test.go # 當前播放測試
│   │   ├── download.go        # /download 指令
│   │   ├── loop.go            # /loop 指令
│   │   ├── loop_test.go       # 循環播放測試
│   │   ├── shuffle.go         # /shuffle 指令
│   │   ├── control_panel.go   # 音樂控制面板
│   │   ├── test_helpers.go    # 測試輔助工具
│   │   └── help.go            # /help 指令
│   │
│   ├── config/                 # 配置管理
│   │   ├── config.go          # 配置載入
│   │   └── config_test.go     # 配置測試
│   │
│   ├── player/                 # 播放器邏輯
│   │   ├── player.go          # GuildPlayer 實現
│   │   ├── player_test.go     # 播放器測試
│   │   ├── queue.go           # 佇列管理
│   │   ├── queue_test.go      # 佇列測試
│   │   ├── loop.go            # 循環模式定義
│   │   ├── loop_test.go       # 循環模式測試
│   │   ├── song.go            # 歌曲結構
│   │   ├── song_test.go       # 歌曲測試
│   │   ├── manager.go         # 播放器管理器
│   │   └── manager_test.go    # 管理器測試
│   │
│   ├── youtube/                # YouTube 整合
│   │   ├── resolver.go        # yt-dlp 解析器
│   │   └── resolver_test.go   # 解析器測試
│   │
│   └── audio/                  # 音訊處理（已棄用）
│       ├── pipeline.go        # 舊的音訊管線
│       └── pipeline_test.go   # 音訊管線測試
│
├── docker-compose.yml          # Docker 服務編排
├── Dockerfile                  # Bot 容器映像
├── .env                        # 環境變數（不提交）
├── .env.example                # 環境變數範例
├── go.mod                      # Go 模組定義
├── go.sum                      # Go 依賴鎖定
│
├── README.md                   # 專案說明
├── CLAUDE.md                   # Claude 開發指南
└── STRUCTURE.md                # 本文件
```

## 架構層級

```
┌─────────────────────────────────────┐
│         Discord Client              │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│        cmd/bot/main.go              │
│    (程式入口點、啟動邏輯)              │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│      internal/bot/bot.go            │
│  (Bot 初始化、事件處理、生命週期)      │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│      internal/command/*             │
│   (Slash Commands 實現)             │
└─────────────────────────────────────┘
        ↓                 ↓
┌──────────────┐  ┌──────────────────┐
│  Lavalink    │  │  internal/player │
│  (音訊服務)   │  │  (播放器邏輯)     │
└──────────────┘  └──────────────────┘
        ↓                 ↓
┌──────────────────────────────────────┐
│      internal/youtube/               │
│      (yt-dlp 整合)                    │
└──────────────────────────────────────┘
```

## 核心模組說明

### 1. cmd/bot - 程式入口

**職責**:
- 載入配置
- 初始化 Bot
- 啟動服務
- 處理優雅關閉

**關鍵檔案**:
- `main.go` - 程式主入口

### 2. internal/bot - Bot 核心

**職責**:
- Bot 生命週期管理
- 與 Discord Gateway 連線
- Lavalink 整合
- 事件監聽與分發

**關鍵檔案**:
- `bot.go` - Bot 結構和初始化
- `lavalink_handlers.go` - Lavalink 事件處理
- `cleanup.go` - 指令清理工具

**主要介面**:
```go
type Bot struct {
    Client       bot.Client
    Lavalink     disgolink.Client
    cfg          *config.Config
    playerManager *player.Manager
}
```

### 3. internal/command - Slash Commands

**職責**:
- 定義所有 Slash Commands
- 處理用戶交互
- 整合播放器和音訊服務

**指令清單**:
- `play.go` - 播放音樂和播放清單
- `pause.go` - 暫停/繼續
- `skip.go` - 跳過
- `stop.go` - 停止
- `queue.go` - 顯示佇列
- `nowplaying.go` - 當前播放
- `download.go` - 下載音訊
- `help.go` - 使用說明

**核心介面**:
```go
type BotCommand struct {
    Command discord.SlashCommandCreate
    Handler InteractionHandler
}

type InteractionHandler func(*events.ApplicationCommandInteractionCreate)
```

### 4. internal/player - 播放器邏輯

**職責**:
- 管理播放佇列
- 追蹤當前播放
- 提供播放器控制介面

**關鍵檔案**:
- `player.go` - GuildPlayer 實現
- `queue.go` - 佇列管理
- `song.go` - 歌曲結構
- `manager.go` - 播放器管理器

**核心結構**:
```go
type GuildPlayer struct {
    guildID     string
    queue       *Queue
    currentSong *Song
    paused      bool
    stopped     bool
}

type Song struct {
    Title       string
    URL         string
    StreamURL   string
    RequestedBy string
}
```

### 5. internal/youtube - YouTube 整合

**職責**:
- 解析 YouTube URL
- 搜尋關鍵字
- 提取播放清單
- 獲取音訊串流 URL

**關鍵檔案**:
- `resolver.go` - yt-dlp 整合

**核心介面**:
```go
type Resolver interface {
    Resolve(ctx context.Context, query string) (player.Song, error)
}
```

### 6. internal/config - 配置管理

**職責**:
- 載入環境變數
- 驗證配置
- 提供配置訪問

**關鍵檔案**:
- `config.go` - 配置結構和載入

**配置結構**:
```go
type Config struct {
    BotToken string
    GuildID  string
}
```

## 資料流

### 播放流程

```
用戶 → /play → PlayCommand → YouTube Resolver
                                    ↓
                             提取音訊 URL
                                    ↓
                             JoinVoiceAndPlayWithYtDlp
                                    ↓
                             Lavalink 播放
                                    ↓
                             更新 GuildPlayer
```

### 自動播放下一首

```
Lavalink TrackEndEvent → onTrackEnd
                              ↓
                         檢查結束原因
                              ↓
                    GuildPlayer.Dequeue()
                              ↓
                     SetCurrentSong()
                              ↓
                 JoinVoiceAndPlayWithYtDlp
```

### 下載流程

```
用戶 → /download → DownloadCommand
                          ↓
                   構建 yt-dlp 參數
                          ↓
                   執行 yt-dlp 下載
                          ↓
                   檢查檔案大小
                          ↓
                   上傳到 Discord
                          ↓
                   清理暫存檔案
```

## Docker 服務

### 服務定義

**bot** - 主 Bot 服務
- 運行主程式
- 連接 Discord
- 依賴 Lavalink

**lavalink** - 音訊處理服務
- 處理音訊串流
- 提供播放 API
- 獨立運行

**workspace** - 開發環境
- Go 開發工具鏈
- 用於編譯和測試
- 掛載專案目錄

## 測試策略

### 測試層級

1. **單元測試** (`*_test.go`)
   - 測試個別函數
   - 使用 mock 和 fake
   - 快速執行

2. **整合測試** (`*_integration_test.go`)
   - 測試模組間互動
   - 使用真實依賴
   - 較慢執行

3. **端到端測試**
   - 測試完整流程
   - 需要 Discord 環境
   - 手動執行

### 測試覆蓋

**當前覆蓋率** (2026-06-20):
- **總體**: ~45%
- **cmd/bot**: 92.9% ✅
- **internal/config**: 83.3% ✅
- **internal/player**: 78.2% ✅
- **internal/youtube**: 75.0% ✅
- **internal/bot**: 6.9% (需要改進)
- **internal/command**: 5.4% (框架測試已完成)

**測試檔案**:
- ✅ `cmd/bot/main_test.go` - 主程式測試
- ✅ `internal/bot/bot_test.go` - Bot 核心測試
- ✅ `internal/bot/lavalink_handlers_test.go` - Lavalink 事件測試
- ✅ `internal/bot/cleanup_test.go` - 清理邏輯測試
- ✅ `internal/command/*_test.go` - 15+ 指令測試文件
- ✅ `internal/command/test_helpers.go` - 測試輔助工具
- ✅ `internal/config/config_test.go` - 配置測試
- ✅ `internal/player/*_test.go` - 播放器測試
- ✅ `internal/youtube/*_test.go` - YouTube 解析測試
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 開發工作流

### 新增功能

1. 在 `internal/` 對應目錄創建檔案
2. 編寫單元測試
3. 實現功能
4. 更新文檔
5. 提交代碼

### 新增指令

1. 在 `internal/command/` 創建 `newcmd.go`
2. 定義 `NewCommand` 變數
3. 在 `registry.go` 註冊
4. 編寫測試 `newcmd_test.go`
5. 測試並提交

### 修改現有功能

1. 檢查相關測試
2. 更新測試以反映新行為
3. 修改實現
4. 確保測試通過
5. 更新文檔

## 依賴關係

### 主要依賴

- `github.com/disgoorg/disgo` - Discord API 客戶端
- `github.com/disgoorg/disgolink/v3` - Lavalink 客戶端
- `github.com/disgoorg/snowflake/v2` - Snowflake ID 處理
- `github.com/joho/godotenv` - 環境變數載入

### 開發依賴

- `golang.org/x/tools` - Go 工具鏈
- Docker 和 Docker Compose

## 配置檔案

### `.env`

```env
BOT_TOKEN=your_bot_token
GUILD_ID=your_test_guild_id
```

### `docker-compose.yml`

定義三個服務：
- `bot` - 生產環境
- `lavalink` - 音訊處理
- `workspace` - 開發環境

## 慣例與風格

### 命名慣例

- **檔案**: 小寫蛇形 `my_file.go`
- **類型**: 大寫駝峰 `MyType`
- **函數**: 駝峰 `myFunction`
- **常數**: 大寫蛇形 `MY_CONST`

### 組織原則

1. **按功能組織** - 不是按類型
2. **內部包** - 使用 `internal/` 防止外部導入
3. **介面隔離** - 小而專注的介面
4. **依賴注入** - 透過參數傳遞依賴

### 錯誤處理

1. **總是檢查錯誤**
2. **向上傳播或記錄**
3. **提供上下文**
4. **友善的用戶訊息**

## 已知限制

1. **私人播放清單** - 無法訪問（需登入）
2. **檔案大小** - 下載限制 25 MB
3. **影片時長** - 下載限制 10 分鐘
4. **地區限制** - 部分影片無法播放

## 未來改進

### 短期

1. 增加測試覆蓋率
2. 改善錯誤處理
3. 優化效能

### 中期

1. 資料庫整合
2. 用戶播放清單
3. 歷史記錄

### 長期

1. 多伺服器負載平衡
2. 監控和告警
3. 管理 Web 介面

## 維護檢查清單

### 每週

- [ ] 檢查日誌錯誤
- [ ] 監控資源使用
- [ ] 更新依賴

### 每月

- [ ] 審查程式碼品質
- [ ] 更新文檔
- [ ] 效能分析

### 每季

- [ ] 安全審計
- [ ] 架構評估
- [ ] 用戶回饋整合

## 相關文檔

- [README.md](README.md) - 專案說明和快速開始
- [CLAUDE.md](CLAUDE.md) - AI 開發指南
- [go.mod](go.mod) - Go 模組定義

## 最後更新

2026-06-20 - 完整重寫，反映當前架構
