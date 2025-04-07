package main

import (
	"context"
	"dagger/terragrunt/internal/dagger"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// ActionResult represents the result of a Ci action.
type ActionResult struct {
	Unit   string
	Output string
	Err    error
}

// TerragruntActionBuilder helps configure and execute a single Terragrunt action.
type TerragruntActionBuilder struct {
	parent         *Terragrunt       // Reference to the main Terragrunt module
	command        string            // The Terragrunt command (e.g., "hclfmt", "validate", "init")
	args           []string          // Arguments for the command
	unitEnv        string            // Terragrunt unit env (e.g., "global")
	unitLayer      string            // Terragrunt unit layer (e.g., "dni")
	unitName       string            // Terragrunt unit name (e.g., "dni_generator")
	src            *dagger.Directory // Source directory for this action (defaults to parent.Src)
	awsAccessKeyID *dagger.Secret    // Specific AWS Key ID for this action
	awsSecretKey   *dagger.Secret    // Specific AWS Secret Key for this action
	awsRegion      string            // Specific AWS region for this action
	envVars        []string          // Additional env vars for this action
	secrets        []*dagger.Secret  // Additional secrets for this action (name -> secret)
	noCache        bool              // Apply cache buster for this action
	loadDotEnvFile bool              // Load .env file for this action
	// TODO: Add fields for tfToken, ghToken, glToken, gitSshSocket if needed per-action
}

func (m *Terragrunt) NewAction(command string) *TerragruntActionBuilder {
	return &TerragruntActionBuilder{
		parent:    m,
		command:   command,
		unitEnv:   defaultRefArchEnv,
		unitLayer: defaulttRefArchLayer,
		unitName:  defaulttRefArchUnit,
		awsRegion: "eu-central-1",
		envVars:   []string{},
		secrets:   []*dagger.Secret{},
		src:       m.Src, // defaults to parent.Src
	}
}

// ForTgUnit sets the unit environment, layer, and name for the action.
//
// Parameters:
//   - env: The environment to run the command in (e.g., "global")
//   - layer: The layer to run the command in (e.g., "dni")
//   - unit: The unit to run the command in (e.g., "dni_generator")
func (b *TerragruntActionBuilder) ForTgUnit(env, layer, unit string) *TerragruntActionBuilder {
	if env == "" {
		env = defaultRefArchEnv
	}

	if layer == "" {
		layer = defaulttRefArchLayer
	}

	if unit == "" {
		unit = defaulttRefArchUnit
	}

	b.unitEnv = env
	b.unitLayer = layer
	b.unitName = unit

	return b
}

// WithArgs adds additional arguments to the command.
//
// Parameters:
//   - args: The arguments to add to the command.
func (b *TerragruntActionBuilder) WithArgs(args ...string) *TerragruntActionBuilder {
	b.args = append(b.args, args...)
	return b
}

// WithSource sets the source directory for the action.
//
// Parameters:
//   - src: The source directory to use for the action.
func (b *TerragruntActionBuilder) WithSource(src *dagger.Directory) *TerragruntActionBuilder {
	if src != nil {
		b.src = src
	}

	return b
}

// WithAWS sets the AWS credentials for the action.
//
// Parameters:
//   - awsAccessKeyID: The AWS access key ID to use for the action.
//   - awsSecretAccessKey: The AWS secret access key to use for the action.
//   - awsRegion: The AWS region to use for the action.
func (b *TerragruntActionBuilder) WithAWS(awsAccessKeyID, awsSecretAccessKey *dagger.Secret, awsRegion string) *TerragruntActionBuilder {
	awsRegion = getDefaultAWSRegionIfNotSet(awsRegion)

	if awsAccessKeyID != nil {
		b.awsAccessKeyID = awsAccessKeyID
	}

	if awsSecretAccessKey != nil {
		b.awsSecretKey = awsSecretAccessKey
	}

	b.awsRegion = awsRegion

	return b
}

// WithEnv adds an environment variable to the action.
//
// Parameters:
//   - key: The key of the environment variable.
//   - value: The value of the environment variable.
func (b *TerragruntActionBuilder) WithEnv(key, value string) *TerragruntActionBuilder {
	b.envVars = append(b.envVars, fmt.Sprintf("%s=%s", key, value))
	return b
}

// WithSecret adds a secret to the action.
//
// Parameters:
//   - key: The key of the secret.
//   - secret: The secret to add to the action.
func (b *TerragruntActionBuilder) WithSecret(ctx context.Context, name string, secret *dagger.Secret) *TerragruntActionBuilder {
	if secret == nil {
		return b
	}

	secretPlainTxtValue, _ := secret.Plaintext(ctx)
	newSecret := dag.SetSecret(name, secretPlainTxtValue)
	b.secrets = append(b.secrets, newSecret)

	return b
}

// WithNoCache sets the cache buster for the action.
func (b *TerragruntActionBuilder) WithNoCache() *TerragruntActionBuilder {
	b.noCache = true
	return b
}

// WithLoadDotEnvFile loads the .env file for the action.
func (b *TerragruntActionBuilder) WithLoadDotEnvFile() *TerragruntActionBuilder {
	b.loadDotEnvFile = true
	return b
}

func (b *TerragruntActionBuilder) Execute(ctx context.Context) (*dagger.Container, error) {
	if b.parent == nil {
		return nil, NewError("cannot execute action: parent Terragrunt module is nil")
	}
	if b.command == "" {
		return nil, NewError("cannot execute action: command is empty")
	}
	if b.src == nil {
		return nil, NewError("cannot execute action: source directory is nil")
	}

	// Set terragrunt execution path
	tgUnitPath := getTerragruntUnitPath(b.unitEnv, b.unitLayer, b.unitName)

	if len(b.envVars) > 0 {
		modDecorate, envVarErr := b.parent.WithEnvVars(b.envVars)

		if envVarErr != nil {
			return nil, WrapErrorf(envVarErr, "failed to set environment variables for the following environment: %s, layer: %s, unit: %s", b.unitEnv, b.unitLayer, b.unitName)
		}

		b.parent = modDecorate
	}

	// Set AWS credentials
	if b.awsAccessKeyID != nil && b.awsSecretKey != nil {
		b.parent = b.parent.WithAWSCredentials(ctx, b.awsAccessKeyID, b.awsSecretKey, b.awsRegion)
	}

	// Set secrets
	if len(b.secrets) > 0 {
		b.parent = b.parent.WithSecrets(ctx, b.secrets)
	}

	// Set cache buster
	if b.noCache {
		b.parent = b.parent.WithCacheBuster()
	}

	// Set loadDotEnvFile
	if b.loadDotEnvFile {
		modDecorate, envVarErr := b.parent.WithDotEnvFile(ctx, b.src)

		if envVarErr != nil {
			return nil, WrapErrorf(envVarErr, "failed to set environment variables for the following environment: %s, layer: %s, unit: %s", b.unitEnv, b.unitLayer, b.unitName)
		}

		b.parent = modDecorate
	}

	execCtr, execErr := b.parent.Exec(
		ctx,
		"terragrunt",
		b.command,
		b.args,
		false, // autoApprove - maybe add to builder?
		b.src, // Use the source specified for this action
		tgUnitPath,
		b.envVars,
		b.secrets,
		nil, // tfToken - Add to builder if needed per-action
		nil, // ghToken - Add to builder if needed per-action
		nil, // glToken - Add to builder if needed per-action
		nil, // gitSshSocket - Add to builder if needed per-action
		false,
		b.loadDotEnvFile, // loadDotEnvFile - assumes Exec handles this contextually
	)

	if execErr != nil {
		return nil, WrapErrorf(execErr, "failed to execute the terragrunt action for the following environment: %s, layer: %s, unit: %s", b.unitEnv, b.unitLayer, b.unitName)
	}

	return execCtr, nil
}

// String returns a string representation of the ActionResult.
func (ar ActionResult) String() string {
	if ar.Err != nil {
		return fmt.Sprintf("Unit [%s]: Error - %v", ar.Unit, ar.Err)
	}
	if ar.Output != "" {
		return fmt.Sprintf("Unit [%s]:\n%s", ar.Unit, ar.Output)
	}
	return fmt.Sprintf("Unit [%s]: OK (No specific output)", ar.Unit) // Or just OK
}

// cloneTerragrunt creates a deep copy of the Terragrunt struct to avoid concurrency issues
func cloneTerragrunt(original *Terragrunt) *Terragrunt {
	if original == nil {
		return nil
	}
	return &Terragrunt{
		Ctr: original.Ctr,
		Src: original.Src,
	}
}

func (m *Terragrunt) JobTerragruntUnitStaticCheck(
	// ctx is the context for the Dagger container.
	// +optional
	ctx context.Context,
	// awsAccessKeyID is the AWS access key ID to use for the command.
	awsAccessKeyID *dagger.Secret,
	// awsSecretAccessKey is the AWS secret access key to use for the command.
	awsSecretAccessKey *dagger.Secret,
	// noCache set the cache buster
	// +optional
	noCache bool,
	// tgArgs are the arbitrary arguments to pass to the terrawgrunt command
	// +optional
	tgArgs []string,
	// loadDotEnvFile is a boolean indicating whether to load .env file
	// +optional
	loadDotEnvFile bool,
) (string, error) {
	units := []string{
		"dni_generator",
		"lastname_generator",
		"name_generator",
		"age_generator",
	}

	// actions to run on this job
	actions := []struct {
		Command string
		Args    []string
	}{
		{
			Command: "hclfmt",
			Args:    []string{"--check", "--diff"},
		},
		{
			Command: "terragrunt-info",
			Args:    []string{},
		},
		{
			Command: "validate-inputs",
			Args:    []string{},
		},
	}

	var wg sync.WaitGroup
	// Buffered channel, so we can collect the results safely.
	resultChan := make(chan ActionResult, len(units)*len(actions))

	// Run each unit check concurrently
	for _, unit := range units {
		for _, action := range actions {
			currentUnit := unit
			currentAction := action

			wg.Add(1)

			go func() {
				defer wg.Done()
				var finalErr error
				var finalOutput string

				// Create a clone of the Terragrunt struct to avoid concurrency issues
				terragruntClone := cloneTerragrunt(m)

				// Builder
				actionBuilder := terragruntClone.NewAction(currentAction.Command)
				actionBuilder.ForTgUnit(defaultRefArchEnv, defaulttRefArchLayer, currentUnit)
				actionBuilder.WithSource(terragruntClone.Src)
				actionBuilder.WithAWS(awsAccessKeyID, awsSecretAccessKey, "eu-central-1")
				actionBuilder.WithNoCache()
				actionBuilder.WithArgs(currentAction.Args...)
				actionBuilder.WithLoadDotEnvFile()

				// Execute
				compiledCtr, buildErr := actionBuilder.Execute(ctx)
				if buildErr != nil {
					finalErr = WrapErrorf(buildErr, "failed to execute the terragrunt action for the following environment: %s, layer: %s, unit: %s", currentUnit, defaulttRefArchLayer, currentUnit)
				} else {
					// get the stdout of the action
					stdOut, stdOutErr := compiledCtr.Stdout(ctx)
					if stdOutErr != nil {
						finalErr = WrapErrorf(stdOutErr, "failed to get the stdout of the terragrunt action for the following environment: %s, layer: %s, unit: %s", currentUnit, defaulttRefArchLayer, currentUnit)
					} else {
						finalOutput = stdOut
					}
				}

				// Collect results
				resultChan <- ActionResult{
					Unit:   fmt.Sprintf("%s.%s", currentUnit, currentAction.Command),
					Output: finalOutput,
					Err:    finalErr,
				}
			}()
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultChan)

	// collectors
	var collectedActionErrors []error
	// I'm opinionated, here the key defined is "unit.command"
	successfulOutputs := make(map[string]string)

	for result := range resultChan {
		if result.Err != nil {
			collectedActionErrors = append(collectedActionErrors, result.Err)
		} else {
			successfulOutputs[result.Unit] = result.Output
		}
	}

	// Handling, and showing errors.
	if len(collectedActionErrors) > 0 {
		return "", JoinErrors(collectedActionErrors...)
	}

	// Handling, and showing outputs
	var outputBuilder strings.Builder
	outputBuilder.WriteString("All static checks passed successfully.\n\nOutput per unit:\n")
	outputBuilder.WriteString("=====================\n")

	// Group outputs by unit
	unitOutputs := make(map[string]map[string]string)
	for key, output := range successfulOutputs {
		parts := strings.Split(key, ".")
		if len(parts) == 2 {
			unit, command := parts[0], parts[1]
			if _, ok := unitOutputs[unit]; !ok {
				unitOutputs[unit] = make(map[string]string)
			}
			unitOutputs[unit][command] = output
		}
	}

	// Sort units for consistent output order
	sortedUnits := make([]string, 0, len(unitOutputs))
	for unit := range unitOutputs {
		sortedUnits = append(sortedUnits, unit)
	}
	sort.Strings(sortedUnits)

	// Display outputs by unit and command
	for _, unit := range sortedUnits {
		outputBuilder.WriteString(fmt.Sprintf("--- Unit: %s ---\n", unit))

		// Sort commands for consistent output order
		commands := make([]string, 0, len(unitOutputs[unit]))
		for cmd := range unitOutputs[unit] {
			commands = append(commands, cmd)
		}
		sort.Strings(commands)

		for _, cmd := range commands {
			outputBuilder.WriteString(fmt.Sprintf("Command: %s\n", cmd))
			stdout := unitOutputs[unit][cmd]
			if stdout == "" {
				outputBuilder.WriteString("(No standard output)\n")
			} else {
				outputBuilder.WriteString(stdout)
				// Ensure a newline after each command's output if not already present
				if !strings.HasSuffix(stdout, "\n") {
					outputBuilder.WriteString("\n")
				}
			}
			outputBuilder.WriteString("\n")
		}
		outputBuilder.WriteString("--------------------\n")
	}

	return outputBuilder.String(), nil // Return combined stdout and nil error
}

// func (m *Terragrunt) ActionTerragruntHclValidate(
// 	// ctx is the context for the Dagger container.
// 	// +optional
// 	ctx context.Context,
// 	// src is the source directory for the Dagger container.
// 	// +optional
// 	// +defaultPath="/"
// 	// +ignore=["*", "!**/*.hcl", "!**/*.tfvars", "!**/.git/**", "!**/*.tfvars.json", "!**/*.tf", "!*.env", "!*.envrc", "!*.envrc"]
// 	src *dagger.Directory,
// 	// env is the environment to run the command in, following the pattern infra/terraform/[env]/[layer]/[unit]
// 	// +optional
// 	env string,
// 	// layer is the terragrunt layer to run the command in (pattern infra/terraform/[env]/[layer]/[unit])
// 	// +optional
// 	layer string,
// 	// unit is the terragrunt unit to run the command in (pattern infra/terraform/[env]/[layer]/[unit])
// 	// +optional
// 	unit string,
// 	// awsAccessKeyID is the AWS access key ID to use for the command.
// 	awsAccessKeyID *dagger.Secret,
// 	// awsSecretAccessKey is the AWS secret access key to use for the command.
// 	awsSecretAccessKey *dagger.Secret,
// 	// awsRegion is the AWS region to use for the command.
// 	// +optional
// 	awsRegion string,
// 	// envVars is a slice of strings including all environment variables to pass to the command.
// 	// +optional
// 	envVars []string,
// 	// noCache set the cache buster
// 	// +optional
// 	noCache bool,
// 	// tgArgs are the arbitrary arguments to pass to the terrawgrunt command
// 	// +optional
// 	tgArgs []string,
// 	// loadDotEnvFile is a boolean indicating whether to load .env file
// 	// +optional
// 	loadDotEnvFile bool,
// ) (*dagger.Container, error) {
// 	if env == "" {
// 		env = defaultRefArchEnv
// 	}

// 	if layer == "" {
// 		layer = defaulttRefArchLayer
// 	}

// 	if unit == "" {
// 		unit = defaulttRefArchUnit
// 	}

// 	if awsAccessKeyID != nil && awsSecretAccessKey != nil {
// 		m = m.WithAWSCredentials(ctx, awsAccessKeyID, awsSecretAccessKey, awsRegion)
// 	}

// 	if len(envVars) > 0 {
// 		modDecorate, envVarErr := m.WithEnvVars(envVars)

// 		if envVarErr != nil {
// 			return nil, WrapErrorf(envVarErr, "failed to set environment variables for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 		}

// 		m = modDecorate
// 	}

// 	if noCache {
// 		m = m.WithCacheBuster()
// 	}

// 	tgPath := getTerragruntUnitPath(env, layer, unit)

// 	// Run 'terragrunt hclvalidate'
// 	tgHclValidateCtr, tgHclValidateCtrErr := m.Exec(
// 		ctx,
// 		"terragrunt",                   // binary
// 		"hclvalidate",                  // command
// 		[]string{"--show-config-path"}, // args
// 		false,                          // autoApprove
// 		src,                            // src (we are operating within the base src context m.Src)
// 		tgPath,                         // tgUnitPath
// 		nil,                            // envVars (already set on m.Ctr by parseDotEnvFiles)
// 		nil,                            // secrets (already set on m.Ctr by parseDotEnvFiles if applicable)
// 		nil,                            // tfToken
// 		nil,                            // ghToken
// 		nil,                            // glToken
// 		nil,                            // gitSshSocket
// 		true,                           // printPaths
// 		loadDotEnvFile,                 // loadEnvFiles
// 	)

// 	if tgHclValidateCtrErr != nil {
// 		return nil, WrapErrorf(tgHclValidateCtrErr, "failed to run the terragrunt action for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 	}

// 	return tgHclValidateCtr, nil
// }

// func (m *Terragrunt) ActionTerragruntHCLFmt(
// 	// ctx is the context for the Dagger container.
// 	// +optional
// 	ctx context.Context,
// 	// src is the source directory for the Dagger container.
// 	// +optional
// 	// +defaultPath="/"
// 	// +ignore=["*", "!**/*.hcl", "!**/*.tfvars", "!**/.git/**", "!**/*.tfvars.json", "!**/*.tf", "!*.env", "!*.envrc", "!*.envrc"]
// 	src *dagger.Directory,
// 	// env is the environment to run the command in, following the pattern infra/terraform/[env]/[layer]/[unit]
// 	// +optional
// 	env string,
// 	// layer is the terragrunt layer to run the command in (pattern infra/terraform/[env]/[layer]/[unit])
// 	// +optional
// 	layer string,
// 	// unit is the terragrunt unit to run the command in (pattern infra/terraform/[env]/[layer]/[unit])
// 	// +optional
// 	unit string,
// 	// awsAccessKeyID is the AWS access key ID to use for the command.
// 	awsAccessKeyID *dagger.Secret,
// 	// awsSecretAccessKey is the AWS secret access key to use for the command.
// 	awsSecretAccessKey *dagger.Secret,
// 	// awsRegion is the AWS region to use for the command.
// 	// +optional
// 	awsRegion string,
// 	// envVars is a slice of strings including all environment variables to pass to the command.
// 	// +optional
// 	envVars []string,
// 	// noCache set the cache buster
// 	// +optional
// 	noCache bool,
// 	// tgArgs are the arbitrary arguments to pass to the terrawgrunt command
// 	// +optional
// 	tgArgs []string,
// 	// loadDotEnvFile is a boolean indicating whether to load .env file
// 	// +optional
// 	loadDotEnvFile bool,
// ) (*dagger.Container, error) {
// 	if env == "" {
// 		env = defaultRefArchEnv
// 	}

// 	if layer == "" {
// 		layer = defaulttRefArchLayer
// 	}

// 	if unit == "" {
// 		unit = defaulttRefArchUnit
// 	}

// 	if awsAccessKeyID != nil && awsSecretAccessKey != nil {
// 		m = m.WithAWSCredentials(ctx, awsAccessKeyID, awsSecretAccessKey, awsRegion)
// 	}

// 	if len(envVars) > 0 {
// 		modDecorate, envVarErr := m.WithEnvVars(envVars)

// 		if envVarErr != nil {
// 			return nil, WrapErrorf(envVarErr, "failed to set environment variables for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 		}

// 		m = modDecorate
// 	}

// 	if noCache {
// 		m = m.WithCacheBuster()
// 	}

// 	tgPath := getTerragruntUnitPath(env, layer, unit)

// 	// Run HCLfmt
// 	tgHclFmtCtr, tgHclFmtCtrErr := m.Exec(
// 		ctx,
// 		"terragrunt",                           // binary
// 		"hclfmt",                               // command
// 		[]string{"--check", "--all", "--diff"}, // args
// 		false,                                  // autoApprove
// 		src,                                    // src (we are operating within the base src context m.Src)
// 		tgPath,                                 // tgUnitPath
// 		nil,                                    // envVars (already set on m.Ctr by parseDotEnvFiles)
// 		nil,                                    // secrets (already set on m.Ctr by parseDotEnvFiles if applicable)
// 		nil,                                    // tfToken
// 		nil,                                    // ghToken
// 		nil,                                    // glToken
// 		nil,                                    // gitSshSocket
// 		true,                                   // printPaths
// 		loadDotEnvFile,                         // loadEnvFiles
// 	)

// 	if tgHclFmtCtrErr != nil {
// 		return nil, WrapErrorf(tgHclFmtCtrErr, "failed to run the terragrunt action for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 	}

// 	return tgHclFmtCtr, nil

// 	// _, tgHclOutErr := tgHclFmtCtr.Stdout(ctx)
// 	// if tgHclOutErr != nil {
// 	// 	return "", WrapErrorf(tgHclOutErr, "failed to get the stdout of the terragrunt action for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 	// }

// 	// // Run 'terragrunt hclvalidate'
// 	// tgHclValidateCtr, tgHclValidateCtrErr := m.Exec(
// 	// 	ctx,
// 	// 	"terragrunt",                   // binary
// 	// 	"hclvalidate",                  // command
// 	// 	[]string{"--show-config-path"}, // args
// 	// 	false,                          // autoApprove
// 	// 	src,                            // src (we are operating within the base src context m.Src)
// 	// 	tgPath,                         // tgUnitPath
// 	// 	nil,                            // envVars (already set on m.Ctr by parseDotEnvFiles)
// 	// 	nil,                            // secrets (already set on m.Ctr by parseDotEnvFiles if applicable)
// 	// 	nil,                            // tfToken
// 	// 	nil,                            // ghToken
// 	// 	nil,                            // glToken
// 	// 	nil,                            // gitSshSocket
// 	// 	true,                           // printPaths
// 	// 	loadDotEnvFile,                 // loadEnvFiles
// 	// )

// 	// if tgHclValidateCtrErr != nil {
// 	// 	return "", WrapErrorf(tgHclValidateCtrErr, "failed to run the terragrunt action for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 	// }

// 	// _, tgHclValidateOutErr := tgHclValidateCtr.Stdout(ctx)
// 	// if tgHclValidateOutErr != nil {
// 	// 	return "", WrapErrorf(tgHclValidateOutErr, "failed to get the stdout of the terragrunt action for the following environment: %s, layer: %s, unit: %s", env, layer, unit)
// 	// }

// 	// return "", nil
// }
