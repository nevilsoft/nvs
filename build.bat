@echo off
setlocal enabledelayedexpansion

REM NVS CLI Cross-Platform Build Script for Windows
REM Builds the NVS CLI tool for Windows, Linux, and macOS

echo 🚀 Building NVS CLI...

REM Set version and build info
set VERSION=%VERSION%
if "%VERSION%"=="" set VERSION=dev

for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT_HASH=%%i
if "%COMMIT_HASH%"=="" set COMMIT_HASH=unknown

for /f "tokens=*" %%i in ('date /t') do set BUILD_DATE=%%i
for /f "tokens=*" %%i in ('time /t') do set BUILD_TIME=%%i

echo 📅 Build Time: %BUILD_DATE% %BUILD_TIME%
echo 🔗 Commit: %COMMIT_HASH%
echo.

REM Create build directory
if not exist "build" mkdir build

REM Build flags
set LDFLAGS=-X "main.Version=%VERSION%" -X "main.BuildTime=%BUILD_DATE% %BUILD_TIME%" -X "main.CommitHash=%COMMIT_HASH%"

REM Function to build for a specific platform
:build_for_platform
set GOOS=%1
set GOARCH=%2
set EXTENSION=%3
set OUTPUT_NAME=nvs%EXTENSION%

echo 🔨 Building for %GOOS%/%GOARCH%...

REM Set environment variables
set GOOS=%GOOS%
set GOARCH=%GOARCH%
set CGO_ENABLED=0

REM Build the binary
go build -ldflags "%LDFLAGS%" -o "build\%OUTPUT_NAME%" .

if %ERRORLEVEL% EQU 0 (
    echo ✅ Successfully built: build\%OUTPUT_NAME%
    
    REM Get file size
    for %%A in ("build\%OUTPUT_NAME%") do set SIZE=%%~zA
    echo 📦 File size: %SIZE% bytes
) else (
    echo ❌ Failed to build for %GOOS%/%GOARCH%
    exit /b 1
)

echo.
goto :eof

REM Main build process
echo 📋 Building for all supported platforms...
echo.

REM Build for different platforms
call :build_for_platform linux amd64 ""
call :build_for_platform linux arm64 ""
call :build_for_platform windows amd64 ".exe"
call :build_for_platform windows arm64 ".exe"
call :build_for_platform darwin amd64 ""
call :build_for_platform darwin arm64 ""

REM Create checksums
echo 🔍 Creating checksums...
cd build
if exist "nvs*" (
    certutil -hashfile nvs* SHA256 > checksums.sha256 2>nul
    echo ✅ Checksums created: checksums.sha256
) else (
    echo ⚠️ No files to create checksums for
)
cd ..
echo.

REM Create archives
echo 📦 Creating release archives...
cd build

REM Create ZIP archives for Windows
if exist "nvs.exe" (
    powershell -Command "Compress-Archive -Path 'nvs.exe', 'checksums.sha256' -DestinationPath 'nvs-%VERSION%-windows-amd64.zip' -Force" 2>nul
    echo ✅ Created: nvs-%VERSION%-windows-amd64.zip
)

REM Note: For other platforms, you would need additional tools like tar/zip
REM or use PowerShell to create archives

cd ..
echo.

REM Summary
echo 🎉 Build completed successfully!
echo.
echo 📁 Build artifacts:
dir build
echo.
echo 📋 Supported platforms:
echo   • Linux (amd64, arm64)
echo   • Windows (amd64, arm64)
echo   • macOS (amd64, arm64)
echo.
echo 💡 Usage:
echo   • Linux/macOS: ./nvs
echo   • Windows: nvs.exe
echo.
echo 🔧 To install globally:
echo   • Copy the appropriate binary to a directory in your PATH
echo   • Or use: go install github.com/nevilsoft/nvs@latest

pause 