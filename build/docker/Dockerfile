# Build arguments for versioning
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Base stage with common dependencies
FROM golang:1.24-bullseye AS base
ARG VERSION
ARG COMMIT
ARG BUILD_TIME

# Build flags for versioning
ENV VERSION=${VERSION}
ENV COMMIT=${COMMIT}
ENV BUILD_TIME=${BUILD_TIME}
ENV BUILD_FLAGS="-ldflags '-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME} -s -w'"

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && apt-get install -y \
    gcc \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxxf86vm-dev \
    mingw-w64 \
    && rm -rf /var/lib/apt/lists/*

# Dependencies stage
FROM base AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    go mod download

# Tools stage for installing development tools
FROM deps AS tools
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    go install github.com/google/wire/cmd/wire@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Wire generation stage
FROM tools AS wire
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    cd internal/container && wire gen -tags linux

# Lint stage
FROM wire AS lint
ENV CGO_ENABLED=1
CMD ["golangci-lint", "run", "./..."]

# Test stage
FROM wire AS test
ENV CGO_ENABLED=1
CMD ["go", "test", "-v", "-tags=docker", "./..."]

# Linux build stage
FROM wire AS linux-builder
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -tags linux -o bin/godo-linux cmd/godo/main.go

# Windows build stage
FROM wire AS windows-builder
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build ${BUILD_FLAGS} -tags windows -o bin/godo.exe cmd/godo/main.go

# Runtime stage for Linux
FROM debian:bullseye-slim AS linux-runtime
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && apt-get install -y gcc && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=linux-builder /app/bin/godo-linux /app/godo
ENTRYPOINT ["/app/godo"]

# Runtime stage for Windows (artifact only)
FROM scratch AS windows-runtime
COPY --from=windows-builder /app/bin/godo.exe /godo.exe 