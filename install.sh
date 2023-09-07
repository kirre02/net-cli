#!/bin/bash

set -euo pipefail

# Set the name of your Go binary
BINARY_NAME="net-cli"

# Build the Go binary
go build -o "${BINARY_NAME}"

# Remove the current binary if it exists
if [ -f "/usr/local/bin/${BINARY_NAME}" ]; then
    echo "Previous installation found. Removing the existing binary..."
    sudo rm "/usr/local/bin/${BINARY_NAME}"
fi

# Move the new binary to /usr/local/bin
sudo mv "./${BINARY_NAME}" "/usr/local/bin/${BINARY_NAME}"

echo "Installation complete! The binary has been updated."
