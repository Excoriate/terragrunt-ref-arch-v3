package main

import (
	"context"
	"dagger/infra/internal/dagger"
	"fmt"
	"sync"
)

// TODO: Change it accordingly to the environments you use (e.g.: var environments = []string{"global", "dev", "prod"})
var environments = []string{"global"}
var unitsPerStack = map[string][][]string{
	"dni": {
		// Units for the domain stack, meaning, under the infra/terragrunt/<environment>/domain/ folder.
		{
			"dni-generator",
			"lastname-generator",
			"name-generator",
			"age-generator",
		},
	},
	"non-distributable": {
		{
			"random-string-generator",
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
	// remoteStateBucket is the name of the bucket to use for the remote backend.
	// +optional
	remoteStateBucket string,
	// remoteStateLockTable is the name of the lock table to use for the remote backend.
	// +optional
	remoteStateLockTable string,
	// remoteStateRegion is the region of the remote state bucket.
	// +optional
	remoteStateRegion string,
	// deploymentRegion is the AWS region to use for the remote backend.
	// +optional
	deploymentRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// awsSessionToken is the AWS session token.
	// +optional
	awsSessionToken *dagger.Secret,
	// tfGitlabToken is the Terraform Gitlab token.
	// +optional
	tfGitlabToken *dagger.Secret,
	// GitHubToken is the github token
	// +optional
	GitHubToken *dagger.Secret,
	// loadDotEnvFile is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnvFile bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables to set in the container.
	// +optional
	envVars []string,
	// tgBinaryVersionOverride is the Terragrunt binary version to use.
	// +optional
	tgBinaryVersionOverride string,
	// tfBinaryVersionOverride is the Terraform binary version to use.
	// +optional
	tfBinaryVersionOverride string,
	// tfVersionFile is the Terraform version file to use. I'll generate a .terraform-version file in the working directory.
	// +optional
	tfVersionFile string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
	// tgLogLevel is the Terragrunt log level to use.
	// +optional
	tgLogLevel string,
	// stack is the stack name to check.
	stack string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
) (string, error) {
	var remoteStateBucketName string
	var remoteStateLockTableName string

	if remoteStateBucket == "" && remoteStateLockTable == "" {
		remoteStateBucketName = fmt.Sprintf("%s-%s", remoteStateDefaultBucketNamingConvention, environment)
		remoteStateLockTableName = fmt.Sprintf("%s-%s", remoteStateDefaultLockTableNamingConvention, environment)
	} else {
		remoteStateBucketName = remoteStateBucket
		remoteStateLockTableName = remoteStateLockTable
	}

	if remoteStateRegion == "" {
		remoteStateRegion = defaultRemoteStateRegion
	}

	baseCtr, baseCtrErr := m.JobTg(ctx,
		remoteStateBucketName,
		remoteStateLockTableName,
		remoteStateRegion,
		deploymentRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		awsSessionToken,
		tfGitlabToken,
		GitHubToken,
		loadDotEnvFile,
		noCache,
		envVars,
		tgBinaryVersionOverride,
		tfBinaryVersionOverride,
		tfVersionFile,
		gitSSH,
		tgLogLevel,
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
					// FIXME: This command is going tom be deprecated in further versions of Terragrunt. Plan to remove it.
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
	// remoteStateBucket is the name of the bucket to use for the remote backend.
	// +optional
	remoteStateBucket string,
	// remoteStateLockTable is the name of the lock table to use for the remote backend.
	// +optional
	remoteStateLockTable string,
	// remoteStateRegion is the region of the remote state bucket.
	// +optional
	remoteStateRegion string,
	// deploymentRegion is the AWS region to use for the remote backend.
	// +optional
	deploymentRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// awsSessionToken is the AWS session token.
	// +optional
	awsSessionToken *dagger.Secret,
	// tfGitlabToken is the Terraform Gitlab token.
	// +optional
	tfGitlabToken *dagger.Secret,
	// GitHubToken is the github token
	// +optional
	GitHubToken *dagger.Secret,
	// loadDotEnvFile is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnvFile bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables to set in the container.
	// +optional
	envVars []string,
	// tgBinaryVersionOverride is the Terragrunt binary version to use.
	// +optional
	tgBinaryVersionOverride string,
	// tfBinaryVersionOverride is the Terraform binary version to use.
	// +optional
	tfBinaryVersionOverride string,
	// tfVersionFile is the Terraform version file to use. I'll generate a .terraform-version file in the working directory.
	// +optional
	tfVersionFile string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
	// tgLogLevel is the Terragrunt log level to use.
	// +optional
	tgLogLevel string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		remoteStateBucket,
		remoteStateLockTable,
		remoteStateRegion,
		deploymentRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		awsSessionToken,
		tfGitlabToken,
		GitHubToken,
		loadDotEnvFile,
		noCache,
		envVars,
		tgBinaryVersionOverride,
		tfBinaryVersionOverride,
		tfVersionFile,
		gitSSH,
		tgLogLevel,
		"non-distributable",
		environment,
	)
}

func (m *Infra) JobCITgStackDniGeneratorStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// remoteStateBucket is the name of the bucket to use for the remote backend.
	// +optional
	remoteStateBucket string,
	// remoteStateLockTable is the name of the lock table to use for the remote backend.
	// +optional
	remoteStateLockTable string,
	// remoteStateRegion is the region of the remote state bucket.
	// +optional
	remoteStateRegion string,
	// deploymentRegion is the AWS region to use for the remote backend.
	// +optional
	deploymentRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// awsSessionToken is the AWS session token.
	// +optional
	awsSessionToken *dagger.Secret,
	// tfGitlabToken is the Terraform Gitlab token.
	// +optional
	tfGitlabToken *dagger.Secret,
	// GitHubToken is the github token
	// +optional
	GitHubToken *dagger.Secret,
	// loadDotEnvFile is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnvFile bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables to set in the container.
	// +optional
	envVars []string,
	// tgBinaryVersionOverride is the Terragrunt binary version to use.
	// +optional
	tgBinaryVersionOverride string,
	// tfBinaryVersionOverride is the Terraform binary version to use.
	// +optional
	tfBinaryVersionOverride string,
	// tfVersionFile is the Terraform version file to use. I'll generate a .terraform-version file in the working directory.
	// +optional
	tfVersionFile string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
	// tgLogLevel is the Terragrunt log level to use.
	// +optional
	tgLogLevel string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		remoteStateBucket,
		remoteStateLockTable,
		remoteStateRegion,
		deploymentRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		awsSessionToken,
		tfGitlabToken,
		GitHubToken,
		loadDotEnvFile,
		noCache,
		envVars,
		tgBinaryVersionOverride,
		tfBinaryVersionOverride,
		tfVersionFile,
		gitSSH,
		tgLogLevel,
		"dni_generator",
		environment,
	)
}

func (m *Infra) JobCITgStackAgeGeneratorStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// remoteStateBucket is the name of the bucket to use for the remote backend.
	// +optional
	remoteStateBucket string,
	// remoteStateLockTable is the name of the lock table to use for the remote backend.
	// +optional
	remoteStateLockTable string,
	// remoteStateRegion is the region of the remote state bucket.
	// +optional
	remoteStateRegion string,
	// deploymentRegion is the AWS region to use for the remote backend.
	// +optional
	deploymentRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// awsSessionToken is the AWS session token.
	// +optional
	awsSessionToken *dagger.Secret,
	// tfGitlabToken is the Terraform Gitlab token.
	// +optional
	tfGitlabToken *dagger.Secret,
	// GitHubToken is the github token
	// +optional
	GitHubToken *dagger.Secret,
	// loadDotEnvFile is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnvFile bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables to set in the container.
	// +optional
	envVars []string,
	// tgBinaryVersionOverride is the Terragrunt binary version to use.
	// +optional
	tgBinaryVersionOverride string,
	// tfBinaryVersionOverride is the Terraform binary version to use.
	// +optional
	tfBinaryVersionOverride string,
	// tfVersionFile is the Terraform version file to use. I'll generate a .terraform-version file in the working directory.
	// +optional
	tfVersionFile string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
	// tgLogLevel is the Terragrunt log level to use.
	// +optional
	tgLogLevel string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		remoteStateBucket,
		remoteStateLockTable,
		remoteStateRegion,
		deploymentRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		awsSessionToken,
		tfGitlabToken,
		GitHubToken,
		loadDotEnvFile,
		noCache,
		envVars,
		tgBinaryVersionOverride,
		tfBinaryVersionOverride,
		tfVersionFile,
		gitSSH,
		tgLogLevel,
		"age_generator",
		environment,
	)
}

func (m *Infra) JobCITgStackNameGeneratorStaticAnalysis(
	// Context is the context for managing the operation's lifecycle
	// +optional
	ctx context.Context,
	// remoteStateBucket is the name of the bucket to use for the remote backend.
	// +optional
	remoteStateBucket string,
	// remoteStateLockTable is the name of the lock table to use for the remote backend.
	// +optional
	remoteStateLockTable string,
	// remoteStateRegion is the region of the remote state bucket.
	// +optional
	remoteStateRegion string,
	// deploymentRegion is the AWS region to use for the remote backend.
	// +optional
	deploymentRegion string,
	// awsAccessKeyID is the AWS access key ID.
	// +optional
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	// +optional
	awsSecretAccessKey *dagger.Secret,
	// awsSessionToken is the AWS session token.
	// +optional
	awsSessionToken *dagger.Secret,
	// tfGitlabToken is the Terraform Gitlab token.
	// +optional
	tfGitlabToken *dagger.Secret,
	// GitHubToken is the github token
	// +optional
	GitHubToken *dagger.Secret,
	// loadDotEnvFile is a flag to enable source .env files from the local directory.
	// +optional
	loadDotEnvFile bool,
	// NoCache is a flag to disable caching of the container.
	// +optional
	noCache bool,
	// envVars are the environment variables to set in the container.
	// +optional
	envVars []string,
	// tgBinaryVersionOverride is the Terragrunt binary version to use.
	// +optional
	tgBinaryVersionOverride string,
	// tfBinaryVersionOverride is the Terraform binary version to use.
	// +optional
	tfBinaryVersionOverride string,
	// tfVersionFile is the Terraform version file to use. I'll generate a .terraform-version file in the working directory.
	// +optional
	tfVersionFile string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
	// tgLogLevel is the Terragrunt log level to use.
	// +optional
	tgLogLevel string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
) (string, error) {
	return m.JobCITgStackStaticAnalysis(
		ctx,
		remoteStateBucket,
		remoteStateLockTable,
		remoteStateRegion,
		deploymentRegion,
		awsAccessKeyID,
		awsSecretAccessKey,
		awsSessionToken,
		tfGitlabToken,
		GitHubToken,
		loadDotEnvFile,
		noCache,
		envVars,
		tgBinaryVersionOverride,
		tfBinaryVersionOverride,
		tfVersionFile,
		gitSSH,
		tgLogLevel,
		"name_generator",
		environment,
	)
}
