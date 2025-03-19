package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	defaultTerraformVersion  = "1.7.1"
	defaultTerragruntVersion = "0.55.1"
	defaultImage             = "alpine"
	defaultImageTag          = "3.21.3"
)

// Terragrunt represents a structure that encapsulates operations related to Terragrunt,
// a tool for managing Terraform configurations. This struct can be extended with methods
// that perform various tasks such as executing commands in containers, managing directories,
// and other functionalities that facilitate the use of Terragrunt in a Dagger pipeline.
type Terragrunt struct {
	// Ctr is a Dagger container that can be used to run Terragrunt commands
	Ctr *dagger.Container
}

func New(
	// ctx is the context for the Dagger container.
	ctx context.Context,

	// tgVersion (image tag) to use from the official Terragrunt image.
	//
	// +optional
	tgVersion string,

	// tfVersion is the Terraform version to use.
	//
	// +optional
	tfVersion string,

	// Container is the custom container to use for Terragrunt operations.
	//
	// +optional
	container *dagger.Container,

	// Secrets are the secrets that will be used to run the Terragrunt commands.
	//
	// +optional
	secrets []*dagger.Secret,

	// EnvVars are the environment variables that will be used to run the Terragrunt commands.
	//
	// +optional
	envVars []string,
) (*Terragrunt, error) {
	if container != nil {
		return &Terragrunt{
			Ctr: container,
		}, nil
	}

	m := &Terragrunt{
		Ctr: dag.Container().
			From(fmt.Sprintf("%s:%s", defaultImage, defaultImageTag)),
	}

	if tgVersion == "" {
		tgVersion = defaultTerragruntVersion
	}

	if tfVersion == "" {
		tfVersion = defaultTerraformVersion
	}

	if len(secrets) > 0 {
		for _, secret := range secrets {
			secretName, secretErr := secret.Name(ctx)

			if secretErr != nil {
				return nil, fmt.Errorf("failed to get secret name: %w", secretErr)
			}

			m.Ctr = m.Ctr.WithSecretVariable(secretName, secret)
		}
	}

	if len(envVars) > 0 {
		for _, envVar := range envVars {
			parts := strings.Split(envVar, "=")

			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid environment variable: %s", envVar)
			}

			m.Ctr = m.Ctr.WithEnvVariable(parts[0], parts[1])
		}
	}

	// m.Ctr = m.WithTerragrunt(tgVersion)
	// m.Ctr = m.WithTerraform(tfVersion)

	return m, nil
}

// OpenTerminal returns a terminal
//
// It returns a terminal for the container.
// Arguments:
// - None.
// Returns:
// - *Terminal: The terminal for the container.
func (m *Terragrunt) OpenTerminal() *dagger.Container {
	return m.Ctr.Terminal()
}

// WithTerraform sets the Terraform version to use and installs it.
// It takes a version string as an argument and returns a pointer to a dagger.Container.
func (m *Terragrunt) WithTerraform(version string) *dagger.Container {
	tfInstallationCmd := getTerraformInstallationCommand(version)
	return m.
		Ctr.
		WithExec([]string{tfInstallationCmd}).
		WithExec([]string{"terraform", "--version"})
}

// WithTerragrunt sets the Terragrunt version to use and installs it.
// It takes a version string as an argument and returns a pointer to a dagger.Container.
func (m *Terragrunt) WithTerragrunt(version string) *dagger.Container {
	tgInstallationCmd := getTerragruntInstallationCommand(version)
	return m.
		Ctr.
		WithExec([]string{tgInstallationCmd}).
		WithExec([]string{"terragrunt", "--version"})
}

func getTerraformInstallationCommand(version string) string {
	installDir := "/usr/local/bin"
	installPath := filepath.Join(installDir, "terraform")
	command := fmt.Sprintf(`set -ex
curl -L https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip -o /tmp/terraform.zip
unzip /tmp/terraform.zip -d /tmp
mv /tmp/terraform %[2]s
chmod +x %[2]s
rm /tmp/terraform.zip`, version, installPath)

	return strings.TrimSpace(command)
}

func getTerragruntInstallationCommand(version string) string {
	installDir := "/usr/local/bin"
	installPath := filepath.Join(installDir, "terragrunt")
	command := fmt.Sprintf(`set -ex
curl -L https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_linux_amd64 -o %s
chmod +x %s`, version, installPath, installPath)

	return strings.TrimSpace(command)
}
