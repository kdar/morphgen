@echo off
rem xcopy ui dist\ui /I /Y
rem xcopy images dist\images /I /Y
rem cd dist
rem go run ..\morph.go ..\battlenet.go ..\wowhead.go ..\generic.go ..\version.go ..\ui.go
rem cd ..

xcopy ui dist\ui /I /Y
mkdir dist\resources
copy /Y resources\icon.png dist\resources
copy /Y resources\icon.ico dist\resources

For /f "tokens=2-4 delims=/ " %%a in ('date /t') do (set mydate=%%c-%%a-%%b)
For /f "tokens=1-3 delims=/:" %%a in ("%TIME%") do (set mytime=%%a%%b%%c)

go build -o "dist\morphgen_%mydate%_%mytime%.exe" && "dist\morphgen_%mydate%_%mytime%.exe"
del "dist\morphgen_%mydate%_%mytime%.exe"