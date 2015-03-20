@echo off
windres -o resource.syso -l 0x409 resources\morphgen.rc
go build -ldflags "-s -H windowsgui" -o dist\morphgen.exe
rem go build -ldflags "-s" -o dist\morphgen.exe
rm resource.syso
xcopy ui dist\ui /I /Y
mkdir dist\resources
copy /Y resources\icon.png dist\resources
copy /Y resources\icon.ico dist\resources