#!/bin/bash

# Exit on any error
set -e

# --- Argument Parsing ---
# Try to read version from local VERSION file, then from GitHub, then fallback
if [ -f "$(dirname "$0")/VERSION" ]; then
    VERSION=$(cat "$(dirname "$0")/VERSION" | tr -d '[:space:]')
elif command -v curl >/dev/null 2>&1; then
    VERSION=$(curl -sSL https://raw.githubusercontent.com/romerramos/dwui/main/dist/VERSION 2>/dev/null | tr -d '[:space:]' || echo "v0.0.1")
else
    VERSION="v0.0.1"  # Fallback when curl is not available
fi
# Save original args to pass them to the binary later
ARGS=("$@")

# Parse arguments to find a custom version
# We only need the version for the download URL. The rest of the args are for the binary.
for i in "${!ARGS[@]}"; do
    if [[ "${ARGS[$i]}" == "--version" ]]; then
        # Found --version, so grab the next element as the value
        VERSION="${ARGS[$i+1]}"
        break
    fi
done


# --- Platform Detection ---
OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
PLATFORM=""

case "$OS_TYPE" in
    linux)
        case "$ARCH" in
            x86_64)
                PLATFORM="linux-amd64"
                ;;
            aarch64 | arm64)
                PLATFORM="linux-arm64"
                ;;
        esac
        ;;
    darwin)
        case "$ARCH" in
            x86_64)
                PLATFORM="mac-amd64"
                ;;
            arm64)
                PLATFORM="mac-arm64"
                ;;
        esac
        ;;
esac

if [ -z "$PLATFORM" ]; then
    echo "Unsupported platform: ${OS_TYPE}-${ARCH}" >&2
    echo "Supported platforms are Linux (amd64, arm64) and macOS (amd64, arm64)." >&2
    exit 1
fi

# --- Download and Execute ---
BINARY_NAME="dwui-${PLATFORM}"
DOWNLOAD_URL="https://github.com/romerramos/dwui/releases/download/${VERSION}/${BINARY_NAME}"
TMP_BINARY="/tmp/${BINARY_NAME}"

echo "Downloading DWUI ${VERSION} for ${PLATFORM} using ${DOWNLOAD_URL}..."
curl -sSL -o "$TMP_BINARY" "$DOWNLOAD_URL"

echo "Making the binary executable..."
chmod +x "$TMP_BINARY"

echo "Starting DWUI..."
# Use exec to replace the script process with the dwui process
# Pass the original arguments to the binary
exec "$TMP_BINARY" "${ARGS[@]}" 