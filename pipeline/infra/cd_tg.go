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
	// stack is the stack to run the Terragrunt commands.
	stack string,
	// environment is the environment to run the Terragrunt commands.
	environment string,
	// gitSSH is a flag to enable SSH for the container.
	// +optional
	gitSSH *dagger.Socket,
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

	if runApply {
		jobStackOut, jobStackErr = m.JobTgStack(ctx,
			stack,
			awsRegion,
			awsAccessKeyID,
			awsSecretAccessKey,
			loadDotEnv,
			noCache,
			envVars,
			environment,
			gitSSH,
			[]string{"apply"},
			[]string{"-auto-approve"},
			nil,
		)
	} else if runDestroy {
		jobStackOut, jobStackErr = m.JobTgStack(ctx,
			stack,
			awsRegion,
			awsAccessKeyID,
			awsSecretAccessKey,
			loadDotEnv,
			noCache,
			envVars,
			environment,
			gitSSH,
			[]string{"destroy"},
			[]string{"-auto-approve"},
			nil,
		)
	} else if runPlan {
		jobStackOut, jobStackErr = m.JobTgStack(ctx,
			stack,
			awsRegion,
			awsAccessKeyID,
			awsSecretAccessKey,
			loadDotEnv,
			noCache,
			envVars,
			environment,
			gitSSH,
			[]string{"plan"},
			[]string{},
			nil,
		)
	}

	return jobStackOut, jobStackErr
}
