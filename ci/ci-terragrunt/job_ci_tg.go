package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"sync"
)

// JobTerragruntUnitStaticCheck performs a series of static checks on Terragrunt
// configurations for a specific unit. It utilizes the provided AWS credentials
// and other parameters to execute the checks in a Dagger container context.
//
// Parameters:
//
//	ctx: A context.Context for managing the lifecycle of the Dagger container.
//	     This is used to control cancellation and timeouts for the operation.
//	awsAccessKeyID: A pointer to a dagger.Secret containing the AWS access key ID
//	                 required for authentication with AWS services.
//	awsSecretAccessKey: A pointer to a dagger.Secret containing the AWS secret access
//	                    key required for authentication with AWS services.
//	awsRoleArn: A string representing the AWS role ARN to assume for the command.
//	             This parameter is optional and can be omitted if not needed.
//	awsOidcToken: A pointer to a dagger.Secret containing the OIDC token (JWT)
//	               obtained from CI (e.g., from CI_JOB_JWT_V2). This is used for
//	               authentication when using OIDC. This parameter is optional.
//	awsRegion: A string specifying the AWS region to use for the command. This
//	           parameter is optional and can be omitted if not needed.
//	awsRoleSessionName: A string representing an optional name for the assumed
//	                    role session. This parameter is optional.
//	noCache: A boolean flag indicating whether to disable caching for the command.
//	         If set to true, the command will not use any cached results.
//	tgArgs: A slice of strings containing arbitrary arguments to pass to the
//	         Terragrunt command. This parameter is optional.
//	loadDotEnvFile: A boolean indicating whether to load a .env file for
//	                 environment variable configuration. This parameter is optional.
//	envVars: A slice of strings including all environment variables to pass to
//	          the command. This parameter is optional.
//
// Returns:
//
//	A string containing the output of the static checks, and an error if any
//	occurred during the execution of the checks.
func (m *Terragrunt) JobTerragruntUnitsStaticCheck(
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
	units := []string{
		"dni_generator",
		"lastname_generator",
		"name_generator",
		"age_generator",
	}

	var wg sync.WaitGroup
	resultChan := make(chan JobResult, len(units))

	if loadDotEnvFile {
		mWithDotEnvLoaded, err := m.WithDotEnvFile(ctx, m.Src)

		if err != nil {
			return "", WrapErrorf(err, "failed to load .env file")
		}

		m = mWithDotEnvLoaded
	}

	// Connect to AWS, with credentials or OIDC.
	m = m.WithAWSKeys(ctx, awsAccessKeyID, awsSecretAccessKey, awsRegion)

	if noCache {
		m = m.WithCacheBuster()
	}

	if len(envVars) > 0 {
		mWithEnvVars, err := m.WithEnvVars(envVars)

		if err != nil {
			return "", WrapErrorf(err, "failed to set environment variables")
		}

		m = mWithEnvVars
	}

	// Run each unit check concurrently
	for _, unit := range units {
		currentUnit := unit

		// mount path dynamically calculated per unit.
		tgUnitPath := getTerragruntExecutionPath(defaultRefArchEnv, defaulttRefArchLayer, currentUnit)
		tgUnitMntPath := fmt.Sprintf("%s/%s", defaultMntPath, tgUnitPath)

		// Global configuration, let's start with the source code and the .env file.
		mWithSrc, err := m.WithSRC(ctx, tgUnitMntPath, m.Src)

		if err != nil {
			return "", WrapErrorf(err, "failed to set source code for unit %s", currentUnit)
		}

		m = mWithSrc

		wg.Add(1)

		go func() {
			defer wg.Done()
			var finalErr error
			var finalOutput string

			// Execution runtime, just a fancy nawe for the container's state per execution.
			execCtr := m.Ctr.
				WithExec([]string{"terragrunt", "init", "--working-dir", tgUnitMntPath}).
				WithExec([]string{"terragrunt", "terragrunt-info", "--working-dir", tgUnitMntPath}).
				WithExec([]string{"terragrunt", "hclfmt", "--check", "--diff", "--working-dir", tgUnitMntPath}).
				WithExec([]string{"terragrunt", "validate-inputs", "--working-dir", tgUnitMntPath}).
				WithExec([]string{"terragrunt", "hclvalidate", "--show-config-path", "--working-dir", tgUnitMntPath})

			stdout, err := execCtr.Stdout(ctx)
			if err != nil {
				finalErr = WrapErrorf(err, "failed to get the stdout of the terragrunt action for the following environment: %s, layer: %s, unit: %s", currentUnit, defaulttRefArchLayer, currentUnit)
			} else {
				finalOutput = stdout
			}

			// Collect results
			resultChan <- JobResult{
				WorkDir: tgUnitPath,
				Output:  finalOutput,
				Err:     finalErr,
			}
		}()
	}

	wg.Wait()
	close(resultChan)

	finalOutput, err := processActionAsyncResults(resultChan)
	if err != nil {
		return "", err
	}

	return finalOutput, nil
}

// JobTerragruntUnitsPlan executes a Terragrunt plan command for multiple units concurrently.
// It takes a context for managing the lifecycle of the Dagger container and various AWS credentials
// and configuration options necessary for the command execution.
//
// Parameters:
//   - ctx: The context for the Dagger container. This is used to control the execution lifecycle.
//     // +optional
//   - awsAccessKeyID: The AWS access key ID to use for the command. This is a secret value.
//     // +optional
//   - awsSecretAccessKey: The AWS secret access key to use for the command. This is a secret value.
//     // +optional
//   - awsRoleArn: The AWS role ARN to use for the command. This is used for assuming a role in AWS.
//     // +optional
//   - awsOidcToken: The OIDC token (JWT) obtained from CI (e.g., from CI_JOB_JWT_V2). Pass as a secret.
//     // +optional
//   - awsRegion: The AWS region to use for the command. This specifies the geographical region for AWS services.
//     // +optional
//   - awsRoleSessionName: An optional name for the assumed role session. This can help identify the session in AWS.
//     // +optional
//   - noCache: A boolean flag that indicates whether to disable caching for the command execution.
//     // +optional
//   - tgArgs: A slice of arbitrary arguments to pass to the Terragrunt command. This allows customization of the command.
//     // +optional
//   - loadDotEnvFile: A boolean indicating whether to load a .env file. If true, environment variables from the file will be loaded.
//     // +optional
//   - envVars: A slice of strings including all environment variables to pass to the command. This allows for dynamic configuration.
//     // +optional
//
// Returns:
// - A string containing the output of the command execution, or an error if the execution fails.
func (m *Terragrunt) JobTerragruntUnitsPlan(
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
	units := []string{
		"dni_generator",
		"lastname_generator",
		"name_generator",
		"age_generator",
	}

	var wg sync.WaitGroup
	resultChan := make(chan JobResult, len(units))

	if loadDotEnvFile {
		mWithDotEnvLoaded, err := m.WithDotEnvFile(ctx, m.Src)

		if err != nil {
			return "", WrapErrorf(err, "failed to load .env file")
		}

		m = mWithDotEnvLoaded
	}

	// Connect to AWS, with credentials or OIDC.
	m = m.WithAWSKeys(ctx, awsAccessKeyID, awsSecretAccessKey, awsRegion)

	if noCache {
		m = m.WithCacheBuster()
	}

	if len(envVars) > 0 {
		mWithEnvVars, err := m.WithEnvVars(envVars)

		if err != nil {
			return "", WrapErrorf(err, "failed to set environment variables")
		}

		m = mWithEnvVars
	}

	// Run each unit check concurrently
	for _, unit := range units {
		currentUnit := unit

		// mount path dynamically calculated per unit.
		tgUnitPath := getTerragruntExecutionPath(defaultRefArchEnv, defaulttRefArchLayer, currentUnit)
		tgUnitMntPath := fmt.Sprintf("%s/%s", defaultMntPath, tgUnitPath)

		// Global configuration, let's start with the source code and the .env file.
		mWithSrc, err := m.WithSRC(ctx, tgUnitMntPath, m.Src)

		if err != nil {
			return "", WrapErrorf(err, "failed to set source code for unit %s", currentUnit)
		}

		m = mWithSrc

		wg.Add(1)

		go func() {
			defer wg.Done()
			var finalErr error
			var finalOutput string

			// Execution runtime, just a fancy nawe for the container's state per execution.
			execCtr := m.Ctr.
				WithExec([]string{"terragrunt", "init", "--working-dir", tgUnitMntPath}).
				WithExec([]string{"terragrunt", "plan", "--working-dir", tgUnitMntPath})

			stdout, err := execCtr.Stdout(ctx)
			if err != nil {
				finalErr = WrapErrorf(err, "failed to get the stdout of the terragrunt action for the following environment: %s, layer: %s, unit: %s", currentUnit, defaulttRefArchLayer, currentUnit)
			} else {
				finalOutput = stdout
			}

			// Collect results
			resultChan <- JobResult{
				WorkDir: tgUnitPath,
				Output:  finalOutput,
				Err:     finalErr,
			}
		}()
	}

	wg.Wait()
	close(resultChan)

	finalOutput, err := processActionAsyncResults(resultChan)
	if err != nil {
		return "", err
	}

	return finalOutput, nil
}
