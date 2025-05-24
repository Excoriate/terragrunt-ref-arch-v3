# Environment Variables Management

## Overview

The Terragrunt Reference Architecture utilizes `.env` files for managing environment variables. This approach provides a straightforward way to configure settings for different environments. Typically, you will copy the provided `.env.example` file to a new file named `.env` (which should be excluded from version control via `.gitignore`) and then customize the variables within your `.env` file as needed.

The `justfile` in this project is configured with `set dotenv-load`, which means that `just` commands will automatically load variables from a `.env` file found in the project root.

## Key Features

- **Simple Key-Value Pairs**: `.env` files use a simple `KEY=value` format.
- **Easy Customization**: Users can easily override or set variables by creating and editing a `.env` file.
- **Gitignored by Default**: User-specific `.env` files are typically gitignored to keep sensitive data or local preferences out of version control.
- **Automatic Loading**: `just` recipes automatically source variables from `.env` files.

## Configuration Structure

The primary configuration for environment variables involves:

```
/
├── .env.example          # Example file showing available variables and default/template values
└── .env                  # User-specific, gitignored file where actual variable values are set
```

## Configuration Principles

Environment variables are typically defined in a `.env` file in the project root. This file is loaded by `just` when recipes are executed.

- **`.env.example`**: Serves as a template and documentation for required or optional environment variables.
- **`.env`**: This is your local, untracked file where you set the actual values for your development or deployment environment.
- **Loading**: Variables are made available to the execution context of `just` recipes.

## Variable Customization

To customize environment variables for your local setup or specific deployments:

1.  **Copy `.env.example`**: If a `.env` file does not already exist in the project root, copy the `.env.example` file to a new file named `.env`.
    ```sh
    cp .env.example .env
    ```
2.  **Edit `.env`**: Open the `.env` file and modify the variable values as needed for your environment. For example:
    ```env
    # .env (example content)
    DEFAULT_REGION=us-west-2
    TG_STACK_APP_PRODUCT_NAME=my-cool-app
    # Add any other custom variables required
    MY_CUSTOM_API_KEY=xxxxxxxxxxxx
    ```
3.  **Ensure `.env` is Gitignored**: The `.env` file typically contains sensitive information or local-specific settings and should not be committed to version control. Make sure `.env` is listed in your project's `.gitignore` file.

Variables defined in the `.env` file will be automatically loaded by `just` when you run recipes defined in the `justfile`.

## Comprehensive Environment Variable List

This section provides a comprehensive list of environment variables used throughout the Terragrunt Reference Architecture, including those commonly defined in `.env` files (based on `.env.example`) and those utilized within HCL configurations.

| Category                    | Variable Name                                       | Description                                                                 | Default Value (`.env.example` or HCL)     | Primary Source/Usage                                  |
|-----------------------------|-----------------------------------------------------|-----------------------------------------------------------------------------|-------------------------------------------|-------------------------------------------------------|
| **Project & Metadata**      | `PROJECT_ROOT`                                      | Absolute path to the project root directory.                                | Dynamically set by scripts if needed      | Shell environment, Scripts                            |
|                             | `TG_STACK_APP_AUTHOR`                               | Author of the configuration.                                                | `Your Name` (from `.env.example`)         | `.env`, `infra/terragrunt/_shared/_config/tags.hcl` |
|                             | `TG_STACK_APP_PRODUCT_NAME`                         | Project/application name for identification.                                | `your-app-name` (from `.env.example`)     | `.env`, `infra/terragrunt/_shared/_config/app.hcl`, `infra/terragrunt/_shared/_config/tags.hcl` |
| **Cloud Provider & Region** | `DEFAULT_REGION`                                    | Default cloud provider region.                                              | `us-east-1` (from `.env.example`)         | `.env`                                                |
|                             | `TG_STACK_DEPLOYMENT_REGION`                        | AWS region for deployments, overrides `DEFAULT_REGION` for Terragrunt.    | `us-east-1` (from HCL)                    | `infra/terragrunt/config.hcl`                         |
| **Terraform & Terragrunt**  |                                                     |                                                                             |                                           |                                                       |
| *General*                   | `TG_NON_INTERACTIVE`                                | Replaces `TF_INPUT`. If true, disables interactive prompts.                 | `true` (from `.env.example`)              | `.env` (Terragrunt Standard)                        |
|                             | `TG_STACK_TF_VERSION`                               | Enforced Terraform version.                                                 | `1.9.0` (from `.env.example`), `1.11.3` (HCL) | `.env`, `infra/terragrunt/config.hcl`             |
| *Terragrunt Performance*    | `TERRAGRUNT_DOWNLOAD_DIR`                           | Directory where Terragrunt downloads remote sources.                        | `${HOME}/.terragrunt-cache/...` (from `.env.example`) | `.env` (Terragrunt Standard)                        |
|                             | `TERRAGRUNT_CACHE_MAX_AGE`                          | Max age for items in Terragrunt cache.                                      | `168h` (from `.env.example`)              | `.env` (Terragrunt Standard)                        |
| *Terragrunt Behavior*       | `TG_LOG_LEVEL`                                      | Terragrunt logging verbosity. Replaces `TERRAGRUNT_LOG_LEVEL`.            | `info` (from `.env.example`)              | `.env` (Terragrunt Standard)                        |
|                             | `TERRAGRUNT_DISABLE_CONSOLE_OUTPUT`                 | If true, disables Terragrunt console output.                                | `false` (from `.env.example`)             | `.env` (Terragrunt Standard)                        |
|                             | `TG_NO_AUTO_INIT`                                   | Replaces `TERRAGRUNT_AUTO_INIT` (inverted). If true, disables auto-init.   | `false` (meaning auto-init is enabled, from `.env.example`) | `.env` (Terragrunt Standard)                        |
|                             | `TG_NO_AUTO_RETRY`                                  | Replaces `TERRAGRUNT_AUTO_RETRY` (inverted). If true, disables auto-retry. | `false` (meaning auto-retry is enabled, from `.env.example`) | `.env` (Terragrunt Standard)                        |
| *Terragrunt Flags (HCL)*    | `TG_STACK_FLAG_ENABLE_PROVIDERS_OVERRIDE`           | Controls dynamic provider file generation.                                  | `true` (from HCL)                         | `infra/terragrunt/config.hcl`                         |
|                             | `TG_STACK_FLAG_ENABLE_VERSIONS_OVERRIDE`            | Controls dynamic version file generation.                                   | `true` (from HCL)                         | `infra/terragrunt/config.hcl`                         |
|                             | `TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE` | Controls `.terraform-version` file generation.                             | `false` (from HCL)                        | `infra/terragrunt/config.hcl`                         |
| **Logging & Diagnostics**   | `LOG_LEVEL`                                         | General logging level for scripts.                                          | `info` (from `.env.example`)              | `.env`                                                |
|                             | `LOG_DIR`                                           | Directory for storing logs.                                                 | `${HOME}/.logs/...` (derived, from `.env.example`) | `.env`                                                |
| **Remote State Config**     | `TG_STACK_REMOTE_STATE_BUCKET_NAME`                 | S3 bucket for Terraform remote state.                                       | `terraform-state-makemyinfra` (from `.env.example`) | `.env`, `infra/terragrunt/_shared/_config/remote_state.hcl` |
|                             | `TG_STACK_REMOTE_STATE_LOCK_TABLE`                  | DynamoDB table for state locking.                                           | `terraform-state-makemyinfra` (from `.env.example`) | `.env`, `infra/terragrunt/_shared/_config/remote_state.hcl` |
|                             | `TG_STACK_REMOTE_STATE_REGION`                      | AWS region for remote state storage.                                        | `us-east-1` (from HCL)                    | `infra/terragrunt/_shared/_config/remote_state.hcl`   |
|                             | `TG_STACK_REMOTE_STATE_OBJECT_BASENAME`             | Basename for the Terraform state object file.                               | `terraform.tfstate.json` (from HCL)       | `infra/terragrunt/_shared/_config/remote_state.hcl`   |
|                             | `TG_STACK_REMOTE_STATE_BACKEND_TF_FILENAME`         | Filename for the generated backend configuration.                           | `backend.tf` (from HCL)                   | `infra/terragrunt/_shared/_config/remote_state.hcl`   |
| **Module Version Overrides**| `TG_STACK_TF_MODULE_DNI_GENERATOR_VERSION_DEFAULT`  | Default version for the DNI generator module.                               | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/dni_generator.hcl`   |
|                             | `TG_STACK_TF_MODULE_NAME_GENERATOR_VERSION_DEFAULT` | Default version for the name generator module.                              | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/name_generator.hcl`  |
|                             | `TG_STACK_TF_MODULE_AGE_GENERATOR_VERSION_DEFAULT`  | Default version for the age generator module.                               | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/age_generator.hcl`   |
|                             | `TG_STACK_TF_MODULE_LASTNAME_GENERATOR_VERSION_DEFAULT` | Default version for the lastname generator module.                          | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/lastname_generator.hcl`|

*Note: Some Terragrunt standard variables like `TF_INPUT`, `TERRAGRUNT_LOG_LEVEL`, `TERRAGRUNT_AUTO_INIT`, `TERRAGRUNT_AUTO_RETRY` have been replaced by their `TG_` prefixed counterparts (e.g., `TG_NON_INTERACTIVE`, `TG_LOG_LEVEL`, `TG_NO_AUTO_INIT`, `TG_NO_AUTO_RETRY`) in the `.env` (via `.env.example`) for consistency and to align with newer Terragrunt practices where applicable. The `TG_NO_AUTO_INIT` and `TG_NO_AUTO_RETRY` variables have inverted logic compared to their predecessors.*

## Best Practices

- Use the `.env.example` file as a template for your `.env` file.
- Store all local environment-specific configurations in your `.env` file.
- Ensure `.env` is included in your project's `.gitignore` file to avoid committing local configurations or sensitive data.
- Keep `.env.example` up-to-date with all variables required or commonly customized for the project.
- Use descriptive variable names in your `.env` files and document their purpose if they are custom additions not found in `.env.example`.

## Troubleshooting

### Common Issues

1.  **Variables Not Loading**
    - Ensure you have a `.env` file in the project root.
    - If you copied from `.env.example`, ensure the file is named exactly `.env`.
    - Check for syntax errors in your `.env` file (e.g., ensure there are no unquoted spaces in values if the value itself shouldn't contain leading/trailing spaces, though `just`'s dotenv loader is generally robust).
    - Verify that `set dotenv-load` is present in your `justfile`.

2.  **Incorrect Variable Values**
    - Double-check the variable names and values in your `.env` file for typos.
    - Remember that variables loaded from the shell environment may override those in `.env` files, depending on the system and shell configuration (though `just`'s `dotenv-load` typically gives `.env` precedence over pre-existing shell environment variables of the same name unless `dotenv-override` is also set).

### Debugging Commands

```bash
# Show all environment variables available to just (run from a just recipe)
# Example just recipe:
# show-env:
# @echo "Current Environment Variables:"
# @env

# Show a specific variable's value (run from a just recipe)
# Example just recipe for a variable MY_VAR:
# show-my-var:
# @echo "MY_VAR is: {{MY_VAR}}"

# Alternatively, to check outside of just, after running a just command that should load the .env:
# (This depends on your shell and if just exports the variables to its parent shell)
# echo $MY_VARIABLE_NAME
```

## Recommended Tools

- [just](https://just.systems/): For command running and automatic `.env` loading.
- [sops](https://github.com/mozilla/sops): For managing encrypted secrets, which can then be loaded into environment variables (requires separate setup).
