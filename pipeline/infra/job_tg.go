package main

import (
	"context"
	"dagger/infra/internal/dagger"
)

// JobTg performs a command on Terragrunt by:
func (m *Infra) JobTg(
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
) (*dagger.Container, error) {
	if len(envVars) > 0 {
		mWithEnvVars, err := m.WithEnvVars(envVars)
		if err != nil {
			return nil, WrapErrorf(err, "failed to set environment variables")
		}

		m = mWithEnvVars
	}

	// No interactivity is mandatory, since it's a CI/CD pipeline.
	m = m.WithTerragruntNonInteractive()

	if deploymentRegion != "" {
		m = m.WithTrragruntDeploymentRegion(deploymentRegion)
	}

	if tgBinaryVersionOverride != "" {
		m = m.WithTerragrunt(tgBinaryVersionOverride)
	}

	if tfBinaryVersionOverride != "" {
		m = m.WithTerraform(tfBinaryVersionOverride)
	}

	if remoteStateBucket != "" && remoteStateLockTable != "" {
		// Passing the region should be always optional, and shouldn't be considered a mandatory condition,
		// hence, why I haven't included it as part of the condition statement.
		m = m.WithRemoteBackendConfiguration(remoteStateBucket, remoteStateLockTable, remoteStateRegion)
	}

	if tfVersionFile != "" {
		m = m.WithDotTerraformVersionFileGeneration(tfVersionFile)
	}

	if noCache {
		m = m.WithCacheBuster()
	}

	if tgLogLevel != "" {
		mDecorated, err := m.WithTerragruntLogLevelProgramatically(tgLogLevel)
		if err != nil {
			return nil, WrapErrorf(err, "failed to run the job with log level %s", tgLogLevel)
		}

		m = mDecorated
	}

	if gitSSH != nil {
		m = m.WithSSHAuthSocket(gitSSH, "", "", false, true)
	}

	if loadDotEnvFile {
		mDecorated, err := m.WithDotEnvFile(ctx, m.Src)
		if err != nil {
			return nil, WrapErrorf(err, "failed to source .env files from the local directory")
		}

		m = mDecorated
	}

	if awsAccessKeyID != nil && awsSecretAccessKey != nil {
		m = m.WithAWSKeys(ctx, awsAccessKeyID, awsSecretAccessKey, deploymentRegion, awsSessionToken)
	}

	if tfGitlabToken != nil {
		m = m.WithTerraformGitlabToken(ctx, tfGitlabToken)
	}

	if GitHubToken != nil {
		m = m.WithGitHubToken(ctx, GitHubToken)
	}

	return m.Ctr, nil
}

// JobTg performs a command on Terragrunt by:
func (m *Infra) JobTgExec(
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
	// cmd is the command to execute on the container.
	cmd []string,
	// environment is the environment to use for the container.
	// +optional
	environment string,
	// layer is the stack or layer to use for the container.
	layer string,
	// unit is the unit to use for the container.
	unit string,
) (string, error) {
	// Getting the base container
	jobTgCtrBase, jobTgErr := m.JobTg(ctx,
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
	)

	if jobTgErr != nil {
		return "", WrapErrorf(jobTgErr, "failed to create base jobTg container for environment %s, stack %s, unit %s", environment, layer, unit)
	}

	if environment == "" {
		environment = defaultRefArchEnv
	}

	// Getting the Terragrunt working directory
	tgWorkDir := getTerragruntExecutionPath(environment, layer, unit)

	tgCmd := []string{"terragrunt"}
	tgCmd = append(tgCmd, cmd...)
	tgCmd = append(tgCmd, "--working-dir", tgWorkDir)

	stdout, err := jobTgCtrBase.
		WithExec(tgCmd).
		Stdout(ctx)

	if err != nil {
		return "", WrapErrorf(err, "failed to execute command %v", tgCmd)
	}

	return stdout, nil
}

// JobTgStack runs the Terragrunt commands for the specified stack.
//
// This function takes the following parameters:
//   - ctx: The context for managing the operation's lifecycle.
//   - stack: The stack name to check.
//   - awsRegion: The AWS region to use for the remote backend.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - loadDotEnv: A flag to enable source .env files from the local directory.
//   - noCache: A flag to disable caching of the container.
//   - envVars: The environment variables that will be used to run the Terragrunt commands.
//   - environment: The environment to run the Terragrunt commands.
//   - gitSSH: A flag to enable SSH for the container.
//   - tgCmd: The commands to run on the container.
//   - tgCmdArgs: The arguments to run on the command.
//   - gitlabToken: The Gitlab token.
//
// Returns:
//   - string: The output of the Terragrunt commands.
func (m *Infra) JobTgStack(
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
	// tgCmd is the command to run on the container.
	tgCmd []string,
	// tgCmdArgs is the arguments to run on the command.
	// +optional
	tgCmdArgs []string,
	// stack is the stack name to check.
	stack string,
	// environment is the environment to use for the container.
	environment string,
) (string, error) {
	baseCtr, baseCtrErr := m.JobTg(ctx,
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
	)

	if len(tgCmd) == 0 {
		return "", WrapErrorf(nil, "no commands to run for stack %s", stack)
	}

	tgCmd = append([]string{"terragrunt", "run-all"}, tgCmd...)

	if baseCtrErr != nil {
		return "", WrapErrorf(baseCtrErr, "failed to create base jobTg container for the job tg-stack %s", stack)
	}

	// Define the working directory for this unit
	tgWorkDir := getTerragruntExecutionPathForStacks(environment, stack)

	if len(tgCmdArgs) > 0 {
		tgCmd = append(tgCmd, tgCmdArgs...)
	}

	tgCmd = append(tgCmd, "--working-dir", tgWorkDir)

	// Add the commands, and the working directory to the container
	baseCtr = baseCtr.
		WithExec(tgCmd)

	tgCmdOut, tgCmdErr := baseCtr.
		Stdout(ctx)

	if tgCmdErr != nil {
		return "", WrapErrorf(tgCmdErr, "failed to run terragrunt commands for stack %s", stack)
	}

	return tgCmdOut, nil
}

func executeDaggerCtrAsync(
	ctx context.Context,
	resultChan chan<- JobResult,
	baseCtr *dagger.Container,
	tgWorkDir string,
	commands [][]string,
) {
	jobRes := JobResult{WorkDir: tgWorkDir, Output: "", Err: nil}

	execCtr := baseCtr
	for _, command := range commands {
		execCtr = execCtr.
			WithExec(command)
	}

	stdout, err := execCtr.Stdout(ctx)
	jobRes.Output = stdout

	if err != nil {
		jobRes.Err = WrapErrorf(err, "dagger command failed on working directory: %s", tgWorkDir)
	}

	resultChan <- jobRes
}
