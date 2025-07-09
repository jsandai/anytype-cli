#!/usr/bin/env bash
set -euo pipefail

REPO="anyproto/anytype-heart"
GITHUB="api.github.com"

MENU_SELECTION=0

select_option() {
    local options=("$@")
    local selected=0
    local key=""
    
    tput civis

    echo "Choose your platform and architecture."
    echo "Use ↑/↓ arrow keys to move, Enter to select:"
    echo ""
    
    while true; do
        # Draw menu
        for i in "${!options[@]}"; do
            if [ "$i" -eq $selected ]; then
                echo "  ▶ ${options[$i]}"
            else
                echo "    ${options[$i]}"
            fi
        done
        
        read -rsn1 key
        
        if [[ $key == $'\x1b' ]]; then
            read -rsn2 key
            case $key in
                '[A') ((selected = selected == 0 ? ${#options[@]} - 1 : selected - 1)) ;;
                '[B') ((selected = (selected + 1) % ${#options[@]})) ;;
            esac
        elif [[ $key == "" ]]; then
            break
        fi
        
        # Move cursor back up to redraw menu
        tput cuu ${#options[@]}
    done
    
    # Clear menu lines
    for ((i=0; i<${#options[@]}; i++)); do
        tput el
        [ $i -lt $((${#options[@]} - 1)) ] && tput cud1
    done
    tput cuu $((${#options[@]} - 1))
    
    tput cnorm
    MENU_SELECTION=$selected
}

if [[ $# -eq 2 ]]; then
    PLATFORM="$1"
    ARCH="$2"
    echo "Using provided platform: $PLATFORM-$ARCH"
else
    options=(
        "Linux AMD64"
        "Linux ARM64"
        "macOS Apple Silicon (ARM64)"
        "macOS Intel (AMD64)"
        "Windows AMD64"
    )
    
    select_option "${options[@]}"
    choice=$MENU_SELECTION
    
    case $choice in
        0)
            PLATFORM="linux"
            ARCH="amd64"
            ;;
        1)
            PLATFORM="linux"
            ARCH="arm64"
            ;;
        2)
            PLATFORM="darwin"
            ARCH="arm64"
            ;;
        3)
            PLATFORM="darwin"
            ARCH="amd64"
            ;;
        4)
            PLATFORM="windows"
            ARCH="amd64"
            ;;
    esac
    
    echo ""
    echo "Selected: ${options[$choice]}"
fi

if [[ "$PLATFORM" == "windows" ]]; then
    EXT="zip"
else
    EXT="tar.gz"
fi

ASSET_NAME="js_.*_${PLATFORM}-${ARCH}.${EXT}"

RESPONSE=$(curl -s \
  -H "Accept: application/vnd.github.v3+json" \
  "https://$GITHUB/repos/$REPO/releases/latest")

TAG=$(echo "$RESPONSE" | jq -r .tag_name 2>/dev/null)

if [[ -z "$TAG" || "$TAG" == "null" ]]; then
  echo "Failed to fetch latest release tag" >&2
  echo "Error: $RESPONSE" >&2
  exit 1
fi

echo ""
echo "Latest release: $TAG"
echo "Downloading: $PLATFORM-$ARCH"
echo ""

ASSET_INFO=$(curl -s \
    -H "Accept: application/vnd.github.v3+json" \
    "https://$GITHUB/repos/$REPO/releases/tags/$TAG" \
  | jq -r --arg pattern "$ASSET_NAME" \
      '.assets[] | select(.name | test($pattern)) | "\(.id) \(.name)"' 2>/dev/null)

if [[ -z "$ASSET_INFO" ]]; then
  echo "No asset found matching pattern: $ASSET_NAME" >&2
  exit 1
fi

ASSET_ID=$(echo "$ASSET_INFO" | cut -d' ' -f1)
ASSET_FILENAME=$(echo "$ASSET_INFO" | cut -d' ' -f2)

mkdir -p dist
OUT_FILE="dist/$ASSET_FILENAME"
curl -L \
  -H "Accept: application/octet-stream" \
  "https://$GITHUB/repos/$REPO/releases/assets/$ASSET_ID" \
  -o "$OUT_FILE"

cd dist
if [[ "$ASSET_FILENAME" == *.zip ]]; then
  unzip -o "$ASSET_FILENAME"
else
  tar -zxf "$ASSET_FILENAME"
fi

rm -f "$ASSET_FILENAME"
cd ..

echo "Downloaded and extracted successfully!"

# Make the server executable
if [[ -f "dist/grpc-server" ]]; then
    chmod +x dist/grpc-server
    SERVER_NAME="grpc-server"
else
    echo "✗ Could not find server executable (grpc-server or mw)"
    SERVER_NAME="unknown"
fi

echo ""
echo "Anytype middleware server downloaded to: dist/$SERVER_NAME"
echo ""
echo "To use the Anytype CLI:"
echo "1. Build the CLI: make build"
echo "2. Install the CLI: make install (or make install-local)"
echo "3. Start the daemon: anytype daemon"
echo "4. Start the server: anytype server start"