---
slug: /api/filters
---

# Directory Filters

When you pass a directory to a Dagger Function as argument, Dagger uploads everything in that directory tree to the Dagger Engine. For large monorepos or directories containing large-sized files, this can significantly slow down your Dagger Function while filesystem contents are transferred. To mitigate this problem, Dagger lets you apply filters to control which files and directories are uploaded.

## Directory arguments

Dagger Functions do not have access to the filesystem of the host you invoke the Dagger Function from (i.e. the host you execute a CLI command like `dagger` from). Instead, host files and directories need to be explicitly passed as command-line arguments to Dagger Functions.

There are two important reasons for this.

- Reproducibility: By providing a call-time mechanism to define and control the files available to a Dagger Function, Dagger guards against creating hidden dependencies on ambient properties of the host filesystem that could change at any moment.
- Security: By forcing you to explicitly specify which host files and directories a Dagger Function "sees" on every call, Dagger ensures that you're always 100% in control. This reduces the risk of third-party Dagger Functions gaining access to your data.

To tell a Dagger Function which directory to use, specify its path as an argument when using `dagger call`. Here's a simple example, which passes a directory from the host (`./example/hello`) to a Dagger Function:

```
git clone https://github.com/golang/example
dagger -m github.com/kpenfound/dagger-modules/golang@v0.2.0 call build --source=./example/hello --args=. directory --p
```

The important thing to know here is that, by default, Dagger will copy and upload everything in the specified directory and its sub-directories to the Dagger Engine. For complex directory trees, directories containing a large number of files, or directories containing large-sized files, this can add minutes to your Dagger Function execution time while the contents are transferred.

Dagger offers pre- and post-call filtering to mitigate this problem and optimize how your directories are handled.

## Why filter?

Filtering improves the performance of your Dagger Functions in three ways:

- It reduces the size of the files being transferred from the host to the Dagger Engine, allowing the upload step to complete faster.
- It ensures that minor unrelated changes in the source directory don't invalidate Dagger's build cache.
- It enables different use-cases, such as setting up component/feature/service-specific pipelines for monorepos.

It is worth noting that Dagger already uses caching to optimize file uploads. Subsequent calls to a Dagger Function will only upload files that have changed since the preceding call. Filtering is an additional optimization that you can apply to improve the performance of your Dagger Function.

## Pre-call filtering

Pre-call filtering means that a directory is filtered before it's uploaded to the Dagger Engine container. This is useful for:

- Large monorepos. Typically your Dagger Function only operates on a subset of the monorepo, representing a specific component or feature. Uploading the entire worktree imposes a prohibitive cost.

- Large files, such as audio/video files and other binary content. These files take time to upload. If they're not directly relevant, you'll usually want your Dagger Function to ignore them.

  > **Tip:**
  > The `.git` directory is a good example of both these cases. It contains a lot of data, including large binary objects, and for projects with a long version history, it can sometimes be larger than your actual source code.

- Dependencies. If you're developing locally, you'll typically have your project dependencies installed locally: `node_modules` (Node.js), `.venv` (Python), `vendor` (PHP) and so on. When you call your Dagger Function locally, Dagger will upload all these installed dependencies as well. This is both bad practice and inefficient. Typically, you'll want your Dagger Function to ignore locally-installed dependencies and only operate on the project source code.

> **Note:**
> At the time of writing, Dagger [does not read exclusion patterns from existing `.dockerignore`/`.gitignore` files](https://github.com/dagger/dagger/issues/6627). If you already use these files, you'll need to manually implement the same patterns in your Dagger Function.

To implement a pre-call filter in your Dagger Function, add an `ignore` parameter to your `Directory` argument. The `ignore` parameter follows the [`.gitignore` syntax](https://git-scm.com/docs/gitignore). Some important points to keep in mind are:

- The order of arguments is significant: the pattern `"**", "!**"` includes everything but `"!**", "**"` excludes everything.
- Prefixing a path with `!` negates a previous ignore: the pattern `"!foo"` has no effect, since nothing is previously ignored, while the pattern `"**", "!foo"` excludes everything except `foo`.

### Go

Here's an example of a Dagger Function that excludes everything in a given directory except Go source code files:

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns the list of Go files in the directory
func (m *MyModule) ListGoFiles(
	ctx context.Context,
	// Source directory
	// +ignore=["*", "!**/*.go"]
	source *dagger.Directory,
) ([]string, error) {
	return source.Entries(ctx)
}

```

### Python

Here's an example of a Dagger Function that excludes everything in a given directory except Python source code files:

```python
from typing import Annotated, List

import dagger
from dagger import Doc, Ignore, function, object_type


@object_type
class MyModule:
    @function
    async def list_python_files(
        self,
        source: Annotated[
            dagger.Directory,
            Doc("Source directory"),
            Ignore(["*", "!**/*.py"]),
        ],
    ) -> List[str]:
        """Returns the list of Python files in the directory"""
        return await source.entries()

```

### TypeScript

Here's an example of a Dagger Function that excludes everything in a given directory except TypeScript source code files:

```typescript
import { Directory, func, object, argument } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns the list of TypeScript files in the directory
   *
   * @param source Source directory
   */
  @func()
  async listTypescriptFiles(
    @argument({ description: "Source directory", ignore: ["*", "!**/*.ts"] })
    source: Directory,
  ): Promise<string[]> {
    return await source.entries()
  }
}

```

### PHP

Here's an example of a Dagger Function that excludes everything in a given directory except PHP source code files:

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Attribute\Ignore;
use Dagger\Client\Directory;

#[DaggerObject]
class MyModule
{
    /**
     * Returns the list of PHP files in the directory
     *
     * @param Directory $source Source directory
     * @return string[]
     */
    #[DaggerFunction]
    public function listPhpFiles(
        #[Ignore('*', '!**/*.php')]
        Directory $source
    ): array {
        return $source->entries();
    }
}

```

### Java

Here's an example of a Dagger Function that excludes everything in a given directory except Java source code files:

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Directory;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;
import io.dagger.module.annotation.Ignore;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns the list of Java files in the directory
   */
  @Function
  public List<String> listJavaFiles(
      @Description("Source directory") @Ignore({"*", "!**/*.java"}) Directory source) throws Exception {
    return source.entries().get();
  }
}

```

Here are a few examples of useful patterns:

### Go
```go
// exclude Go tests and test data
+ignore=["**_test.go", "**/testdata/**"]

// exclude binaries
+ignore=["bin"]

// exclude Python dependencies
+ignore=["**/.venv", "**/__pycache__"]

// exclude Node.js dependencies
+ignore=["**/node_modules"]

// exclude Git metadata
+ignore=[".git", "**/.gitignore"]
```

### Python
```python
# exclude Pytest tests and test data
Ignore(["tests/", ".pytest_cache"])

# exclude binaries
Ignore(["bin"])

# exclude Python dependencies
Ignore(["**/.venv", "**/__pycache__"])

# exclude Node.js dependencies
Ignore(["**/node_modules"])

# exclude Git metadata
Ignore([".git", "**/.gitignore"])
```

### TypeScript
```typescript
// exclude Mocha tests
@argument({ ignore: ["**.spec.ts"] })

// exclude binaries
@argument({ ignore: ["bin"] })

// exclude Python dependencies
@argument({ ignore: ["**/.venv", "**/__pycache__"] })

// exclude Node.js dependencies
@argument({ ignore: ["**/node_modules"] })

// exclude Git metadata
@argument({ ignore: [".git", "**/.gitignore"] })
```

### PHP
```php
// exclude PHPUnit tests and test data
#[Ignore('tests/', '.phpunit.cache', '.phpunit.result.cache')]

// exclude binaries
#[Ignore('bin')]

// exclude Composer dependencies
#[Ignore('vendor/')]

// exclude Node.js dependencies
#[Ignore('**/node_modules')]

// exclude Git metadata
#[Ignore('.git/', '**/.gitignore')]
```

### Java
```java
// exclude Java tests and test data
@Ignore({"src/test"})

// exclude binaries
@Ignore({"bin"})

// exclude Python dependencies
@Ignore({"**/.venv", "**/__pycache__"})

// exclude Node.js dependencies
@Ignore({"**/node_modules"})

// exclude Git metadata
@Ignore({".git", "**/.gitignore"})
```

## Post-call filtering

Post-call filtering means that a directory is filtered after it's uploaded to the Dagger Engine.

This is useful when working with directories that are modified "in place" by a Dagger Function. When building an application, your Dagger Function might modify the source directory during the build by adding new files to it. A post-call filter allows you to use that directory in another operation, only fetching the new files and ignoring the old ones.

A good example of this is a multi-stage build. Imagine a Dagger Function that reads and builds an application from source, placing the compiled binaries in a new sub-directory (stage 1). Instead of then transferring everything to the final container image for distribution (stage 2), you could use a post-call filter to transfer only the compiled files.

### Go

To implement a post-call filter in your Dagger Function, use the `DirectoryWithDirectoryOpts` or `ContainerWithDirectoryOpts` structs, which support `Include` and `Exclude` patterns for `Directory` objects. Here's an example:

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Build the application and return the build directory
func (m *MyModule) Build(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "build/app"}).
		Directory("build"). // return the build directory
		Sync(ctx)
}

// Build and package the application in a container
func (m *MyModule) Package(ctx context.Context, source *dagger.Directory) (*dagger.Container, error) {
	buildDir, err := m.Build(ctx, source)
	if err != nil {
		return nil, err
	}

	return dag.Container().
		From("alpine:latest").
		// copy only the compiled binary from the build directory
		WithDirectory("/usr/local/bin", buildDir, dagger.ContainerWithDirectoryOpts{
			Include: []string{"app"},
		}).
		Sync(ctx)
}

```

### Python

To implement a post-call filter in your Dagger Function, use the `include` and `exclude` parameters when working with `Directory` objects. Here's an example:

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def build(self, source: dagger.Directory) -> dagger.Directory:
        """Build the application and return the build directory"""
        return await (
            dag.container()
            .from_("golang:latest")
            .with_directory("/src", source)
            .with_workdir("/src")
            .with_exec(["go", "build", "-o", "build/app"])
            .directory("build")  # return the build directory
        )

    @function
    async def package(self, source: dagger.Directory) -> dagger.Container:
        """Build and package the application in a container"""
        build_dir = await self.build(source)

        return (
            dag.container()
            .from_("alpine:latest")
            # copy only the compiled binary from the build directory
            .with_directory(
                "/usr/local/bin",
                build_dir,
                include=["app"],
            )
        )

```

### TypeScript

To implement a post-call filter in your Dagger Function, use the `include` and `exclude` parameters when working with `Directory` objects. Here's an example:

```typescript
import { dag, Directory, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Build the application and return the build directory
   */
  @func()
  async build(source: Directory): Promise<Directory> {
    return await dag
      .container()
      .from("golang:latest")
      .withDirectory("/src", source)
      .withWorkdir("/src")
      .withExec(["go", "build", "-o", "build/app"])
      .directory("build") // return the build directory
  }

  /**
   * Build and package the application in a container
   */
  @func()
  async package(source: Directory): Promise<Container> {
    const buildDir = await this.build(source)

    return dag
      .container()
      .from("alpine:latest")
      // copy only the compiled binary from the build directory
      .withDirectory("/usr/local/bin", buildDir, {
        include: ["app"],
      })
  }
}

```

### PHP

To implement a post-call filter in your Dagger Function, use the `include` and `exclude` parameters when working with `Directory` objects. Here's an example:

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Client\Container;
use Dagger\Client\Directory;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Build the application and return the build directory
     */
    #[DaggerFunction]
    public function build(Directory $source): Directory
    {
        return dag()
            ->container()
            ->from('golang:latest')
            ->withDirectory('/src', $source)
            ->withWorkdir('/src')
            ->withExec(['go', 'build', '-o', 'build/app'])
            ->directory('build'); // return the build directory
    }

    /**
     * Build and package the application in a container
     */
    #[DaggerFunction]
    public function package(Directory $source): Container
    {
        $buildDir = $this->build($source);

        return dag()
            ->container()
            ->from('alpine:latest')
            // copy only the compiled binary from the build directory
            ->withDirectory(
                '/usr/local/bin',
                $buildDir,
                include: ['app']
            );
    }
}

```

### Java

To implement a post-call filter in your Dagger Function, use the `Container.WithDirectoryArguments` class which support `withInclude` and `withExclude` functions when working with `Directory` objects. Here's an example:

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.client.Directory;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Build the application and return the build directory
   */
  @Function
  public Directory build(Directory source) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("golang:latest")
          .withDirectory("/src", source)
          .withWorkdir("/src")
          .withExec(List.of("go", "build", "-o", "build/app"))
          .directory("build") // return the build directory
          .sync();
    }
  }

  /**
   * Build and package the application in a container
   */
  @Function
  public Container package_(Directory source) throws Exception {
    Directory buildDir = this.build(source);
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("alpine:latest")
          // copy only the compiled binary from the build directory
          .withDirectory(
              "/usr/local/bin",
              buildDir,
              new Container.WithDirectoryArguments().withInclude(List.of("app")))
          .sync();
    }
  }
}

```

Here are a few examples of useful patterns:

### Go
```go
// exclude all Markdown files
dirOpts := dagger.ContainerWithDirectoryOpts{
  Exclude: "*.md*",
}

// include only the build output directory
dirOpts := dagger.ContainerWithDirectoryOpts{
  Include: "build",
}

// include only ZIP files
dirOpts := dagger.DirectoryWithDirectoryOpts{
  Include: "*.zip",
}

// exclude Git metadata
dirOpts := dagger.DirectoryWithDirectoryOpts{
  Exclude: "*.git",
}
```

### Python
```python
# exclude all Markdown files
dir_opts = {"exclude": ["*.md*"]}

# include only the build output directory
dir_opts = {"include": ["build"]}

# include only ZIP files
dir_opts = {"include": ["*.zip"]}

# exclude Git metadata
dir_opts = {"exclude": ["*.git"]}
```

### TypeScript
```typescript
// exclude all Markdown files
const dirOpts = { exclude: ["*.md*"] }

// include only the build output directory
const dirOpts = { include: ["build"] }

// include only ZIP files
const dirOpts = { include: ["*.zip"] }

// exclude Git metadata
const dirOpts = { exclude: ["*.git"] }
```

### PHP
```php
// exclude all Markdown files
$dirOpts = ['exclude' => ['*.md*']];

// include only the build output directory
$dirOpts = ['include' => ['build']];

// include only ZIP files
$dirOpts = ['include' => ['*.zip']];

// exclude Git metadata
$dirOpts = ['exclude' => ['*.git']];
```

### Java
```java
// exclude all Markdown files
var dirOpts = new Container.WithDirectoryArguments()
    .withExclude(List.of("*.md*"));

// include only the build output directory
var dirOpts = new Container.WithDirectoryArguments()
    .withInclude(List.of("build"));

// include only ZIP files
var dirOpts = new Container.WithDirectoryArguments()
    .withInclude(List.of("*.zip"));

// exclude Git metadata
var dirOpts = new Container.WithDirectoryArguments()
    .withExclude(List.of("*.git"));
```

## Debugging

### Using logs

Both Dagger Cloud and the Dagger TUI provide detailed information on the patterns Dagger uses to filter your directory uploads - look for the upload step in the TUI logs or Trace:

![Dagger TUI](/img/current_docs/api/fs-filters-tui.png)

![Dagger Cloud Trace](/img/current_docs/api/fs-filters-trace.png)

### Inspecting directory contents

Another way to debug how directories are being filtered is to create a function that receives a `Directory` as input, and returns the same `Directory`:

### Go
```go
func (m *MyModule) Debug(
  ctx context.Context,
  // +ignore=["*", "!analytics"]
  source *dagger.Directory,
) *dagger.Directory {
  return source
}
```

### Python
```python
@function
async def foo(
    self,
    source: Annotated[
        dagger.Directory, Ignore(["*", "!analytics"])
    ],
) -> dagger.Directory:
    return source
```

### TypeScript
```typescript
@func()
debug(
   @argument({ ignore: ["*", "!analytics"] }) source: Directory,
): Directory {
  return source
}
```

### PHP
```php
    #[DaggerFunction]
    public function debug(
        #[Ignore('*'/, '!analytics')]
        Directory $source,
    ): Directory {
        return $source;
    }
```

### Java
```java
@Function
public Directory debug(@Ignore({"*", "!analytics"}) Directory source) {
    return source;
}
```

Calling the function will show you the directory’s digest and top level entries. The digest is content addressed, so it changes if there are changes in the contents of the directory. Looking at the entries field you may be able to spot an interloper:

### System shell
```shell
dagger -c 'debug .'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
debug .
```

### Dagger CLI
```shell
dagger call debug --source=.
```

You can also list all files, recursively to check it more deeply:

### System shell
```shell
dagger -c 'debug . | glob "**/*"'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
debug . | glob "**/*"
```

### Dagger CLI
```shell
dagger call debug --source=. glob --pattern="**/*"
```

You can open the directory in an interactive terminal to inspect the filesystem:

### System shell
```shell
dagger -c 'debug . | terminal'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
debug . | terminal
```

### Dagger CLI
```shell
dagger call debug --source=. terminal
```

You can export the filtered directory to your host and check it with local tools:

### System shell
```shell
dagger -c 'debug . | export audit'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
debug . | export audit
```

### Dagger CLI
```shell
dagger call debug --source=. export --path=audit