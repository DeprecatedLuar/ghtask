#!/bin/bash

# ghtask Install Script
# Usage: ./install-local.sh

set -e

cd "$(dirname "$0")/.."

echo "Building ghtask..."
go build -o ghtask cmd/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"

    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"

    echo "Installing to $INSTALL_DIR/..."
    cp ghtask "$INSTALL_DIR/"
    echo "Done! Shortcuts (gt, gt0-gt3) will be auto-created on first run."
else
    echo "Build failed!"
    exit 1
fi
