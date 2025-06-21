#!/usr/bin/env bash
# lib/validate.sh – Go version and dependency validation

validate_versions() {
    local REQUIRED_GO_VERSION="${_VERSION_GO:-1.20.0}"
    local GO_VERSION
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [[ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]]; then
        log error "Go version must be >= $REQUIRED_GO_VERSION. Detected: $GO_VERSION"
        exit 1
    fi
    log success "Valid Go version: $GO_VERSION"
    go mod tidy || return 1
}

check_dependencies() {
    for dep in "$@"; do
        if ! command -v "$dep" > /dev/null; then
            log error "$dep is not installed."
            exit 1
        else
            log success "$dep is installed."
        fi
    done
}

export -f validate_versions
export -f check_dependencies
