---
slug: /api/terminal
---

# Interactive Terminal

Dagger provides an interactive terminal that can help greatly when trying to debug a pipeline failure.

To use this, set one or more explicit breakpoints in your Dagger pipeline with the `Container.terminal()` method. Dagger then starts an interactive terminal session at each breakpoint. This lets you inspect a `Directory` or a `Container` at any point in your pipeline run, with all the necessary context available to you.

Here is a simple example, which opens an interactive terminal in an `alpine` container:

### System shell
```shell
dagger -c 'container | from alpine | terminal'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container | from alpine | terminal
```

### Dagger CLI
```shell
dagger core container from --address=alpine terminal
```

Here is an example of a Dagger Function which opens an interactive terminal at two different points in the Dagger pipeline to inspect the built container:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Build a container and open a terminal at two different points
func (m *MyModule) Build(ctx context.Context) (*dagger.Container, error) {
	ctr := dag.Container().From("alpine:latest")

	// open terminal before adding packages
	ctr = ctr.Terminal()

	// add packages
	ctr = ctr.WithExec([]string{"apk", "add", "curl"})

	// open terminal after adding packages
	ctr = ctr.Terminal()

	return ctr, nil
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def build(self) -> dagger.Container:
        """Build a container and open a terminal at two different points"""
        ctr = dag.container().from_("alpine:latest")

        # open terminal before adding packages
        ctr = await ctr.terminal()

        # add packages
        ctr = ctr.with_exec(["apk", "add", "curl"])

        # open terminal after adding packages
        ctr = await ctr.terminal()

        return ctr

```

### TypeScript

```typescript
import { dag, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Build a container and open a terminal at two different points
   */
  @func()
  async build(): Promise<Container> {
    let ctr = dag.container().from("alpine:latest")

    // open terminal before adding packages
    ctr = await ctr.terminal()

    // add packages
    ctr = ctr.withExec(["apk", "add", "curl"])

    // open terminal after adding packages
    ctr = await ctr.terminal()

    return ctr
  }
}

```

The `Container.terminal()` method can be chained. It returns a `Container`, so it can be injected at any point in a pipeline (in this example, between `Container.from()` and `Container.withExec()` methods).

> **Tip:**
> Multiple terminals are supported in the same Dagger Function; they will open in sequence.

It's also possible to inspect a directory using the `Container.terminal()` method. Here is an example of a Dagger Function which opens an interactive terminal on a directory:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Open a terminal on a directory
func (m *MyModule) DebugDir(ctx context.Context, dir *dagger.Directory) (*dagger.Container, error) {
	return dir.Terminal(ctx)
}

```

### Python
```python
import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    async def debug_dir(self, dir: dagger.Directory) -> dagger.Container:
        """Open a terminal on a directory"""
        return await dir.terminal()

```

### TypeScript

```typescript
import { Directory, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Open a terminal on a directory
   */
  @func()
  async debugDir(dir: Directory): Promise<Container> {
    return await dir.terminal()
  }
}

```

Under the hood, this creates a new container (defaults to `alpine`) and starts a shell, mounting the directory inside. This container can be customized using additional options. Here is a more complex example, which produces the same result as the previous one but this time using an `ubuntu` container image and `bash` shell instead of the default `alpine` container image and `sh` shell:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Open a terminal on a directory using a custom container and shell
func (m *MyModule) DebugDirCustom(ctx context.Context, dir *dagger.Directory) (*dagger.Container, error) {
	return dir.Terminal(ctx, dagger.DirectoryTerminalOpts{
		Cmd:       []string{"bash"},
		Container: dag.Container().From("ubuntu:latest"),
	})
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def debug_dir_custom(self, dir: dagger.Directory) -> dagger.Container:
        """Open a terminal on a directory using a custom container and shell"""
        return await dir.terminal(
            cmd=["bash"],
            container=dag.container().from_("ubuntu:latest"),
        )

```

### TypeScript

```typescript
import { dag, Directory, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Open a terminal on a directory using a custom container and shell
   */
  @func()
  async debugDirCustom(dir: Directory): Promise<Container> {
    return await dir.terminal({
      cmd: ["bash"],
      container: dag.container().from("ubuntu:latest"),
    })
  }
}
