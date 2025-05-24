#!/bin/sh
#
# Script to set up ssh-agent for GitLab CI jobs.
# Ensures SSH key is added and agent variables are available.
# Creates ssh_agent_vars.sh in CI_PROJECT_DIR for sourcing by subsequent steps.
#
# Adheres to principles from Google's Shell Style Guide.

# Exit on error (-e), exit on unset variable (-u).
set -eu

# --- Constants and Globals ---
# CI_PROJECT_DIR is a predefined GitLab CI variable.
# If this script is run outside CI, CI_PROJECT_DIR might need to be PWD or similar.
project_dir="${CI_PROJECT_DIR:-$(pwd)}" # Default to current dir if not in CI
ssh_agent_vars_file="${project_dir}/ssh_agent_vars.sh"
ssh_dir="${HOME}/.ssh"
ssh_key_file="${ssh_dir}/id_rsa"
known_hosts_file="${ssh_dir}/known_hosts"

# --- Helper Functions ---
# Standardized logging functions.
log_info() {
  echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') - setup_ssh_agent.sh: ${*}" >&2
}

log_error() {
  echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') - setup_ssh_agent.sh: ${*}" >&2
}

# --- Main Logic ---
main() {
  log_info "Starting ssh-agent setup..."

  # Validate essential GitLab CI variables
  if [ -z "${CI_PROJECT_DIR:-}" ]; then
    # This specific check is for when the script *relies* on CI_PROJECT_DIR.
    # Since we default project_dir above, this might be redundant if PWD is an acceptable fallback.
    # However, for GitLab CI context, CI_PROJECT_DIR is expected.
    log_info "CI_PROJECT_DIR is not set. Using current directory: ${project_dir}"
  fi

  if [ -z "${GITLAB_SSH_PRIVATE_KEY:-}" ]; then
    log_error "GITLAB_SSH_PRIVATE_KEY CI/CD variable is not set. Cannot proceed with SSH key setup."
    exit 1
  fi

  log_info "Ensuring openssh dependencies are installed..."
  if command -v apk >/dev/null 2>&1; then
    log_info "Alpine system detected. Installing openssh via apk."
    apk update >/dev/null
    apk add --no-cache openssh
  elif command -v apt-get >/dev/null 2>&1; then
    log_info "Debian/Ubuntu system detected. Installing openssh-client via apt-get."
    apt-get update -y >/dev/null
    apt-get install -y openssh-client
  elif command -v yum >/dev/null 2>&1; then
    log_info "CentOS/RHEL system detected. Installing openssh-clients via yum."
    yum install -y openssh-clients
  else
    log_info "Could not detect apk, apt-get, or yum. Assuming openssh/ssh-agent is already available or installed by other means."
  fi

  if ! command -v ssh-agent >/dev/null 2>&1; then
    log_error "'ssh-agent' command not found. Please ensure OpenSSH client tools are installed and in PATH."
    exit 1
  fi
  if ! command -v ssh-add >/dev/null 2>&1; then
    log_error "'ssh-add' command not found. Please ensure OpenSSH client tools are installed and in PATH."
    exit 1
  fi
  if ! command -v ssh-keyscan >/dev/null 2>&1; then
    log_error "'ssh-keyscan' command not found. Please ensure OpenSSH client tools are installed and in PATH."
    exit 1
  fi

  log_info "Starting ssh-agent and creating variable file at: ${ssh_agent_vars_file}"
  # The output of ssh-agent -s typically includes:
  # SSH_AUTH_SOCK=/tmp/ssh-XXXXXXXXXX/agent.XXXXX; export SSH_AUTH_SOCK;
  # SSH_AGENT_PID=XXXXX; export SSH_AGENT_PID;
  # echo Agent pid XXXXX;
  # Ensure we capture all of this.
  if ! ssh-agent -s > "${ssh_agent_vars_file}"; then
    log_error "ssh-agent -s command failed."
    exit 1
  fi

  if [ ! -f "${ssh_agent_vars_file}" ] || [ ! -s "${ssh_agent_vars_file}" ]; then
    log_error "Failed to create or populate ${ssh_agent_vars_file} from ssh-agent."
    exit 1
  fi
  log_info "${ssh_agent_vars_file} created. Contents:"
  cat "${ssh_agent_vars_file}" >&2 # Log to stderr for visibility

  log_info "Sourcing ${ssh_agent_vars_file} to set agent environment variables for this script..."
  # Use . (dot) for POSIX-compliant sourcing
  # shellcheck disable=SC1090
  . "${ssh_agent_vars_file}"

  log_info "Creating SSH directory: ${ssh_dir}"
  mkdir -p "${ssh_dir}"
  log_info "Writing private key to: ${ssh_key_file}"
  echo "${GITLAB_SSH_PRIVATE_KEY}" > "${ssh_key_file}"
  chmod 600 "${ssh_key_file}"

  log_info "Adding gitlab.com and github.com to known_hosts file: ${known_hosts_file}"
  # Ensure the known_hosts file exists with correct permissions before appending
  touch "${known_hosts_file}"
  chmod 644 "${known_hosts_file}"
  ssh-keyscan -H gitlab.com github.com >> "${known_hosts_file}"

  log_info "Configuring SSH to force use of the correct key for github.com (CI workaround)"
  ssh_config_file="${ssh_dir}/config"
  cat <<EOF > "${ssh_config_file}"
Host github.com
  HostName github.com
  User git
  IdentityFile ${ssh_key_file}
  IdentitiesOnly yes
  StrictHostKeyChecking no
EOF
  chmod 600 "${ssh_config_file}"

  log_info "Adding private key to ssh-agent..."
  if ! ssh-add "${ssh_key_file}"; then
    log_error "Failed to add SSH key to agent. Check key format and agent status."
    # Sensitive: Avoid logging key content in production
    # log_info "Key file (${ssh_key_file}) permissions: $(ls -l "${ssh_key_file}")"
    exit 1
  fi

  log_info "Listing loaded SSH keys in agent (for debug):"
  ssh-add -l || log_info "ssh-add -l failed (no keys loaded?)"

  log_info "Testing SSH authentication to github.com (ignore exit code, for debug only):"
  ssh -T git@github.com || log_info "SSH test to github.com failed (expected if key is not authorized)"

  log_info "Verifying SSH_AUTH_SOCK is set and socket exists..."
  if [ -z "${SSH_AUTH_SOCK:-}" ] || [ ! -S "${SSH_AUTH_SOCK}" ]; then
    log_error "SSH_AUTH_SOCK is not set correctly or the agent socket file does not exist."
    log_error "SSH_AUTH_SOCK current value: '${SSH_AUTH_SOCK:-}'"
    # List /tmp for potential ssh-agent sockets if debugging is desperate
    # ls -la /tmp | grep ssh
    exit 1
  fi

  log_info "SSH_AUTH_SOCK is set to: ${SSH_AUTH_SOCK}"
  log_info "ssh-agent setup completed successfully. ${ssh_agent_vars_file} is ready to be sourced by subsequent CI script steps."
}

# Script entry point
main "$@"
