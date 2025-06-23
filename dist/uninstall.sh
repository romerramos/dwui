#!/bin/bash

# Exit on any error
set -e

# Check for root user
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please use sudo." >&2
  exit 1
fi

INSTALL_PATH="/usr/local/bin/dwui"
SERVICE_FILE="/etc/systemd/system/dwui.service"

echo "Starting DWUI uninstallation..."

# Stop the service
echo "Stopping DWUI service..."
systemctl stop dwui

# Disable the service
echo "Disabling DWUI service..."
systemctl disable dwui

# Remove the service file
echo "Removing systemd service file..."
rm -f "$SERVICE_FILE"

# Reload systemd
echo "Reloading systemd..."
systemctl daemon-reload

# Remove the binary
if [ -f "$INSTALL_PATH" ]; then
    echo "Removing DWUI binary from $INSTALL_PATH..."
    rm -f "$INSTALL_PATH"
fi

echo "DWUI has been uninstalled successfully." 