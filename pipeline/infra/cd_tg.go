package main

import (
	"context"
	"dagger/infra/internal/dagger"
	"fmt"
)

func (m *Infra) JobCDTgStack(
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
	// stack is the stack to run the Terragrunt commands.
	stack string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
	// runApply is a flag to run the apply command.
	// +optional
	runApply bool,
	// runDestroy is a flag to run the destroy command.
	// +optional
	runDestroy bool,
	// runPlan is a flag to run the plan command.
	// +optional
	runPlan bool,
) (string, error) {
	var jobStackOut string
	var jobStackErr error

	if !runApply && !runDestroy && !runPlan {
		return "", fmt.Errorf("either --run-apply (runApply) or --run-destroy (runDestroy) or --run-plan (runPlan) must be set to true")
	}

	if runApply && runDestroy {
		return "", fmt.Errorf("cannot set both --run-apply (runApply) and --run-destroy (runDestroy) to true")
	}

	if runPlan && runApply {
		return "", fmt.Errorf("cannot set both --run-plan (runPlan) and --run-apply (runApply) to true")
	}

	if runPlan && runDestroy {
		return "", fmt.Errorf("cannot set both --run-plan (runPlan) and --run-destroy (runDestroy) to true")
	}

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

	if runApply {
		jobStackOut, jobStackErr = m.JobTgStack(ctx,
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
			[]string{"apply"},
			[]string{"-auto-approve"},
			stack,
			environment,
		)
	} else if runDestroy {
		jobStackOut, jobStackErr = m.JobTgStack(ctx,
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
			[]string{"destroy"},
			[]string{"-auto-approve"},
			stack,
			environment,
		)
	} else if runPlan {
		jobStackOut, jobStackErr = m.JobTgStack(ctx,
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
			[]string{"plan"},
			[]string{},
			stack,
			environment,
		)
	}
	// End of Selection

	return jobStackOut, jobStackErr
}

func (m *Infra) JobCDTgStackNonDistributable(
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
	// runApply is a flag to run the apply command.
	// +optional
	runApply bool,
	// runDestroy is a flag to run the destroy command.
	// +optional
	runDestroy bool,
	// runPlan is a flag to run the plan command.
	// +optional
	runPlan bool,
) (string, error) {
	return m.JobCDTgStack(ctx,
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
		runApply,
		runDestroy,
		runPlan,
	)
}
