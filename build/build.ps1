# Create dist directory if it doesn't exist
New-Item -ItemType Directory -Force -Path dist

# Build the Docker image
Write-Host "Building Docker image..."
docker build -t godo-builder -f Dockerfile.build .

# Build Linux version
Write-Host "`nBuilding Linux version..."
docker run --rm -v ${PWD}/dist:/go/src/app/dist godo-builder go build -tags linux -ldflags "-s -w" -o dist/godo ./cmd/godo

# Build Windows version
Write-Host "`nBuilding Windows version..."
docker run --rm -v ${PWD}/dist:/go/src/app/dist -e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=x86_64-w64-mingw32-gcc godo-builder go build -tags windows -ldflags "-s -w" -o dist/godo.exe ./cmd/godo

Write-Host "`nBuild complete! Check the dist directory for the binaries." 