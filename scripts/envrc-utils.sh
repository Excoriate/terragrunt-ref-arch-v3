#!/usr/bin/env bash
# =============================================================================
# Terragrunt Reference Architecture - Environment Utilities
# =============================================================================
# Common shell functions for .envrc files across the project
# This file contains reusable functions that can be sourced by any .envrc file

# =============================================================================
# CORE LOGGING FUNCTIONS
# =============================================================================

# Simple logging function with timestamp and log level
# Usage: _log "INFO" "Your message here"
_log() {
  local log_level="${1:-INFO}"
  local message="${2}"
  local timestamp
  timestamp="$(date +"%Y-%m-%d %H:%M:%S")"
  echo "[${log_level}] ${timestamp} - ${message}" >&2
}

# Layer-specific logging with prefix
# Usage: _layer_log "INFO" "Your message" "LAYER_NAME"
_layer_log() {
  local log_level="${1:-INFO}"
  local message="${2}"
  local prefix="${3:-}"

  if [[ -n "${prefix}" ]]; then
    _log "${log_level}" "[${prefix}] ${message}"
  else
    _log "${log_level}" "${message}"
  fi
}

# =============================================================================
# ENVIRONMENT VARIABLE MANAGEMENT
# =============================================================================

# Secure environment variable export with validation
# Usage: _safe_export "VARIABLE_NAME" "variable_value"
_safe_export() {
  local var_name="${1}"
  local var_value="${2}"
  local caller_info="${3:-unknown}"

  # Validate variable name
  if [[ ! "${var_name}" =~ ^[A-Z][A-Z0-9_]*$ ]]; then
    _log "ERROR" "Invalid environment variable name: ${var_name}. Must be uppercase and start with a letter."
    return 1
  fi

  # Check for empty values
  if [[ -z "${var_value// }" ]]; then
    _log "WARN" "Attempted to export empty or whitespace-only value for '${var_name}'"
    return 1
  fi

  # Sanitize and export
  var_value="$(echo "${var_value}" | xargs)"
  export "${var_name}"="${var_value}"
  _log "TRACK" "${var_name} = [REDACTED] (set from ${caller_info})"
}

# Layer-specific variable export with custom logging
# Usage: _layer_export "VARIABLE_NAME" "variable_value" "LAYER_NAME"
_layer_export() {
  local var_name="${1}"
  local var_value="${2}"
  local layer_name="${3:-}"

  # Use the existing _safe_export with caller information
  _safe_export "${var_name}" "${var_value}" "$(caller | awk '{print $2}')"

  if [[ -n "${layer_name}" ]]; then
    _layer_log "TRACK" "${layer_name} config: ${var_name} = [REDACTED]" "${layer_name}"
  fi
}

# =============================================================================
# PROJECT STRUCTURE DETECTION
# =============================================================================

# Project root detection
# Sets PROJECT_ROOT environment variable
_detect_project_root() {
  local current_dir="${PWD}"
  local root_markers=(".git" "justfile")

  while [ "${current_dir}" != "/" ]; do
    for marker in "${root_markers[@]}"; do
      if [ -e "${current_dir}/${marker}" ]; then
        _safe_export PROJECT_ROOT "${current_dir}"
        _log "INFO" "Project root detected: ${current_dir}"
        return 0
      fi
    done
    current_dir="$(dirname "${current_dir}")"
  done

  _safe_export PROJECT_ROOT "${PWD}"
  _log "WARN" "No specific project root marker found. Using current directory."
}

# Terragrunt root detection
# Sets TERRAGRUNT_ROOT environment variable
_detect_terragrunt_root() {
  local current_dir="${PWD}"
  local terragrunt_markers=(
    "terragrunt.hcl"
    "root.hcl"
    "config.hcl"
  )

  while [ "${current_dir}" != "/" ]; do
    for marker in "${terragrunt_markers[@]}"; do
      if [ -e "${current_dir}/${marker}" ]; then
        _layer_export TERRAGRUNT_ROOT "${current_dir}" "TERRAGRUNT"
        _layer_log "INFO" "Terragrunt root detected: ${current_dir}" "TERRAGRUNT"
        return 0
      fi
    done
    current_dir="$(dirname "${current_dir}")"
  done

  _layer_export TERRAGRUNT_ROOT "${PWD}" "TERRAGRUNT"
  _layer_log "WARN" "No specific Terragrunt root marker found. Using current directory." "TERRAGRUNT"
}

# =============================================================================
# LAYER CONFIGURATION FUNCTIONS
# =============================================================================

# Validate configuration for a specific layer
# Usage: _validate_layer_config "LAYER_NAME" "VAR1" "VAR2" ...
_validate_layer_config() {
  local layer_name="${1}"
  shift
  local required_vars=("$@")
  local missing_vars=()

  for var in "${required_vars[@]}"; do
    # Use indirect expansion safely
    if [[ -z "${!var:-}" ]]; then
      missing_vars+=("${var}")
    fi
  done

  if [[ ${#missing_vars[@]} -gt 0 ]]; then
    _layer_log "ERROR" "Missing critical configurations:" "${layer_name}"
    printf '%s\n' "${missing_vars[@]}"
    return 1
  fi

  _layer_log "INFO" "Configuration validated successfully" "${layer_name}"
}

# Display layer environment information
# Usage: _layer_env_info "LAYER_NAME" "VAR1:Description" "VAR2:Description" ...
_layer_env_info() {
  local layer_name="${1}"
  shift
  local var_descriptions=("$@")

  _layer_log "INFO" "🔹 ${layer_name} Layer Environment Initialized" "${layer_name}"

  for var_desc in "${var_descriptions[@]}"; do
    IFS=':' read -r var_name var_description <<< "${var_desc}"
    # Use indirect expansion safely
    if [[ -n "${!var_name:-}" ]]; then
      local display_value="${!var_name}"
      # Mask sensitive values if needed
      if [[ "${var_name}" == *"SECRET"* || "${var_name}" == *"PASSWORD"* || "${var_name}" == *"KEY"* ]]; then
        display_value="[REDACTED]"
      fi
      _layer_log "INFO" "${var_description:-${var_name}}: ${display_value}" "${layer_name}"
    fi
  done
}

# =============================================================================
# HELPER FUNCTIONS FOR ENVIRONMENT INITIALIZATION
# =============================================================================

# Initialize core environment
_core_init() {
  # Detect project root
  _detect_project_root

  # Set up basic environment variables
  _log "INFO" "Initializing core environment"
}

# Initialize Terragrunt environment
_terragrunt_init() {
  # Detect Terragrunt root
  _detect_terragrunt_root

  _layer_log "INFO" "Initializing Terragrunt environment" "TERRAGRUNT"
}

# =============================================================================
# USAGE EXAMPLES
# =============================================================================
#
# In root .envrc:
#   source "${PWD}/scripts/envrc-utils.sh"
#   _core_init
#   _safe_export VARIABLE_NAME "value"
#
# In any layer .envrc:
#   source_up || true
#   source "${PROJECT_ROOT}/scripts/envrc-utils.sh"
#
#   # Define layer name (can be any name that makes sense for your structure)
#   LAYER_NAME="MY_LAYER"
#
#   # Initialize layer variables
#   _layer_export VAR_NAME "value" "$LAYER_NAME"
#
#   # Display layer information
#   _layer_env_info "$LAYER_NAME" \
#     "VAR_NAME:Description of VAR_NAME" \
#     "ANOTHER_VAR:Description of ANOTHER_VAR"

# Display all exported variables with optional filtering
# Usage: _display_exported_vars [FILTER_PREFIX]
_display_exported_vars() {
  local filter_prefix="${1:-}"
  local exported_vars=()

  # Collect exported variables
  while IFS='=' read -r var value; do
    # Only include variables that match the optional prefix
    if [[ -z "${filter_prefix}" ]] || [[ "${var}" == "${filter_prefix}"* ]]; then
      # Mask sensitive values if needed
      if [[ "${var}" == *"SECRET"* || "${var}" == *"PASSWORD"* || "${var}" == *"KEY"* || "${var}" == *"TOKEN"* ]]; then
        value="[REDACTED]"
      fi
      exported_vars+=("${var}: ${value}")
    fi
  done < <(env | grep -E '^[A-Z_]+=' | sort)

  # Display header and variables
  if [[ ${#exported_vars[@]} -gt 0 ]]; then
    _log "INFO" "Exported Environment Variables:"
    printf '%s\n' "${exported_vars[@]}" >&2
  else
    _log "WARN" "No exported variables found$([ -n "${filter_prefix}" ] && echo " with prefix '${filter_prefix}'")"
  fi
}

# Shellcheck directives for sourcing and shell compatibility
# shellcheck shell=bash
# shellcheck disable=SC2155
true
