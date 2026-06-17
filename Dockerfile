FROM golang:1.22

WORKDIR /workspace

# 複製 go.mod/go.sum 並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 安裝開發時可以使用的工具（可選）
RUN go install github.com/githubnemo/CompileDaemon@latest

# 預設工作目錄為 /workspace
COPY . .

# 預設啟動一個空閒容器，方便進入 shell 或執行指令
CMD ["tail", "-f", "/dev/null"]
