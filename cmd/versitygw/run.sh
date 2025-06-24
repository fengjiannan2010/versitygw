#!/bin/bash

# ============================================
# 运行 netzongw 正常模式（不含调试）
# ============================================

# 可执行文件名
BIN_NAME="netzongw"

# 检查是否存在可执行文件
if [ ! -f "$BIN_NAME" ]; then
  echo "❌ 可执行文件 '$BIN_NAME' 不存在，请先编译：go build -o $BIN_NAME"
  exit 1
fi

# 启动程序
echo "🚀 正在启动 $BIN_NAME ..."
./$BIN_NAME \
  --port :11000 \
  --access admin \
  --secret admin \
  posix --nometa /home/vgw

