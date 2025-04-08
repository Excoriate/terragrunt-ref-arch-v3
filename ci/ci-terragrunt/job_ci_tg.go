package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"sync"
)

func (m *Terragrunt) JobTerragruntUnitStaticCheck(
	// ctx is the context for the Dagger container.
	// +optional
	ctx context.Context,
	// awsAccessKeyID is the AWS access key ID to use for the command.
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key to use for the command.
	awsSecretAccessKey *dagger.Secret,
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
	// secrets is a slice of secrets to pass to the command.
	// +optional
	secrets []*dagger.Secret,
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
				actionBuilder.WithAWS(awsAccessKeyID, awsSecretAccessKey, "eu-central-1")
				actionBuilder.WithNoCache()
				actionBuilder.WithArgs(currentAction.Args...)
				actionBuilder.WithLoadDotEnvFile()
				actionBuilder.WithEnvVars(envVars)
				// actionBuilder.WithSecret(ctx, secrets)

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
					Unit:   fmt.Sprintf("%s.%s", currentUnit, currentAction.Command),
					Output: finalOutput,
					Err:    finalErr,
				}
			}()
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultChan)

	// Process results using the reusable helper function
	finalOutput, err := processActionResults(resultChan)
	if err != nil {
		return "", err // Return the aggregated error
	}

	return finalOutput, nil // Return the formatted success report
}
