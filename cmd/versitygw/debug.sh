#!/bin/bash

# ============================================
# 启动 netzongw 的 Delve 远程调试服务（含清理逻辑）
# ============================================

BIN_NAME="netzongw"

# 构建带调试信息的可执行文件
echo "🔧 正在编译调试版本..."
go build -gcflags="all=-N -l" -o "$BIN_NAME"
if [ $? -ne 0 ]; then
  echo "❌ 编译失败，退出"
  exit 1
fi

# 清理函数：在脚本退出或 Ctrl+C 时执行
cleanup() {
  echo "🧹 正在清理调试进程..."
  if [ -n "$DLV_PID" ] && kill -0 "$DLV_PID" 2>/dev/null; then
    kill "$DLV_PID"
    echo "✅ 已终止 dlv 进程 (PID: $DLV_PID)"
  fi
  exit 0
}

# 捕捉 Ctrl+C (SIGINT) 和 脚本退出(SIGTERM/EXIT)
trap cleanup INT TERM EXIT

# 启动 Delve 调试器（后台）
echo "🚀 启动 Delve 远程调试服务..."

export ROOT_ACCESS_KEY="admin"
export ROOT_SECRET_KEY="admin"
export VGW_PORT=":11000"
export NOTIFY_BASE_URL="http://192.168.11.17:8080"
export NOTIFY_ENDPOINT_PATH="/oss/rest/restServer/creatArcTask"
export VGW_META_NONE="true"
export SKIP_DIRS="VisualDiscSpace,IsoBuffer"

dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./$BIN_NAME posix /data/vol &


# 记录 PID 并等待
DLV_PID=$!
echo "🆔 Delve PID: $DLV_PID"

# 等待子进程（Delve）退出
wait $DLV_PID
