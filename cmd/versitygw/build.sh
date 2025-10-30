#!/bin/bash

# === 输出目录 ===
OUTPUT_DIR="build"

# === 创建输出目录 ===
if [ ! -d "$OUTPUT_DIR" ]; then
    mkdir "$OUTPUT_DIR"
fi

# === 构建 linux/arm64 (64位 ARM，如树莓派 4/5、ARM 云主机等) ===

export GOARCH=arm64
export GOOS=linux
export CGO_ENABLED=1

echo "Building for linux/$GOARCH..."
go build -o "$OUTPUT_DIR/netzongw-$GOOS-$GOARCH"

export GOARCH=amd64
echo "Building for linux/$GOARCH..."
go build -o "$OUTPUT_DIR/netzongw-$GOOS-$GOARCH"

cp start_process.sh $OUTPUT_DIR
cp stop_process.sh $OUTPUT_DIR
echo "Build complete. Output in $OUTPUT_DIR"
tree $OUTPUT_DIR
tar -czvf "netzongw-$GOOS.tar.gz" $OUTPUT_DIR
