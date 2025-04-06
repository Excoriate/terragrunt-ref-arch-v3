---
slug: /ci/integrations/gitlab
---

# GitLab CI

Dagger provides a programmable container engine that allows you to replace your YAML pipeline definitions in GitLab with Dagger Functions written in a regular programming language. This allows you to execute your pipeline the same locally and in GitLab, with the additional benefit of intelligent caching.

## How it works

When running a CI pipeline with Dagger using GitLab CI, the general workflow looks like this:

1. GitLab receives a trigger based on a repository event.
1. GitLab begins processing the stages and jobs in the `.gitlab-ci.yml` file.
1. GitLab downloads the Dagger CLI.
1. GitLab executes one (or more) Dagger CLI commands, such as `dagger call ...`.
1. The Dagger CLI attempts to find an existing Dagger Engine or spins up a new one inside the GitLab runner.
1. The Dagger CLI calls the specified Dagger Function and sends telemetry to Dagger Cloud if the `DAGGER_CLOUD_TOKEN` environment variable is set.
1. The pipeline completes with success or failure. Logs appear in GitLab as usual.

## Prerequisites

- A GitLab repository
- Any one of the following:
  - [GitLab-hosted runners](https://docs.gitlab.com/ee/ci/runners/index.html) using the (default) [Docker Machine executor](https://docs.gitlab.com/runner/executors/docker_machine.html)
  - [Self-managed GitLab Runners](https://docs.gitlab.com/runner/install/index.html) using the [Docker executor](https://docs.gitlab.com/runner/executors/docker.html).
  - [Self-managed GitLab Runners](https://docs.gitlab.com/runner/install/index.html) in a Kubernetes cluster and using the [Kubernetes executor](https://docs.gitlab.com/runner/executors/kubernetes/).

## Examples

### Docker executor

The following example demonstrates how to call a Dagger Function in a GitLab CI/CD pipeline using the (default) [Docker Machine executor](https://docs.gitlab.com/runner/executors/docker_machine.html) or the [Docker executor](https://docs.gitlab.com/runner/executors/docker.html). In both these cases, the Dagger Engine is provisioned "just in time" using a Docker-in-Docker (`dind`) service.

```yaml title=".gitlab-ci.yml"
# .gitlab-ci.yml
stages:
  - hello

hello:
  stage: hello
  image: alpine:latest
  services:
    - docker:dind
  variables:
    # When using dind service we need to instruct docker to talk with the daemon started by the service.
    # The daemon is available with a network connection instead of the default /var/run/docker.sock socket.
    # See https://docs.gitlab.com/ee/ci/docker/using_docker_build.html#use-docker-in-docker-executor
    DOCKER_HOST: tcp://docker:2375
    # Instruct Docker not to start over TLS.
    DOCKER_TLS_CERTDIR: ""
    # Improve performance with overlayfs.
    DOCKER_DRIVER: overlay2
  before_script:
    # Install Dagger CLI
    - apk add --update curl && curl -L https://dl.dagger.io/dagger/install.sh | sh
    - mv bin/dagger /usr/local/bin
    - dagger version
  script:
    # Replace with your Dagger Function call
    # Example: dagger call container-echo --string-arg="Hello from Dagger!" stdout
    # Example: dagger call test build publish --tag=ttl.sh/my-image
    - dagger call container-echo --string-arg="Hello from Dagger!" stdout
```

The following is a more complex example demonstrating how to create a GitLab pipeline that checks out source code, calls a Dagger Function to test the project, and then calls another Dagger Function to build and publish a container image of the project. This example uses a simple [Go application](https://github.com/kpenfound/greetings-api) and assumes that you have already forked it in your own GitLab repository.

```yaml title=".gitlab-ci.yml"
# .gitlab-ci.yml
stages:
  - test
  - build

variables:
  # When using dind service we need to instruct docker to talk with the daemon started by the service.
  # The daemon is available with a network connection instead of the default /var/run/docker.sock socket.
  # See https://docs.gitlab.com/ee/ci/docker/using_docker_build.html#use-docker-in-docker-executor
  DOCKER_HOST: tcp://docker:2375
  # Instruct Docker not to start over TLS.
  DOCKER_TLS_CERTDIR: ""
  # Improve performance with overlayfs.
  DOCKER_DRIVER: overlay2

before_script:
  # Install Dagger CLI
  - apk add --update curl && curl -L https://dl.dagger.io/dagger/install.sh | sh
  - mv bin/dagger /usr/local/bin
  - dagger version

test:
  stage: test
  image: alpine:latest
  services:
    - docker:dind
  script:
    # Replace with your Dagger Function call
    # Example: dagger call test
    - dagger call test

build:
  stage: build
  image: alpine:latest
  services:
    - docker:dind
  script:
    # Replace with your Dagger Function call
    # Example: dagger call build publish --tag=ttl.sh/my-image
    - dagger call build publish --tag=ttl.sh/my-image
```

### Kubernetes executor

The following example demonstrates how to call a Dagger Function in a GitLab CI/CD pipeline using the [Kubernetes executor](https://docs.gitlab.com/runner/executors/kubernetes/).

```yaml title=".gitlab-ci.yml"
# .gitlab-ci.yml
stages:
  - hello

hello:
  stage: hello
  image: alpine:latest
  tags:
    - dagger-node
  before_script:
    # Install Dagger CLI
    - apk add --update curl && curl -L https://dl.dagger.io/dagger/install.sh | sh
    - mv bin/dagger /usr/local/bin
    - dagger version
  script:
    # Replace with your Dagger Function call
    # Example: dagger call container-echo --string-arg="Hello from Dagger!" stdout
    # Example: dagger call test build publish --tag=ttl.sh/my-image
    - dagger call container-echo --string-arg="Hello from Dagger!" stdout
```

The following is a more complex example demonstrating how to create a GitLab pipeline that checks out source code, calls a Dagger Function to test the project, and then calls another Dagger Function to build and publish a container image of the project. This example uses a simple [Go application](https://github.com/kpenfound/greetings-api) and assumes that you have already forked it in your own GitLab repository.


```yaml title=".gitlab-ci.yml"
# .gitlab-ci.yml
stages:
  - test
  - build

before_script:
  # Install Dagger CLI
  - apk add --update curl && curl -L https://dl.dagger.io/dagger/install.sh | sh
  - mv bin/dagger /usr/local/bin
  - dagger version

test:
  stage: test
  image: alpine:latest
  tags:
    - dagger-node
  script:
    # Replace with your Dagger Function call
    # Example: dagger call test
    - dagger call test

build:
  stage: build
  image: alpine:latest
  tags:
    - dagger-node
  script:
    # Replace with your Dagger Function call
    # Example: dagger call build publish --tag=ttl.sh/my-image
    - dagger call build publish --tag=ttl.sh/my-image
```

In both cases, each GitLab Runner must be configured to only run on nodes with pre-provisioned instances of the Dagger Engine. This is achieved using taints and tolerations on the nodes, and pod affinity.

The following code listings illustrate the configuration to be applied to each GitLab Runner, with taints, tolerations and pod affinity set via the `dagger-node` key. For an example of the corresponding node configuration, refer to the [OpenShift](./openshift.md) integration page.

To use this configuration, replace the YOUR-GITLAB-URL placeholder with the URL of your GitLab instance and replace the YOUR-GITLAB-RUNNER-TOKEN-REFERENCE placeholder with your [GitLab Runner authentication token](https://docs.gitlab.com/ee/ci/runners/runners_scope.html#create-a-shared-runner-with-a-runner-authentication-token).

```yaml title="runner-config.yml"
# runner-config.yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gitlab-runner-config
  namespace: gitlab-runner
data:
  config.toml: |
    concurrent = 10
    check_interval = 30
    log_level = "info"
    listen_address = '[::]:9252'

    [session_server]
      session_timeout = 1800

    [[runners]]
      name = "dagger-runner"
      url = "YOUR-GITLAB-URL"
      token = "YOUR-GITLAB-RUNNER-TOKEN-REFERENCE"
      executor = "kubernetes"
      [runners.kubernetes]
        namespace = "gitlab-runner"
        privileged = true
        image = "alpine:latest"
        [runners.kubernetes.node_selector]
          "dagger/enabled" = "true"
        [runners.kubernetes.node_tolerations]
          "dagger=enabled:NoSchedule" = ""
        [runners.kubernetes.affinity]
          [runners.kubernetes.affinity.nodeAffinity]
            [runners.kubernetes.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution]
              [[runners.kubernetes.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms]]
                [[runners.kubernetes.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.matchExpressions]]
                  key = "dagger/enabled"
                  operator = "In"
                  values = ["true"]
```

```yaml title="runner.yml"
# runner.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gitlab-runner
  namespace: gitlab-runner
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitlab-runner
  template:
    metadata:
      labels:
        app: gitlab-runner
    spec:
      serviceAccountName: gitlab-runner
      containers:
        - name: gitlab-runner
          image: gitlab/gitlab-runner:latest
          args:
            - run
          volumeMounts:
            - name: config
              mountPath: /etc/gitlab-runner/config.toml
              subPath: config.toml
            - name: dagger-socket
              mountPath: /var/run/dagger
      volumes:
        - name: config
          configMap:
            name: gitlab-runner-config
        - name: dagger-socket
          hostPath:
            path: /var/run/dagger
            type: DirectoryOrCreate
```


## Resources

If you have any questions about additional ways to use GitLab with Dagger, join our [Discord](https://discord.gg/dagger-io) and ask your questions in our [GitLab channel](https://discord.com/channels/707636530424053791/1122940615806685296).

## About GitLab

[GitLab](https://gitlab.com/) is a popular Web-based platform used for version control and collaboration. It allows developers to store and manage their code in repositories, track changes over time, and collaborate with other developers on projects.