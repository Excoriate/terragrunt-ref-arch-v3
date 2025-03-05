# 🏗️ Terragrunt Reference Architecture - Justfile
# This Justfile provides a streamlined interface for managing Terragrunt-based infrastructure
# Designed to simplify complex infrastructure workflows and provide consistent, reproducible deployments

# 📍 Path configurations
# Centralize path management to ensure consistent directory references across recipes
TERRAGRUNT_DIR := "./infra/terragrunt"

# 🌍 Load environment variables
# Automatically load environment-specific configurations from .env file
# Enables flexible, environment-aware infrastructure management
set dotenv-load

# 🐚 Shell configuration
# Use bash with strict error handling to prevent silent failures
# -u: Treat unset variables as an error
# -e: Exit immediately if a command exits with a non-zero status
set shell := ["bash", "-uce"]

# 📋 Default recipe: List all available commands
# Provides a quick overview of available infrastructure management commands
default:
    @just --list

# 🗑️ Clean macOS system files
# Removes .DS_Store files that can cause unnecessary version control noise
# Helps maintain a clean repository across different operating systems
clean-ds:
    @echo "🧹 Cleaning .DS_Store files"
    @find . -name '.DS_Store' -type f -delete

# 🧹 Terragrunt and Terraform cache cleanup
# Removes cached Terragrunt and Terraform directories to ensure clean state
# Useful for troubleshooting and preventing stale cache-related issues
tg-clean:
    @echo "🧹 Cleaning Terragrunt cache for all environments and .terraform directories"
    @cd infra/terragrunt && find . -type d -name ".terragrunt-cache" -exec rm -rf {} +
    @cd infra/terragrunt && find . -type d -name ".terraform" -exec rm -rf {} +

# 🚀 Run Terragrunt on a specific infrastructure unit
# Flexible recipe for running Terragrunt commands on individual units
# Parameters:
# - env: Environment (default: global)
# - stack: Infrastructure stack (default: dni)
# - unit: Specific infrastructure unit (default: dni_generator)
# - cmd: Terragrunt command (default: plan)
# Example: `just tg-run env=staging stack=network unit=vpc cmd=apply`
tg-run env="global" stack="dni" unit="dni_generator" cmd="plan":
    @cd infra/terragrunt/{{env}}/{{stack}}/{{unit}} && terragrunt {{cmd}}

# 🌐 Run Terragrunt plan across all units in a stack
# Provides a comprehensive view of potential infrastructure changes
# Useful for pre-deployment validation and impact assessment
tg-run-all-plan env="global" stack="dni":
    @cd infra/terragrunt/{{env}}/{{stack}} && terragrunt run-all plan

# 🚀 Apply infrastructure changes across all units in a stack
# Automated, non-interactive deployment of infrastructure
# Includes auto-approval to streamline deployment processes
tg-run-all-apply env="global" stack="dni":
    @cd infra/terragrunt/{{env}}/{{stack}} && terragrunt run-all apply --auto-approve --terragrunt-non-interactive

# 💥 Destroy infrastructure across all units in a stack
# Provides a safe, controlled method for infrastructure teardown
# Non-interactive with auto-approval for scripting and automation
tg-run-all-destroy env="global" stack="dni":
    @cd infra/terragrunt/{{env}}/{{stack}} && terragrunt run-all destroy --terragrunt-non-interactive --auto-approve

# 🛠️ Allow direnv to run
# Ensures that direnv is allowed to run in the current directory
# Useful for managing environment variables and configurations
allow-direnv:
    @echo "🔒 Allow direnv to run..."
    @direnv allow

# 🛠️ Enter Nix development shell
# Provides a consistent development environment with all required tools
# Uses flake.nix to ensure reproducible tool versions across all developers
dev: allow-direnv
    @echo "🔍 Verifying Git and Nix configuration..."
    @if [ ! -f "flake.nix" ]; then echo "❌ flake.nix not found!"; exit 1; fi
    @if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then echo "❌ Not in a Git repository!"; exit 1; fi
    @if [ -n "$(git status --porcelain)" ]; then \
        echo "⚠️  Git repository has uncommitted changes. This may cause issues with Nix flakes."; \
        echo "Consider committing or stashing changes first."; \
        read -p "Continue anyway? [y/N] " -n 1 -r; \
        echo; \
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then \
            echo "Aborting."; \
            exit 1; \
        fi \
    fi
    @echo "🚀 Entering Nix development shell..."
    @NIXPKGS_ALLOW_UNFREE=1 nix develop --impure --extra-experimental-features "nix-command flakes"



