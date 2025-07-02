#!/usr/bin/env bash
set -euo pipefail

REPO="anyproto/anytype-heart"
GITHUB="api.github.com"

if [[ $# -eq 2 ]]; then
    PLATFORM="$1"
    ARCH="$2"
    echo "Using provided platform: $PLATFORM-$ARCH"
else
    echo "Anytype Heart Download Script"
    echo "============================="
    echo ""
    echo "Select platform:"
    echo "1) Linux AMD64 (default)"
    echo "2) Linux ARM64"
    echo "3) macOS Apple Silicon (ARM64)"
    echo "4) macOS Intel (AMD64)"
    echo "5) Windows AMD64"
    echo ""
    read -rp "Enter choice [1-5] (default: 1): " choice

    case ${choice:-1} in
        1)
            PLATFORM="linux"
            ARCH="amd64"
            ;;
        2)
            PLATFORM="linux"
            ARCH="arm64"
            ;;
        3)
            PLATFORM="darwin"
            ARCH="arm64"
            ;;
        4)
            PLATFORM="darwin"
            ARCH="amd64"
            ;;
        5)
            PLATFORM="windows"
            ARCH="amd64"
            ;;
        *)
            echo "Invalid choice. Using default (Linux AMD64)"
            PLATFORM="linux"
            ARCH="amd64"
            ;;
    esac
fi

if [[ "$PLATFORM" == "windows" ]]; then
    EXT="zip"
else
    EXT="tar.gz"
fi

ASSET_NAME="js_.*_${PLATFORM}-${ARCH}.${EXT}"

TAG=$(curl -s \
  -H "Accept: application/vnd.github.v3+json" \
  "https://$GITHUB/repos/$REPO/releases/latest" \
  | jq -r .tag_name)

ASSET_INFO=$(curl -s \
    -H "Accept: application/vnd.github.v3+json" \
    "https://$GITHUB/repos/$REPO/releases/tags/$TAG" \
  | jq -r --arg pattern "$ASSET_NAME" \
      '.assets[] | select(.name | test($pattern)) | "\(.id) \(.name)"')

if [[ -z "$ASSET_INFO" ]]; then
  echo "No asset found matching pattern: $ASSET_NAME" >&2
  exit 1
fi

ASSET_ID=$(echo "$ASSET_INFO" | cut -d' ' -f1)
ASSET_FILENAME=$(echo "$ASSET_INFO" | cut -d' ' -f2)

mkdir -p dist
OUT_FILE="dist/$ASSET_FILENAME"
echo ""
echo "Downloading: $OUT_FILE..."

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