package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
)

func getTFInstallCmd(tfVersion string) string {
	installDir := "/usr/local/bin/terraform"
	command := fmt.Sprintf(`apk add --no-cache curl unzip &&
	curl -L https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip -o /tmp/terraform.zip &&
	unzip -o /tmp/terraform.zip -d /tmp &&
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

func getTerragruntExecutionPath(env, layer, unit string) string {
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

func getTerraformModulesExecutionPath(moduleName string) string {
	return filepath.Join(configRefArchATerraformModulesRootPath, moduleName)
}

type EnvVarDagger struct {
	Key   string
	Value string
}

func getEnvVarsDaggerFromSlice(envVars []string) ([]EnvVarDagger, error) {
	envVarsDagger := []EnvVarDagger{}
	for _, envVar := range envVars {
		trimmedEnvVar := strings.TrimSpace(envVar)
		if trimmedEnvVar == "" {
			return nil, NewError("environment variable cannot be empty")
		}

		if !strings.Contains(trimmedEnvVar, "=") {
			return nil, NewError(fmt.Sprintf("environment variable must be in the format ENVARKEY=VALUE: %s", trimmedEnvVar))
		}

		parts := strings.Split(trimmedEnvVar, "=")
		if len(parts) != 2 {
			return nil, NewError(fmt.Sprintf("environment variable must be in the format ENVARKEY=VALUE: %s", trimmedEnvVar))
		}

		envVarsDagger = append(envVarsDagger, EnvVarDagger{
			Key:   parts[0],
			Value: parts[1],
		})
	}

	return envVarsDagger, nil
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
