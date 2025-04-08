package main

import (
	"context"
	"fmt"
)

func (m *Terragrunt) JobTerraformModulesStaticCheck(ctx context.Context) (string, error) {
	modules := []string{
		"dni-generator",
		"lastname-generator",
		"name-generator",
		"age-generator",
	}

	results := []JobResult{}

	for _, module := range modules {
		tfModuleMntPath := fmt.Sprintf("%s/%s", defaultMntPath, getTerraformModulesExecutionPath(module))
		m, err := m.WithSRC(ctx, tfModuleMntPath, m.Src)

		if err != nil {
			return "", WrapErrorf(err, "failed to set source code for module %s mounted in path %s", module, tfModuleMntPath)
		}

		execCtr := m.Ctr.
			WithExec([]string{"terraform", "init", "-backend=false"}).
			WithExec([]string{"terraform", "validate"}).
			WithExec([]string{"terraform", "fmt", "-recursive"})

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
