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
func (m *Terragrunt) JobTerragruntUnitStaticCheck(
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

	// actions to run on this job
	actions := []ActionCmd{
		{
			Command: "hclfmt",
			Args:    []string{"--check", "--diff"},
		},
		{
			Command: "terragrunt-info",
			Args:    []string{},
		},
		{
			Command: "validate-inputs",
			Args:    []string{},
		},
		{
			Command: "hclvalidate",
			Args:    []string{"--show-config-path"},
		},
	}

	var wg sync.WaitGroup
	// Buffered channel, so we can collect the results safely.
	resultChan := make(chan ActionResult, len(units)*len(actions))

	// Run each unit check concurrently
	for _, unit := range units {
		for _, action := range actions {
			currentUnit := unit
			currentAction := action

			wg.Add(1)

			go func() {
				defer wg.Done()
				var finalErr error
				var finalOutput string

				// Create a clone of the Terragrunt struct to avoid concurrency issues
				terragruntClone := cloneTerragrunt(m)

				// Builder
				actionBuilder := terragruntClone.NewAction(currentAction.Command)
				actionBuilder.ForTgUnit(defaultRefArchEnv, defaulttRefArchLayer, currentUnit)
				actionBuilder.WithSource(terragruntClone.Src)
				actionBuilder.ConnectToAWSWithCredentials(awsAccessKeyID, awsSecretAccessKey, awsRegion)
				// TODO: Up to you ;)
				// actionBuilder.ConnectToAWSWithOIDC(awsRoleArn, awsOidcToken, awsRegion, awsRoleSessionName)
				actionBuilder.WithNoCache()
				actionBuilder.WithArgs(currentAction.Args...)
				actionBuilder.WithLoadDotEnvFile()
				actionBuilder.WithEnvVars(envVars)

				// Execute
				compiledCtr, buildErr := actionBuilder.Execute(ctx)
				if buildErr != nil {
					finalErr = WrapErrorf(buildErr, "failed to execute the terragrunt action for the following environment: %s, layer: %s, unit: %s", currentUnit, defaulttRefArchLayer, currentUnit)
				} else {
					// get the stdout of the action
					stdOut, stdOutErr := compiledCtr.Stdout(ctx)
					if stdOutErr != nil {
						finalErr = WrapErrorf(stdOutErr, "failed to get the stdout of the terragrunt action for the following environment: %s, layer: %s, unit: %s", currentUnit, defaulttRefArchLayer, currentUnit)
					} else {
						finalOutput = stdOut
					}
				}

				// Collect results
				resultChan <- ActionResult{
					WorkDir: fmt.Sprintf("%s.%s", currentUnit, currentAction.Command),
					Output:  finalOutput,
					Err:     finalErr,
				}
			}()
		}
	}

	wg.Wait()
	close(resultChan)

	finalOutput, err := processActionResults(resultChan)
	if err != nil {
		return "", err
	}

	return finalOutput, nil
}
