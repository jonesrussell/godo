FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    dbus-x11 \
    libdbus-1-dev \
    libx11-dev \
    xorg-dev \
    libxtst-dev \
    libxkbcommon-dev \
    libxkbcommon-x11-dev \
    libxcb1-dev \
    libxcb-xkb-dev \
    libxcb-keysyms1-dev \
    libxcb-xtest0-dev \
    libayatana-appindicator3-dev \
    libayatana-indicator3-dev \
    x11-apps \
    && rm -rf /var/lib/apt/lists/*

# Install task with proper permissions
RUN go install github.com/go-task/task/v3/cmd/task@latest && \
    mkdir -p /go/pkg/mod && \
    chown -R vscode:vscode /go

# Set up display for WSL2
ENV DISPLAY=:0
ENV WAYLAND_DISPLAY=wayland-0
ENV XDG_RUNTIME_DIR=/tmp

# Switch to non-root user
USER vscode