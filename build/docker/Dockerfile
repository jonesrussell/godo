FROM golang:1.23

# Install dependencies for Linux build
RUN apt-get update && apt-get install -y \
    gcc \
    libgl1-mesa-dev \
    xorg-dev \
    mingw-w64 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app

# Copy only necessary files for dependency download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build command will be provided when running the container 