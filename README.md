# Discord Bot (Go)

使用 `discordgo` 開發的 Discord Bot 專案。

## 安裝與設定

1. 安裝 Go 1.22 或以上版本
2. 在專案根目錄下載依賴：

```bash
go mod tidy
```

3. 建立 `.env` 檔案：

```bash
copy .env.example .env
```

4. 在 `.env` 中填入你的 Bot Token：

```env
BOT_TOKEN=your_discord_bot_token_here
```

（可選）若想在私人伺服器註冊 Slash Command，可加入 Guild ID：

```env
GUILD_ID=your_test_guild_id
```

5. 本地啟動 Bot：

```bash
go run ./cmd/bot
```

## Docker 開發環境

本專案已建立 Go 開發環境，使用 `Dockerfile` 建置開發映像，並由 `docker-compose.yml` 啟動。

### 建置 Go 開發映像

```bash
docker compose build
```

### 啟動開發工作區容器

```bash
docker compose up -d workspace
```

### 進入開發容器

```bash
docker compose exec workspace bash
```

### 在容器內執行程式

```bash
go run ./cmd/bot
```

### 直接啟動 Bot

```bash
docker compose up bot
```

### 停止容器

```bash
docker compose down
```

> 請先確認 `.env` 已正確設定 `BOT_TOKEN`，否則容器啟動時會失敗。

## 使用方式

目前 Bot 支援以下指令：

- `/ping` → 回應 `Pong!`
