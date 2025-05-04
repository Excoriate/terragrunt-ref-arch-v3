package main

import (
	"context"
	"fmt"
	"sync"
)

var tfModules = []string{
	"dni-generator",
	"lastname-generator",
	"name-generator",
	"age-generator",
}

type TfModulesMatrixConfig struct {
	Module      string
	IsCIEnabled bool
	TFversions  []string
}

var tfModulesMatrixConfig = []TfModulesMatrixConfig{
	{
		Module:      "read-aws-metadata",
		IsCIEnabled: true,
		TFversions:  []string{"1.11.3", "1.11.1", "1.11.0"},
	},
	{
		Module:      "random-string-generator",
		IsCIEnabled: true,
		TFversions:  []string{"1.11.3", "1.11.1", "1.11.0"},
	},
}

// JobTfModulesStaticCheck performs static checks on Terraform modules by:
// - Initializing Terraform without backend configuration
// - Validating the module configuration
// - Formatting the module code recursively
//
// It processes all modules defined in tfModules and returns a combined result
// of all checks. If any module fails validation, the function returns an error.
//
// Parameters:
//   - ctx: Context for managing the operation's lifecycle
//
// Returns:
//   - string: Combined output of all module checks
//   - error: Any error encountered during the checks, nil if successful
func (m *Infra) JobTfModulesStaticCheck(ctx context.Context) (string, error) {
	results := []JobResult{}

	for _, module := range tfModules {
		tfModuleMntPath := fmt.Sprintf("%s/%s", defaultMntPath, getTerraformModulesExecutionPath(module))
		m, err := m.WithSRC(ctx, tfModuleMntPath, m.Src)

		if err != nil {
			return "", WrapErrorf(err, "failed to set source code for module %s mounted in path %s", module, tfModuleMntPath)
		}

		m = m.WithCacheBuster()

		execCtr := m.Ctr.
			WithExec([]string{"terraform", "init", "-backend=false"}).
			WithExec([]string{"terraform", "validate"}).
			WithExec([]string{"terraform", "fmt", "-recursive", "-check", "-diff"})

		stdout, err := execCtr.Stdout(ctx)

		if err != nil {
			return "", WrapErrorf(err, "failed to validate module %s", module)
		}

		results = append(results, JobResult{
			WorkDir: tfModuleMntPath,
			Output:  stdout,
			Err:     nil,
		})
	}

	return ProcessActionSyncResults(results)
}

// JobTfModulesCompatibilityCheck performs compatibility checks on Terraform modules asynchronously by:
// - Testing each module against multiple Terraform versions specified in the module's configuration in parallel.
// - Initializing Terraform without backend configuration for each version.
// - Validating the module configuration.
// - Formatting the module code recursively.
// - Generating a JSON representation of the module's configuration.
//
// The function processes all modules defined in tfModulesMatrixConfig that have CI enabled.
// For each module, it runs the checks against all specified Terraform versions using goroutines.
// The results are collected via a channel and processed asynchronously.
//
// Parameters:
//   - ctx: Context for managing the operation's lifecycle.
//
// Returns:
//   - string: Combined output of all module compatibility checks.
//   - error: Any error encountered during the checks, nil if successful.
func (m *Infra) JobTfModulesCompatibilityCheck(ctx context.Context) (string, error) {
	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	// Estimate buffer size: number of modules * typical number of versions
	// Adjust buffer size if needed based on actual config length
	bufferSize := len(tfModulesMatrixConfig) * 3
	if bufferSize == 0 {
		bufferSize = 10 // Default buffer if config is empty
	}
	resultChan := make(chan JobResult, bufferSize)

	// Base Infra struct for this job - avoid modifying this directly in goroutines
	baseInfra := m

	for _, moduleCfg := range tfModulesMatrixConfig {
		if !moduleCfg.IsCIEnabled {
			continue
		}

		for _, tfVersion := range moduleCfg.TFversions {
			// Increment WaitGroup counter for each task
			wg.Add(1)

			// Launch a goroutine for each module/version combination
			// Capture loop variables correctly by passing them as arguments
			go func(modCfg TfModulesMatrixConfig, version string) {
				// Decrement counter when goroutine completes
				defer wg.Done()

				// Create a result struct to hold the outcome for this specific check
				jobRes := JobResult{
					// Use a more descriptive WorkDir combining module and version
					WorkDir: fmt.Sprintf("%s/tf-%s", getTerraformModulesExecutionPath(modCfg.Module), version),
					Output:  "",
					Err:     nil,
				}

				// Configure Terraform version and plugin cache starting from the base state
				// This ensures each goroutine works with an isolated container config state
				infraForVersion := baseInfra.WithTerraform(version).
					WithTerraformPluginCache().
					WithCacheBuster() // Add cache buster per check

				tfModuleMntPath := fmt.Sprintf("%s/%s", defaultMntPath, getTerraformModulesExecutionPath(modCfg.Module))
				// Mount source code for this specific check
				infraWithSrc, err := infraForVersion.WithSRC(ctx, tfModuleMntPath, baseInfra.Src) // Use baseInfra.Src to avoid mounting issues

				if err != nil {
					jobRes.Err = WrapErrorf(err, "goroutine failed to set source code for module %s (TF %s) mounted in path %s", modCfg.Module, version, tfModuleMntPath)
					resultChan <- jobRes
					return // Exit goroutine on setup error
				}

				infraWithSrc = infraWithSrc.WithCacheBuster() // Add cache buster again after SRC mount might be needed

				// Define and execute Terraform commands
				execCtr := infraWithSrc.Ctr.
					WithExec([]string{"terraform", "init", "-backend=false"}).
					WithExec([]string{"terraform", "validate"}).
					WithExec([]string{"terraform", "fmt", "-recursive", "-check", "-diff"}).
					WithExec([]string{"terraform", "show", "-json"})

				// Execute commands and capture output/error
				stdout, err := execCtr.Stdout(ctx)
				// Capture output even if there's an error
				jobRes.Output = stdout
				if err != nil {
					jobRes.Err = WrapErrorf(err, "goroutine failed executing checks for module %s with TF %s", modCfg.Module, version)
					resultChan <- jobRes
					return // Exit goroutine on execution error
				}

				// Send the successful result to the channel
				resultChan <- jobRes

			}(moduleCfg, tfVersion) // Pass copies of loop variables to the goroutine
		}
	}

	// Start a goroutine to close the channel once all checks are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results collected from the channel asynchronously
	// Use the async processor function from job.go
	return processActionAsyncResults(resultChan)
}
