---
slug: /ci/integrations/circleci
---

# CircleCI

Dagger provides a programmable container engine that allows you to replace your YAML workflows in CircleCI with Dagger Functions written in a regular programming language. This allows you to execute your pipeline the same locally and in CI, with the additional benefit of intelligent caching.

## How it works

When running a CI pipeline with Dagger using CircleCI, the general workflow looks like this:

1. CircleCI receives a trigger based on a repository event.
1. CircleCI begins processing the jobs and steps in the `.circleci/config.yml` workflow file.
1. CircleCI downloads the Dagger CLI.
1. CircleCI executes one (or more) Dagger CLI commands, such as `dagger call ...`.
1. The Dagger CLI attempts to find an existing Dagger Engine or spins up a new one inside the CircleCI runner.
1. The Dagger CLI calls the specified Dagger Function and sends telemetry to Dagger Cloud if the `DAGGER_CLOUD_TOKEN` environment variable is set.
1. The pipeline completes with success or failure. Logs appear in CircleCI as usual.

> **Note:**
> In a Dagger context, you won't have access to CircleCI's test splitting functionality. You will need to implement your own test distribution logic or run all tests in a single execution.

## Prerequisites

- A CircleCI project
- A GitHub, Bitbucket or GitLab repository connected to the CircleCI project
- Docker, if using a [CircleCI execution environment](https://circleci.com/docs/executor-intro) other than `docker`

## Examples

The examples below use the `docker` executor, which come with a Docker execution environment preconfigured. If using a [different executor](https://circleci.com/docs/executor-intro), such as `machine`, you must install Docker in the execution environment before proceeding with the examples.

The following example demonstrates how to call a Dagger Function in a CircleCI workflow.

```yaml title=".circleci/config.yml"
# .circleci/config.yml
version: 2.1

jobs:
  hello:
    docker:
      - image: cimg/base:stable
    steps:
      - setup_remote_docker:
          version: 20.10.14
          docker_layer_caching: true
      - run:
          name: Install Dagger CLI
          command: |
            cd /usr/local
            curl -L https://dl.dagger.io/dagger/install.sh | sh
            cd bin
            sudo mv dagger /usr/local/bin
            dagger version
      - run:
          name: Run Dagger Function
          command: |
            # Replace with your Dagger Function call
            # Example: dagger call container-echo --string-arg="Hello from Dagger!" stdout
            # Example: dagger call test build publish --tag=ttl.sh/my-image
            dagger call container-echo --string-arg="Hello from Dagger!" stdout

workflows:
  dagger:
    jobs:
      - hello
```

The following is a more complex example demonstrating how to create a CircleCI workflow that checks out source code, calls a Dagger Function to test the project, and then calls another Dagger Function to build and publish a container image of the project. This example uses a simple [Go application](https://github.com/kpenfound/greetings-api) and assumes that you have already forked it in the repository connected to the CircleCI project.

```yaml title=".circleci/config.yml"
# .circleci/config.yml
version: 2.1

jobs:
  test:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - setup_remote_docker:
          version: 20.10.14
          docker_layer_caching: true
      - run:
          name: Install Dagger CLI
          command: |
            cd /usr/local
            curl -L https://dl.dagger.io/dagger/install.sh | sh
            cd bin
            sudo mv dagger /usr/local/bin
            dagger version
      - run:
          name: Run Dagger Function
          command: |
            # Replace with your Dagger Function call
            # Example: dagger call test
            dagger call test

  build:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - setup_remote_docker:
          version: 20.10.14
          docker_layer_caching: true
      - run:
          name: Install Dagger CLI
          command: |
            cd /usr/local
            curl -L https://dl.dagger.io/dagger/install.sh | sh
            cd bin
            sudo mv dagger /usr/local/bin
            dagger version
      - run:
          name: Run Dagger Function
          command: |
            # Replace with your Dagger Function call
            # Example: dagger call build publish --tag=ttl.sh/my-image
            dagger call build publish --tag=ttl.sh/my-image

workflows:
  dagger:
    jobs:
      - test
      - build:
          requires:
            - test
```

## Resources

If you have any questions about additional ways to use CircleCI with Dagger, join our [Discord](https://discord.gg/dagger-io) and ask your questions in our [help channel](https://discord.com/channels/707636530424053791/1030538312508776540).

## About CircleCI

[CircleCI](https://circleci.com/) is a popular CI/CD platform to test, build and deploy software applications.