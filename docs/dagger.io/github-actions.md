---
slug: /ci/integrations/github-actions
---

# GitHub Actions

Dagger provides a GitHub Action that can be used in any GitHub Actions workflow to call one or more Dagger Functions on specific events, such as a new commit.

Dagger has also partnered with [Depot](https://depot.dev/) to provide managed, Dagger Powered GitHub Actions runners. These runners, which serve as drop-in replacements for GitHub's own runners, come with Dagger pre-installed and pre-configured to best practices, automatic persistent layer caching, and multi-architecture support. They make it faster and easier to run Dagger pipelines in GitHub repositories.

## How it works

Depending on whether you choose a standard GitHub Actions runner or a Dagger Powered Depot GitHub Actions Runner, here's how it works.

### GitHub

When running a CI pipeline with Dagger using a standard GitHub Actions runner, the general workflow looks like this:

1. GitHub receives a workflow trigger based on a repository event.
1. GitHub begins processing the jobs and steps in the workflow.
1. GitHub finds the "Dagger for GitHub" action and passes control to it.
1. The "Dagger for GitHub" action calls the Dagger CLI with the specified sub-command, module, function name, and arguments.
1. The Dagger CLI attempts to find an existing Dagger Engine or spins up a new one inside the GitHub Actions runner.
1. The Dagger CLI executes the specified sub-command and sends telemetry to Dagger Cloud if the `DAGGER_CLOUD_TOKEN` environment variable is set.
1. The workflow completes with success or failure. Logs appear in GitHub as usual.

### Depot

When running a CI pipeline on a Dagger Powered Depot GitHub Actions Runner, the general workflow looks like this:

1. GitHub emits a webhook event which triggers a job in a GitHub Actions workflow
1. Depot receives this event, processes it, and launches an isolated Depot GitHub Actions Runner.
1. Each Depot GitHub Actions Runner is automatically connected to an already running Dagger Engine (this exists outside of the Depot GitHub Actions Runner).
    * The external Dagger Engine comes with persistent caching pre-configured.
1. The Depot GitHub Actions Runner has the Dagger CLI pre-installed.
    * The Dagger CLI running inside the Depot GitHub Actions Runner gets automatically connected to Dagger Cloud.
    * The Dagger CLI is also pre-configured to connect to the external Dagger Engine with persistent caching
1. The "Dagger for GitHub" GitHub Action calls the Dagger CLI with the specified sub-command, module, function name, and arguments.
1. The Dagger CLI executes the specified sub-command and sends telemetry to Dagger Cloud.
1. The workflow completes with success or failure. Logs appear in GitHub as usual. They also appear in Dagger Cloud.

## Prerequisites

### GitHub

- A GitHub repository

### Depot

- A GitHub repository in a GitHub organization
- A Depot account and [organization](https://depot.dev/docs/github-actions/quickstart#create-an-organization)
- The [Depot organization connected to the GitHub organization](https://depot.dev/docs/github-actions/quickstart#connect-to-github)

## Examples

### GitHub

The following example demonstrates how to call a Dagger Function on a standard GitHub runner in a GitHub Actions workflow.

```yaml title=".github/workflows/dagger.yml"
# .github/workflows/dagger.yml
name: Dagger

on:
  push:
    branches: ["main"]

jobs:
  hello:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Dagger CLI
        uses: dagger/setup-dagger@v5

      - name: Run Dagger Function
        uses: dagger/dagger-for-github@v5
        with:
          # Replace with your Dagger Function call
          # Example: call container-echo --string-arg="Hello from Dagger!" stdout
          # Example: call test build publish --tag=ttl.sh/my-image
          args: "call container-echo --string-arg='Hello from Dagger!' stdout"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}
```

The following is a more complex example demonstrating how to create a GitHub Actions workflow that checks out source code, calls a Dagger Function to test the project, and then calls another Dagger Function to build and publish a container image of the project. This example uses a simple [Go application](https://github.com/kpenfound/greetings-api) and assumes that you have already forked it in your own GitHub repository.

```yaml title=".github/workflows/dagger.yml"
# .github/workflows/dagger.yml
name: Dagger

on:
  push:
    branches: ["main"]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Dagger CLI
        uses: dagger/setup-dagger@v5

      - name: Run Dagger Function
        uses: dagger/dagger-for-github@v5
        with:
          # Replace with your Dagger Function call
          # Example: call test
          args: "call test"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}

  build:
    runs-on: ubuntu-latest
    needs: [test]
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Dagger CLI
        uses: dagger/setup-dagger@v5

      - name: Run Dagger Function
        uses: dagger/dagger-for-github@v5
        with:
          # Replace with your Dagger Function call
          # Example: call build publish --tag=ttl.sh/my-image
          args: "call build publish --tag=ttl.sh/my-image"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}
```

More information is available in the [Dagger for GitHub page](https://github.com/marketplace/actions/dagger-for-github).

### Depot

The following example demonstrates how to call a Dagger Function on a Dagger Powered Depot runner in a GitHub Actions workflow. The `runs-on` clause specifies the runner and Dagger versions.

```yaml title=".github/workflows/dagger.yml"
# .github/workflows/dagger.yml
name: Dagger

on:
  push:
    branches: ["main"]

jobs:
  hello:
    runs-on: depot-ubuntu-2204-m # Replace with your Depot runner label
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Run Dagger Function
        uses: dagger/dagger-for-github@v5
        with:
          # Replace with your Dagger Function call
          # Example: call container-echo --string-arg="Hello from Dagger!" stdout
          # Example: call test build publish --tag=ttl.sh/my-image
          args: "call container-echo --string-arg='Hello from Dagger!' stdout"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}
```

The following is a more complex example demonstrating how to create a GitHub Actions workflow that checks out source code, calls a Dagger Function to test the project, and then calls another Dagger Function to build and publish a container image of the project. This example uses a simple [Go application](https://github.com/kpenfound/greetings-api) and assumes that you have already forked it in your own GitHub repository.

```yaml title=".github/workflows/dagger.yml"
# .github/workflows/dagger.yml
name: Dagger

on:
  push:
    branches: ["main"]

jobs:
  test:
    runs-on: depot-ubuntu-2204-m # Replace with your Depot runner label
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Run Dagger Function
        uses: dagger/dagger-for-github@v5
        with:
          # Replace with your Dagger Function call
          # Example: call test
          args: "call test"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}

  build:
    runs-on: depot-ubuntu-2204-m # Replace with your Depot runner label
    needs: [test]
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Run Dagger Function
        uses: dagger/dagger-for-github@v5
        with:
          # Replace with your Dagger Function call
          # Example: call build publish --tag=ttl.sh/my-image
          args: "call build publish --tag=ttl.sh/my-image"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}
```

> **Warning:**
> Always ensure that the same Dagger versions are specified in the `runs-on` and `version` clauses of the workflow definition.

More information is available in the [Dagger for GitHub page](https://github.com/marketplace/actions/dagger-for-github) and in the [Depot documentation](https://depot.dev/docs/).

### SSH configuration

When using SSH keys in GitHub Actions, ensure proper SSH agent setup:

```yaml
- name: Set up SSH
  run: |
    eval "$(ssh-agent -s)"
    ssh-add - <<< '${{ secrets.SSH_PRIVATE_KEY }}'
```

Replace `${{ secrets.SSH_PRIVATE_KEY }}` with your provider secret containing the private key.

## Resources

Below are some resources from the Dagger community that may help as well. If you have any questions about additional ways to use GitHub with Dagger, join our [Discord](https://discord.gg/dagger-io) and ask your questions in our [GitHub channel](https://discord.com/channels/707636530424053791/1117139064274034809).

- [Video: Boost Dagger Performance on GitHub Actions with Depot](https://www.youtube.com/watch?v=55fQqPnJ4XU): This video demonstrates how to run your Dagger pipelines on Depot's Dagger Powered runners.
- [Video: Set up Dagger Cloud with Depot](https://www.youtube.com/watch?v=Y59gmLLOOT8): This video provides a step-by-step walkthrough of connecting Dagger Cloud to Depot and GitHub Actions.

## About GitHub

[GitHub](https://github.com/) is a popular Web-based platform used for version control and collaboration. It allows developers to store and manage their code in repositories, track changes over time, and collaborate with other developers on projects.

## About Depot

[Depot](https://depot.dev/) is a build acceleration platform focused on making builds more efficient and performant.