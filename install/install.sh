#!/bin/bash

set -e

echo "Welcome to MoLing MCP Server initialization script."
echo "Home page: https://gojue.cc/moling"
echo "Github: https://github.com/gojue/moling"

# Determine the OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Determine the download URL
VERSION=$(curl -s https://api.github.com/repos/gojue/moling/releases/latest | grep 'tag_name' | cut -d\" -f4)
BASE_URL="https://github.com/gojue/moling/releases/download/${VERSION}"
FILE_NAME="moling-${VERSION}-${OS}-${ARCH}.tar.gz"

DOWNLOAD_URL="${BASE_URL}/${FILE_NAME}"

# Download the installation package
echo "Downloading ${DOWNLOAD_URL}..."
curl -LO "${DOWNLOAD_URL}"
echo "Download completed. filename: ${FILE_NAME}"
# Extract the package
tar -xzf "${FILE_NAME}"

# Move the binary to /usr/local/bin
mv moling /usr/local/bin/moling
chmod +x /usr/local/bin/moling

# Clean up
rm -rf moling "${FILE_NAME}"

# Check if the installation was successful
if command -v moling &> /dev/null; then
    echo "MoLing installation was successful!"
else
    echo "MoLing installation failed."
    exit 1
fi

# initialize the configuration
echo "Initializing MoLing configuration..."
moling config --init
echo "MoLing configuration initialized successfully!"

echo "setup MCP Server configuration into MCP Client"
moling client -i -d
echo "MCP Client configuration setup successfully!"