echo off 

go build -ldflags "-w -s" omip_updater.go
IF '%ERRORLEVEL%'=='0' GOTO OK
GOTO:EOF
:OK
"C:\Program Files (x86)\Microsoft SDKs\ClickOnce\SignTool\signtool.exe" sign /n "Open Source Developer, Christian Wilmes" /t http://time.certum.pl/ /fd sha256 "E:\upload\go_test_folder\go_scripts\omip\exec\omip_updater.exe"