# Godo Build Script
# This script builds both Windows and Linux versions of Godo using Docker
# Requirements:
#   - Docker
#   - PowerShell
#   - Git (for development)

# Enable strict mode
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# Enable BuildKit for better caching
$env:DOCKER_BUILDKIT = "1"

# Configuration
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$ROOT_DIR = Split-Path -Parent $SCRIPT_DIR
$BUILD_DIR = Join-Path -Path $ROOT_DIR -ChildPath "dist"
$DOCKER_DIR = Join-Path -Path $SCRIPT_DIR -ChildPath "docker"
$DOCKERFILE = Join-Path -Path $DOCKER_DIR -ChildPath "Dockerfile"

# Get version information
$VERSION = git describe --tags --always 2>$null
if (-not $VERSION) { $VERSION = "dev" }
$COMMIT = git rev-parse --short HEAD 2>$null
if (-not $COMMIT) { $COMMIT = "unknown" }
$BUILD_TIME = [DateTime]::UtcNow.ToString("o")

# Build args for versioning
$BUILD_ARGS = @(
    "--build-arg", "VERSION=$VERSION",
    "--build-arg", "COMMIT=$COMMIT",
    "--build-arg", "BUILD_TIME=$BUILD_TIME"
)

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

# Build and test in Docker
Write-Host "Running tests and building..." -ForegroundColor Green
try {
    # Run linting
    Write-Host "`nRunning linters..." -ForegroundColor Yellow
    docker build $BUILD_ARGS --target lint -t godo-lint -f $DOCKERFILE $ROOT_DIR
    docker run --rm godo-lint
    if ($LASTEXITCODE -ne 0) { throw "Linting failed" }

    # Run tests
    Write-Host "`nRunning tests..." -ForegroundColor Yellow
    docker build $BUILD_ARGS --target test -t godo-test -f $DOCKERFILE $ROOT_DIR
    docker run --rm godo-test
    if ($LASTEXITCODE -ne 0) { throw "Tests failed" }

    # Build both versions in parallel
    Write-Host "`nBuilding Linux and Windows versions..." -ForegroundColor Yellow
    $linuxJob = Start-Job -ScriptBlock {
        docker build $using:BUILD_ARGS --target linux-runtime -t godo-linux -f $using:DOCKERFILE $using:ROOT_DIR
        if ($LASTEXITCODE -ne 0) { throw "Linux build failed" }
        docker create --name godo-linux-temp godo-linux
        docker cp godo-linux-temp:/app/godo $using:BUILD_DIR/godo
        docker rm godo-linux-temp
    }

    $windowsJob = Start-Job -ScriptBlock {
        docker build $using:BUILD_ARGS --target windows-runtime -t godo-windows -f $using:DOCKERFILE $using:ROOT_DIR
        if ($LASTEXITCODE -ne 0) { throw "Windows build failed" }
        docker create --name godo-windows-temp godo-windows
        docker cp godo-windows-temp:/godo.exe $using:BUILD_DIR/godo.exe
        docker rm godo-windows-temp
    }

    # Wait for builds to complete
    $null = Wait-Job $linuxJob, $windowsJob
    Receive-Job $linuxJob, $windowsJob

    # Check for errors
    if ($linuxJob.State -eq "Failed" -or $windowsJob.State -eq "Failed") {
        throw "One or more builds failed"
    }

} catch {
    Handle-Error "Build process failed: $_"
} finally {
    Remove-Job $linuxJob, $windowsJob -Force -ErrorAction SilentlyContinue
}

Write-Host "`nBuild complete! Binaries are in the dist directory:" -ForegroundColor Green
Write-Host "  - Windows: dist/godo.exe (version: $VERSION, commit: $COMMIT)"
Write-Host "  - Linux:   dist/godo (version: $VERSION, commit: $COMMIT)" 