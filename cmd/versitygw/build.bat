@echo off
REM === 设置程序主入口文件名 ===
set MAIN_FILE=main.go

REM === 输出目录 ===
set OUTPUT_DIR=build

REM === 创建输出目录 ===
if not exist %OUTPUT_DIR% (
    mkdir %OUTPUT_DIR%
)

set CGO_ENABLED=1
set LDFLAGS=-ldflags "-s -w -X main.Version=1.0.14 -X main.Build=netzon -X main.BuildTime=20250621"

REM === 构建 linux/arm64 (64位 ARM，如树莓派 4/5、ARM 云主机等) ===
echo Building for linux/arm64...
set GOARCH=arm64

go build %LDFLAGS% -o %OUTPUT_DIR%\netzongw-linux-arm64.exe %MAIN_FILE%
echo Build complete. Output in %OUTPUT_DIR%
pause

