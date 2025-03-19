#!/usr/bin/env bash
# =============================================================================
# Terragrunt Reference Architecture - Justfile Utilities
# =============================================================================
# Utility functions for use with Justfile recipes
# This file contains functions that can be called from Justfile to keep
# recipes clean and maintainable

# =============================================================================
# CORE LOGGING FUNCTIONS
# =============================================================================

# Simple logging function with timestamp and level
# Usage: log_message "INFO" "Your message here"
log_message() {
  local log_level="${1:-INFO}"
  local message="$2"
  local timestamp
  timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  echo "[${log_level}] ${timestamp} - ${message}" >&2
}

# =============================================================================
# TERRAGRUNT UTILITIES
# =============================================================================

# Format Terragrunt HCL files
# Usage: terragrunt_format "/path/to/terragrunt/dir" "check" "diff" "exclude_pattern"
terragrunt_format() {
  local terragrunt_dir="$1"
  local check_mode="$2"
  local diff_mode="$3"
  local exclude_pattern="$4"

  log_message "INFO" "🔍 Advanced Terragrunt HCL Formatting"

  # Validate inputs
  if [[ ! -d "${terragrunt_dir}" ]]; then
    log_message "ERROR" "Terragrunt directory does not exist: ${terragrunt_dir}"
    return 1
  fi

  # Set up the command base
  local cmd="terragrunt hclfmt"

  # Add options based on parameters
  if [[ "${check_mode}" == "true" ]]; then
    cmd="${cmd} --check"
    log_message "INFO" "ℹ️ Running in check-only mode (no changes will be made)"
  fi

  if [[ "${diff_mode}" == "true" ]]; then
    cmd="${cmd} --diff"
    log_message "INFO" "ℹ️ Showing diffs between original and formatted files"
  fi

  # Build exclude patterns for grep
  local exclude_grep_pattern=""
  if [[ -n "${exclude_pattern}" ]]; then
    log_message "INFO" "ℹ️ Excluding directories: ${exclude_pattern}"
    IFS=',' read -ra exclude_dirs <<< "${exclude_pattern}"
    for dir in "${exclude_dirs[@]}"; do
      exclude_grep_pattern="${exclude_grep_pattern} -e ${dir}"
    done
  fi

  # Change to the terragrunt directory
  cd "${terragrunt_dir}" || {
    log_message "ERROR" "Failed to change to directory: ${terragrunt_dir}"
    return 1
  }

  # Find all HCL files, excluding specified directories
  local hcl_files
  if [[ -n "${exclude_grep_pattern}" ]]; then
    hcl_files=$(find . -name "*.hcl" | grep -v "${exclude_grep_pattern}" || true)
  else
    hcl_files=$(find . -name "*.hcl" || true)
  fi

  # Count total HCL files for reporting
  local total_files
  total_files=$(echo "${hcl_files}" | grep -c -v '^$')
  log_message "INFO" "📊 Found ${total_files} HCL files in ${terragrunt_dir}"

  # Exit early if no files found
  if [[ "${total_files}" -eq 0 ]]; then
    log_message "WARN" "No HCL files found to process"
    return 0
  fi

  # Process each file individually
  local formatted_count=0
  local failed_count=0
  local unchanged_count=0

  log_message "INFO" "🔄 Formatting HCL files..."
  while IFS= read -r file; do
    if [[ -z "${file}" ]]; then
      continue
    fi

    echo "  Processing: ${file}"
    if ${cmd} --file "${file}" 2>/dev/null; then
      if grep -q "was updated" <<< "$(terragrunt hclfmt --check --file "${file}" 2>&1)"; then
        formatted_count=$((formatted_count+1))
        echo "    ✅ File updated: ${file}"
      else
        unchanged_count=$((unchanged_count+1))
        echo "    ℹ️ Already formatted: ${file}"
      fi
    else
      failed_count=$((failed_count+1))
      echo "    ❌ Failed to format: ${file}"
    fi
  done <<< "${hcl_files}"

  # Show success message with stats
  echo ""
  log_message "INFO" "📊 Formatting Statistics:"
  echo "   - Total files processed: ${total_files}"
  if [[ "${check_mode}" != "true" ]]; then
    echo "   - Files updated: ${formatted_count}"
  fi
  echo "   - Files already formatted: ${unchanged_count}"
  echo "   - Files failed: ${failed_count}"

  if [[ "${check_mode}" == "true" ]]; then
    log_message "INFO" "✅ HCL format check completed"
  else
    log_message "INFO" "✅ HCL formatting completed"
  fi

  return 0
}

# =============================================================================
# Command-line interface when script is executed directly
# =============================================================================

# Execute terraform_format when script is called directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  # Default values
  terragrunt_dir="${1:-.}"
  check_mode="${2:-false}"
  diff_mode="${3:-false}"
  exclude_pattern="${4:-}"

  terragrunt_format "${terragrunt_dir}" "${check_mode}" "${diff_mode}" "${exclude_pattern}"
fi

# Shellcheck directives for sourcing and shell compatibility
# shellcheck shell=bash
# shellcheck disable=SC2155
true
