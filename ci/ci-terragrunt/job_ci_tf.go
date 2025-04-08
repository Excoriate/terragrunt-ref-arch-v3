package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
)

func (m *Terragrunt) JobTerraformModulesStaticCheck(
	// ctx is the context for the Dagger container.
	// +optional
	ctx context.Context,
	// awsAccessKeyID is the AWS access key ID to use for the command.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key to use for the command.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// awsRoleArn is the AWS role ARN to use for the command.
	// +optional
	awsRoleArn string,
	// awsOidcToken is the OIDC token (JWT) obtained from  CI (e.g., from CI_JOB_JWT_V2). Pass as a secret.
	// +optional
	awsOidcToken *dagger.Secret,
	// awsRegion is the AWS region to use for the command.
	// +optional
	awsRegion string,
	// awsRoleSessionName is an optional name for the assumed role session.
	// +optional
	awsRoleSessionName string,
	// noCache set the cache buster
	// +optional
	noCache bool,
	// tgArgs are the arbitrary arguments to pass to the terrawgrunt command
	// +optional
	tgArgs []string,
	// loadDotEnvFile is a boolean indicating whether to load .env file
	// +optional
	loadDotEnvFile bool,
	// envVars is a slice of strings including all environment variables to pass to the command.
	// +optional
	envVars []string,
) (string, error) {
	// Slice to collect results from each action
	var results []ActionResult

	modules := []string{
		"dni-generator",
		"lastname-generator",
		"name-generator",
		"age-generator",
	}

	// actions to run on this job
	syncActions := []ActionCmd{
		{
			Command: "init",
			Args:    []string{"-backend=false", "-no-color"},
		},
		{
			Command: "validate",
			Args:    []string{"-no-color"},
		},
	}

	// Iterate through each module and action synchronously
	for _, module := range modules {
		for _, action := range syncActions {
			var actionErr error
			var actionOutput string

			actionBuilder := m.NewAction(action.Command)
			actionBuilder.ForTgModule(module) // Use ForTgModule for TF module context
			actionBuilder.WithTerraformInsteadOfTerragrunt()
			actionBuilder.WithNoCache()            // Apply common flags if needed
			actionBuilder.WithArgs(action.Args...) // Add specific args for the command
			actionBuilder.WithSourceCodeMounted(m.Src)

			// Execute the action
			compiledCtr, buildErr := actionBuilder.Execute(ctx)
			if buildErr != nil {
				// Capture the build error
				actionErr = WrapErrorf(buildErr, "failed to execute terraform action '%s' for module '%s'", action.Command, module)
			} else {
				// If execution succeeded, try to get stdout
				stdOut, stdOutErr := compiledCtr.Stdout(ctx)
				if stdOutErr != nil {
					// Capture the stdout error
					actionErr = WrapErrorf(stdOutErr, "failed to get stdout for terraform action '%s' for module '%s'", action.Command, module)
				} else {
					actionOutput = stdOut
				}
			}

			// Append the result (including potential error) to the results slice
			results = append(results, ActionResult{
				WorkDir: fmt.Sprintf("%s.%s", module, action.Command), // Use module.command as identifier
				Output:  actionOutput,
				Err:     actionErr,
			})
		}
	}

	// Process all collected results using the reusable function
	return ProcessActionSyncResults(results)
}
