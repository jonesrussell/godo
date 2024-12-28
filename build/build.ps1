# Godo Build Script
# This script builds both Windows and Linux versions of Godo using Docker
# Requirements:
#   - Docker
#   - PowerShell
#   - Git (for development)

# Configuration
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$ROOT_DIR = Split-Path -Parent $SCRIPT_DIR
$BUILD_DIR = Join-Path -Path $ROOT_DIR -ChildPath "dist"
$DOCKER_DIR = Join-Path -Path $SCRIPT_DIR -ChildPath "docker"

# Create build directory if it doesn't exist
if (-not (Test-Path -Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null
}

# Build the Docker image
Write-Host "Building Docker image..."
docker build -t godo-builder -f (Join-Path -Path $DOCKER_DIR -ChildPath "Dockerfile") $ROOT_DIR

# Build Linux version
Write-Host "`nBuilding Linux version..."
docker run --rm -v "${BUILD_DIR}:/go/src/app/dist" godo-builder go build -tags linux -ldflags "-s -w" -o dist/godo ./cmd/godo

# Build Windows version
Write-Host "`nBuilding Windows version..."
docker run --rm -v "${BUILD_DIR}:/go/src/app/dist" -e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=x86_64-w64-mingw32-gcc godo-builder go build -tags windows -ldflags "-s -w" -o dist/godo.exe ./cmd/godo

Write-Host "`nBuild complete! Binaries are in the dist directory:"
Write-Host "  - Windows: dist/godo.exe"
Write-Host "  - Linux:   dist/godo" 