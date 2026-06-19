FROM golang:1.22

WORKDIR /workspace

# 安裝系統依賴
# - ffmpeg: 音訊處理與格式轉換
# - python3: yt-dlp 運行環境
# - curl: 下載工具
# - libopus0/libopus-dev: Opus 音訊編碼庫（語音通訊必需）
# - pkg-config: 編譯時查找庫依賴
# 同時安裝 yt-dlp（YouTube 音訊/影片下載工具）
RUN apt-get update && apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    curl \
    libopus0 \
    libopus-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/* \
    && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
        -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

# 複製 Go 模組定義並下載依賴
# 這樣可以利用 Docker 層快取，只在依賴變更時重新下載
COPY go.mod go.sum ./
RUN go mod download \
    && go list -m github.com/disgoorg/disgo \
    && go list -m layeh.com/gopus

# 複製專案檔案
COPY . .

# 預編譯 bot 並安裝開發工具
RUN go build -o bin/bot ./cmd/bot \
    && test -f bin/bot && echo "✅ Bot binary compiled successfully" \
    && go install github.com/githubnemo/CompileDaemon@latest

# 預設指令：使用 CompileDaemon 進行熱重載開發
# 生產環境建議直接運行預編譯的二進制：./bin/bot
CMD ["CompileDaemon", "-build=go build -o bin/bot ./cmd/bot", "-command=./bin/bot"]
