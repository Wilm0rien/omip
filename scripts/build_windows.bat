@echo off
setlocal
cd /d %~dp0
set CURDIR=%CD%
cd..

go build -ldflags "-s -w -X main.CmdLineOpt=default_cmd"
move omip.exe omip_cmd.exe
go build -ldflags "-H=windowsgui -s -w"

IF '%ERRORLEVEL%'=='0' GOTO OK
GOTO:EOF
:OK
fyne package -appBuild 108 -os windows -icon logo.png -appID omip.exe -appVersion 1.0.8  -executable omip.exe -name "omip" -release -tags 1.0.8

"C:\Program Files (x86)\Windows Kits\10\App Certification Kit\signtool.exe" sign /n "Open Source Developer, Christian Wilmes" /t http://time.certum.pl/ /fd sha256 "C:\upload\git\omip\omip.exe"
"C:\Program Files (x86)\Windows Kits\10\App Certification Kit\signtool.exe" sign /n "Open Source Developer, Christian Wilmes" /t http://time.certum.pl/ /fd sha256 "C:\upload\git\omip\omip_cmd.exe"
::go test .\... -coverprofile=c.out
::go tool cover -html=c.out -o coverage.html
cd %CURDIR%