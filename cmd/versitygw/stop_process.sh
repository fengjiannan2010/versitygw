#!/bin/bash

BasePath=$(cd $(dirname $0); pwd)
cd $BasePath


# ============================================
# 通过进程名称停止 netzongw 进程
# ============================================

# 可执行文件名（对应启动脚本中的文件名）
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

# 获取进程名
PROCESS_NAME=$(basename "$BIN_NAME")

# 查找进程 PID
PIDS=$(pgrep -f "$PROCESS_NAME")

if [ -z "$PIDS" ]; then
  echo "未找到名为 '$PROCESS_NAME' 的进程。"
  exit 1
fi

# 输出找到的进程 PID
echo "找到以下进程："
echo "$PIDS"

# 通过 PID 杀死进程
echo "正在停止进程 '$PROCESS_NAME' ..."
kill $PIDS

# 确认是否成功
if [ $? -eq 0 ]; then
  echo "进程 '$PROCESS_NAME' 已成功停止."
else
  echo "停止进程 '$PROCESS_NAME' 失败."
fi
