@echo off
setlocal
cd /d %~dp0
set CURDIR=%CD%
cd..
go build -ldflags -H=windowsgui
fyne package -os windows -icon logo.png
::go test .\... -coverprofile=c.out
::go tool cover -html=c.out -o coverage.html
cd %CURDIR%