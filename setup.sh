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
    
    while true; do
        clear
        echo "Anytype Heart Download Script"
        echo "============================="
        echo ""
        echo "Use ↑/↓ arrow keys to move, Enter to select:"
        echo ""
        
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
                '[A')
                    if [ $selected -eq 0 ]; then
                        selected=$((${#options[@]} - 1))
                    else
                        ((selected--))
                    fi
                    ;;
                '[B')
                    if [ $selected -eq $((${#options[@]} - 1)) ]; then
                        selected=0
                    else
                        ((selected++))
                    fi
                    ;;
            esac
        elif [[ $key == "" ]]; then
            break
        fi
    done
    
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
    
    clear
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