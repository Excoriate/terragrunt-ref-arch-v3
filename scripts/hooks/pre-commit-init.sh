#!/usr/bin/env bash
# shellcheck disable=SC2317

# Pre-Commit Hook Management Script
#
# This script provides functionality to manage pre-commit hooks for the repository.
# It follows Google's Bash Style Guide and provides reliable hook management.

# Strict error handling
set -euo pipefail

# Logging function with timestamp and color
log() {
    local -r level="${1}"
    local -r message="${2}"
    local -r timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    local color=""

    case "${level}" in
        INFO)    color="\033[0;32m" ;;  # Green
        WARNING) color="\033[0;33m" ;;  # Yellow
        ERROR)   color="\033[0;31m" ;;  # Red
        *)       color="\033[0m" ;;     # Default
    esac

    # shellcheck disable=SC2059
    printf "${color}[${level}] ${timestamp}: ${message}\033[0m\n" >&2
}

# Ensure pre-commit is installed
ensure_pre_commit_installed() {
    if ! command -v pre-commit &> /dev/null; then
        log ERROR "pre-commit is not installed. Installing via pip..."
        if ! pip3 install pre-commit; then
            log ERROR "Failed to install pre-commit. Please install manually."
            return 1
        fi
    fi
    log INFO "pre-commit is installed and ready."
}

# Verify hook installation
verify_hook_installation() {
    local hook_types=("pre-commit" "pre-push")
    local git_dir
    git_dir=$(git rev-parse --git-dir)

    for hook_type in "${hook_types[@]}"; do
        if [ ! -f "${git_dir}/hooks/${hook_type}" ]; then
            log ERROR "Hook ${hook_type} not installed correctly"
            return 1
        fi
    done

    log INFO "All Git hooks verified successfully"
}

# Function to safely unset core.hooksPath from Git configuration
# This resolves conflicts with pre-commit installation
unset_core_hooks_path() {
    local hooks_path_value

    # Check if core.hooksPath is set and get its value
    if hooks_path_value=$(git config core.hooksPath 2>/dev/null); then
        log INFO "Found core.hooksPath set to: ${hooks_path_value}"
        log INFO "Pre-commit requires exclusive control over Git hooks directory"
        log INFO "Attempting to unset core.hooksPath to allow pre-commit installation..."

        # Try to unset from local repository first
        if git config --local --unset-all core.hooksPath 2>/dev/null; then
            log INFO "Successfully unset core.hooksPath from local repository configuration"
            return 0
        fi

        # Try to unset from global configuration
        if git config --global --unset-all core.hooksPath 2>/dev/null; then
            log INFO "Successfully unset core.hooksPath from global Git configuration"
            return 0
        fi

        # If both fail, provide helpful guidance
        log WARNING "Could not automatically unset core.hooksPath"
        log WARNING "This may be set at system level or require different permissions"
        log WARNING "You may need to manually run one of these commands:"
        log WARNING "  git config --local --unset-all core.hooksPath"
        log WARNING "  git config --global --unset-all core.hooksPath"
        log WARNING "  git config --system --unset-all core.hooksPath (requires admin)"
        log WARNING "Continuing with pre-commit installation attempt..."
        return 1
    else
        log INFO "core.hooksPath is not set - pre-commit installation should proceed normally"
        return 0
    fi
}

# Install pre-commit hooks
# Exposed function 1: Initialize repository hooks
pc_init() {
    ensure_pre_commit_installed

    log INFO "Updating pre-commit hooks to the latest version..."
    if ! pre-commit autoupdate; then
        log WARNING "Failed to update pre-commit hooks to the latest version. Continuing with existing hooks."
    fi

    # Handle core.hooksPath configuration that conflicts with pre-commit
    unset_core_hooks_path

    log INFO "Installing pre-commit hooks..."
    if ! pre-commit install; then
        log ERROR "Failed to install pre-commit hooks"
        log ERROR "If you see 'Cowardly refusing to install hooks with core.hooksPath set':"
        log ERROR "  Run: git config --unset-all core.hooksPath"
        log ERROR "  Then retry: just hooks-install"
        return 1
    fi

    if ! pre-commit install --hook-type pre-commit; then
        log ERROR "Failed to install pre-commit hooks for commit stage"
        return 1
    fi

    if ! pre-commit install --hook-type pre-push; then
        log ERROR "Failed to install pre-commit hooks for pre-push stage"
        return 1
    fi

    # Verify hook installation
    if ! verify_hook_installation; then
        log ERROR "Hook verification failed. Please check your Git configuration."
        return 1
    fi

    log INFO "🎉 All pre-commit hooks installed and verified successfully!"
    log INFO "Your repository is now configured with automated code quality checks"
}

# Run pre-commit hooks on all files
# Exposed function 2: Run hooks across all files
pc_run() {
    log INFO "Running pre-commit hooks on all files..."
    if ! pre-commit run --all-files; then
        log ERROR "Pre-commit hooks failed on some files"
        return 1
    fi
    log INFO "Pre-commit hooks completed successfully"
}

# Main function for script execution
main() {
    local command="${1:-}"

    case "${command}" in
        init)
            pc_init
            ;;
        run)
            pc_run
            ;;
        *)
            log ERROR "Invalid command. Use 'init' or 'run'."
            exit 1
            ;;
    esac
}

# Allow sourcing for function access or direct script execution
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
