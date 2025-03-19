# 🏗️ Terragrunt Reference Architecture - Justfile
# This Justfile provides a streamlined interface for managing Terragrunt-based infrastructure
# Designed to simplify complex infrastructure workflows and provide consistent, reproducible deployments

# 📍 Path configurations
# Centralize path management to ensure consistent directory references across recipes
TERRAGRUNT_DIR := "./infra/terragrunt"

# 🐚 Shell configuration
# Use bash with strict error handling to prevent silent failures
# -u: Treat unset variables as an error
# -e: Exit immediately if a command exits with a non-zero status
set shell := ["bash", "-uce"]
set dotenv-load

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

# 🧹 Terragrunt format, run hclfmt on all Terragrunt files
# Example: `just tg-format check=true diff=true exclude=".terragrunt-cache,modules"`
tg-format check="false" diff="false" exclude="":
    #!/usr/bin/env bash
    echo "🔍 Advanced Terragrunt HCL Formatting"

    # Set up the command base
    CMD="terragrunt hclfmt"

    # Add options based on parameters
    if [[ "{{check}}" == "true" ]]; then
        CMD="$CMD --check"
        echo "ℹ️ Running in check-only mode (no changes will be made)"
    fi

    if [[ "{{diff}}" == "true" ]]; then
        CMD="$CMD --diff"
        echo "ℹ️ Showing diffs between original and formatted files"
    fi

    # Build exclude patterns for grep
    EXCLUDE_PATTERN=""
    if [[ -n "{{exclude}}" ]]; then
        echo "ℹ️ Excluding directories: {{exclude}}"
        IFS=',' read -ra EXCLUDE_DIRS <<< "{{exclude}}"
        for dir in "${EXCLUDE_DIRS[@]}"; do
            EXCLUDE_PATTERN="$EXCLUDE_PATTERN -e $dir"
        done
    fi

    # Find all HCL files, excluding specified directories
    cd {{TERRAGRUNT_DIR}}
    if [[ -n "$EXCLUDE_PATTERN" ]]; then
        HCL_FILES=$(find . -name "*.hcl" | grep -v $EXCLUDE_PATTERN)
    else
        HCL_FILES=$(find . -name "*.hcl")
    fi

    # Count total HCL files for reporting
    TOTAL_FILES=$(echo "$HCL_FILES" | wc -l | xargs)
    echo "📊 Found $TOTAL_FILES HCL files in {{TERRAGRUNT_DIR}}"

    # Process each file individually
    FORMATTED_COUNT=0
    FAILED_COUNT=0
    UNCHANGED_COUNT=0

    echo "🔄 Formatting HCL files..."
    while read -r file; do
        if [[ -z "$file" ]]; then
            continue
        fi

        echo "  Processing: $file"
        if $CMD --file "$file" 2>/dev/null; then
            if grep -q "was updated" <<< "$(terragrunt hclfmt --check --file "$file" 2>&1)"; then
                FORMATTED_COUNT=$((FORMATTED_COUNT+1))
                echo "    ✅ File updated: $file"
            else
                UNCHANGED_COUNT=$((UNCHANGED_COUNT+1))
                echo "    ℹ️ Already formatted: $file"
            fi
        else
            FAILED_COUNT=$((FAILED_COUNT+1))
            echo "    ❌ Failed to format: $file"
        fi
    done <<< "$HCL_FILES"

    # Show success message with stats
    echo ""
    echo "📊 Formatting Statistics:"
    echo "   - Total files processed: $TOTAL_FILES"
    if [[ "{{check}}" != "true" ]]; then
        echo "   - Files updated: $FORMATTED_COUNT"
    fi
    echo "   - Files already formatted: $UNCHANGED_COUNT"
    echo "   - Files failed: $FAILED_COUNT"

    if [[ "{{check}}" == "true" ]]; then
        echo "✅ HCL format check completed"
    else
        echo "✅ HCL formatting completed"
    fi

# 🚀 Run Terragrunt on a specific infrastructure unit
# Flexible recipe for running Terragrunt commands on individual units
# Parameters:
# - env: Environment (default: global)
# - stack: Infrastructure stack (default: dni)
# - unit: Specific infrastructure unit (default: dni_generator)
# - cmd: Terragrunt command (default: plan)
# Example: `just tg-run env=staging stack=network unit=vpc cmd=apply`
tg-run env="global" stack="dni" unit="dni_generator" cmd="init":
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

# 🔄 Reload direnv environment
# Manually reload the direnv environment when needed
reload-env:
    @echo "🔄 Manually reloading direnv environment..."
    @direnv reload

# 🧹 Clean direnv cache
# Removes the direnv cache to force a fresh environment build
# Useful when experiencing issues with the development environment
clean-direnv:
    @echo "🧹 Cleaning direnv cache..."
    @rm -rf .direnv
    @direnv allow
    @echo "✅ direnv cache cleaned. Environment will rebuild on next shell activation."
