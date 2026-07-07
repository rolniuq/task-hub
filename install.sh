#!/usr/bin/env bash
set -e

# Task Hub Desktop — macOS Installer
# Usage: curl -sfL https://raw.githubusercontent.com/rolniuq/task-hub/main/install.sh | sh

BOLD='\033[1m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo ""
echo -e "${BOLD}🚀 Task Hub Desktop — macOS Installer${NC}"
echo ""

# ---- Check prerequisites ----

if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo "Install it first: https://go.dev/dl/"
    echo "Or via Homebrew: brew install go"
    exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
echo -e "  ✓ Go ${GO_VERSION} found"

if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: git is not installed.${NC}"
    echo "Install Xcode Command Line Tools: xcode-select --install"
    exit 1
fi
echo -e "  ✓ git found"

# ---- Create temp directory ----
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT
echo ""

# ---- Clone repo ----
echo -e "${YELLOW}Downloading Task Hub...${NC}"
git clone --depth 1 https://github.com/rolniuq/task-hub.git "$TEMP_DIR" 2>/dev/null
echo -e "  ✓ Source downloaded"

# ---- Build ----
echo -e "${YELLOW}Building desktop application...${NC}"
cd "$TEMP_DIR"
go build -o task-hub-desktop ./cmd/desktop/main.go
echo -e "  ✓ Build complete"

# ---- Install ----
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    # Fall back to user-local bin
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

mv task-hub-desktop "$INSTALL_DIR/task-hub-desktop"
chmod +x "$INSTALL_DIR/task-hub-desktop"

echo ""
echo -e "${GREEN}✅ Task Hub Desktop installed to ${INSTALL_DIR}/task-hub-desktop${NC}"
echo ""

# ---- Check PATH ----
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠️  ${INSTALL_DIR} is not in your PATH.${NC}"
    echo "   Add it to your shell profile:"
    echo "   echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.zshrc"
    echo "   source ~/.zshrc"
    echo ""
fi

echo -e "${BOLD}Usage:${NC}"
echo "  1. Create a .env file with your database config (see .env.example)"
echo "  2. Run: task-hub-desktop"
echo ""
echo -e "${BOLD}Note:${NC} PostgreSQL and NATS must be running. Use Docker:"
echo "  docker compose -f $TEMP_DIR/docker-compose.yml up -d db nats"
echo ""
