# GitLab CI/CD Configuration for Terragrunt Reference Architecture

This directory contains the GitLab CI/CD pipeline configuration for the Terragrunt Reference Architecture project. It defines various workflows, jobs, and utility configurations to automate testing, building, and deployment processes.

## Directory Structure

```
.gitlab/
├── .gitlab-ci.yml                # Main GitLab CI/CD pipeline configuration
├── issue_templates/              # Templates for GitLab issues (currently empty)
├── merge_request_templates/      # Templates for GitLab merge requests (currently empty)
├── scripts/                      # Utility scripts for CI jobs
│   └── setup_ssh_agent.sh        # Script to set up SSH agent
├── utils/                        # Reusable CI job configurations and utilities
│   ├── pipeline_auth.yml         # Authentication related CI components
│   ├── pipeline_infra.yml        # Infrastructure related CI components
│   └── pipeline_tooling.yml      # Tooling and setup related CI components
└── workflows/                    # Individual workflow definition files
    ├── workflow_infra_stack_mr_nondist.yml
    ├── workflow_infra_stack_onmaster_nondist.yml
    ├── workflow_infra_stacks_ci.yml
    ├── workflow_infra_terraform_ci.yml
    ├── workflow_pipeline_build.yml
    ├── workflow_setup_aws.yml
    └── workflow_static_analysis.yml
```

## Main Configuration: `.gitlab-ci.yml`

The [`.gitlab-ci.yml`](../.gitlab-ci.yml) file is the entry point for all CI/CD pipelines.

### Key Components:

*   **Variables**: Defines global and job-specific variables used throughout the pipelines. Notable variables include:
    *   `ENVIRONMENT`: Target environment (e.g., `dev`, `prod`).
    *   `DEPLOYMENT_REGION`: AWS region for deployment.
    *   `GIT_REFERENCE`: Branch/tag for the Dagger module.
    *   `ADDITIONAL_ENV_VARS`: For injecting custom environment variables.
    *   `NOCACHE`: To bypass Dagger cache.
    *   `DAGGER_VERSION`: Specifies the Dagger version.
    *   `TRIGGER_WORKFLOW_*`: A set of boolean flags to manually trigger specific workflows.
*   **Includes**: The main configuration file uses `include:local` to incorporate various workflow files from the [`.gitlab/workflows/`](./workflows) directory. This modular approach keeps the main file clean and organizes workflows logically.
*   **Stages**: Defines the execution order of jobs within the pipeline. Common stages include `setup`, `dagger-ci`, `infra-ci-terraform`, `infra-ci-stacks-plan`, `stack-nondist-plan`, `stack-nondist-apply`, and `pipeline-infra`.
*   **Default Job**: A placeholder `default_job` is included to ensure pipeline validity even if no other jobs are triggered.

### Included Workflows:

The [`.gitlab-ci.yml`](../.gitlab-ci.yml) file conditionally includes the following workflows based on changes in specific paths or manual trigger variables:

1.  **Pipeline Build Workflow** ([`workflows/workflow_pipeline_build.yml`](./workflows/workflow_pipeline_build.yml))
    *   **Triggered by**: Changes in `pipeline/infra/**/*` or if `TRIGGER_WORKFLOW_PIPELINE_BUILD` is `true`.
    *   **Purpose**: Handles CI for the Dagger pipeline code itself.

2.  **Infrastructure Terraform CI Workflow** ([`workflows/workflow_infra_terraform_ci.yml`](./workflows/workflow_infra_terraform_ci.yml))
    *   **Triggered by**: Changes in `infra/terraform/**/*` or if `TRIGGER_WORKFLOW_INFRA_CI_TERRAFORM` is `true`.
    *   **Purpose**: Runs CI jobs specifically for Terraform module changes.

3.  **Infrastructure Stacks CI Workflow** ([`workflows/workflow_infra_stacks_ci.yml`](./workflows/workflow_infra_stacks_ci.yml))
    *   **Triggered by**: If `TRIGGER_WORKFLOW_INFRA_CI_STACKS` is `true`.
    *   **Purpose**: Runs CI jobs for Terragrunt stack changes (likely a broader trigger).

4.  **AWS Setup Workflow** ([`workflows/workflow_setup_aws.yml`](./workflows/workflow_setup_aws.yml))
    *   **Triggered by**: Changes in `.gitlab/workflows/**/*.yml`, `.gitlab-ci.yml`, or if `TRIGGER_WORKFLOW_SETUP_AWS` is `true`.
    *   **Purpose**: Tests AWS OIDC setup and related configurations.

5.  **Infrastructure Stack Non-Distributable (Merge Request) Workflow** ([`workflows/workflow_infra_stack_mr_nondist.yml`](./workflows/workflow_infra_stack_mr_nondist.yml))
    *   **Triggered by**: Changes in `infra/terragrunt/dev/non-distributable/**/*`, `infra/terragrunt/prod/non-distributable/**/*`, `infra/terragrunt/_shared/_units/**/*` on merge requests, or if `TRIGGER_WORKFLOW_INFRA_CI_MR_STACK_NON_DIST` is `true`.
    *   **Purpose**: Runs CI (plan) for the "non-distributable" stack components when changes are proposed in a merge request.

6.  **Infrastructure Stack Non-Distributable (Master/Default Branch) Workflow** ([`workflows/workflow_infra_stack_onmaster_nondist.yml`](./workflows/workflow_infra_stack_onmaster_nondist.yml))
    *   *(Note: This workflow is present in your open tabs but its include rule in `.gitlab-ci.yml` was not explicitly provided in the initial file content. Assuming it's triggered by `TRIGGER_WORKFLOW_INFRA_ONMASTER_STACK_NONDIST` or pushes to the default branch with relevant path changes.)*
    *   **Purpose**: Handles apply operations for the "non-distributable" stack components when changes are merged to the main branch.

7.  **Static Analysis Workflow** ([`workflows/workflow_static_analysis.yml`](./workflows/workflow_static_analysis.yml))
    *   *(Note: This workflow is present in your open tabs but its include rule in `.gitlab-ci.yml` was not explicitly provided in the initial file content. Assuming it's triggered by `TRIGGER_WORKFLOW_STATIC_ANALYSIS` or relevant code changes.)*
    *   **Purpose**: Performs static analysis checks on the codebase (e.g., linting, security scans).

## Workflows (`.gitlab/workflows/`)

This directory contains YAML files defining specific CI/CD workflows. Each workflow typically groups related jobs and defines its own rules for execution. The primary workflows are listed above as included by the main [`.gitlab-ci.yml`](../.gitlab-ci.yml).

## Utility Scripts (`.gitlab/scripts/`)

This directory is intended for shell scripts or other utility programs used by CI jobs.
*   [`setup_ssh_agent.sh`](./scripts/setup_ssh_agent.sh): A script likely used to initialize an SSH agent within CI jobs, which can be necessary for operations requiring SSH key authentication (e.g., cloning private Git repositories).

## Reusable CI Components (`.gitlab/utils/`)

This directory contains YAML files with reusable CI job definitions or configurations that can be imported into various workflows using GitLab's `include` or `extends` keywords. This promotes DRY (Don't Repeat Yourself) principles in CI configuration.
*   [`pipeline_auth.yml`](./utils/pipeline_auth.yml): Likely contains common job configurations for authentication, such as setting up AWS OIDC, handling secrets, or configuring access to other services.
*   [`pipeline_infra.yml`](./utils/pipeline_infra.yml): Probably defines reusable job templates or snippets for infrastructure-related tasks, such as running `terragrunt plan` or `terragrunt apply` with consistent settings.
*   [`pipeline_tooling.yml`](./utils/pipeline_tooling.yml): May contain job definitions for setting up necessary tools (e.g., Terraform, Terragrunt, Dagger, linters) within the CI environment.

## Issue and Merge Request Templates

The directories [`.gitlab/issue_templates/`](./issue_templates) and [`.gitlab/merge_request_templates/`](./merge_request_templates) are standard GitLab locations for defining templates that pre-fill issue or merge request descriptions. While currently empty, they can be populated with Markdown files to standardize these processes.

## How to Use and Extend

*   **Modifying Workflows**: To change existing CI processes, edit the relevant YAML file in the [`.gitlab/workflows/`](./workflows) directory.
*   **Adding New Workflows**: Create a new YAML file in [`.gitlab/workflows/`](./workflows) and include it in the main [`.gitlab-ci.yml`](../.gitlab-ci.yml) with appropriate `rules`.
*   **Custom Variables**: Adjust CI behavior by modifying the default values of variables in [`.gitlab-ci.yml`](../.gitlab-ci.yml) or by setting them during manual pipeline runs or in project/group CI/CD settings.
*   **Manual Triggers**: Use the `TRIGGER_WORKFLOW_*` variables (e.g., by running a pipeline with custom variables) to force specific workflows to run, which is useful for debugging or ad-hoc tasks.