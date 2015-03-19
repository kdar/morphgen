@echo off
xcopy ui dist\ui /I /Y
mkdir dist\resources
copy /Y resources\icon.png dist\resources
copy /Y resources\icon.ico dist\resources