#!/bin/bash

BasePath=$(cd $(dirname $0); pwd)
cd $BasePath

# ============================================
# 运行 netzongw 正常模式（不含调试）
# ============================================

# 可执行文件名
BIN_NAME_AMD64="netzongw-linux-amd64"
BIN_NAME_ARM="netzongw-linux-arm64"

BIN_NAME=""

# 检查平台并选择对应的可执行文件
if [ -f "$BIN_NAME_AMD64" ] && [ -f "$BIN_NAME_ARM" ]; then
  # 检测系统架构
  ARCH=$(uname -m)
  if [ "$ARCH" == "x86_64" ]; then
    BIN_NAME="$BIN_NAME_AMD64"
  elif [ "$ARCH" == "aarch64" ]; then
    BIN_NAME="$BIN_NAME_ARM"
  else
    echo "不支持的系统架构: $ARCH"
    exit 1
  fi
elif [ -f "$BIN_NAME_AMD64" ]; then
  BIN_NAME="$BIN_NAME_AMD64"
elif [ -f "$BIN_NAME_ARM" ]; then
  BIN_NAME="$BIN_NAME_ARM"
else
  echo "可执行文件 ($BIN_NAME_AMD64 或 $BIN_NAME_ARM) 不存在，请先编译：go build -o $BIN_NAME"
  exit 1
fi

export ROOT_ACCESS_KEY="admin"
export ROOT_SECRET_KEY="admin"
export VGW_PORT=":11000"
export NOTIFY_BASE_URL="http://192.168.11.17:8080"
export NOTIFY_ENDPOINT_PATH="/oss/rest/restServer/creatArcTask"
export VGW_META_NONE="true"
export SKIP_DIRS="VisualDiscSpace,IsoBuffer"

# 启动程序
echo "正在启动 $BIN_NAME ..."
nohup ./$BIN_NAME posix /data/vol >/dev/null 2>&1 &

# 记录 PID 并等待
NETZONGW_PID=$!
echo "$BIN_NAME PID: $NETZONGW_PID"
