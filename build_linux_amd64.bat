@echo off
set GOOS=linux
set GOARCH=amd64
set OUTPUT=pkgup

go build -o %OUTPUT% github.com/realheyu/pkgup

if %ERRORLEVEL% neq 0 (
    echo Build failed.
    exit /b %ERRORLEVEL%
) else (
    echo Build succeeded. Output file: %OUTPUT%
)
