#!/bin/sh
# te_demo installer for Linux / macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/YOUR_ORG/te_demo/main/install.sh | sh

set -e

REPO="YOUR_ORG/te_demo"
BIN="te_demo"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)       ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

ASSET="${BIN}-${OS}-${ARCH}"

echo "Installing te_demo (${OS}/${ARCH})..."

API="https://api.github.com/repos/${REPO}/releases/latest"
URL=$(curl -fsSL "$API" | grep "browser_download_url" | grep "$ASSET" | cut -d'"' -f4)
VER=$(curl -fsSL "$API" | grep '"tag_name"' | cut -d'"' -f4)

echo "Version: $VER"

TMP=$(mktemp)
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

if [ -w "/usr/local/bin" ]; then
  mv "$TMP" "/usr/local/bin/$BIN"
  echo "Installed to /usr/local/bin/$BIN"
else
  mkdir -p "$HOME/.local/bin"
  mv "$TMP" "$HOME/.local/bin/$BIN"
  echo "Installed to $HOME/.local/bin/$BIN"
  echo "Add to PATH: export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Done! Run: te_demo help"
