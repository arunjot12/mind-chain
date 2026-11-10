#!/bin/bash

# Download mc-deployer pre-built binary
# Based on the tutorial instructions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect platform
detect_platform() {
    case "$(uname -s)" in
        Darwin)
            case "$(uname -m)" in
                arm64)
                    echo "darwin-arm64"
                    ;;
                x86_64)
                    echo "darwin-amd64"
                    ;;
                *)
                    log_error "Unsupported macOS architecture: $(uname -m)"
                    exit 1
                    ;;
            esac
            ;;
        Linux)
            echo "linux-amd64"
            ;;
        *)
            log_error "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac
}

# Download mc-deployer
download_op_deployer() {
    # Detect host platform
    local os
    local arch
    case "$(uname -s)" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        *)
            log_error "Unsupported OS: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        aarch64|arm64)
            arch="arm64"
            ;;
        x86_64|amd64)
            arch="amd64"
            ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
    local platform="$os-$arch"
    local releases_url="https://api.github.com/repos/ethereum-optimism/optimism/releases"

    log_info "Detecting platform: $platform"

    # Find the latest mc-deployer release
    log_info "Finding latest mc-deployer release..."

    # Get all releases and find the latest one with mc-deployer assets
    local latest_deployer_release
    latest_deployer_release=$(curl -s "$releases_url?per_page=50" | jq -r '.[] | select(.tag_name | startswith("mc-deployer/")) | .tag_name' | sort -V | tail -1)

    if [ -z "$latest_deployer_release" ]; then
        log_error "Could not find any mc-deployer releases"
        exit 1
    fi

    local tag_name="$latest_deployer_release"
    log_info "Found latest mc-deployer release: $tag_name"

    # Get the release info
    local release_info
    release_info=$(curl -s "$releases_url/tags/$tag_name")

    if [ -z "$release_info" ] || [ "$release_info" = "null" ] || echo "$release_info" | jq -e '.message' >/dev/null 2>&1; then
        log_error "Could not fetch release information for $tag_name"
        exit 1
    fi

    local has_deployer_assets
    has_deployer_assets=$(echo "$release_info" | jq -r '.assets[] | select(.name | contains("mc-deployer")) | .name' | wc -l)

    if [ "$has_deployer_assets" -eq 0 ]; then
        log_error "Release $tag_name does not have mc-deployer assets"
        exit 1
    fi

    log_info "Release $tag_name has mc-deployer assets"

    # Check if this release has mc-deployer assets for our platform
    local asset_name
    asset_name=$(echo "$release_info" | jq -r ".assets[] | select(.name | contains(\"mc-deployer\") and contains(\"$platform\")) | .name")

    # If arm64 is not available, fall back to amd64 (Rosetta compatibility)
    if [ -z "$asset_name" ] && [ "$platform" = "linux-arm64" ]; then
        log_info "linux-arm64 not available, trying linux-amd64 (Rosetta compatible)"
        platform="linux-amd64"
        asset_name=$(echo "$release_info" | jq -r ".assets[] | select(.name | contains(\"mc-deployer\") and contains(\"$platform\")) | .name")
    fi

    if [ -z "$asset_name" ]; then
        log_error "Could not find mc-deployer asset for platform $platform in release $tag_name"
        log_info "Available mc-deployer assets in this release:"
        echo "$release_info" | jq -r '.assets[] | select(.name | contains("mc-deployer")) | .name'
        log_info "Please check: https://github.com/ethereum-optimism/optimism/releases/tag/$tag_name"
        exit 1
    fi

    local download_url="https://github.com/ethereum-optimism/optimism/releases/download/$tag_name/$asset_name"

    log_info "Downloading mc-deployer $tag_name for $platform..."
    log_info "From: $download_url"

    # Download the asset
    if ! curl -L -o "mc-deployer.tar.gz" "$download_url"; then
        log_error "Failed to download mc-deployer"
        exit 1
    fi

    # Extract (assuming it's a tar.gz)
    log_info "Extracting mc-deployer..."
    tar -xzf mc-deployer.tar.gz

    # Find the binary (it might be in a subdirectory)
    local binary_path

    # Look for the extracted directory first
    local extracted_dir
    extracted_dir=$(find . -name "mc-deployer-*" -type d | head -1)

    if [ -n "$extracted_dir" ] && [ -f "$extracted_dir/mc-deployer" ]; then
        binary_path="$extracted_dir/mc-deployer"
    elif [ -f "mc-deployer" ]; then
        binary_path="mc-deployer"
    elif [ -f "mc-deployer-$platform" ]; then
        binary_path="mc-deployer-$platform"
    else
        # Look for it anywhere (macOS compatible)
        binary_path=$(find . -name "mc-deployer*" -type f -perm +111 | head -1)
    fi

    if [ -z "$binary_path" ]; then
        log_error "Could not find mc-deployer binary in extracted files"
        ls -la
        exit 1
    fi

    # Make executable and move to current directory
    chmod +x "$binary_path"
    mv "$binary_path" ./mc-deployer

    # Cleanup
    rm -f mc-deployer.tar.gz
    rm -rf "${binary_path%/*}" 2>/dev/null || true

    # Test the binary (only if we're downloading for the same platform)
    local current_platform=$(detect_platform)
    if [ "$platform" = "$current_platform" ]; then
        if ! ./mc-deployer --version >/dev/null 2>&1; then
            log_error "Downloaded mc-deployer binary appears to be invalid"
            exit 1
        fi
        log_success "mc-deployer downloaded and ready: $(./mc-deployer --version)"
    else
        # Cross-platform download - just verify the file exists and is executable
        if [ ! -x "./mc-deployer" ]; then
            log_error "Downloaded mc-deployer binary is not executable"
            exit 1
        fi
        log_success "mc-deployer downloaded for $platform platform (cannot test on $current_platform host)"
    fi
}

# Alternative: Manual download instructions
show_manual_instructions() {
    log_warning "Automatic download failed. Please download manually:"
    echo ""
    echo "1. Go to: https://github.com/ethereum-optimism/optimism/releases"
    echo "2. Find the latest release that includes mc-deployer"
    echo "3. Download the appropriate binary for your system:"
    echo "   - Linux: mc-deployer-linux-amd64"
    echo "   - macOS Intel: mc-deployer-darwin-amd64"
    echo "   - macOS Apple Silicon: mc-deployer-darwin-arm64"
    echo "4. Extract and place the binary as './mc-deployer'"
    echo "5. Run: chmod +x ./mc-deployer"
    echo ""
    exit 1
}

# Main execution
main() {
    log_info "Downloading mc-deployer..."

    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        log_error "jq is required but not installed"
        exit 1
    fi

    # Try automatic download
    if download_op_deployer; then
        log_success "mc-deployer setup complete!"
        exit 0
    else
        show_manual_instructions
    fi
}

# Run main function
main "$@"
