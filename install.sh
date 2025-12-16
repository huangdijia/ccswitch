#!/bin/bash

# ccswitch installation script
# Usage: curl https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default settings
REPO="huangdijia/ccswitch"
INSTALL_DIR="$HOME/.local/bin"
LATEST_RELEASE="true"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--directory)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            LATEST_RELEASE="false"
            shift 2
            ;;
        -h|--help)
            echo "ccswitch Installation Script"
            echo ""
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  -d, --directory DIR    Install directory (default: $HOME/.local/bin)"
            echo "  -v, --version VERSION  Install specific version"
            echo "  -h, --help            Show this help message"
            echo ""
            echo "Examples:"
            echo "  curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash"
            echo "  curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash -s -- -d /usr/local/bin"
            echo "  curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash -s -- -v v1.0.0"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Print colored output
print_info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

print_success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1"
}

print_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

# Detect platform and architecture
detect_platform() {
    local platform
    local arch

    case "$(uname -s)" in
        Darwin)
            platform="Darwin"
            ;;
        Linux)
            platform="Linux"
            ;;
        Windows|CYGWIN*|MINGW*|MSYS*)
            platform="Windows"
            ;;
        *)
            print_error "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)
            arch="x86_64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        arm*)
            arch="armv7"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${platform}_${arch}"
}

# Get latest release version
get_latest_version() {
    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    local version

    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "${api_url}" | grep '"tag_name":' | sed -E 's/.*"tag_name": ?"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "${api_url}" | grep '"tag_name":' | sed -E 's/.*"tag_name": ?"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi

    if [[ -z "$version" ]]; then
        print_error "Failed to fetch latest version from GitHub"
        print_error "Please check your internet connection or specify a version with -v"
        exit 1
    fi

    echo "$version"
}

# Download ccswitch binary
download_binary() {
    local version="$1"
    local platform=$(detect_platform)
    local filename="ccswitch_${version#v}_${platform}"
    local archive_ext="tar.gz"

    # For Windows, use zip format
    if [[ $platform == Windows_* ]]; then
        archive_ext="zip"
    fi

    local archive_name="${filename}.${archive_ext}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"
    local temp_dir=$(mktemp -d)
    local archive_path="${temp_dir}/${archive_name}"

    print_info "Downloading ccswitch ${version} for ${platform}..."

    if command -v curl >/dev/null 2>&1; then
        curl -sSL "${download_url}" -o "${archive_path}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${download_url}" -O "${archive_path}"
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi

    # Verify download was successful
    if [[ ! -f "${archive_path}" ]] || [[ ! -s "${archive_path}" ]]; then
        print_error "Failed to download ccswitch binary"
        print_error "URL: ${download_url}"
        print_error "Please check your internet connection and try again"
        exit 1
    fi

    # Extract the archive
    print_info "Extracting archive..."
    if [[ $archive_ext == "tar.gz" ]]; then
        tar -xzf "${archive_path}" -C "${temp_dir}"
    elif [[ $archive_ext == "zip" ]]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip "${archive_path}" -d "${temp_dir}"
        else
            print_error "unzip is required for Windows archives"
            exit 1
        fi
    fi

    # Find the ccswitch binary in the extracted files
    local binary_path=$(find "${temp_dir}" -name "ccswitch" -type f -perm +111 | head -n 1)
    if [[ -z "${binary_path}" ]]; then
        # Fallback for systems without -perm support
        binary_path=$(find "${temp_dir}" -name "ccswitch" -type f | head -n 1)
    fi

    if [[ -z "${binary_path}" ]]; then
        print_error "Failed to find ccswitch binary in archive"
        print_error "Contents of archive:"
        ls -la "${temp_dir}"
        exit 1
    fi

    # Make binary executable
    chmod +x "${binary_path}"

    # Move to install directory
    if [[ ! -d "${INSTALL_DIR}" ]]; then
        print_info "Creating installation directory: ${INSTALL_DIR}"
        mkdir -p "${INSTALL_DIR}"
    fi

    if mv "${binary_path}" "${INSTALL_DIR}/ccswitch"; then
        print_success "Binary installed to ${INSTALL_DIR}/ccswitch"
    else
        print_error "Failed to move binary to ${INSTALL_DIR}"
        print_error "You may need to run with sudo or choose a different directory"
        exit 1
    fi

    # Clean up temp directory
    rm -rf "${temp_dir}"
}

# Check if install directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        print_warning "${INSTALL_DIR} is not in your PATH"
        echo ""
        echo "To use ccswitch, add the following line to your shell profile:"
        echo ""
        echo -e "${BLUE}export PATH=\"${INSTALL_DIR}:\$PATH\"${NC}"
        echo ""
        echo "Then restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"
    fi
}

# Verify installation
verify_installation() {
    local ccswitch_path="${INSTALL_DIR}/ccswitch"

    if [[ -f "${ccswitch_path}" ]]; then
        local version_output
        version_output=$("${ccswitch_path}" --version 2>/dev/null || echo "version unknown")
        print_success "ccswitch ${version_output} installed successfully!"
        echo ""
        echo "Quick start:"
        echo "  ccswitch profiles   # List all available profiles"
        echo "  ccswitch use <name> # Switch to a specific profile"
        echo "  ccswitch init       # Initialize configuration"
        echo ""
        echo "Documentation: https://github.com/${REPO}"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Main installation logic
main() {
    print_info "Starting ccswitch installation..."

    # Check if running on supported platform
    if ! detect_platform >/dev/null 2>&1; then
        exit 1
    fi

    # Determine version to install
    if [[ "$LATEST_RELEASE" == "true" ]]; then
        VERSION=$(get_latest_version)
        print_info "Installing latest version: ${VERSION}"
    else
        print_info "Installing version: ${VERSION}"
    fi

    # Download binary
    download_binary "$VERSION"

    # Check PATH
    check_path

    # Verify installation
    verify_installation
}

# Run main function
main "$@"