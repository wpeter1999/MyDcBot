# Discord Bot 專案結構與開發流程

## 目錄結構

```
discord-bot/
├── cmd/
│   └── bot/
│       └── main.go              # 應用程式入口點
├── internal/
│   ├── bot/
│   │   └── bot.go               # Bot 生命週期與事件處理
│   ├── command/
│   │   ├── command.go           # BotCommand 結構與 Registrar 介面
│   │   ├── registry.go          # 指令註冊表
│   │   ├── registrar.go         # RegisterCommands / HandleInteraction
│   │   └── ping.go              # /ping 指令實作
│   ├── config/
│   │   └── config.go            # 配置載入
│   └── pkg/
│       └── logger/              # 共用日誌工具（預留）
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── STRUCTURE.md                 # 本文件
└── README.md
```

## 設計原則

- 使用 `cmd/` + `internal/` 標準 Go 專案佈局
- 所有商業邏輯放在 `internal/`，避免被外部引用
- 每個指令獨立管理，新增指令只需修改 `registry.go`
- 配置與 Bot 生命週期分離

## 標準開發流程

### 1. 本地開發

```bash
# 啟動 Bot
go run ./cmd/bot

# 執行測試
go test ./...

# 建置
go build -o bin/bot ./cmd/bot
```

### 2. Docker 開發

```bash
# 建置映像
docker compose build

# 啟動開發容器
docker compose up -d workspace

# 進入容器
docker compose exec workspace bash

# 在容器內執行
go run ./cmd/bot
```

### 3. 新增指令流程

1. 先撰寫或更新對應的單元測試
2. 再開始實作新指令的程式碼
3. 在 `internal/command/` 建立新檔案（如 `weather.go`）
4. 定義 `BotCommand` 變數與指令處理函式
5. 在 `registry.go` 的 `CommandRegistry` 中加入新指令
6. 執行 `go test ./...` 確認測試通過
7. 再執行 `go run ./cmd/bot` 測試實際行為

### 4. 測試策略

- **Unit Test**：測試單一函式邏輯（使用 fake）
- **Integration Test**：測試模組之間整合（可建立 `*_integration_test.go`）
- **E2E Test**：使用真實 Discord 伺服器測試完整流程（選用）

## 注意事項

- 所有 import 路徑已更新為 `discordbot/internal/...`
- 建議之後加入 `Makefile` 統一開發指令
