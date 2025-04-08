# Terragrunt CI Dagger Module

This Dagger module provides a reusable environment and functions for running Terragrunt and Terraform commands, primarily intended for Continuous Integration (CI) workflows within the context of the `terragrunt-ref-arch-v3` repository structure.

## Core Concepts

*   **Dagger Module:** This is a self-contained unit of Dagger configuration and code, written in Go, that exposes functions callable by Dagger clients (CLI, other SDKs). It encapsulates the setup and execution logic for Terragrunt/Terraform tasks.
*   **`Terragrunt` Struct:** The central component of this module (`main.go`). It holds the Dagger `Container` (`Ctr`) where commands are executed and a reference to the source `Directory` (`Src`).
*   **Container Setup:** The `New()` constructor initializes the module. It can:
    *   Use a pre-built container image (`imageURL`).
    *   Use a custom container (`ctr`).
    *   Build a default container based on Alpine (`main.go:defaultImage`), installing specific versions of Git, Terraform (`main.go:defaultTerraformVersion`), and Terragrunt (`main.go:defaultTerragruntVersion`).
    *   The `CommonSetup()` method (`main.go`) applies standard configurations like installing tools and setting up caching.
*   **Caching:** The module configures Dagger cache volumes for Terraform plugins (`main.go:configterraformPluginCachePath`) and Terragrunt artifacts (`main.go:configterragruntCachePath`) to speed up subsequent runs. It also enables the Terragrunt provider cache server by default (`WithTerragruntProvidersCacheServerEnabled`). Caching can be bypassed using the `noCache` flag in job functions or the `WithCacheBuster` method.

## Features & Configuration (`main.go`)

The module offers various methods (primarily on the `Terragrunt` struct) to configure the execution environment:

*   **Source Code:** `WithSRC()` mounts your project source code into the container.
*   **Tool Versions:** `WithTerraform()`, `WithTerragrunt()` install specific tool versions if not using a pre-built image.
*   **Environment Variables:**
    *   `WithEnvVars()`: Sets multiple environment variables.
    *   `WithDotEnvFile()`: Loads variables from `.env` files found in the source directory.
*   **Secrets:**
    *   `WithSecrets()`: Mounts multiple Dagger secrets as environment variables.
    *   `WithToken()`: A specific helper for mounting a single secret.
    *   `WithAWSOIDC()`: Configures AWS authentication via OIDC, mounting the token as a secret file.
    *   `WithAWSKeys()`: Configures AWS authentication using access keys (secrets).
    *   `WithGitlabToken()`, `WithGitHubToken()`, `WithTerraformToken()`: Specific helpers for common tokens (mounted as env vars).
*   **Authentication:**
    *   `WithNewNetrcFile*()`: Creates a `.netrc` file for Git authentication (GitHub/GitLab). Use `WithNewNetrcFileAsSecret*` for secure password handling.
    *   `WithSSHAuthSocket()`: Mounts an SSH agent socket for Git SSH authentication.
    *   AWS Auth: Covered by `WithAWSKeys` and `WithAWSOIDC`.
*   **Caching:**
    *   `WithTerraformPluginCache()`: Configures TF plugin cache dir and env var.
    *   `WithTerragruntCache()`: Configures TG cache volume.
    *   `WithTerragruntProvidersCacheServerEnabled() / Disabled()`: Toggles the TG provider cache server.
    *   `WithRegistriesToCacheProvidersFrom()`: Adds custom registries to the TG provider cache.
    *   `WithCacheBuster()`: Adds an environment variable to break build cache layers if needed.
*   **Terragrunt Options:**
    *   `WithTerragruntLogLevel()`: Sets `TERRAGRUNT_LOG_LEVEL`.
    *   `WithTerragruntNonInteractive()`: Sets `TERRAGRUNT_NON_INTERACTIVE`.
    *   `WithTerragruntNoColor()`: Sets `TERRAGRUNT_NO_COLOR`.
*   **Terminal:** `OpenTerminal()` provides an interactive terminal within the fully configured container.

## Available Jobs

The module provides pre-defined functions for common CI tasks:

*   **`JobTerraformModulesStaticCheck(ctx)` (`job_ci_tf.go`):**
    *   **Purpose:** Performs static checks (`init -backend=false`, `validate`, `fmt -recursive`) on Terraform modules.
    *   **Note:** Currently iterates over a hardcoded list of module names specific to this reference architecture.
    *   **Execution:** Synchronous execution across modules.
*   **`JobTerragruntUnitsStaticCheck(ctx, ...)` (`job_ci_tg.go`):**
    *   **Purpose:** Performs static checks (`init`, `terragrunt-info`, `hclfmt --check`, `validate-inputs`, `hclvalidate`) on Terragrunt units (directories containing `terragrunt.hcl`).
    *   **Parameters:** Accepts AWS credentials/OIDC config, cache options, custom arguments (`tgArgs`), env vars.
    *   **Note:** Currently iterates over a hardcoded list of unit names specific to this reference architecture (`dni_generator`, etc.) within a fixed env/layer (`global`/`dni`).
    *   **Execution:** Asynchronous execution across units using goroutines.
*   **`JobTerragruntUnitsPlan(ctx, ...)` (`job_ci_tg.go`):**
    *   **Purpose:** Runs `terragrunt plan` on Terragrunt units.
    *   **Parameters:** Accepts AWS credentials/OIDC config, cache options, custom arguments (`tgArgs`), env vars.
    *   **Note:** Currently iterates over a hardcoded list of unit names specific to this reference architecture.
    *   **Execution:** Asynchronous execution across units using goroutines.

*   **Generic Execution:** `Exec(ctx, binary, command, args...)` (`main.go`):
    *   **Purpose:** Allows executing arbitrary `terragrunt` or `terraform` commands with fine-grained control over arguments, environment variables, secrets, tokens, AWS auth, SSH sockets, etc.
    *   **Use Case:** Useful for running commands not covered by the pre-defined jobs or for custom workflows.

## Usage Example (Dagger Go Client)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger/terragrunt/internal/dagger"
)

func main() {
	ctx := context.Background()

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatalf("Error connecting to Dagger engine: %v", err)
	}
	defer client.Close()

	// Get reference to source code directory
	srcDir := client.Host().Directory(".") // Assumes run from repo root

	// Initialize the Terragrunt module
	// Uses default versions and Alpine base image
	tgModule, err := client.Terragrunt().New(ctx,
		// No imageURL, tgVersion, tfVersion specified, will use defaults
		dagger.TerragruntNewOpts{
			SrcDir: srcDir,
			// Add any necessary EnvVars here if not using .env files
			// EnvVars: []string{"MY_VAR=value"},
		},
	)
	if err != nil {
		log.Fatalf("Error initializing Terragrunt module: %v", err)
	}

	// Example: Run Terragrunt static checks
	// Assumes AWS credentials are set as environment variables or via OIDC in the CI environment
	// Secrets need to be created using client.SetSecret()
	// awsAccessKey = client.SetSecret(...)
	// awsSecretKey = client.SetSecret(...)
	// oidcToken = client.SetSecret(...)

	output, err := tgModule.JobTerragruntUnitsStaticCheck(ctx,
		dagger.TerragruntJobTerragruntUnitsStaticCheckOpts{
			// Pass AWS secrets if needed:
			// AwsAccessKeyID: awsAccessKey,
			// AwsSecretAccessKey: awsSecretKey,
			// Or OIDC token:
			// AwsOidcToken: oidcToken,
			// AwsRoleArn: "arn:aws:iam::ACCOUNT:role/MyCIRole",
			LoadDotEnvFile: true, // Example: Load .env files from srcDir root
		},
	)

	if err != nil {
		log.Fatalf("Error running Terragrunt static checks: %v", err)
	}

	fmt.Println("Terragrunt static checks completed successfully:")
	fmt.Println(output)

	// Example: Run Terraform module checks
	tfModOutput, err := tgModule.JobTerraformModulesStaticCheck(ctx)
	if err != nil {
		log.Fatalf("Error running Terraform module checks: %v", err)
	}
	fmt.Println("Terraform module checks completed successfully:")
	fmt.Println(tfModOutput)
}

```

## Extending the Module

This module can be extended to add new functionality or adapt existing jobs:

1.  **Adding New Jobs:**
    *   Create a new file (e.g., `job_my_custom_task.go`).
    *   Define a new method on the `*Terragrunt` struct (e.g., `func (m *Terragrunt) JobMyCustomTask(...)`).
    *   Within the method, use the existing `m.Ctr` and `With*` methods (`main.go`) to configure the container as needed.
    *   Use `m.Ctr.WithExec([]string{"terragrunt", "your-command", ...})` to run the desired command(s).
    *   Return the output string and error, potentially using the result processing helpers from `job.go` if running multiple steps or concurrently.
2.  **Modifying Existing Jobs:**
    *   Update the hardcoded lists of modules/units in `job_ci_tf.go` and `job_ci_tg.go` or modify the logic to dynamically discover them (e.g., by scanning the source directory).
    *   Adjust the specific `WithExec` commands within the job functions.
3.  **Adding Configuration:**
    *   Add new `With*` methods to `main.go` to support additional configuration options (e.g., different cloud providers, new tools).
    *   These methods should typically modify `m.Ctr` by adding environment variables, mounting files/secrets, or installing packages.

Refer to the Dagger Go SDK documentation (concepts similar to those in `docs/dagger.io`) for more details on interacting with containers, directories, secrets, and executing commands.