#!/bin/bash

# Script to calculate SHA256 for Homebrew formula
# Usage: ./scripts/homebrew-sha.sh [version]

set -e

VERSION=${1:-$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")}
VERSION_NO_V=${VERSION#v}

echo "Calculating SHA256 for ccswitch version: $VERSION"
echo "URL: https://github.com/huangdijia/ccswitch/archive/refs/tags/$VERSION.tar.gz"
echo ""

# Download the source
echo "Downloading source..."
curl -sL "https://github.com/huangdijia/ccswitch/archive/refs/tags/$VERSION.tar.gz" -o source.tar.gz

# Calculate SHA256
echo "Calculating SHA256..."
SHA256=$(shasum -a 256 source.tar.gz | awk '{print $1}')

echo "SHA256: $SHA256"
echo ""

# Generate formula snippet
cat << EOF
For Homebrew formula:

url "https://github.com/huangdijia/ccswitch/archive/refs/tags/$VERSION.tar.gz"
sha256 "$SHA256"
EOF

echo ""
echo "Cleaning up..."
rm -f source.tar.gz

echo ""
echo "Formula ready for version $VERSION_NO_V"