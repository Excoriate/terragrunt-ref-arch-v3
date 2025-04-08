---
slug: /ci/integrations/github
---

# GitHub

Dagger can directly interact with GitHub pull requests, making it easy to test the functionality of specific forks or branches of a GitHub repository.

Dagger also supports publishing container images to [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry).

## How it works

GitHub contains a shorthand redirect at the `/merge` endpoint that allows you to reference the correct branch of a repository from a pull request (PR), without needing to know anything about the fork or branch where the PR came from.

By default, the Dagger `Directory` type works with both local directories and [remote Git repositories](../../cookbook/cookbook.md#copy-a-directory-or-remote-repository-to-a-container). This makes it possible to work with the directory tree at the root of a Git repository or a given branch.

By combining these two features, Dagger users can write Dagger Functions that directly use GitHub pull requests as arguments.

## Prerequisites

- A GitHub repository

## Examples

Given a Dagger Function called `foo` that accepts a `Directory` as argument, you can pass it a GitHub pull request URL as argument like this:

```shell
dagger call foo --directory=https://github.com/ORGANIZATION/REPOSITORY#pull/NUMBER/merge
```

If your GitHub repository contains a Dagger module, you can test the functionality of a specific branch by calling the Dagger module with the corresponding pull request URL, as shown below:

```shell
dagger call -m github.com/ORGANIZATION/REPOSITORY@pull/NUMBER/merge --help
```

You can also use a Dagger Function in a GitHub Actions workflow to publish a container image to GitHub Container Registry.

```yaml title=".github/workflows/dagger.yml"
# .github/workflows/dagger.yml
name: Dagger

on:
  push:
    branches: ["main"]

jobs:
  build:
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
          args: "call build publish --tag=ghcr.io/${{ github.repository }}:${{ github.sha }}"
          # Optional: Dagger Cloud token
          # cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }}
        env:
          # Required: GitHub token to publish container image to GitHub Container Registry
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Resources

If you have any questions about additional ways to use GitHub with Dagger, join our [Discord](https://discord.gg/dagger-io) and ask your questions in our [GitHub channel](https://discord.com/channels/707636530424053791/1117139064274034809).

## About GitHub

[GitHub](https://github.com/) is a popular Web-based platform used for version control and collaboration. It allows developers to store and manage their code in repositories, track changes over time, and collaborate with other developers on projects.