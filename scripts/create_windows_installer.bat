@echo off
setlocal
cd /d %~dp0
:: download wix binaries from https://github.com/wixtoolset/wix3/releases

:: install wix binaries into WIX_PATH

SET WIX_PATH=E:\upload\go_test_folder\go_scripts\wix311-binaries
set CURDIR=%CD%

call build_windows.bat

if not exist ..\omip.exe goto:ERROR
set path=%path%;%WIX_PATH%
candle -ext WixUIExtension -ext WixUtilExtension -arch x64 omip.wxs
light -ext WixUIExtension -ext WixUtilExtension omip.wixobj
pause
goto:EOF
:ERROR
echo ERROR could not find OMIP.exe
pause
:EOF