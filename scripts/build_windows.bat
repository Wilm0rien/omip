@echo off
setlocal
cd /d %~dp0
set CURDIR=%CD%
cd..
go build -ldflags "-H=windowsgui -s -w"
IF '%ERRORLEVEL%'=='0' GOTO OK
GOTO:EOF
:OK
fyne package -appBuild 15 -os windows -icon logo.png -appID omip.exe -appVersion 1.0.1  -executable omip.exe -name "omip v1.0.1" -release -tags 1.0.1

"C:\Program Files (x86)\Microsoft SDKs\ClickOnce\SignTool\signtool.exe" sign /n "Open Source Developer, Christian Wilmes" /t http://time.certum.pl/ /fd sha256 "E:\upload\go_test_folder\go_scripts\omip\omip.exe"
::go test .\... -coverprofile=c.out
::go tool cover -html=c.out -o coverage.html
cd %CURDIR%