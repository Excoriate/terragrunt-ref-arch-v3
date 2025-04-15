# 🏗️ Terragrunt Reference Architecture - Justfile
# This Justfile provides a streamlined interface for managing Terragrunt-based infrastructure
# Designed to simplify complex infrastructure workflows and provide consistent, reproducible deployments

# 📍 Path configurations
# Centralize path management to ensure consistent directory references across recipes
TERRAGRUNT_DIR := "./infra/terragrunt"
TERRAFORM_MODULES_DIR := "./infra/terraform/modules"

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
    @./scripts/justfile-utils.sh "{{TERRAGRUNT_DIR}}" "{{check}}" "{{diff}}" "{{exclude}}"

# ✅ Terragrunt validate, run hclvalidate on all Terragrunt files
# Example: `just tg-hclvalidate`
tg-hclvalidate:
    @echo "✅ Running Terragrunt HCL validation via utility script"
    @./scripts/justfile-utils.sh terragrunt_hclvalidate "{{TERRAGRUNT_DIR}}"

tg_env := "global"
tg_stack := "dni"
tg_unit := "dni_generator"

# 🚀 Run Terragrunt CI checks (hclvalidate and format)
tg-ci: (tg-hclvalidate) (tg-format)

# 🚀 Run Terragrunt on a specific infrastructure unit
# Flexible recipe for running Terragrunt commands on individual units
# Example: `just tg-run cmd=init`
[working-directory:'infra/terragrunt']
tg-run cmd="init":
    @cd {{tg_env}}/{{tg_stack}}/{{tg_unit}} && terragrunt {{cmd}}

# 🌐 Run Terragrunt plan across all units in a stack
# Provides a comprehensive view of potential infrastructure changes
# Useful for pre-deployment validation and impact assessment
[working-directory:'infra/terragrunt']
tg-run-all-plan :
    @cd {{tg_env}}/{{tg_stack}} && terragrunt run-all plan

# 🚀 Apply infrastructure changes across all units in a stack
# Automated, non-interactive deployment of infrastructure
# Includes auto-approval to streamline deployment processes
[working-directory:'infra/terragrunt']
tg-run-all-apply :
    @cd {{tg_env}}/{{tg_stack}} && terragrunt run-all apply --auto-approve --terragrunt-non-interactive

# 💥 Destroy infrastructure across all units in a stack
# Provides a safe, controlled method for infrastructure teardown
# Non-interactive with auto-approval for scripting and automation
tg-run-all-destroy:
    @cd infra/terragrunt/{{tg_env}}/{{tg_stack}} && terragrunt run-all destroy --terragrunt-non-interactive --auto-approve



# 🔍 Open Dagger CI terminal. E.g.: just ci-terminal --help
[working-directory:'ci/ci-terragrunt']
ci-terminal args="":
    @echo "🔍 Open Dagger CI terminal"
    @echo "🔍 Building the dagger module"
    @dagger develop
    @echo "🔍 Inspecting the available functions"
    @dagger functions
    @echo "🔍 Running the function"
    @dagger call open-terminal {{args}}

# 🔍 Run Dagger CI function
[working-directory:'ci/ci-terragrunt']
ci-shell:
    @echo "🔍 Running Dagger CI for terragrunt"
    @echo "🔍 Building the dagger module"
    @dagger develop
    @echo "🔍 Inspecting the available functions"
    @dagger functions
    @echo "🔍 Running the function"
    @dagger

# aws_access_key_id := env("AWS_ACCESS_KEY_ID")
# aws_secret_access_key := env("AWS_SECRET_ACCESS_KEY")

# 🔍 Run Dagger CI function
[working-directory:'ci/ci-terragrunt']
ci-job-units-static-check env="global" layer="dni" unit="dni_generator":
    @echo "🔍 Building the dagger module"
    @dagger develop
    @echo "🔍 Inspecting the available functions"
    @dagger functions
    @echo "🔍 Running the function"
    @dagger call job-terragrunt-units-static-check \
      --load-dot-env-file \
      --no-cache \
      --aws-access-key-id env://AWS_ACCESS_KEY_ID \
      --aws-secret-access-key env://AWS_SECRET_ACCESS_KEY

# 🔍 Run Dagger CI function
[working-directory:'ci/ci-terragrunt']
ci-job-units-plan env="global" layer="dni" unit="dni_generator":
    @echo "🔍 Building the dagger module"
    @dagger develop
    @echo "🔍 Inspecting the available functions"
    @dagger functions
    @echo "🔍 Running the function"
    @dagger call job-terragrunt-units-plan \
      --load-dot-env-file \
      --no-cache \
      --aws-access-key-id env://AWS_ACCESS_KEY_ID \
      --aws-secret-access-key env://AWS_SECRET_ACCESS_KEY

[working-directory:'ci/ci-terragrunt']
ci-job-tfmodules-static-check:
    @echo "🔍 Building the dagger module"
    @dagger develop
    @echo "🔍 Inspecting the available functions"
    @dagger functions
    @echo "🔍 Running the function"
    @dagger call job-terraform-modules-static-check

dev:
    @echo "🌿 Starting Nix Development Shell for Terraform Registry Module Template 🏷️"
    @nix develop . --impure --extra-experimental-features nix-command --extra-experimental-features flakes
