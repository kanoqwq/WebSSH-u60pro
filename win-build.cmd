@echo off
setlocal EnableDelayedExpansion
set "now_dir=%~dp0"

cd "%now_dir%webssh"
call pnpm install
call npm run buildOnly

:: 删除 webroot 目录下所有文件和子目录（保留 webroot 目录本身）
if exist "%now_dir%gossh\webroot" (
    :: 删除文件
	del /f /s /q "%now_dir%gossh\webroot\*"
	:: 删除子目录
	for /d %%i in ("%now_dir%gossh\webroot\*") do rd /s /q "%%i"
)

:: 复制 webssh\dist 下所有文件到 gossh\webroot
:: robocopy "%now_dir%webssh\dist" "%now_dir%gossh\webroot" /E
:: 复制 dist 下所有文件
xcopy "%now_dir%webssh\dist\*" "%now_dir%gossh\webroot\" /E /Y /I

cd "%now_dir%gossh"

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64

go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy

set build_time=%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set build_time=%build_time: =0%

go build -trimpath -ldflags="-s -w" -o webssh_%build_time%

::upx --best --lzma webssh
%now_dir%upx.exe --best --lzma webssh_%build_time%

pause