FROM golang:1.22

WORKDIR /workspace

# 安裝系統依賴：ffmpeg, python3, curl, libopus
RUN apt-get update && apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    curl \
    libopus0 \
    libopus-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# 安裝 yt-dlp（YouTube 下載工具）
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

# 複製 go.mod/go.sum 並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 安裝開發時可以使用的工具（可選）
RUN go install github.com/githubnemo/CompileDaemon@latest

# 預設工作目錄為 /workspace
COPY . .

# 預設啟動一個空閒容器，方便進入 shell 或執行指令
CMD ["tail", "-f", "/dev/null"]

