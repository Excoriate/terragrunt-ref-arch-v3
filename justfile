# 🏗️ Terragrunt Reference Architecture - Justfile
# This Justfile provides a streamlined interface for managing Terragrunt-based infrastructure
# Designed to simplify complex infrastructure workflows and provide consistent, reproducible deployments

# 📍 Path configurations
# Centralize path management to ensure consistent directory references across recipes
TERRAGRUNT_DIR := "./infra/terragrunt"
TERRAFORM_MODULES_DIR := "./infra/terraform/modules"

# 📍 Default values for Terragrunt environment, stack, and unit
# TODO: Change it accordingly to the environment, stack, and unit you are working on
# 📍 Terragrunt default(s)
TG_DEFAULT_STACK := "non-distributable"
TG_DEFAULT_UNIT := "random-string-generator"

# 🐚 Shell configuration
# Use bash with strict error handling to prevent silent failures
# -u: Treat unset variables as an error
# -e: Exit immediately if a command exits with a non-zero status
set shell := ["bash", "-uce"]
set dotenv-load

# Avoid reporting traces to Dagger
# TODO: Uncomment if traces aren't needed.
# export NOTHANKS := "1"

# 📋 Default recipe: List all available commands
# Provides a quick overview of available infrastructure management commands
default:
    @just --list

# 🌿 Start Nix Development Shell
dev:
    @echo "🌿 Starting Nix Development Shell for Terragrunt Reference Architecture 🚀"
    @nix develop . --impure --extra-experimental-features nix-command --extra-experimental-features flakes

# 🗑️ Clean macOS system files
# Removes .DS_Store files that can cause unnecessary version control noise
# Helps maintain a clean repository across different operating systems
clean-ds:
    @echo "🧹 Cleaning .DS_Store files"
    @find . -name '.DS_Store' -type f -delete

# 🔧 Install pre-commit hooks in local environment for code consistency
hooks-install:
    @echo "🧰 Installing pre-commit hooks locally..."
    @./scripts/hooks/pre-commit-init.sh init

# 🕵️ Run pre-commit hooks across all files in local environment
hooks-run:
    @echo "🔍 Running pre-commit hooks from .pre-commit-config.yaml..."
    @./scripts/hooks/pre-commit-init.sh run

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

# 🔍 Run Terraform command for a specific module
[working-directory:'infra/terraform/modules']
tf-run module='random-string-generator' cmd='init' args='':
    @echo "🔍 Preparing to run Terraform command..."
    @echo "📂 Module Path: {{module}}"
    @echo "⚙️ Command: {{cmd}}"
    @echo "📋 Arguments: {{args}}"
    @cd {{module}} && terraform {{cmd}} {{args}}
    @echo "✅ Terraform {{cmd}} executed successfully for module: {{module}}"

# 🌿 Format all Terraform files across modules, examples, and tests directories
tf-format-all:
    @echo "🌿 Formatting all Terraform files across the repository..."
    @echo "📂 Scanning directories: {{TERRAFORM_MODULES_DIR}}/"

    @echo "\n🔍 Formatting files in modules/"
    @pushd {{TERRAFORM_MODULES_DIR}} > /dev/null && \
    find . -type f \( -name "*.tf" -o -name "*.tfvars" \) | sort | while read -r file; do \
        echo "   📄 Processing: $file"; \
    done && \
    terraform fmt -recursive && \
    popd > /dev/null

    @echo "\n✅ All Terraform files have been formatted!"

# 🧹 Clean Terraform cache for all modules
[working-directory:'infra/terraform/modules']
tf-clean-all:
    @echo "🧹 Cleaning Terraform cache for all modules"
    @if [ -n "$(find . -maxdepth 4 -type d -name ".terraform" 2>/dev/null)" ]; then \
        echo "🔍 Found .terraform directories to clean"; \
        find . -maxdepth 4 -type d -name ".terraform" -exec rm -rf {} +; \
        echo "✅ Removed .terraform directories"; \
    else \
        echo "ℹ️ No .terraform directories found"; \
    fi
    @if [ -n "$(find . -maxdepth 4 -type f -name ".terraform.lock.hcl" 2>/dev/null)" ]; then \
        echo "🔍 Found .terraform.lock.hcl files to clean"; \
        find . -maxdepth 4 -type f -name ".terraform.lock.hcl" -exec rm -rf {} +; \
        echo "✅ Removed .terraform.lock.hcl files"; \
    else \
        echo "ℹ️ No .terraform.lock.hcl files found"; \
    fi
    @echo "🧹 Cleaning completed"

# 🧹 Terragrunt and Terraform cache cleanup
[working-directory:'infra/terragrunt']
tg-clean-all:
    @echo "🧹 Cleaning Terragrunt cache for all environments and .terraform directories"
    @find . -maxdepth 4 -type d \( -name ".terragrunt-cache" -o -name ".terraform" \) -exec rm -rf {} +
    @find . -maxdepth 4 -type f -name ".terraform.lock.hcl" -exec rm -rf {} +
    @find . -maxdepth 4 -type f -name ".terraform.lock.hcl" -exec rm -rf {} +

# 🧹 Terragrunt and Terraform cache cleanup for a specific path
[working-directory:'infra/terragrunt']
tg-clean tgpath:
    @echo "🧹 Cleaning Terragrunt cache for specific path: {{tgpath}}"
    @if [ -d {{tgpath}} ]; then \
        cd {{tgpath}} && \
        find . -maxdepth 4 -type d \( -name ".terragrunt-cache" -o -name ".terraform" \) -exec rm -rf {} + && \
        find . -maxdepth 4 -type f -name ".terraform.lock.hcl" -exec rm -rf {} +; \
    else \
        echo "❌ Directory {{tgpath}} does not exist."; \
    fi

# 🧹 Terragrunt format, run hclfmt on all Terragrunt files
# Example: `just tg-format check=true diff=true exclude=".terragrunt-cache,modules"`
tg-format check="false" diff="false" exclude="":
    @echo "🔍 Running Terragrunt HCL formatting via utility script"
    @./scripts/justfile-utils.sh terragrunt_format "{{TERRAGRUNT_DIR}}" "{{check}}" "{{diff}}" "{{exclude}}"

# ✅ Terragrunt validate, run hclvalidate on all Terragrunt files
# Example: `just tg-hclvalidate`
tg-hclvalidate:
    @echo "✅ Running Terragrunt HCL validation via utility script"
    @./scripts/justfile-utils.sh terragrunt_hclvalidate "{{TERRAGRUNT_DIR}}"

# 🚀 Run Terragrunt CI checks (hclvalidate and format)
tg-ci: (tg-hclvalidate) (tg-format)

# 🚀 Run Terragrunt on a specific infrastructure unit
# Example: `just tg-run env=dev stack=non-distributable unit=random-string-generator cmd=init`
[working-directory:'infra/terragrunt']
tg-run env stack unit cmd="init":
    @cd {{env}}/{{stack}}/{{unit}} && \
    if [ "{{cmd}}" = "apply" ] || [ "{{cmd}}" = "destroy" ]; then \
        echo "🔒 Running Terragrunt {{cmd}} with auto-approve flag" && \
        terragrunt {{cmd}} --auto-approve; \
    else \
        echo "🔍 Running Terragrunt {{cmd}}" && \
        terragrunt {{cmd}}; \
    fi

# 🚀 Run Terragrunt on all infrastructure units in a stack
# Example: `just tg-run-all env=dev stack=non-distributable unit=random-string-generator cmd=init`
[working-directory:'infra/terragrunt']
tg-run-all env stack cmd="init":
    @cd {{env}}/{{stack}} && terragrunt run-all {{cmd}} $(if [ "{{cmd}}" = "apply" ] || [ "{{cmd}}" = "destroy" ]; then echo "--auto-approve"; fi)


# 🔨 Build the Dagger pipeline
[working-directory:'pipeline/infra']
pipeline-infra-build:
    @echo "🔨 Initializing Dagger development environment"
    @dagger develop
    @echo "📋 Building Dagger pipeline"
    @dagger functions

# 🔨 Help for Dagger job
[working-directory:'pipeline/infra']
pipeline-job-help fn: (pipeline-infra-build)
    @echo "🔨 Help for Dagger job: {{fn}}"
    @dagger call {{fn}} --help

# 🔨 Open an interactive development shell for the Infra pipeline
[working-directory:'pipeline/infra']
pipeline-infra-shell args="": (pipeline-infra-build)
    @echo "🚀 Launching interactive terminal"
    @dagger call open-terminal {{args}}

# 🔨 Validate Terraform modules for best practices and security
[working-directory:'pipeline/infra']
pipeline-infra-tf-modules-static-check args="": (pipeline-infra-build)
    @echo " Analyzing Terraform modules for security and best practices"
    @echo "⚡ Running static analysis checks"
    @dagger call job-tf-modules-static-check {{args}}
    @echo "✅ Static analysis completed successfully"

# 🔨 Test Terraform modules against multiple provider versions
[working-directory:'pipeline/infra']
pipeline-infra-tf-modules-versions args="": (pipeline-infra-build)
    @echo " Testing module compatibility across provider versions"
    @echo "⚡ Running version compatibility checks"
    @dagger call job-tf-modules-compatibility-check {{args}}
    @echo "✅ Version compatibility testing completed"

# 🔨 Run comprehensive CI checks for Terraform modules
pipeline-infra-tf-ci args="": (pipeline-infra-tf-modules-static-check) (pipeline-infra-tf-modules-versions)

# 🔨 Run a Terragrunt job with custom arguments
[working-directory:'pipeline/infra']
pipeline-infra-tg-exec args="": (pipeline-infra-build)
    @echo "🔨 Running Terragrunt single command through Dagger"
    @dagger call job-tg-exec {{args}}

# 🔨 Run a Terragrunt job with custom arguments
[working-directory:'pipeline/infra']
pipeline-infra-tg-ci-static env="global" stack="dni": (pipeline-infra-build)
    @echo "🔄 Running Terragrunt CI checks through Dagger"
    @echo "🌍 Environment: {{env}} | 📚 Stack: {{stack}}"
    @dagger call job-citg-stack-static-analysis \
        --aws-region eu-west-1 \
        --aws-access-key-id env:AWS_ACCESS_KEY_ID \
        --aws-secret-access-key env:AWS_SECRET_ACCESS_KEY \
        --load-dot-env \
        --no-cache \
        --environment "{{env}}" \
        --stack "{{stack}}"

    @echo "✅ Terragrunt CI checks completed successfully on environment: {{env}} | 📚 Stack: {{stack}}"

# 🔨 Run a Terragrunt job on stacks, with custom arguments
[working-directory:'pipeline/infra']
pipeline-infra-tg-stack env="dev" stack="non-distributable" tg-cmd="validate" tg-cmd-args="--terragrunt-ignore-external-dependencies": (pipeline-infra-build)
    @echo "🔄 Running Terragrunt stack through Dagger"
    @echo "🌍 Environment: {{env}} | 📚 Stack: {{stack}}"
    @echo "🔨 Terragrunt command: {{tg-cmd}}"
    @echo "🔨 Terragrunt command arguments: {{tg-cmd-args}}"
    @dagger call job-tg-stack \
        --aws-region eu-west-1 \
        --aws-access-key-id env:AWS_ACCESS_KEY_ID \
        --aws-secret-access-key env:AWS_SECRET_ACCESS_KEY \
        --load-dot-env \
        --no-cache \
        --environment "{{env}}" \
        --stack "{{stack}}" \
        --git-ssh $SSH_AUTH_SOCK \
        --tg-cmd "{{tg-cmd}}" \
        --tg-cmd-args "{{tg-cmd-args}}"

    @echo "✅ Terragrunt stack completed successfully on environment: {{env}} | 📚 Stack: {{stack}}"
