@echo off

:: 声明采用UTF-8编码
chcp 65001

:: 打开批处理所在目录
cd /d %~dp0

go build -ldflags="-H windowsgui" -o ../../build/tools

REM 出错暂停
if %errorlevel% neq 0 pause exit 