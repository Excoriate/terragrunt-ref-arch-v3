---
slug: /api/module-dependencies
---

# Module Dependencies

## Installation

You can call Dagger Functions from any other Dagger module in your own Dagger module simply by adding it as a module dependency with `dagger install`, as in the following example:

```shell
dagger install github.com/shykes/daggerverse/hello@v0.3.0
```

This module will be added to your `dagger.json`:

```json
...
"dependencies": [
  {
    "name": "hello",
    "source": "github.com/shykes/daggerverse/hello@54d86c6002d954167796e41886a47c47d95a626d"
  }
]
```

When you add a dependency to your module with `dagger install`, the dependent module will be added to the code-generation routines and can be accessed from your own module's code.

The entrypoint to accessing dependent modules from your own module's code is `dag`, the Dagger client, which is pre-initialized. It contains all the core types (like `Container`, `Directory`, etc.), as well as bindings to any dependencies your module has declared.

Here is an example of accessing the installed `hello` module from your own module's code:

### Go

```go
func (m *MyModule) Greeting(ctx context.Context) (string, error) {
  return dag.Hello().Hello(ctx)
}
```

### Python

```python
@function
async def greeting(self) -> str:
  return await dag.hello().hello()
```

### TypeScript

```typescript
@func()
async greeting(): Promise<string> {
  return await dag.hello().hello()
}
```

### PHP

```php
#[DaggerFunction]
public function greeting(): string
{
    return dag()->hello()->hello();
}
```

### Java

```java
@Function
public String greeting() throws ExecutionException, DaggerQueryException, InterruptedException {
    return dag().hello().hello();
}
```

Here is a more complex example. It is a Dagger Function that utilizes a module from the Daggerverse to build a Go project, then chains a Dagger API method to open an interactive terminal session in the build directory.

First, install the module:

```shell
dagger install github.com/kpenfound/dagger-modules/golang@v0.2.0
```

Next, create a new Dagger Function:

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Build a Go project and open a terminal in the build directory
func (m *MyModule) Example(
	ctx context.Context,
	// Source directory
	buildSrc *dagger.Directory,
	// Build arguments
	buildArgs []string,
) *dagger.Container {
	// Call the Golang module's build function
	buildDir := dag.Golang().Build(buildSrc, buildArgs)

	// Open a terminal in the build directory
	return buildDir.Terminal(ctx)
}

```

### Python
```python
from typing import List

import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def example(
        self,
        build_src: dagger.Directory,
        build_args: List[str],
    ) -> dagger.Container:
        """Build a Go project and open a terminal in the build directory"""
        # Call the Golang module's build function
        build_dir = dag.golang().build(build_src, build_args)

        # Open a terminal in the build directory
        return await build_dir.terminal()

```

### TypeScript
```typescript
import { dag, Directory, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Build a Go project and open a terminal in the build directory
   */
  @func()
  async example(
    buildSrc: Directory,
    buildArgs: string[],
  ): Promise<Container> {
    // Call the Golang module's build function
    const buildDir = dag.golang().build(buildSrc, buildArgs)

    // Open a terminal in the build directory
    return await buildDir.terminal()
  }
}

```

### PHP
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
     * Build a Go project and open a terminal in the build directory
     *
     * @param Directory $buildSrc Source directory
     * @param string[] $buildArgs Build arguments
     */
    #[DaggerFunction]
    public function example(
        Directory $buildSrc,
        array $buildArgs
    ): Container {
        // Call the Golang module's build function
        $buildDir = dag()->golang()->build($buildSrc, $buildArgs);

        // Open a terminal in the build directory
        return $buildDir->terminal();
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.client.Directory;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Build a Go project and open a terminal in the build directory
   */
  @Function
  public Container example(
      @Description("Source directory") Directory buildSrc,
      @Description("Build arguments") List<String> buildArgs) throws Exception {
    try (Client client = Dagger.connect()) {
      // Call the Golang module's build function
      Directory buildDir = client.golang().build(buildSrc, buildArgs);

      // Open a terminal in the build directory
      return buildDir.terminal();
    }
  }
}

```

This Dagger Function accepts two arguments - the source directory and a list of build arguments - and does the following:

- It invokes the Golang module via the `dag` Dagger client.
- It calls a Dagger Function from the module to build the source code and return a just-in-time directory with the compiled binary.
- It chains a core Dagger API method to open an interactive terminal session in the returned directory.

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'example https://github.com/golang/example#master:/hello .'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
example https://github.com/golang/example#master:/hello .
```

### Dagger CLI
```shell
dagger call example --build-src=https://github.com/golang/example#master:/hello --build-args=.
```

You can also use local modules as dependencies. However, they must be stored in a sub-directory of your module. For example:

```shell
dagger install ./path/to/module
```

> **Note:**
> Installing a module using a local path (relative or absolute) is only possible if your module is within the repository root (for Git repositories) or the directory containing the `dagger.json` file (for all other cases).

## Uninstallation

To remove a dependency from your Dagger module, use the `dagger uninstall` command. The `dagger uninstall` command can be passed either a remote repository reference or a local module name.

The commands below are equivalent:

```shell
dagger uninstall hello
dagger uninstall github.com/shykes/daggerverse/hello
```

## Update

To update a dependency in your Dagger module to the latest version (or the version specified), use the `dagger update` command. The target module must be local.

The commands below are equivalent:

```shell
dagger update hello
dagger update github.com/shykes/daggerverse/hello
```

> **Note:**
> Given a dependency like `github.com/path/name@branch/tag`:
> - `dagger update github.com/path/name` updates the dependency to the latest commit of the branch/tag.
> - `dagger update github.com/path/name@version` updates the dependency to the latest commit for the `version` branch/tag.