#!/bin/bash
# 驗證 Docker 容器中的音樂 bot 執行環境

echo "=== 驗證 yt-dlp ==="
yt-dlp --version
if [ $? -eq 0 ]; then
    echo "✅ yt-dlp 已安裝"
else
    echo "❌ yt-dlp 未安裝"
    exit 1
fi

echo ""
echo "=== 驗證 ffmpeg ==="
ffmpeg -version | head -n 1
if [ $? -eq 0 ]; then
    echo "✅ ffmpeg 已安裝"
else
    echo "❌ ffmpeg 未安裝"
    exit 1
fi

echo ""
echo "=== 驗證 Python 3 ==="
python3 --version
if [ $? -eq 0 ]; then
    echo "✅ python3 已安裝"
else
    echo "❌ python3 未安裝"
    exit 1
fi

echo ""
echo "=== 驗證 Go 環境 ==="
go version
if [ $? -eq 0 ]; then
    echo "✅ Go 已安裝"
else
    echo "❌ Go 未安裝"
    exit 1
fi

echo ""
echo "✅ 所有依賴驗證通過"
