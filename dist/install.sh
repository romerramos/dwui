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
# Try to read version from VERSION file in same directory, fallback to hardcoded
if [ -f "$(dirname "$0")/VERSION" ]; then
    VERSION=$(cat "$(dirname "$0")/VERSION" | tr -d '[:space:]')
else
    VERSION="v0.0.1"  # Fallback for standalone usage
fi
PASSWORD=""
PASSWORD_FILE="/var/lib/dwui/password"

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

DWUI_BINARY="dwui-linux-${GO_ARCH}"
DWUI_URL="https://github.com/romerramos/dwui/releases/download/${VERSION}/${DWUI_BINARY}"
INSTALL_PATH="/usr/local/bin/dwui"
SERVICE_FILE="/etc/systemd/system/dwui.service"

echo "Starting DWUI installation..."

# Download the binary
echo "Downloading DWUI binary version ${VERSION} from ${DWUI_URL}..."
curl -L -o "$DWUI_BINARY" "$DWUI_URL"

# Make it executable
chmod +x "$DWUI_BINARY"

# Move it to a directory in your PATH
echo "Installing DWUI to $INSTALL_PATH..."
mv "$DWUI_BINARY" "$INSTALL_PATH"

# Create a service file
echo "Creating systemd service file..."

EXEC_START_COMMAND="$INSTALL_PATH --port $PORT"
if [ -n "$PASSWORD" ]; then
  EXEC_START_COMMAND="$EXEC_START_COMMAND --password $PASSWORD"
  echo "Password provided. Configuring service with a fixed password."
else
  echo "No password provided. A random password will be generated on each start."
  mkdir -p "$(dirname "$PASSWORD_FILE")"
  touch "$PASSWORD_FILE"
  chown root:root "$(dirname "$PASSWORD_FILE")"
  chown root:root "$PASSWORD_FILE"
  chmod 600 "$PASSWORD_FILE"
  EXEC_START_COMMAND="$EXEC_START_COMMAND --password-file $PASSWORD_FILE"
  echo "You can find the password by running: sudo cat $PASSWORD_FILE"
fi

cat > "$SERVICE_FILE" <<EOL
[Unit]
Description=DWUI - Docker Web UI
After=docker.service
Requires=docker.service

[Service]
ExecStart=$EXEC_START_COMMAND
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