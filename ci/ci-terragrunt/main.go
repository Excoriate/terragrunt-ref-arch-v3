package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	defaultTerraformVersion  = "1.10.1"
	defaultTerragruntVersion = "0.77.7"
	defaultImage             = "alpine"
	defaultImageTag          = "3.21.3"
	defaultMntPath           = "/mnt"
	defaultNetrcRootPath     = "/root/.netrc"
	// specific terragrunt configuration details
	terraformPluginCachePath = "/root/.terraform.d/plugin-cache"
	terragruntCachePath      = "/root/.terragrunt-cache"
)

// Terragrunt represents a structure that encapsulates operations related to Terragrunt,
// a tool for managing Terraform configurations. This struct can be extended with methods
// that perform various tasks such as executing commands in containers, managing directories,
// and other functionalities that facilitate the use of Terragrunt in a Dagger pipeline.
type Terragrunt struct {
	// Ctr is a Dagger container that can be used to run Terragrunt commands
	Ctr *dagger.Container

	// Src is the source code for the Terragrunt project.
	Src *dagger.Directory
}

func New(
	// ctx is the context for the Dagger container.
	ctx context.Context,

	// imageURL is the URL of the image to use as the base container.
	// It should includes tags. E.g. "ghcr.io/devops-infra/docker-terragrunt:tf-1.9.5-ot-1.8.2-tg-0.67.4"
	// +optional
	imageURL string,

	// tgVersion (image tag) to use from the official Terragrunt image.
	//
	// +optional
	tgVersion string,

	// tfVersion is the Terraform version to use.
	//
	// +optional
	tfVersion string,

	// Ctr is the custom container to use for Terragrunt operations.
	//
	// +optional
	ctr *dagger.Container,

	// srcDir is the directory to mount as the source code.
	// +optional
	// +defaultPath="/"
	// +ignore=["*", "!**/*.hcl", "!**/*.tfvars", "!**/*.tfvars.json", "!**/*.tf"]
	srcDir *dagger.Directory,

	// Secrets are the secrets that will be used to run the Terragrunt commands.
	//
	// +optional
	secrets []*dagger.Secret,

	// EnvVars are the environment variables that will be used to run the Terragrunt commands.
	//
	// +optional
	envVars []string,
) (*Terragrunt, error) {
	if ctr != nil {
		mod := &Terragrunt{Ctr: ctr}
		mod, enVarError := mod.WithEnvVars(envVars)

		if enVarError != nil {
			return nil, enVarError
		}

		mod.WithSRC(ctx, defaultMntPath, srcDir)

		mod = mod.
			WithTerraformPluginCache().
			WithTerragruntCache().
			WithTerragruntProvidersCacheServerEnabled()

		return mod, nil
	}

	if imageURL != "" {
		mod := &Terragrunt{}
		mod.Ctr = dag.Container().From(imageURL)
		mod.WithSecrets(ctx, secrets)
		mod.WithSRC(ctx, defaultMntPath, srcDir)
		mod, enVarError := mod.WithEnvVars(envVars)

		if enVarError != nil {
			return nil, enVarError
		}

		return mod.
			WithTerraformPluginCache().
			WithTerragruntCache().
			WithTerragruntProvidersCacheServerEnabled(), nil
	}

	// We'll use the binary that should be downloaded from its source, or github repository.
	mod := &Terragrunt{}
	if tfVersion == "" {
		tfVersion = defaultTerraformVersion
	}

	if tgVersion == "" {
		tgVersion = defaultTerragruntVersion
	}

	defaultImageWithTag := fmt.Sprintf("%s:%s", defaultImage, defaultImageTag)

	mod.Ctr = dag.Container().From(defaultImageWithTag)

	mod.WithTerraform(tfVersion)
	mod.WithTerragrunt(tgVersion)

	mod.WithSecrets(ctx, secrets)
	mod.WithSRC(ctx, defaultMntPath, srcDir)
	mod, enVarError := mod.WithEnvVars(envVars)

	if enVarError != nil {
		return nil, enVarError
	}

	mod = mod.
		WithTerraformPluginCache().
		WithTerragruntCache().
		WithTerragruntProvidersCacheServerEnabled()

	return mod, nil
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

// WithSecrets adds secrets to the Terragrunt container, making them available as environment variables.
//
// This method allows secure injection of sensitive information into the container.
//
// Parameters:
//   - secrets: A slice of Dagger secrets to be mounted in the container
//
// Returns:
//   - The updated Terragrunt instance with secrets mounted
func (m *Terragrunt) WithSecrets(ctx context.Context, secrets []*dagger.Secret) *Terragrunt {
	for _, secret := range secrets {
		secretName, _ := secret.Name(ctx)
		m.Ctr = m.Ctr.WithSecretVariable(secretName, secret)
	}

	return m
}

// WithSRC mounts a source directory into the Terragrunt container.
//
// This method sets the working directory and mounts the provided directory,
// preparing the container for source code operations.
//
// Parameters:
//   - dir: A Dagger directory to be mounted in the container
//
// Returns:
//   - The updated Terragrunt instance with source directory mounted
func (m *Terragrunt) WithSRC(ctx context.Context, workdir string, dir *dagger.Directory) (*Terragrunt, error) {
	if workdir == "" {
		workdir = defaultMntPath
	} else {
		if workdir != defaultMntPath {
			workdir = filepath.Join(defaultMntPath, workdir)
		}
	}

	if err := isNonEmptyDaggerDir(ctx, dir); err != nil {
		return nil, fmt.Errorf("failed to validate the src/ directory passed: %w", err)
	}

	m.Ctr = m.Ctr.
		WithWorkdir(workdir).
		WithMountedDirectory(workdir, dir)

	m.Src = dir

	return m, nil
}

// WithTFPluginCache mounts a cache volume for Terraform plugins.
//
// This method sets up a cache directory for Terraform plugins and mounts it into the container.
//
// Returns:
//   - The updated Terragrunt instance with the plugin cache mounted
func (m *Terragrunt) WithTerraformPluginCache() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"mkdir", "-p", terraformPluginCachePath}).
		WithExec([]string{"chmod", "755", terraformPluginCachePath}).
		WithMountedCache(terraformPluginCachePath, dag.CacheVolume("terraform-plugin-cache")).
		WithEnvVariable("TF_PLUGIN_CACHE_DIR", terraformPluginCachePath)

	return m
}

// WithTerragruntProvidersCacheServerEnabled enables the Terragrunt providers cache server.
//
// This method sets the environment variable TG_PROVIDER_CACHE to "1" to enable the Terragrunt providers cache server.
//
// Returns:
//   - The updated Terragrunt instance with the providers cache server enabled
func (m *Terragrunt) WithTerragruntProvidersCacheServerEnabled() *Terragrunt {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PROVIDER_CACHE", "1")

	return m
}

// WithTerragruntCache mounts a cache volume for Terragrunt.
//
// This method sets up a cache directory for Terragrunt and mounts it into the container.
//
// Returns:
//   - The updated Terragrunt instance with the cache mounted
func (m *Terragrunt) WithTerragruntCache() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"mkdir", "-p", terragruntCachePath}).
		WithExec([]string{"chmod", "755", terragruntCachePath}).
		WithMountedCache(terragruntCachePath, dag.CacheVolume("terragrunt-cache"))

	return m
}

// WithEnvVars adds environment variables to the Terraformci container.
//
// This method allows setting multiple environment variables in key=value format.
// It performs validation to ensure each environment variable is correctly formatted.
//
// Parameters:
//   - envVars: A slice of environment variables in "KEY=VALUE" format
//
// Returns:
//   - The updated Terragrunt instance with environment variables set
//   - An error if any environment variable is incorrectly formatted
func (m *Terragrunt) WithEnvVars(envVars []string) (*Terragrunt, error) {
	if len(envVars) > 0 {
		for _, envVar := range envVars {
			trimmedEnvVar := strings.TrimSpace(envVar)
			if !strings.Contains(trimmedEnvVar, "=") {
				return nil, fmt.Errorf("environment variable must be in the format ENVARKEY=VALUE: %s", trimmedEnvVar)
			}

			parts := strings.Split(trimmedEnvVar, "=")
			if len(parts) != 2 {
				return nil, fmt.Errorf("environment variable must be in the format ENVARKEY=VALUE: %s", trimmedEnvVar)
			}

			m.Ctr = m.Ctr.WithEnvVariable(parts[0], parts[1])
		}
	}

	return m, nil
}

// WithToken adds a token to the Terragrunt container.
//
// This method adds a token to the container, making it available as an environment variable.
//
// Parameters:
//   - ctx: The context for the Dagger container.
//   - tokenValue: The value of the token to add to the container.
//
// Returns:
//   - The updated Terragrunt instance with the token added
func (m *Terragrunt) WithToken(ctx context.Context, tokenValue *dagger.Secret) *Terragrunt {
	return m.WithSecrets(ctx, []*dagger.Secret{tokenValue})
}

// WithTerraform sets the Terraform version to use and installs it.
// It takes a version string as an argument and returns a pointer to a dagger.Container.
func (m *Terragrunt) WithTerraform(version string) *Terragrunt {
	tfInstallationCmd := getTFInstallCmd(version)
	m.Ctr = m.Ctr.
		WithExec([]string{"/bin/sh", "-c", tfInstallationCmd}).
		WithExec([]string{"terraform", "--version"})

	return m
}

// WithTerragrunt sets the Terragrunt version to use and installs it.
// It takes a version string as an argument and returns a pointer to a dagger.Container.
func (m *Terragrunt) WithTerragrunt(version string) *Terragrunt {
	tgInstallationCmd := getTerragruntInstallationCommand(version)
	m.Ctr = m.Ctr.
		WithExec([]string{"/bin/sh", "-c", tgInstallationCmd}).
		WithExec([]string{"terragrunt", "--version"})

	return m
}

// WithNewNetrcFileGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Terragrunt) WithNewNetrcFileGitHub(
	username string,
	password string,
) *Terragrunt {
	machineCMD := "machine github.com\nlogin " + username + "\npassword " + password + "\n"

	m.Ctr = m.Ctr.WithNewFile(defaultNetrcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileAsSecretGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
// The argument 'password' is a secret that is not exposed in the logs.
func (m *Terragrunt) WithNewNetrcFileAsSecretGitHub(username string, password *dagger.Secret) *Terragrunt {
	passwordTxtValue, _ := password.Plaintext(context.Background())
	machineCMD := fmt.Sprintf("machine github.com\nlogin %s\npassword %s\n", username, passwordTxtValue)
	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile(defaultNetrcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Terragrunt) WithNewNetrcFileGitLab(
	username string,
	password string,
) *Terragrunt {
	machineCMD := "machine gitlab.com\nlogin " + username + "\npassword " + password + "\n"

	m.Ctr = m.Ctr.WithNewFile(defaultNetrcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileAsSecretGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
// The argument 'password' is a secret that is not exposed in the logs.
func (m *Terragrunt) WithNewNetrcFileAsSecretGitLab(username string, password *dagger.Secret) *Terragrunt {
	passwordTxtValue, _ := password.Plaintext(context.Background())
	machineCMD := fmt.Sprintf("machine gitlab.com\nlogin %s\npassword %s\n", username, passwordTxtValue)

	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile(defaultNetrcRootPath, machineCMD)

	return m
}

// WithSSHAuthSocket configures SSH authentication for Terraform modules with Git SSH sources.
//
// This function mounts an SSH authentication socket into the container, enabling Terraform to authenticate
// when fetching modules from Git repositories using SSH URLs (e.g., git@github.com:org/repo.git).
//
// Parameters:
//   - sshAuthSocket: The SSH authentication socket to mount in the container.
//   - socketPath: The path where the SSH socket will be mounted in the container.
//   - owner: Optional. The owner of the mounted socket in the container.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt instance with SSH authentication configured for Terraform modules.
func (m *Terragrunt) WithSSHAuthSocket(
	// sshAuthSocket is the SSH socket to use for authentication.
	sshAuthSocket *dagger.Socket,
	// socketPath is the path where the SSH socket will be mounted in the container.
	// +optional
	socketPath string,
	// owner is the owner of the mounted socket in the container. Optional parameter.
	// +optional
	owner string,
) *Terragrunt {
	// Default the socket path if not provided
	if socketPath == "" {
		socketPath = "/ssh-agent.sock"
	}

	socketOpts := dagger.ContainerWithUnixSocketOpts{}

	if owner != "" {
		socketOpts.Owner = owner
	}

	m.Ctr = m.Ctr.
		WithUnixSocket(socketPath, sshAuthSocket, socketOpts).
		WithEnvVariable("SSH_AUTH_SOCK", socketPath)

	return m
}

// Exec executes a given command within a dagger container.
// It returns the output of the command or an error if the command is invalid or fails to execute.
//
//nolint:lll,cyclop // It's okay, since the ignore pattern is included
func (m *Terragrunt) Exec(
	// ctx is the context to use when executing the command.
	// +optional
	//nolint:contextcheck // It's okay, since the ignore pattern is included.
	ctx context.Context,
	// command is the terragrunt command to execute. It's the actual command that comes after 'terragrunt'
	command string,
	// args are the arguments to pass to the command.
	// +optional
	args []string,
	// autoApprove is the flag to auto approve the command.
	// +optional
	autoApprove bool,
	// src is the source directory that includes the source code.
	// +defaultPath="/"
	// +ignore=["*", "!**/*.hcl", "!**/*.tfvars", "!**/*.tfvars.json", "!**/*.tf"]
	src *dagger.Directory,
	// module is the module to execute or the terragrunt configuration where the terragrunt.hcl file is located.
	// +optional
	module string,
	// envVars is the environment variables to pass to the container.
	// +optional
	envVars []string,
	// secrets is the secrets to pass to the container.
	// +optional
	secrets []*dagger.Secret,
	// tfToken is the Terraform registry token to pass to the container.
	// +optional
	tfToken *dagger.Secret,
	// ghToken is the GitHub token to pass to the container.
	// +optional
	ghToken *dagger.Secret,
	// glToken is the GitLab token to pass to the container.
	// +optional
	glToken *dagger.Secret,
	// gitSshSocket is the Git SSH socket that can be forwarded from the host, to allow private git support through SSH
	// +optional
	gitSshSocket *dagger.Socket,
) (*dagger.Container, error) {
	tgCmd := []string{command}

	if len(args) > 0 {
		tgCmd = append(tgCmd, args...)
	}

	if autoApprove && (command == "apply" || command == "destroy") {
		tgCmd = append(tgCmd, "--auto-approve")
	}

	if src != nil {
		modWithSrc, err := m.WithSRC(ctx, module, src)

		if err != nil {
			return nil, WrapError(err, "failed to mount the source directory")
		}

		m = modWithSrc
	}

	if envVars != nil {
		modWithEnvVars, err := m.WithEnvVars(envVars)

		if err != nil {
			return nil, WrapError(err, "failed to set the environment variables")
		}

		m = modWithEnvVars
	}

	if secrets != nil {
		m.WithSecrets(ctx, secrets)
	}

	if tfToken != nil {
		m.WithToken(ctx, tfToken)
	}

	if ghToken != nil {
		m.WithToken(ctx, ghToken)
	}

	if glToken != nil {
		m.WithToken(ctx, glToken)
	}

	if gitSshSocket != nil {
		m.WithSSHAuthSocket(gitSshSocket, "/ssh-agent.sock", "root")
	}

	return m.Ctr.WithExec(tgCmd), nil
}

func getTFInstallCmd(tfVersion string) string {
	installDir := "/usr/local/bin/terraform"
	command := fmt.Sprintf(`apk add --no-cache curl unzip &&
	curl -L https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip -o /tmp/terraform.zip &&
	unzip /tmp/terraform.zip -d /tmp &&
	mv /tmp/terraform %[2]s &&
	chmod +x %[2]s &&
	rm /tmp/terraform.zip`, tfVersion, installDir)

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

func isNonEmptyDaggerDir(ctx context.Context, dir *dagger.Directory) error {
	if dir == nil {
		return fmt.Errorf("dagger directory cannot be nil")
	}

	entries, err := dir.Entries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get entries from the dagger directory passed: %w", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("no entries found in the dagger directory passed")
	}

	return nil
}
