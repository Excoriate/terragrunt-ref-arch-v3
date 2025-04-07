package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	// Default version for binaries
	defaultTerraformVersion  = "1.10.1"
	defaultTerragruntVersion = "0.77.7"
	defaultImage             = "alpine"
	defaultImageTag          = "3.21.3"
	defaultMntPath           = "/mnt"
	defaultBinary            = "terragrunt"
	// Default Ref-Arch configuration details
	// TODO: change this later, according to your specific architecture.
	defaultRefArchEnv    = "global"
	defaulttRefArchLayer = "dni"
	defaulttRefArchUnit  = "dni_generator"
	// Default for AWS
	defaultAWSRegion = "us-east-1"
	// Configuration
	configRefArchRootPath          = "infra/terragrunt"
	configterraformPluginCachePath = "/root/.terraform.d/plugin-cache"
	configterragruntCachePath      = "/root/.terragrunt-cache"
	configNetrcRootPath            = "/root/.netrc"
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
	// +ignore=["*", "!**/*.hcl", "!**/*.tfvars", "!**/.git/**", "!**/*.tfvars.json", "!**/*.tf", "!*.env", "!*.envrc", "!*.envrc"]
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

		mod.WithSRC(ctx, defaultMntPath, srcDir, false)

		mod = mod.CommonSetup(tfVersion, tgVersion)

		return mod, nil
	}

	if imageURL != "" {
		mod := &Terragrunt{}
		mod.Ctr = dag.Container().From(imageURL)
		mod.WithSecrets(ctx, secrets)
		mod.WithSRC(ctx, defaultMntPath, srcDir, false)
		mod, enVarError := mod.WithEnvVars(envVars)

		if enVarError != nil {
			return nil, enVarError
		}

		mod = mod.CommonSetup(tfVersion, tgVersion)
		return mod, nil
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

	mod.Ctr = dag.
		Container().
		From(defaultImageWithTag)

	mod.WithSecrets(ctx, secrets)
	mod.WithSRC(ctx, defaultMntPath, srcDir, false)
	mod, enVarError := mod.WithEnvVars(envVars)

	if enVarError != nil {
		return nil, enVarError
	}

	mod = mod.CommonSetup(tfVersion, tgVersion)

	return mod, nil
}

// CommonSetup configures the Terragrunt container with common dependencies and settings.
// It installs Git, sets up specified Terraform and Terragrunt versions, and configures
// cache volumes for Terraform plugins and Terragrunt operations. It also enables the
// Terragrunt provider cache server.
//
// Parameters:
//   - tfVersion: The version of Terraform to install.
//   - tgVersion: The version of Terragrunt to install.
//
// Returns:
//   - The updated Terragrunt instance with common setup applied.
func (m *Terragrunt) CommonSetup(tfVersion, tgVersion string) *Terragrunt {
	m = m.
		WithGitPkgInstalled().
		WithTerraform(tfVersion).
		WithTerragrunt(tgVersion).
		WithTerraformPluginCache().
		WithTerragruntCache().
		WithTerragruntProvidersCacheServerEnabled()

	return m
}

// OpenTerminal returns a terminal
//
// It returns a terminal for the container.
// Arguments:
// - ctx: The context for the operation.
// - srcDir: The source directory to be mounted in the container. If nil, the default source directory is used.
// - loadEnvFiles: A boolean to load the environment files.
// Returns:
// - *dagger.Container: The terminal for the container.
func (m *Terragrunt) OpenTerminal(
	// ctx is the context for the operation.
	// +optional
	ctx context.Context,
	// srcDir is the source directory to be mounted in the container.
	// +optional
	srcDir *dagger.Directory,
	// loadEnvFiles is a boolean to load the environment files.
	// +optional
	loadEnvFiles bool,
) *dagger.Container {
	if srcDir == nil {
		srcDir = m.Src
	}

	if loadEnvFiles {
		m.WithDotEnvFile(ctx, m.Src)
	}

	return m.
		Ctr.
		Terminal()
}

// WithSRC mounts a source directory into the Terragrunt container.
//
// This method sets the working directory and mounts the provided directory,
// preparing the container for source code operations.
//
// Parameters:
//   - dir: A Dagger directory to be mounted in the container
//   - printPaths: A boolean to print the paths of the mounted directories.
//
// Returns:
//   - The updated Terragrunt instance with source directory mounted
func (m *Terragrunt) WithSRC(
	// ctx is the context for the Dagger container.
	ctx context.Context,
	// workdir is the working directory to set in the container.
	// +optional
	workdir string,
	// dir is the directory to mount in the container.
	dir *dagger.Directory,
	// printPaths is a boolean to print the paths of the mounted directories.
	// +optional
	printPaths bool) (*Terragrunt, error) {
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

	if printPaths {
		m.Ctr = m.Ctr.
			WithExec([]string{"echo", "Workdir: " + workdir}).
			WithExec([]string{"echo", "Mounted directory: " + workdir}).
			WithExec([]string{"ls", "-la", workdir})
	}

	return m, nil
}

// WithGitPkgInstalled installs the Git package in the container.
//
// This method adds the Git package to the container's package manager.
//
// Returns:
//   - The updated Terragrunt instance with Git installed
func (m *Terragrunt) WithGitPkgInstalled() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "git"})

	return m
}

// WithTFPluginCache mounts a cache volume for Terraform plugins.
//
// This method sets up a cache directory for Terraform plugins and mounts it into the container.
//
// Returns:
//   - The updated Terragrunt instance with the plugin cache mounted
func (m *Terragrunt) WithTerraformPluginCache() *Terragrunt {
	m.Ctr = m.Ctr.
		WithExec([]string{"mkdir", "-p", configterraformPluginCachePath}).
		WithExec([]string{"chmod", "755", configterraformPluginCachePath}).
		WithMountedCache(configterraformPluginCachePath, dag.CacheVolume("terraform-plugin-cache")).
		WithEnvVariable("TF_PLUGIN_CACHE_DIR", configterraformPluginCachePath)

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

// WithCacheBuster enables the cache buster for the container.
//
// This method sets the environment variable DAGGER_APT_CACHE_BUSTER to a unique value based on the current time.
//
// Returns:
//   - The updated Terragrunt instance with the cache buster enabled
func (m *Terragrunt) WithCacheBuster() *Terragrunt {
	m.Ctr = m.Ctr.
		WithEnvVariable("DAGGER_OPT_CACHE_BUSTER", fmt.Sprintf("%d", time.Now().Truncate(24*time.Hour).Unix()))

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
		WithExec([]string{"mkdir", "-p", configterragruntCachePath}).
		WithExec([]string{"chmod", "755", configterragruntCachePath}).
		WithMountedCache(configterragruntCachePath, dag.CacheVolume("terragrunt-cache"))

	return m
}

// WithSecrets mounts secrets into the container.
//
// This method mounts secrets into the container for use by Terragrunt.
//
// Parameters:
//   - ctx: The context for the operation
//   - secrets: A slice of dagger.Secret instances to be mounted
//
// Returns:
//   - The updated Terragrunt instance with secrets mounted
func (m *Terragrunt) WithSecrets(ctx context.Context, secrets []*dagger.Secret) *Terragrunt {
	for _, secret := range secrets {
		// FIXME: This is suitable when secrets are created within dagger. Assumming there's a name set.
		secretName, _ := secret.Name(ctx)
		m.Ctr = m.Ctr.WithSecretVariable(secretName, secret)
	}

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

	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

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
	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

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

	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

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
	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

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

// WithAWSCredentials sets the AWS credentials and region in the container.
//
// This method sets the AWS credentials and region in the container, making them available as environment variables.
// It also mounts the AWS credentials as secrets into the container.
//
// Parameters:
//   - ctx: The context for the Dagger container.
//   - awsAccessKeyID: The AWS access key ID.
//   - awsSecretAccessKey: The AWS secret access key.
//   - awsRegion: The AWS region.
//
// Returns:
//   - *Terragrunt: The updated Terragrunt instance with AWS credentials and region set
func (m *Terragrunt) WithAWSCredentials(
	// ctx is the context for the Dagger container.
	// +optional
	ctx context.Context,
	// awsAccessKeyID is the AWS access key ID.
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key.
	awsSecretAccessKey *dagger.Secret,
	// awsRegion is the AWS region.
	// +optional
	awsRegion string,
) *Terragrunt {
	awsRegion = getDefaultAWSRegionIfNotSet(awsRegion)

	m.Ctr = m.Ctr.
		WithEnvVariable("AWS_REGION", awsRegion).
		WithSecretVariable("AWS_ACCESS_KEY_ID", awsAccessKeyID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", awsSecretAccessKey)

	// m = m.WithSecrets(ctx, []*dagger.Secret{awsAccessKeyID, awsSecretAccessKey})

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
	// binary is the binary to execute, Possible and valid options are 'terragrunt', 'terraform'
	// +optional
	binary string,
	// command is the terragrunt command to execute. It's the actual command that comes after 'terragrunt'
	command string,
	// args are the arguments to pass to the command.
	// +optional
	args []string,
	// autoApprove is the flag to auto approve the command.
	// +optional
	autoApprove bool,
	// src is the source directory that includes the source code.
	src *dagger.Directory,
	// tgUnitPath is the path to the terragrunt unit to execute.
	// +optional
	tgUnitPath string,
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
	// printPaths is a boolean to print the paths of the mounted directories.
	// +optional
	printPaths bool,
	// loadEnvFiles is a boolean to load the environment files.
	// +optional
	loadEnvFiles bool,
) (*dagger.Container, error) {
	cmdEntryPoint := getDefaultBinaryIfANotSet(binary)

	tgCmd := []string{cmdEntryPoint, command}

	if len(args) > 0 {
		tgCmd = append(tgCmd, args...)
	}

	if autoApprove && (command == "apply" || command == "destroy") {
		tgCmd = append(tgCmd, "--auto-approve")
	}

	if tgUnitPath != "" {
		tgCmd = append(tgCmd, "--working-dir", tgUnitPath)
	}

	_, srcErr := m.WithSRC(ctx, "", src, printPaths)

	if srcErr != nil {
		return nil, WrapError(srcErr, "failed to mount the source directory")
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

	if loadEnvFiles {
		// At this point, m.Src is either the directory passed, or the one set at the constructor.
		ctrWithDotEnvLoaded, err := m.WithDotEnvFile(ctx, m.Src)

		if err != nil {
			return nil, WrapError(err, "failed to load .env files")
		}

		m = ctrWithDotEnvLoaded
	}

	return m.Ctr.WithExec(tgCmd), nil
}

// WithDotEnvFile loads and processes environment variables from .env files in the provided directory.
//
// This method finds all .env files in the given directory, reads their contents, and sets
// environment variables in the Terragrunt container. Files containing "secret" in their name
// will have their values added as secret variables rather than regular environment variables.
//
// The method supports standard .env file formats with KEY=VALUE pairs on each line.
// Comments (lines starting with #) and empty lines are ignored. Values can be optionally
// quoted with single or double quotes, which will be automatically removed.
//
// Parameters:
//   - ctx: Context for the Dagger operations
//   - src: Directory containing the .env files to process
//
// Returns:
//   - *Terragrunt: The updated Terragrunt instance with environment variables set
//   - error: An error if file reading or parsing fails
func (m *Terragrunt) WithDotEnvFile(ctx context.Context, src *dagger.Directory) (*Terragrunt, error) {
	if src == nil {
		return nil, NewError("failed to load .env file, the source directory is nil")
	}

	// Check if there's any dotenv file on the source directory passed, or set.
	entries, err := src.Entries(ctx)
	if err != nil {
		return nil, WrapErrorf(err, "failed to list files in source directory")
	}

	foundDotEnvFiles := []string{}
	for _, entry := range entries {
		if strings.HasSuffix(entry, ".env") {
			foundDotEnvFiles = append(foundDotEnvFiles, entry)
		}
	}

	if len(foundDotEnvFiles) == 0 {
		return nil, NewError("No .env files found when inspecting the source directory")
	}

	dotEnvFilesInSrc, srcError := src.Glob(ctx, "*.env")

	if srcError != nil {
		return nil, WrapErrorf(srcError, "failed to glob dot env files")
	}

	ctrWithDotEnvFiles, dotEnvFilesParseErr := parseDotEnvFiles(ctx, m.Ctr, src, dotEnvFilesInSrc)

	if dotEnvFilesParseErr != nil {
		return nil, WrapErrorf(dotEnvFilesParseErr, "failed to parse dot env files")
	}

	m.Ctr = ctrWithDotEnvFiles

	return m, nil
}

// parseDotEnvFiles processes .env files found by WithDotEnvFile.
// It handles basic .env syntax including comments (#), empty lines,
// KEY=VALUE pairs, whitespace trimming, and basic quote removal (' or ").
func parseDotEnvFiles(ctx context.Context, container *dagger.Container, src *dagger.Directory, envFiles []string) (*dagger.Container, error) {
	for _, file := range envFiles {
		fileContent, err := src.File(file).Contents(ctx)
		if err != nil {
			// Wrap error for better context
			return nil, fmt.Errorf("failed to read dot env file '%s': %w", file, err)
		}

		lines := strings.Split(fileContent, "\n")

		for lineNum, line := range lines {
			trimmedLine := strings.TrimSpace(line)

			// Skip empty lines and comments
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
				continue
			}

			// Split line into key/value pair by the first '='
			parts := strings.SplitN(trimmedLine, "=", 2)
			if len(parts) != 2 {
				// Return error for lines without '='
				return nil, fmt.Errorf("invalid format in file '%s' on line %d: '%s'", file, lineNum+1, trimmedLine)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Check for empty key
			if key == "" {
				return nil, fmt.Errorf("empty key found in file '%s' on line %d: '%s'", file, lineNum+1, trimmedLine)
			}

			// Trim surrounding quotes (basic handling)
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}

			// Determine if it's a secret based on filename
			isSecret := strings.Contains(file, "secret")

			if isSecret {
				// Use a distinct name for the Dagger secret object itself
				secretName := fmt.Sprintf("%s_secret_%s", key, file)
				container = container.WithSecretVariable(key, dag.SetSecret(secretName, value))
			} else {
				container = container.WithEnvVariable(key, value)
			}
		}
	}

	return container, nil
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

func getDefaultAWSRegionIfNotSet(awsRegion string) string {
	if awsRegion == "" {
		return defaultAWSRegion
	}

	return awsRegion
}

func getDefaultBinaryIfANotSet(binary string) string {
	if binary == "" {
		return defaultBinary
	}

	return binary
}

func getTerragruntUnitPath(env, layer, unit string) string {
	if env == "" {
		env = defaultRefArchEnv
	}

	if layer == "" {
		layer = defaulttRefArchLayer
	}

	if unit == "" {
		unit = defaulttRefArchUnit
	}

	return filepath.Join(configRefArchRootPath, env, layer, unit)
}
