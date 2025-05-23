---
slug: /api
---

# Dagger API

The Dagger API is an unified interface for composing Dagger pipelines. It provides a set of core functions and core types for creating and managing application delivery pipelines, including types for containers, files, directories, network services, secrets, and more. You can call these functions from `dagger`, the Dagger CLI, or from custom functions written in any supported programming language.

## Functions

To create a pipeline with the Dagger API, you call multiple functions, combining them together in sequence to form a pipeline. Here's an example:

### System shell
```shell
dagger <<EOF
container |
  from alpine |
  file /etc/os-release |
  contents
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from alpine |
  file /etc/os-release |
  contents
```

### Dagger CLI
```shell
dagger core container \
  from --address=alpine \
  file --path=/etc/os-release \
  contents
```

This example calls multiple functions from the Dagger API in sequence to create a pipeline that builds an image from an `alpine` container and returns the contents of the `/etc/os-release` file to the caller.

## Types

Like regular functions, Dagger functions can accept arguments. You can see this in the previous example: the `from` function accepts a container image address (`--address`) as argument, the file function accepts a filesystem location (`--path`) as argument, and so on.

In addition to supporting basic types (string, boolean, integer, array...), the Dagger API also provides [powerful core types](./types.md) for working with pipelines, such as `Container`, `GitRepository`, `Directory`, `Service`, `Secret`. You can use these core types as arguments to Dagger functions.

Here's an example:

### System shell
```shell
dagger <<EOF
container |
  with-directory /src https://github.com/dagger/dagger |
  directory /src |
  entries
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  with-directory /src https://github.com/dagger/dagger |
  directory /src |
  entries
```

### Dagger CLI
```shell
dagger core container \
  with-directory --path=/src --directory=https://github.com/dagger/dagger \
  directory --path=/src \
  entries
```

This example creates a scratch container, adds Dagger's GitHub repository to it, and returns the directory contents. In this case, the address of the GitHub repository is passed to the function as a `Directory` argument. The `Directory` loaded in this manner is not merely a string, but it is the actual filesystem state of the target directory, managed by the Dagger Engine and handled in code just like any another variable.

## Chaining

Each of Dagger's core types comes with functions of its own, which can be used to interact with the corresponding object.

When calling a Dagger function that returns a core type, the Dagger API lets you follow up by calling one of that type's functions, which itself can return another type, and so on. This is called "function chaining", and is a core feature of Dagger.

For example, if a Dagger function returns a `Directory`, the caller can continue the chain by calling a function from the `Directory` type to export it to the local filesystem, modify it, mount it into a container, and so on.

Although you may not have realized it, you've already seen function chaining in action. Both the previous examples chain functions together into pipelines. Here is one more example to illustrate the concept:

### System shell
```shell
dagger <<EOF
container |
  from golang:latest |
  with-directory /src https://github.com/dagger/dagger#main |
  with-workdir /src/cmd/dagger |
  with-exec -- go build -o dagger . |
  file ./dagger |
  export ./dagger.bin
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from golang:latest |
  with-directory /src https://github.com/dagger/dagger#main |
  with-workdir /src/cmd/dagger |
  with-exec -- go build -o dagger . |
  file ./dagger |
  export ./dagger.bin
```

### Dagger CLI
```shell
dagger core container from --address="golang:latest" \
  with-directory --path="/src" --directory="https://github.com/dagger/dagger#main" \
  with-workdir --path="/src/cmd/dagger" \
  with-exec --args="go","build","-o","dagger","." \
  file --path="./dagger" \
  export --path="./dagger.bin"
```

This example chains multiple function calls into a pipeline that builds the Dagger CLI from source and exports it to the Dagger host:
- `from` returns a `golang` container image as a `Container` type
- `with-directory` adds the Dagger open source repository to the container image filesystem
- `with-workdir` sets the working directory to the Dagger repository
- `with-exec` compiles the Dagger CLI
- `file` returns the built binary as a `File` type
- `export` exports the binary artifact to the Dagger host as `./dagger.bin`

Functions can be chained with the CLI, or programmatically in a [custom Dagger function](./custom-functions.md) using a Dagger SDK. The following are equivalent:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

func (m *MyModule) Publish(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithEntrypoint([]string{"cat", "/etc/os-release"}).
		Publish(ctx, "ttl.sh/my-alpine")
}

```

### Python

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def publish(self) -> str:
        """Publish the container"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_entrypoint(["cat", "/etc/os-release"])
            .publish("ttl.sh/my-alpine")
        )

```

### TypeScript

```typescript
import { dag, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async publish(): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withEntrypoint(["cat", "/etc/os-release"])
      .publish("ttl.sh/my-alpine")
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

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    #[DaggerFunction]
    public function publish(): string
    {
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withEntrypoint(['cat', '/etc/os-release'])
            ->publish('ttl.sh/my-alpine');
    }
}

```

### System shell
```shell
dagger <<EOF
container |
  from alpine:latest |
  with-entrypoint  cat /etc/os-release |
  publish ttl.sh/my-alpine
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from alpine:latest |
  with-entrypoint  cat /etc/os-release |
  publish ttl.sh/my-alpine
```

### Dagger CLI
```shell
dagger core container from --address="alpine:latest" \
  with-entrypoint --args="cat","/etc/os-release" \
  publish --address="ttl.sh/my-alpine"
```

As these example illustrate, function chaining is extremely powerful. See [more examples](./chaining.md) of it in action.

## Modules

Dagger [modules](../features/modules.md) are collections of Dagger functions, packaged together for easy sharing and consumption. They are portable and reusable across languages.

## Calling the Dagger API

You can call the Dagger API from code using a custom Dagger function created with a type-safe Dagger SDK, or from the command line using the Dagger CLI or the Dagger Shell.

The different ways to call the Dagger API are:

- From the Dagger CLI
  - [`dagger`, `dagger core` and `dagger call`](./clients-cli.md)
- From a custom application created with a Dagger SDK
  - [`dagger run`](./clients-sdk.md#custom-applications)
- From a [language-native GraphQL client](./clients-http.md#language-native-http-clients)
- From a command-line HTTP or GraphQL client
  - [`curl`](./clients-http.md#command-line-http-clients)
  - [`dagger query`](./clients-http.md#dagger-cli)

> **Note:**
> The Dagger CLI can call any Dagger function, either from your local filesystem or from a remote Git repository. They can be used interactively, from a shell script, or from a CI configuration. Apart from the Dagger CLI, no other infrastructure or tooling is required to call Dagger functions.

## Extending the Dagger API

The Dagger API is extensible and shareable by design. You can extend the API with [custom functions](./custom-functions.md), which are loaded via Dagger [modules](../features/modules.md). You are encouraged to write your own Dagger modules and share them with others.

Dagger also lets you import and reuse modules developed by your team, your organization or the broader Dagger community. The [Daggerverse](https://daggerverse.dev) is a free service run by Dagger, which indexes all publicly available Dagger modules and Dagger functions, and lets you easily search and consume them.

When a Dagger module is loaded, the Dagger API is [dynamically extended](./internals.md#api-extension-with-dagger-functions) with new Dagger functions served by that module. So, after loading a Dagger module, an API client can now call all of the original core functions _plus_ the new Dagger functions provided by that module.