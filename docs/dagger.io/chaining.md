---
slug: /api/chaining
---

# Chaining

Function chaining is one of Dagger's most powerful features, as it allows you to dynamically compose complex pipelines by connecting one Dagger Function with another. The following sections demonstrate a few more examples of function chaining with the Dagger CLI.

## Execute commands in containers

The Dagger CLI can add follow-up processing to a just-in-time container, essentially enabling you to continue the pipeline directly from the command-line. `Container` objects expose a `withExec()` API method, which lets you execute a command in the corresponding container.

Here is an example of chaining a `Container.withExec()` call to a container returned by a Wolfi container builder Dagger Function, to execute a command that displays the contents of the `/etc/` directory:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-exec ls /etc/ |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-exec ls /etc/ |
  stdout
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container \
  with-exec --args="ls","/etc/" \
  stdout
```

Here is an example of chaining a `Container.withExec()` function call to a container returned by a Wolfi container builder Dagger Function, to execute a command that displays the contents of the `/etc/os-release` file:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-exec cat /etc/os-release |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-exec cat /etc/os-release |
  stdout
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container \
  with-exec --args="cat","/etc/os-release" \
  stdout
```

Here is an example of chaining a `Container.withExec()` function call to do the reverse: modify a container returned by a Wolfi container builder Dagger Function, by removing the `/etc/os-release` file from the container filesysytem:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-exec rm /etc/os-release |
  with-exec ls /etc |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-exec rm /etc/os-release |
  with-exec ls /etc |
  stdout
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container \
  with-exec --args="rm","/etc/os-release" \
  with-exec --args="ls","/etc" \
  stdout
```

Here is another example which chains multiple `Container.withExec()` calls to install the `curl` package in a Wolfi container, send an HTTP request, and return the output:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container --packages curl |
  with-exec -- curl -L dagger.io |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container --packages curl |
  with-exec -- curl -L dagger.io |
  stdout
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container --packages="curl" \
  with-exec --args="curl","-L","dagger.io" \
  stdout
```

## Live-debug container builds

`Container` objects expose a `terminal()` API method, which lets you starting an ephemeral interactive terminal session in the corresponding container. This feature is very useful for debugging and experimenting since it allows you to inspect containers directly and at any stage of your Dagger Function execution.

Here is an example of chaining a `Container.terminal()` function call to start an interactive terminal in the container returned by a Wolfi container builder Dagger Function:

### System shell
```shell
dagger -c 'github.com/dagger/dagger/modules/wolfi@v0.16.2 | container --packages=cowsay | terminal'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 | container --packages=cowsay | terminal
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container --packages="cowsay" \
  terminal
```

By default, the terminal is started with the `sh` shell, although this can be overridden by adding the `--cmd` argument. To start the same terminal with the `zsh` shell, use:

### System shell
```shell
dagger -c 'github.com/dagger/dagger/modules/wolfi@v0.16.2 | container --packages=cowsay,zsh | terminal --cmd=zsh'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 | container --packages=cowsay,zsh | terminal --cmd=zsh
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container --packages="cowsay,zsh" \
  terminal --cmd=zsh
```

## Export directories, files and containers

When a host directory or file is copied or mounted to a container's filesystem, modifications made to it in the container do not automatically transfer back to the host. Data flows only one way between Dagger operations, because they are connected in a DAG. To transfer modifications back to the local host, you must explicitly export the directory or file back to the host filesystem.

Just-in-time artifacts such as containers, directories and files can be exported to the host filesystem from the Dagger Function that produced them using the `export` function. The destination path on the host is specified using the `--path` argument.

Here is an example of exporting the build directory returned by a Go builder Dagger Function to the `./my-build` directory on the host:

### System shell
```shell
dagger <<EOF
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build ./cmd/dagger --source=https://github.com/dagger/dagger |
  export ./my-build
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build ./cmd/dagger --source=https://github.com/dagger/dagger |
  export ./my-build
```

### Dagger CLI
```shell
dagger -m github.com/kpenfound/dagger-modules/golang@v0.2.1 call \
  build --source=https://github.com/dagger/dagger --args=./cmd/dagger \
  export --path=./my-build
```

By default, the `Directory.export()` method exports the files that exist in the returned directory to the host, but it does not modify or delete any files that already exist at that host path. To replace the contents of the target host directory, such that it exactly matches the directory being exported, add the optional `--wipe` argument.

Here is an example of exporting the build directory returned by the same Dagger Function above, deleting and replacing files as needed in the `./my-build` directory on the host:

### System shell
```shell
dagger <<EOF
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build ./cmd/dagger --source=https://github.com/dagger/dagger |
  export ./my-build --wipe
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build ./cmd/dagger --source=https://github.com/dagger/dagger |
  export ./my-build --wipe
```

### Dagger CLI
```shell
dagger -m github.com/kpenfound/dagger-modules/golang@v0.2.1 call \
  build --source=https://github.com/dagger/dagger --args=./cmd/dagger \
  export --path=./my-build --wipe
```

Instead of exporting an entire directory, you can also export a file. Here is an example of exporting the compiled binary file from the build directory, as `./my-file` on the host:

### System shell
```shell
dagger <<EOF
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build ./cmd/dagger --source=https://github.com/dagger/dagger |
  file ./dagger |
  export ./my-file
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build ./cmd/dagger --source=https://github.com/dagger/dagger |
  file ./dagger |
  export ./my-file
```

### Dagger CLI
```shell
dagger -m github.com/kpenfound/dagger-modules/golang@v0.2.1 call \
  build --source=https://github.com/dagger/dagger --args=./cmd/dagger \
  file ./dagger \
  export ./my-file
```

Here is another example, this time exporting the results of a `ruff` linter Dagger Function as `/tmp/report.json` on the host:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/ruff |
  lint https://github.com/dagger/dagger |
  report |
  export /tmp/report.json
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/ruff |
  lint https://github.com/dagger/dagger |
  report |
  export /tmp/report.json
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/ruff call \
  lint --source https://github.com/dagger/dagger \
  report \
  export --path=/tmp/report.json
```

Here is an example of exporting a container returned by a Wolfi container builder Dagger Function as an OCI tarball named `/tmp/tarball.tar` on the host:

### System shell
```shell
dagger -c 'github.com/dagger/dagger/modules/wolfi@v0.16.2 | container | export ./tarball.tar'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 | container | export ./tarball.tar
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call container export --path=./tarball.tar
```

## Inspect directories, files and containers

Here is an example of listing the contents of a directory returned by a Dagger Function, by chaining a call to the `Directory.entries()` method:

### System shell
```shell
dagger <<EOF
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build . --source=https://github.com/golang/example#master:/hello |
  directory . |
  entries
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/golang@v0.2.1 |
  build . --source=https://github.com/golang/example#master:/hello |
  directory . |
  entries
```

### Dagger CLI
```shell
dagger -m github.com/kpenfound/dagger-modules/golang@v0.2.1 call \
  build --source=https://github.com/golang/example#master:/hello --args=. \
  directory --path=. \
  entries
```

Here is an example of using the `File.contents()` method to print the JSON report of a linter run:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/ruff |
  lint https://github.com/dagger/dagger |
  report |
  contents
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/ruff |
  lint https://github.com/dagger/dagger |
  report |
  contents
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/ruff call \
  lint --source=https://github.com/dagger/dagger \
  report \
  contents
```

Similar, the `File` object exposes a method to return the size of the corresponding file. Here is an example of obtaining the size of the ZIP file returned by a file archiving Dagger Function, by chaining a call to `File.size()` method:

### System shell
```shell
dagger <<EOF
github.com/sagikazarmark/daggerverse/arc@40057665476af62e617cc8def9ef5a87735264a9 |
  archive-directory dagger-cli https://github.com/dagger/dagger#main:cmd/dagger |
  create zip |
  size
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/sagikazarmark/daggerverse/arc@40057665476af62e617cc8def9ef5a87735264a9 |
  archive-directory dagger-cli https://github.com/dagger/dagger#main:cmd/dagger |
  create zip |
  size
```

### Dagger CLI
```shell
dagger -m github.com/sagikazarmark/daggerverse/arc@40057665476af62e617cc8def9ef5a87735264a9 call \
  archive-directory --name=dagger-cli '--directory=https://github.com/dagger/dagger#main:cmd/dagger' \
  create --format=zip \
  size
```

## Publish containers

Every `Container` object exposes a `Container.publish()` API method, which publishes the container as a new image to a specified container registry. The registry address is passed to the function using the `--address` argument, and the return value is a string referencing the container image address in the registry.

Here is an example of publishing the container returned by a Wolfi container builder Dagger Function to the `ttl.sh` registry, by chaining a `Container.publish()` call:

### System shell
```shell
dagger -c 'github.com/dagger/dagger/modules/wolfi@v0.16.2 | container | publish ttl.sh/my-wolfi'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 | container | publish ttl.sh/my-wolfi
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call container publish --address=ttl.sh/my-wolfi
```

## Start containers as services

Every `Container` object exposes a `Container.asService()` API method, which turns the container into a `Service`. These services can then be spun up for use by other Dagger Functions or by clients on the Dagger host by forwarding their ports. This is akin to a "programmable docker-compose".

To start a `Service` returned by a Dagger Function and have it forward traffic to a specified address via the host, chain a call to the `Service.up()` API method.

Here is an example of starting an NGINX service on host port 80 by chaining calls to `Container.asService()` and `Service.up()`:

### System shell
```shell
dagger -c 'github.com/kpenfound/dagger-modules/nginx@v0.1.0 | container | as-service | up'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/nginx@v0.1.0 | container | as-service | up
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call container as-service up
```

By default, each port maps to the same port on the host. To specify a different mapping, use the additional `--ports` argument with a list of host/service port mappings. To bind ports randomly, use the `--random` argument.

To start the same service and map NGINX port 80 to host port 8080, use:

### System shell
```shell
dagger -c 'github.com/kpenfound/dagger-modules/nginx@v0.1.0 | container | as-service | up --ports=8080:80'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/nginx@v0.1.0 | container | as-service | up --ports=8080:80
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call container as-service up --ports=8080:80
```

The service can now be accessed on the specified port. For example, in another terminal, execute the following command to receive the default NGINX welcome page:

```shell
curl localhost:8080
```

To start the same service and map NGINX port 80 to a random port on the host, use:

### System shell
```shell
dagger -c 'github.com/kpenfound/dagger-modules/nginx@v0.1.0 | container | as-service | up --random'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/nginx@v0.1.0 | container | as-service | up --random
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call container as-service up --random
```

## Modify container filesystems

Here is an example of modifying a container by adding the current directory from the host to the container filesysytem at `/src`, by chaining a call to the `Container.withDirectory()` method:

> **Warning:**
> The example below uploads the entire current directory to the container filesystem. This can take a significant amount of time with large directories. To reduce the time spent on upload, run this example from a directory containing only a few small files.

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-directory /src . |
  with-exec ls /src |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-directory /src . |
  with-exec ls /src |
  stdout
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container \
  with-directory --path=/src --directory=. \
  with-exec --args="ls","/src" \
  stdout
```

Here is an example of passing a `README.md` file from the host to a container builder Dagger Function, by chaining a call to the `Container.withFile()` function:

### System shell
```shell
dagger <<EOF
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-file /README.md $(host | file ./README.md) |
  with-exec cat /README.md |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/dagger/dagger/modules/wolfi@v0.16.2 |
  container |
  with-file /README.md $(host | file ./README.md) |
  with-exec cat /README.md |
  stdout
```

### Dagger CLI
```shell
dagger -m github.com/dagger/dagger/modules/wolfi@v0.16.2 call \
  container \
  with-file --path=/README.md --source=./README.md \
  with-exec --args="cat","/README.md" \
  stdout