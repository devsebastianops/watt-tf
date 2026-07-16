#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

OWNER="devsebastianops"
REPO="watt-tf"
BINARY="wtf"

# Color constants for nice output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0;3m' # No Color
NC_BOLD='\033[1;0m'

log_info() {
    echo -e "${BLUE}[info]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[success]${NC} $1"
}

log_error() {
    echo -e "${RED}[error]${NC} $1" >&2
}

# 1. Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Normalize OS names to match your matrix [linux, darwin]
case "${OS}" in
    linux)  PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    *)
        log_error "Unsupported Operating System: ${OS}. Currently only Linux and macOS are supported by this installer."
        exit 1
        ;;
esac

# Normalize CPU architectures to match your matrix [amd64, arm64]
case "${ARCH}" in
    x86_64|amd64) ARCH_NORM="amd64" ;;
    arm64|aarch64) ARCH_NORM="arm64" ;;
    *)
        log_error "Unsupported CPU Architecture: ${ARCH}. Currently only amd64 and arm64 are supported."
        exit 1
        ;;
esac

# 2. Determine target version (Defaults to latest release)
if [ -z "${VERSION}" ]; then
    log_info "Fetching latest release version from GitHub..."
    # Query GitHub API to get the latest tag name
    LATEST_RELEASE_URL="https://api.github.com/repos/${OWNER}/${REPO}/releases/latest"
    VERSION=$(curl -sSL "${LATEST_RELEASE_URL}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "${VERSION}" ]; then
        log_error "Could not fetch latest release version. Are you offline or rate-limited by GitHub?"
        exit 1
    fi
fi

log_info "Target version set to: ${VERSION}"

# 3. Build download URL based on your naming convention:
# wtf-${goos}-${goarch}-${VERSION}.tar.gz
ARCHIVE_NAME="${BINARY}-${PLATFORM}-${ARCH_NORM}-${VERSION}.tar.gz"
DOWNLOAD_URL="https://github.com/../../releases/download/${VERSION}/${ARCHIVE_NAME}"

# Fallback-URL falls Releases per Redirect geladen werden müssen
STABLE_DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

# Create a temporary directory for downloading and unpacking
TMP_DIR=$(mktemp -d)
clean_up() {
    rm -rf "${TMP_DIR}"
}
trap clean_up EXIT

log_info "Downloading ${ARCHIVE_NAME}..."
if ! curl -sSL --fail -o "${TMP_DIR}/${ARCHIVE_NAME}" "${STABLE_DOWNLOAD_URL}"; then
    log_error "Failed to download release asset from: ${STABLE_DOWNLOAD_URL}"
    log_error "Please verify that the release version ${VERSION} exists."
    exit 1
fi

# 4. Unpack the binary
log_info "Extracting binary..."
tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"

# 5. Determine installation path
# We prefer standard user binary paths that don't require sudo if possible
if [ -d "${HOME}/.local/bin" ]; then
    INSTALL_DIR="${HOME}/.local/bin"
    USE_SUDO=false
elif [ -d "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
    # Check if we have write access to /usr/local/bin without sudo
    if [ -w "/usr/local/bin" ]; then
        USE_SUDO=false
    else
        USE_SUDO=true
    fi
else
    # Fallback to creating a local bin directory
    INSTALL_DIR="${HOME}/.local/bin"
    mkdir -p "${INSTALL_DIR}"
    USE_SUDO=false
fi

log_info "Installing ${BINARY} to ${INSTALL_DIR}..."

if [ "${USE_SUDO}" = true ]; then
    log_info "Write permissions required for ${INSTALL_DIR}. Requesting sudo..."
    sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY}"
else
    mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    chmod +x "${INSTALL_DIR}/${BINARY}"
fi

log_success "Successfully installed Watt TF (${BINARY}) to ${INSTALL_DIR}/${BINARY}"

# 6. Check if install directory is in user's PATH
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo -e "\n${RED}[warning]${NC} ${INSTALL_DIR} is not in your PATH."
    echo -e "You might need to add it to your shell configuration (e.g., ~/.bashrc or ~/.zshrc):"
    echo -e "  export PATH=\"\$PATH:${INSTALL_DIR}\""
else
    echo -e "\nTry it out by running:"
    echo -e "  ${GREEN}wtf --help${NC}"
fi