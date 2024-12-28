# Godo Build Script
# This script builds both Windows and Linux versions of Godo using Docker
# Requirements:
#   - Docker
#   - PowerShell
#   - Git (for development)

# Enable strict mode
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# Configuration
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$ROOT_DIR = Split-Path -Parent $SCRIPT_DIR
$BUILD_DIR = Join-Path -Path $ROOT_DIR -ChildPath "dist"
$DOCKER_DIR = Join-Path -Path $SCRIPT_DIR -ChildPath "docker"
$DOCKERFILE = Join-Path -Path $DOCKER_DIR -ChildPath "Dockerfile"

# Function to handle errors
function Handle-Error {
    param($ErrorMessage)
    Write-Host "Error: $ErrorMessage" -ForegroundColor Red
    exit 1
}

# Verify Docker is installed
if (-not (Get-Command "docker" -ErrorAction SilentlyContinue)) {
    Handle-Error "Docker is not installed or not in PATH"
}

# Create build directory if it doesn't exist
try {
    if (-not (Test-Path -Path $BUILD_DIR)) {
        New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null
    }
} catch {
    Handle-Error "Failed to create build directory: $_"
}

# Build the Docker image
Write-Host "Building Docker image..." -ForegroundColor Green
try {
    docker build -t godo-builder -f $DOCKERFILE $ROOT_DIR
    if ($LASTEXITCODE -ne 0) { throw "Docker build failed" }
} catch {
    Handle-Error "Failed to build Docker image: $_"
}

# Build Linux version
Write-Host "`nBuilding Linux version..." -ForegroundColor Green
try {
    docker run --rm -v "${BUILD_DIR}:/go/src/app/dist" godo-builder go build -tags linux -ldflags "-s -w" -o dist/godo ./cmd/godo
    if ($LASTEXITCODE -ne 0) { throw "Linux build failed" }
} catch {
    Handle-Error "Failed to build Linux version: $_"
}

# Build Windows version
Write-Host "`nBuilding Windows version..." -ForegroundColor Green
try {
    docker run --rm -v "${BUILD_DIR}:/go/src/app/dist" -e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=x86_64-w64-mingw32-gcc godo-builder go build -tags windows -ldflags "-s -w" -o dist/godo.exe ./cmd/godo
    if ($LASTEXITCODE -ne 0) { throw "Windows build failed" }
} catch {
    Handle-Error "Failed to build Windows version: $_"
}

Write-Host "`nBuild complete! Binaries are in the dist directory:" -ForegroundColor Green
Write-Host "  - Windows: dist/godo.exe"
Write-Host "  - Linux:   dist/godo" 