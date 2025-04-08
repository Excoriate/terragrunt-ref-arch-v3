package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"sort"
	"strings"
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

	// collectors
	var collectedActionErrors []error
	// Here the key defined is "unit.command"
	successfulOutputs := make(map[string]string)

	for result := range resultChan {
		if result.Err != nil {
			collectedActionErrors = append(collectedActionErrors, result.Err)
		} else {
			successfulOutputs[result.Unit] = result.Output
		}
	}

	// Handling, and showing errors.
	if len(collectedActionErrors) > 0 {
		return "", JoinErrors(collectedActionErrors...)
	}

	// Handling, and showing outputs
	var outputBuilder strings.Builder
	outputBuilder.WriteString("All static checks passed successfully.\n\nOutput per unit:\n")
	outputBuilder.WriteString("=====================\n")

	// Group outputs by unit
	unitOutputs := make(map[string]map[string]string)
	for key, output := range successfulOutputs {
		parts := strings.Split(key, ".")
		if len(parts) == 2 {
			unit, command := parts[0], parts[1]
			if _, ok := unitOutputs[unit]; !ok {
				unitOutputs[unit] = make(map[string]string)
			}
			unitOutputs[unit][command] = output
		}
	}

	// Sort units for consistent output order
	sortedUnits := make([]string, 0, len(unitOutputs))
	for unit := range unitOutputs {
		sortedUnits = append(sortedUnits, unit)
	}
	sort.Strings(sortedUnits)

	// Display outputs by unit and command
	for _, unit := range sortedUnits {
		outputBuilder.WriteString(fmt.Sprintf("--- Unit: %s ---\n", unit))

		// Sort commands for consistent output order
		commands := make([]string, 0, len(unitOutputs[unit]))
		for cmd := range unitOutputs[unit] {
			commands = append(commands, cmd)
		}
		sort.Strings(commands)

		for _, cmd := range commands {
			outputBuilder.WriteString(fmt.Sprintf("Command: %s\n", cmd))
			stdout := unitOutputs[unit][cmd]
			if stdout == "" {
				outputBuilder.WriteString("(No standard output)\n")
			} else {
				outputBuilder.WriteString(stdout)
				// Ensure a newline after each command's output if not already present
				if !strings.HasSuffix(stdout, "\n") {
					outputBuilder.WriteString("\n")
				}
			}
			outputBuilder.WriteString("\n")
		}
		outputBuilder.WriteString("--------------------\n")
	}

	return outputBuilder.String(), nil // Return combined stdout and nil error
}
