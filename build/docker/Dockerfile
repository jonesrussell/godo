FROM golang:1.23-bookworm AS builder

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    libayatana-appindicator3-dev \
    libgtk-3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -v -o /go/bin/godo cmd/godo/main.go

FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libayatana-appindicator3-1 \
    libgtk-3-0 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/godo /usr/local/bin/godo

ENTRYPOINT ["godo"] 