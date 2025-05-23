package main

import (
	"context"
	"dagger/infra/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// Default version for binaries
	defaultTerraformVersion  = "1.11.3"
	defaultTerragruntVersion = "0.80.2"
	defaultImage             = "alpine"
	defaultImageTag          = "3.21.3"
	defaultMntPath           = "/mnt"
	defaultBinary            = "terragrunt"
	// Default Ref-Arch configuration details
	defaultRefArchEnv    = "dev"
	defaulttRefArchLayer = "non-distributable"
	defaulttRefArchUnit  = "random-string-generator"
	// Default for AWS
	defaultAWSRegion              = "eu-west-1"
	defaultAWSOidcTokenSecretName = "AWS_OIDC_TOKEN"
	// TODO: Change this to the actual region based on your own convention
	defaultRemoteStateRegion = "us-east-1"
	// Configuration
	configRefArchRootPath                  = "infra/terragrunt"
	configRefArchATerraformModulesRootPath = "infra/terraform/modules"
	configterraformPluginCachePath         = "/root/.terraform.d/plugin-cache"
	configterragruntCachePath              = "/root/.terragrunt-cache"
	configNetrcRootPath                    = "/root/.netrc"
	// tfConfig
	// TODO: Change this to the actual bucket and lock table names
	remoteStateDefaultBucketNamingConvention    = "terraform-state-makemyinfra"
	remoteStateDefaultLockTableNamingConvention = "terraform-state-makemyinfra"
)

// Terragrunt represents a structure that encapsulates operations related to Terragrunt,
// a tool for managing Terraform configurations. This struct can be extended with methods
// that perform various tasks such as executing commands in containers, managing directories,
// and other functionalities that facilitate the use of Terragrunt in a Dagger pipeline.
type Infra struct {
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
	// +ignore=["*", "!**/*.hcl", "!**/*.tfvars", "!**/.git/**", "!**/*.tfvars.json", "!**/*.tf", "!*.env"]
	srcDir *dagger.Directory,

	// EnvVars are the environment variables that will be used to run the Terragrunt commands.
	//
	// +optional
	envVars []string,
) (*Infra, error) {
	if ctr != nil {
		mod := &Infra{Ctr: ctr}
		mod, enVarError := mod.WithEnvVars(envVars)
		if enVarError != nil {
			return nil, WrapErrorf(enVarError, "failed to initialise dagger module with environment variables")
		}

		modWithSRC, modWithSRCError := mod.WithSRC(ctx, defaultMntPath, srcDir)
		if modWithSRCError != nil {
			return nil, WrapErrorf(modWithSRCError, "failed to initialise dagger module with source directory")
		}

		mod = modWithSRC
		mod = mod.CommonSetup(tfVersion, tgVersion)

		return mod, nil
	}

	if imageURL != "" {
		mod := &Infra{}
		mod.Ctr = dag.Container().From(imageURL)
		modWithSRC, modWithSRCError := mod.WithSRC(ctx, defaultMntPath, srcDir)
		if modWithSRCError != nil {
			return nil, WrapErrorf(modWithSRCError, "failed to initialise dagger module with source directory")
		}

		mod = modWithSRC
		mod, enVarError := mod.WithEnvVars(envVars)

		if enVarError != nil {
			return nil, WrapErrorf(enVarError, "failed to initialise dagger module with environment variables")
		}

		mod = mod.CommonSetup(tfVersion, tgVersion)
		return mod, nil
	}

	// We'll use the binary that should be downloaded from its source, or github repository.
	mod := &Infra{}
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

	modWithSRC, modWithSRCError := mod.WithSRC(ctx, defaultMntPath, srcDir)
	if modWithSRCError != nil {
		return nil, WrapErrorf(modWithSRCError, "failed to initialise dagger module with source directory")
	}

	mod = modWithSRC
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
func (m *Infra) CommonSetup(tfVersion, tgVersion string) *Infra {
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
func (m *Infra) OpenTerminal(
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
func (m *Infra) WithSRC(
	// ctx is the context for the Dagger container.
	ctx context.Context,
	// workdir is the working directory to set in the container.
	// +optional
	workdir string,
	// dir is the directory to mount in the container.
	dir *dagger.Directory,
) (*Infra, error) {
	if workdir == "" {
		workdir = defaultMntPath
	} else {
		if workdir != defaultMntPath {
			workdir = filepath.Join(defaultMntPath, workdir)
		}
	}

	if err := isNonEmptyDaggerDir(ctx, dir); err != nil {
		return nil, WrapErrorf(err, "failed to validate the src/ directory passed")
	}

	m.Ctr = m.Ctr.
		WithWorkdir(workdir).
		WithMountedDirectory(workdir, dir)

	m.Src = dir

	return m, nil
}

// WithGitPkgInstalled installs the Git package in the container.
//
// This method adds the Git package to the container's package manager.
//
// Returns:
//   - The updated Terragrunt instance with Git installed
func (m *Infra) WithGitPkgInstalled() *Infra {
	m.Ctr = m.Ctr.
		WithExec([]string{"apk", "add", "git"}).
		WithExec([]string{"apk", "add", "openssh"})

	return m
}

// WithTFPluginCache mounts a cache volume for Terraform plugins.
//
// This method sets up a cache directory for Terraform plugins and mounts it into the container.
//
// Returns:
//   - The updated Terragrunt instance with the plugin cache mounted
func (m *Infra) WithTerraformPluginCache() *Infra {
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
func (m *Infra) WithTerragruntProvidersCacheServerEnabled() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PROVIDER_CACHE", "1")

	return m
}

// WithTerragruntProvidersCacheServerDisabled disables the Terragrunt providers cache server.
//
// This method removes the environment variable TG_PROVIDER_CACHE from the container.
//
// Returns:
//   - The updated Terragrunt instance with the providers cache server disabled
func (m *Infra) WithTerragruntProvidersCacheServerDisabled() *Infra {
	m.Ctr = m.Ctr.
		WithoutEnvVariable("TG_PROVIDER_CACHE").
		WithEnvVariable("TG_PROVIDER_CACHE", "0")

	return m
}

// WithRegistriesToCacheProvidersFrom adds extra registries to cache providers from.
//
// This function appends the provided registries to the default list of registries in the Terragrunt configuration.
// By default, the Terragrunt provider's cache only caches registry.terraform.io and registry.opentofu.org.
//
// Parameters:
//   - registries: A slice of strings representing the registries to cache providers from.
//
// Returns:
//   - *Infra: The updated Infra instance with the extra registries to cache providers from.
func (m *Infra) WithRegistriesToCacheProvidersFrom(
	// registries is a slice of strings representing the registries to cache providers from.
	registries []string,
) *Infra {
	defaultRegistries := []string{
		"registry.terraform.io",
		"registry.opentofu.org",
	}

	registries = append(defaultRegistries, registries...)
	registryNames := strings.Join(registries, ",")

	m.Ctr = m.Ctr.
		WithoutEnvVariable("TG_PROVIDER_CACHE_REGISTRY_NAMES").
		WithEnvVariable("TG_PROVIDER_CACHE_REGISTRY_NAMES", registryNames)

	return m
}

// WithCacheBuster enables the cache buster for the container.
//
// This method sets the environment variable DAGGER_APT_CACHE_BUSTER to a unique value based on the current time.
//
// Returns:
//   - The updated Infra instance with the cache buster enabled
func (m *Infra) WithCacheBuster() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("DAGGER_OPT_CACHE_BUSTER", fmt.Sprintf("%d", time.Now().Truncate(24*time.Hour).Unix()))

	return m
}

// WithTerragruntCache mounts a cache volume for Terragrunt.
//
// This method sets up a cache directory for Terragrunt and mounts it into the container.
//
// Returns:
//   - The updated Infra instance with the cache mounted
func (m *Infra) WithTerragruntCache() *Infra {
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
//   - The updated Infra instance with secrets mounted
func (m *Infra) WithSecrets(ctx context.Context, secrets []*dagger.Secret) *Infra {
	for _, secret := range secrets {
		// FIXME: This is suitable when secrets are created within dagger. Assumming there's a name set.
		secretName, _ := secret.Name(ctx)

		m.Ctr = m.
			Ctr.
			WithSecretVariable(secretName, secret)
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
//   - The updated Infra instance with environment variables set
//   - An error if any environment variable is incorrectly formatted
func (m *Infra) WithEnvVars(envVars []string) (*Infra, error) {
	envVarsDagger, err := getEnvVarsDaggerFromSlice(envVars)

	if err != nil {
		return nil, err
	}

	for _, envVar := range envVarsDagger {
		m.Ctr = m.Ctr.WithEnvVariable(envVar.Key, envVar.Value)
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
//   - The updated Infra instance with the token added
func (m *Infra) WithToken(ctx context.Context, tokenValue *dagger.Secret) *Infra {
	return m.WithSecrets(ctx, []*dagger.Secret{tokenValue})
}

// WithTerraform sets the Terraform version to use and installs it.
// It takes a version string as an argument and returns a pointer to a dagger.Container.
func (m *Infra) WithTerraform(version string) *Infra {
	tfInstallationCmd := getTFInstallCmd(version)
	m.Ctr = m.Ctr.
		WithExec([]string{"/bin/sh", "-c", tfInstallationCmd}).
		WithExec([]string{"terraform", "--version"})

	return m
}

// WithTerragrunt sets the Terragrunt version to use and installs it.
// It takes a version string as an argument and returns a pointer to a dagger.Container.
func (m *Infra) WithTerragrunt(version string) *Infra {
	tgInstallationCmd := getTerragruntInstallationCommand(version)
	m.Ctr = m.Ctr.
		WithExec([]string{"/bin/sh", "-c", tgInstallationCmd}).
		WithExec([]string{"terragrunt", "--version"})

	return m
}

// WithNewNetrcFileGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Infra) WithNewNetrcFileGitHub(
	username string,
	password string,
) *Infra {
	machineCMD := "machine github.com\nlogin " + username + "\npassword " + password + "\n"

	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileAsSecretGitHub creates a new .netrc file with the GitHub credentials.
//
// The .netrc file is created in the root directory of the container.
// The argument 'password' is a secret that is not exposed in the logs.
func (m *Infra) WithNewNetrcFileAsSecretGitHub(username string, password *dagger.Secret) *Infra {
	passwordTxtValue, _ := password.Plaintext(context.Background())
	machineCMD := fmt.Sprintf("machine github.com\nlogin %s\npassword %s\n", username, passwordTxtValue)
	//nolint:exhaustruct // This is a method that is used to set the base image and version.
	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
func (m *Infra) WithNewNetrcFileGitLab(
	username string,
	password string,
) *Infra {
	machineCMD := "machine gitlab.com\nlogin " + username + "\npassword " + password + "\n"

	m.Ctr = m.Ctr.WithNewFile(configNetrcRootPath, machineCMD)

	return m
}

// WithNewNetrcFileAsSecretGitLab creates a new .netrc file with the GitLab credentials.
//
// The .netrc file is created in the root directory of the container.
// The argument 'password' is a secret that is not exposed in the logs.
func (m *Infra) WithNewNetrcFileAsSecretGitLab(username string, password *dagger.Secret) *Infra {
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
//   - *Infra: The updated Infra instance with SSH authentication configured for Terraform modules.
func (m *Infra) WithSSHAuthSocket(
	// sshAuthSocket is the SSH socket to use for authentication.
	sshAuthSocket *dagger.Socket,
	// socketPath is the path where the SSH socket will be mounted in the container.
	// +optional
	socketPath string,
	// owner is the owner of the mounted socket in the container. Optional parameter.
	// +optional
	owner string,
	// enableGitlabKnownHosts adds the Gitlab known hosts to the container.
	// +optional
	enableGitlabKnownHosts bool,
	// enableGithubKnownHosts adds the Github known hosts to the container.
	// +optional
	enableGithubKnownHosts bool,
) *Infra {
	// Default the socket path if not provided
	if socketPath == "" {
		socketPath = "/var/run/host.sock"
	}

	socketOpts := dagger.ContainerWithUnixSocketOpts{}

	if owner != "" {
		socketOpts.Owner = owner
	}

	// Ensure .ssh directory exists before running ssh-keyscan
	m.Ctr = m.Ctr.WithExec([]string{"mkdir", "-p", "/root/.ssh"})

	if enableGitlabKnownHosts {
		m.Ctr = m.Ctr.
			WithExec([]string{"sh", "-c", "ssh-keyscan gitlab.com >> /root/.ssh/known_hosts"})
	}

	if enableGithubKnownHosts {
		m.Ctr = m.Ctr.
			WithExec([]string{"sh", "-c", "ssh-keyscan github.com >> /root/.ssh/known_hosts"})
	}

	m.Ctr = m.Ctr.
		WithExec([]string{"chmod", "600", "/root/.ssh/known_hosts"})

	m.Ctr = m.Ctr.WithUnixSocket(socketPath, sshAuthSocket, socketOpts).
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
//   - *Infra: The updated Infra instance with AWS credentials and region set
func (m *Infra) WithAWSKeys(
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
) *Infra {
	awsRegion = getDefaultAWSRegionIfNotSet(awsRegion)

	m.Ctr = m.Ctr.
		WithEnvVariable("AWS_REGION", awsRegion).
		WithSecretVariable("AWS_ACCESS_KEY_ID", awsAccessKeyID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", awsSecretAccessKey)

	return m
}

// WithAWSOIDC sets the AWS OIDC credentials in the container.
//
// This method sets the AWS OIDC credentials in the container, making them available as environment variables.
// It also mounts the AWS OIDC credentials as secrets into the container.
func (m *Infra) WithAWSOIDC(
	// roleARN is the ARN of the IAM role to assume.
	roleARN string,
	// oidcToken is the Dagger Secret containing the OIDC JWT token from GitLab.
	oidcToken *dagger.Secret,
	// oidcTokenName is the name of the secret containing the OIDC JWT token from GitLab.
	// +optional
	oidcTokenName string,
	// awsRegion is the AWS region.
	// +optional
	awsRegion string,
	// awsRoleSessionName is an optional name for the assumed role session.
	// +optional
	awsRoleSessionName string,
) *Infra {
	awsRegion = getDefaultAWSRegionIfNotSet(awsRegion)

	if oidcTokenName == "" {
		oidcTokenName = defaultAWSOidcTokenSecretName
	}

	if awsRoleSessionName == "" {
		awsRoleSessionName = fmt.Sprintf("terragrunt-dagger-%s", uuid.New().String())
	}

	oidcTokenPath := "run/secrets/" + oidcTokenName

	m.Ctr = m.Ctr.
		WithEnvVariable("AWS_REGION", awsRegion).
		WithEnvVariable("AWS_ROLE_ARN", roleARN).
		WithEnvVariable("AWS_ROLE_SESSION_NAME", awsRoleSessionName).
		WithEnvVariable("AWS_WEB_IDENTITY_TOKEN_FILE", oidcTokenPath).
		// cleaning —if set— aws keys.
		WithoutEnvVariable("AWS_ACCESS_KEY_ID").
		WithoutEnvVariable("AWS_SECRET_ACCESS_KEY").
		WithoutEnvVariable("AWS_SESSION_TOKEN").
		WithSecretVariable(oidcTokenName, oidcToken)

	return m
}

// WithGitlabToken sets the GitLab token in the container.
//
// This method sets the GitLab token in the container, making it available as an environment variable.
//
// Parameters:
//   - ctx: The context for the Dagger container.
func (m *Infra) WithGitlabToken(ctx context.Context, token *dagger.Secret) *Infra {
	m.Ctr = m.Ctr.
		WithSecretVariable("GITLAB_TOKEN", token)

	return m
}

// WithGitHubToken sets the GitHub token in the container.
//
// This method sets the GitHub token in the container, making it available as an environment variable.
//
// Parameters:
//   - ctx: The context for the Dagger container.
func (m *Infra) WithGitHubToken(ctx context.Context, token *dagger.Secret) *Infra {
	m.Ctr = m.Ctr.
		WithSecretVariable("GITHUB_TOKEN", token)

	return m
}

// WithTerraformToken sets the Terraform token in the container.
//
// This method sets the Terraform token in the container, making it available as an environment variable.
//
// Parameters:
//   - ctx: The context for the Dagger container.
//   - token: The Terraform token to set.
func (m *Infra) WithTerraformToken(ctx context.Context, token *dagger.Secret) *Infra {
	m.Ctr = m.Ctr.
		WithSecretVariable("TF_TOKEN", token)

	return m
}

// WithTerragruntLogLevel sets the Terragrunt log level in the container.
//
// This method sets the Terragrunt log level in the container, making it available as an environment variable.
//
// Parameters:
//   - level: The log level to set.
func (m *Infra) WithTerragruntLogLevel(level string) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_LOG_LEVEL", level)

	return m
}

// WithTerragruntNonInteractive sets the Terragrunt non interactive option in the container.
//
// This method sets the Terragrunt non interactive option in the container, making it available as an environment variable.
//
// Parameters:
//   - level: The log level to set.
func (m *Infra) WithTerragruntNonInteractive() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_NON_INTERACTIVE", "true")

	return m
}

// WithTerragruntNoColor sets the Terragrunt no color option in the container.
//
// This method sets the Terragrunt no color option in the container, making it available as an environment variable.
//
// Parameters:
//   - level: The log level to set.
func (m *Infra) WithTerragruntNoColor() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_NO_COLOR", "true")

	return m
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
//   - *Infra: The updated Infra instance with environment variables set
//   - error: An error if file reading or parsing fails
func (m *Infra) WithDotEnvFile(ctx context.Context, src *dagger.Directory) (*Infra, error) {
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

// WithRemoteBackendConfiguration sets the remote backend configuration in the container.
//
// This method sets the remote backend configuration in the container, making it available as environment variables.
//
// Parameters:
//   - bucket: The name of the bucket to use for the remote backend.
//   - locktable: The name of the lock table to use for the remote backend.
func (m *Infra) WithRemoteBackendConfiguration(bucket, locktable string) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_STACK_REMOTE_STATE_BUCKET_NAME", bucket).
		WithEnvVariable("TG_STACK_REMOTE_STATE_LOCK_TABLE", locktable)

	return m
}

// WithoutTracingToDagger disables tracing to Dagger in the container.
//
// This method disables tracing to Dagger in the container, making it available as an environment variable.
func (m *Infra) WithoutTracingToDagger() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("NOTHANKS", "1")

	return m
}

// WithTerragruntNoInteractive sets the Terragrunt no interactive option in the container.
//
// This method sets the Terragrunt no interactive option in the container, making it available as an environment variable.
func (m *Infra) WithTerragruntNoInteractive() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_NON_INTERACTIVE", "true")

	return m
}

// WithTerraformGitlabToken sets the Terraform Gitlab token in the container.
//
// This method sets the Terraform Gitlab token in the container, making it available as an environment variable.
//
// Parameters:
//   - ctx: The context for the Dagger container.
func (m *Infra) WithTerraformGitlabToken(ctx context.Context, token *dagger.Secret) *Infra {
	m.Ctr = m.Ctr.
		WithSecretVariable("TF_TOKEN_gitlab_com", token)

	return m
}

// WithTerragruntParallelism sets the parallelism level for Terragrunt operations.
//
// This method controls the number of units that are run concurrently during *-all commands.
//
// Parameters:
//   - parallelism: The number of concurrent operations to allow
func (m *Infra) WithTerragruntParallelism(parallelism int) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PARALLELISM", fmt.Sprintf("%d", parallelism))

	return m
}

// WithTerragruntNoAutoInit disables automatic initialization of Terraform.
//
// This method prevents Terragrunt from automatically running terraform init
// when other commands are executed.
func (m *Infra) WithTerragruntNoAutoInit() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_NO_AUTO_INIT", "true")

	return m
}

// WithTerragruntNoAutoApprove disables automatic approval for Terragrunt operations.
//
// This method prevents Terragrunt from automatically appending -auto-approve
// to underlying Terraform commands in *-all operations.
func (m *Infra) WithTerragruntNoAutoApprove() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_NO_AUTO_APPROVE", "true")

	return m
}

// WithTerragruntNoAutoRetry disables automatic retry of failed commands.
//
// This method prevents Terragrunt from automatically retrying commands
// that fail with transient errors.
func (m *Infra) WithTerragruntNoAutoRetry() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_NO_AUTO_RETRY", "true")

	return m
}

// WithTerragruntInputsDebug enables inputs debug mode.
//
// This method enables debug mode for Terragrunt inputs, creating tfvars files
// that can be used to invoke Terraform modules directly for debugging.
func (m *Infra) WithTerragruntInputsDebug() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_INPUTS_DEBUG", "true")

	return m
}

// WithTerragruntLogFormat sets the log format for Terragrunt output.
//
// This method sets the logging format for Terragrunt. Supported formats:
// "pretty", "bare", "json", "key-value"
//
// Parameters:
//   - format: The log format to use
func (m *Infra) WithTerragruntLogFormat(format string) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_LOG_FORMAT", format)

	return m
}

// WithTerragruntLogDisable disables all Terragrunt logging.
//
// This method disables logging output from Terragrunt and automatically
// enables Terraform stdout forwarding.
func (m *Infra) WithTerragruntLogDisable() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_LOG_DISABLE", "true")

	return m
}

// WithTerragruntLogShowAbsPaths enables absolute paths in log output.
//
// This method configures Terragrunt to show absolute paths in logs
// instead of relative paths to the working directory.
func (m *Infra) WithTerragruntLogShowAbsPaths() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_LOG_SHOW_ABS_PATHS", "true")

	return m
}

// WithTerragruntBackendRequireBootstrap requires backend resources to be bootstrapped.
//
// This method configures Terragrunt to fail if remote state bucket creation
// is necessary, requiring explicit bootstrap operations.
func (m *Infra) WithTerragruntBackendRequireBootstrap() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_BACKEND_REQUIRE_BOOTSTRAP", "true")

	return m
}

// WithTerragruntDisableBucketUpdate disables remote state bucket updates.
//
// This method prevents Terragrunt from updating the remote state bucket,
// useful when the state bucket is managed by a third party.
func (m *Infra) WithTerragruntDisableBucketUpdate() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_DISABLE_BUCKET_UPDATE", "true")

	return m
}

// WithTerragruntDisableCommandValidation disables Terraform command validation.
//
// This method disables Terragrunt's validation of Terraform commands,
// useful when using non-standard commands in hooks.
func (m *Infra) WithTerragruntDisableCommandValidation() *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_DISABLE_COMMAND_VALIDATION", "true")

	return m
}

// WithTerragruntProviderCacheDir sets a custom provider cache directory.
//
// This method sets the path to the Terragrunt provider cache directory.
//
// Parameters:
//   - cacheDir: The path to the provider cache directory
func (m *Infra) WithTerragruntProviderCacheDir(cacheDir string) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PROVIDER_CACHE_DIR", cacheDir)

	return m
}

// WithTerragruntProviderCacheHostname sets the provider cache server hostname.
//
// This method sets the hostname for the Terragrunt provider cache server.
//
// Parameters:
//   - hostname: The hostname for the provider cache server
func (m *Infra) WithTerragruntProviderCacheHostname(hostname string) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PROVIDER_CACHE_HOSTNAME", hostname)

	return m
}

// WithTerragruntProviderCachePort sets the provider cache server port.
//
// This method sets the port for the Terragrunt provider cache server.
//
// Parameters:
//   - port: The port number for the provider cache server
func (m *Infra) WithTerragruntProviderCachePort(port int) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PROVIDER_CACHE_PORT", fmt.Sprintf("%d", port))

	return m
}

// WithTerragruntProviderCacheToken sets the provider cache server authentication token.
//
// This method sets the authentication token for the Terragrunt provider cache server.
//
// Parameters:
//   - token: The authentication token for the provider cache server
func (m *Infra) WithTerragruntProviderCacheToken(token string) *Infra {
	m.Ctr = m.Ctr.
		WithEnvVariable("TG_PROVIDER_CACHE_TOKEN", token)

	return m
}
