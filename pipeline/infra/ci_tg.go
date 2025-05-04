package main

import (
	"context"
	"dagger/infra/internal/dagger"
	"sync"
)

// TODO: Change it accordingly to the environments you use (e.g.: var environments = []string{"global", "dev", "prod"})
var environments = []string{"global"}
var unitsPerStack = map[string][][]string{
	"dni": {
		// Units for the domain stack, meaning, under the infra/terragrunt/<environment>/domain/ folder.
		{
			"dni_generator",
			"lastname_generator",
			"name_generator",
			"age_generator",
		},
	},
	// TODO: Add the stacks you use, and their respective units.
}

func getTotalUnitsPerStack(stack string) int {
	units := unitsPerStack[stack]
	totalUnits := 0

	for _, unit := range units {
		totalUnits += len(unit)
	}

	return totalUnits
}

// JobCITgStackStaticAnalysis runs the Terragrunt CI checks for the specified stack.
//
// This function takes the following parameters:
//   - ctx: The context for managing the operation's lifecycle.
//   - stack: The stack name to check (e.g., "non-distributable", "domain", "landing-zone", "repositories").
//   - awsRegion: The AWS region to use for the remote backend.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - loadDotEnv: A flag to enable source .env files from the local directory.
func (m *Infra) JobCITgStackStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// stack is the stack name to check.
	stack string,
	// awsRegion is the AWS region to use for the remote backend.
	// +optional
	awsRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// loadDotEnv is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnv bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables that will be used to run the Terragrunt commands.
	// +optional
	envVars []string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
) (string, error) {
	baseCtr, baseCtrErr := m.JobTg(ctx,
		"",
		"",
		awsRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		nil,
		loadDotEnv,
		noCache,
		envVars,
		"",
		"",
		gitSSH,
	)

	if baseCtrErr != nil {
		return "", WrapErrorf(baseCtrErr, "failed to create base jobTg container for stack %s", stack)
	}

	// Get the units for the specified stack
	unitsList, exists := unitsPerStack[stack]
	if !exists {
		return "", WrapErrorf(nil, "stack %s not found in unitsPerStack", stack)
	}

	// Concurrency setup
	var wg sync.WaitGroup

	// Calculate total number of units for channel buffer size
	totalUnits := getTotalUnitsPerStack(stack)
	if totalUnits == 0 {
		return "", WrapErrorf(nil, "no units found for stack %s", stack)
	}

	resultChan := make(chan JobResult, totalUnits)

	// Process each unit group
	for _, units := range unitsList {
		// Process each unit in the group
		for _, unit := range units {
			// Increment WaitGroup counter
			wg.Add(1)

			// Launch goroutine for each unit
			go func(unitName string) {
				defer wg.Done()

				// Define the working directory for this unit
				tgWorkDir := getTerragruntExecutionPath(environment, stack, unitName)

				// Define the commands to execute
				commands := [][]string{
					{"terragrunt", "init", "--working-dir", tgWorkDir},
					{"terragrunt", "terragrunt-info", "--working-dir", tgWorkDir},
					{"terragrunt", "hclfmt", "--check", "--diff", "--working-dir", tgWorkDir},
					{"terragrunt", "validate-inputs", "--working-dir", tgWorkDir},
					{"terragrunt", "hclvalidate", "--show-config-path", "--working-dir", tgWorkDir},
				}

				// Execute commands asynchronously using the helper function
				executeDaggerCtrAsync(ctx, resultChan, baseCtr, tgWorkDir, commands)

			}(unit)
		}
	}

	// Start a goroutine to close the channel once all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return processActionAsyncResults(resultChan)
}

// JobCITgStackNonDistributableStaticAnalysis runs the Terragrunt CI checks for the non-distributable stack.
//
// This function takes the following parameters:
//   - ctx: The context for managing the operation's lifecycle.
//   - awsRegion: The AWS region to use for the remote backend.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - loadDotEnv: A flag to enable source .env files from the local directory.
func (m *Infra) JobCITgStackNonDistributableStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// awsRegion is the AWS region to use for the remote backend.
	// +optional
	awsRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// loadDotEnv is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnv bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables that will be used to run the Terragrunt commands.
	// +optional
	envVars []string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		"non-distributable",
		awsRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		loadDotEnv,
		noCache,
		envVars,
		environment,
		nil,
	)
}

// JobCITgStackDomainStaticAnalysis runs the Terragrunt CI checks for the domain stack.
//
// This function takes the following parameters:
//   - ctx: The context for managing the operation's lifecycle.
//   - awsRegion: The AWS region to use for the remote backend.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - loadDotEnv: A flag to enable source .env files from the local directory.
func (m *Infra) JobCITgStackDomainStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// awsRegion is the AWS region to use for the remote backend.
	// +optional
	awsRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// loadDotEnv is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnv bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables that will be used to run the Terragrunt commands.
	// +optional
	envVars []string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		"domain",
		awsRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		loadDotEnv,
		noCache,
		envVars,
		environment,
		gitSSH,
	)
}

// JobCITgStackLandingZoneStaticAnalysis runs the Terragrunt CI checks for the landing-zone stack.
//
// This function takes the following parameters:
//   - ctx: The context for managing the operation's lifecycle.
//   - awsRegion: The AWS region to use for the remote backend.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - loadDotEnv: A flag to enable source .env files from the local directory.
func (m *Infra) JobCITgStackLandingZoneStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// awsRegion is the AWS region to use for the remote backend.
	// +optional
	awsRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// loadDotEnv is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnv bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables that will be used to run the Terragrunt commands.
	// +optional
	envVars []string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		"landing-zone",
		awsRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		loadDotEnv,
		noCache,
		envVars,
		environment,
		gitSSH,
	)
}

// JobCITgStackRepositoriesStaticAnalysis runs the Terragrunt CI checks for the repositories stack.
//
// This function takes the following parameters:
//   - ctx: The context for managing the operation's lifecycle.
//   - awsRegion: The AWS region to use for the remote backend.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - loadDotEnv: A flag to enable source .env files from the local directory.
func (m *Infra) JobCITgStackRepositoriesStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// awsRegion is the AWS region to use for the remote backend.
	// +optional
	awsRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// loadDotEnv is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnv bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables that will be used to run the Terragrunt commands.
	// +optional
	envVars []string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		"repositories",
		awsRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		loadDotEnv,
		noCache,
		envVars,
		environment,
		gitSSH,
	)
}
