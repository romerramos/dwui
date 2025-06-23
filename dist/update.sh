#!/bin/bash

# Exit on any error
set -e

# Try to read version from local VERSION file, then from GitHub, then fallback
if [ -f "$(dirname "$0")/VERSION" ]; then
    VERSION=$(cat "$(dirname "$0")/VERSION" | tr -d '[:space:]')
elif command -v curl >/dev/null 2>&1; then
    VERSION=$(curl -sSL https://raw.githubusercontent.com/romerramos/dwui/main/dist/VERSION 2>/dev/null | tr -d '[:space:]' || echo "v0.0.1")
else
    VERSION="v0.0.1"  # Fallback when curl is not available
fi

# Check for root user
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please use sudo." >&2
  exit 1
fi

# Default values
FORCE_UPDATE=false

# Parse named arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --force)
            FORCE_UPDATE=true
            shift
            ;;
        *)
            echo "Unknown parameter: $1" >&2
            echo "Usage: $0 [--version <version>] [--force]" >&2
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

# Check if DWUI is installed
if [ ! -f "$INSTALL_PATH" ]; then
    echo "DWUI is not installed at $INSTALL_PATH" >&2
    echo "Please install DWUI first using the install.sh script." >&2
    exit 1
fi

# Get current version if possible
CURRENT_VERSION=""
if [ -x "$INSTALL_PATH" ]; then
    # Try to get version from the binary (if it supports --version flag)
    CURRENT_VERSION=$("$INSTALL_PATH" --version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
fi

echo "Current DWUI version: ${CURRENT_VERSION:-unknown}"
echo "Target version: $VERSION"

# Check if update is needed
if [ "$CURRENT_VERSION" = "$VERSION" ] && [ "$FORCE_UPDATE" = false ]; then
    echo "DWUI is already at version $VERSION. Use --force to reinstall."
    exit 0
fi

echo "Starting DWUI update to version $VERSION..."

# Stop the service if it's running
if systemctl is-active --quiet dwui; then
    echo "Stopping DWUI service..."
    systemctl stop dwui
    SERVICE_WAS_RUNNING=true
else
    SERVICE_WAS_RUNNING=false
fi

# Create a backup of the current binary
if [ -f "$INSTALL_PATH" ]; then
    BACKUP_PATH="${INSTALL_PATH}.backup.$(date +%Y%m%d_%H%M%S)"
    echo "Creating backup at $BACKUP_PATH..."
    cp "$INSTALL_PATH" "$BACKUP_PATH"
fi

# Download the new binary
echo "Downloading DWUI binary version ${VERSION} from ${DWUI_URL}..."
TMP_BINARY="/tmp/${DWUI_BINARY}"
if ! curl -L -o "$TMP_BINARY" "$DWUI_URL"; then
    echo "Failed to download DWUI binary. Rolling back..." >&2
    if [ -f "$BACKUP_PATH" ]; then
        mv "$BACKUP_PATH" "$INSTALL_PATH"
    fi
    exit 1
fi

# Make it executable
chmod +x "$TMP_BINARY"

# Replace the existing binary
echo "Updating DWUI binary at $INSTALL_PATH..."
mv "$TMP_BINARY" "$INSTALL_PATH"

# Start the service if it was running before
if [ "$SERVICE_WAS_RUNNING" = true ]; then
    echo "Starting DWUI service..."
    systemctl start dwui
fi

echo "DWUI has been successfully updated to version $VERSION."

# Clean up old backup files (keep only the 3 most recent)
if [ -d "$(dirname "$INSTALL_PATH")" ]; then
    BACKUP_DIR="$(dirname "$INSTALL_PATH")"
    BACKUP_COUNT=$(ls -1 "${BACKUP_DIR}/"*.backup.* 2>/dev/null | wc -l || echo 0)
    if [ "$BACKUP_COUNT" -gt 3 ]; then
        echo "Cleaning up old backup files (keeping 3 most recent)..."
        ls -1t "${BACKUP_DIR}/"*.backup.* 2>/dev/null | tail -n +4 | xargs rm -f
    fi
fi

echo "You can check the status with: systemctl status dwui"
echo "If you experience issues, you can restore from backup: sudo mv $BACKUP_PATH $INSTALL_PATH" 