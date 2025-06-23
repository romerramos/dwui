#!/bin/bash

# Exit on any error
set -e

# Check for root user
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please use sudo." >&2
  exit 1
fi

# Default values
PORT="8300"
VERSION="v0.0.1"
PASSWORD=""

# Parse named arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        --password)
            PASSWORD="$2"
            shift 2
            ;;
        --port)
            PORT="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        *)
            echo "Unknown parameter: $1" >&2
            exit 1
            ;;
    esac
done

# Detect architecture
ARCH=$(uname -m)
GO_ARCH=""
case "$ARCH" in
    x86_64)
        GO_ARCH="amd64"
        ;;
    aarch64 | arm64)
        GO_ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

# Check for password argument
if [ -z "$PASSWORD" ]; then
  echo "Usage: $0 --password <password> [--port <port>] [--version <version>]" >&2
  exit 1
fi

DWUI_BINARY="dwui-linux-${GO_ARCH}"
DWUI_URL="https://github.com/romerramos/dwui/releases/download/${VERSION}/${DWUI_BINARY}"
INSTALL_PATH="/usr/local/bin/dwui"
SERVICE_FILE="/etc/systemd/system/dwui.service"

echo "Starting DWUI installation..."

# Download the binary
echo "Downloading DWUI binary version ${VERSION}..."
curl -L -o "$DWUI_BINARY" "$DWUI_URL"

# Make it executable
chmod +x "$DWUI_BINARY"

# Move it to a directory in your PATH
echo "Installing DWUI to $INSTALL_PATH..."
mv "$DWUI_BINARY" "$INSTALL_PATH"

# Create a service file
echo "Creating systemd service file..."
cat > "$SERVICE_FILE" <<EOL
[Unit]
Description=DWUI - Docker Web UI
After=docker.service
Requires=docker.service

[Service]
ExecStart=$INSTALL_PATH --password ${PASSWORD} --port ${PORT}
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOL

# Enable and start the service
echo "Enabling and starting DWUI service..."
systemctl daemon-reload
systemctl enable dwui
systemctl start dwui

echo "DWUI has been installed and started successfully."
echo "You can check the status with: systemctl status dwui" 