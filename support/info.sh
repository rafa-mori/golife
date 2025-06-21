#!/usr/bin/env bash
# lib/info.sh – Functions to display banners and installation summary

show_about() {
    printf '%s\n\n' "${_ABOUT:-}"
}

show_banner() {
    printf '\n%s\n\n' "${_BANNER:-}"
}

show_headers() {
    show_banner || return 1
    show_about || return 1
}

summary() {
    local install_dir="$_BINARY"
    log success "Build and installation completed!"
    log success "Binary: $_BINARY"
    log success "Installed in: ${install_dir}"
    check_path "$install_dir"
}

export -f show_about
export -f show_banner
export -f show_headers
export -f summary

