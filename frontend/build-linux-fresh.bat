@echo off
set GOOS=linux
set GOARCH=amd64
cd /d C:\Users\YKing\.openclaw\workspace\mdm-project\backend
go clean -cache
go clean -modcache
go build -o mdm-server-linux.exe .
echo Build done
