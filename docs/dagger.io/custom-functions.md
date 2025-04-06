---
slug: /api/custom-functions
---

# Custom Functions

In addition to providing a set of core functions and types, the Dagger API can also be extended with custom Dagger Functions and custom types. These custom Dagger Functions are just regular code, written in your usual language using a type-safe Dagger SDK, and packaged and shared in Dagger [modules](../features/modules.md).

When a Dagger module is loaded into a Dagger session, the Dagger API is [dynamically extended](./internals.md#api-extension-with-dagger-functions) with new functions served by that module. So, after loading a Dagger module, an API client can now call all of the original core functions plus the new functions provided by that module.

## Initialize a Dagger module

<!-- Content from ../partials/_dagger_module_init.mdx would be inserted here if available -->
<!-- Assuming removal as per confirmed rules -->

## Create a Dagger Function

Here's an example of a Dagger Function which calls a remote API method and returns the result:

### Go

Update the `main.go` file with the following code:

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns users matching the provided gender
func (m *MyModule) GetUser(ctx context.Context) (string, error) {
	// NOTE: this uses the http service included in the container, see
	// https://docs.dagger.io/manuals/developer/services/
	ctr := dag.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "curl", "jq"}).
		WithExec([]string{
			"curl",
			"-s",
			"https://randomuser.me/api/",
		})

	// Use jq to extract the first name from the JSON response
	return ctr.
		WithExec([]string{"jq", "-r", ".results[0].name | {title, first, last}"}).
		Stdout(ctx)
}

```

This Dagger Function includes the context as input and error as return in its signature.

### Python

Update the `src/my_module/main.py` file with the following code:

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def get_user(self) -> str:
        """Returns users matching the provided gender"""
        # NOTE: this uses the http service included in the container, see
        # https://docs.dagger.io/manuals/developer/services/
        ctr = (
            dag.container()
            .from_("alpine")
            .with_exec(["apk", "add", "curl", "jq"])
            .with_exec(
                [
                    "curl",
                    "-s",
                    "https://randomuser.me/api/",
                ]
            )
        )

        # Use jq to extract the first name from the JSON response
        return await (
            ctr.with_exec(["jq", "-r", ".results[0].name | {title, first, last}"])
            .stdout()
        )

```

Dagger Functions are implemented as [@dagger.function][@function] decorated
methods, of a [@dagger.object_type][@object_type] decorated class.

It's possible for a module to implement multiple classes (*object types*), but
**the first one needs to have a name that matches the module's name**, in
*PascalCase*. This object is sometimes referred to as the *main object*.

For example, for a module initialized with `dagger init --name=my-module`,
the main object needs to be named `MyModule`.

[@function]: https://dagger-io.readthedocs.io/en/latest/module.html#dagger.function
[@object_type]: https://dagger-io.readthedocs.io/en/latest/module.html#dagger.object_type

### TypeScript

Update the `src/index.ts` file with the following code:

```typescript
import { dag, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns users matching the provided gender
   */
  @func()
  async getUser(): Promise<string> {
    // NOTE: this uses the http service included in the container, see
    // https://docs.dagger.io/manuals/developer/services/
    const ctr = dag
      .container()
      .from("alpine")
      .withExec(["apk", "add", "curl", "jq"])
      .withExec(["curl", "-s", "https://randomuser.me/api/"])

    // Use jq to extract the first name from the JSON response
    return await ctr
      .withExec(["jq", "-r", ".results[0].name | {title, first, last}"])
      .stdout()
  }
}

```

### PHP

Update the `src/MyModule.php` file with the following code:

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
    /**
     * Returns users matching the provided gender
     */
    #[DaggerFunction]
    public function getUser(): string
    {
        // NOTE: this uses the http service included in the container, see
        // https://docs.dagger.io/manuals/developer/services/
        $ctr = dag()
            ->container()
            ->from('alpine')
            ->withExec(['apk', 'add', 'curl', 'jq'])
            ->withExec([
                'curl',
                '-s',
                'https://randomuser.me/api/',
            ]);

        // Use jq to extract the first name from the JSON response
        return $ctr
            ->withExec(['jq', '-r', '.results[0].name | {title, first, last}'])
            ->stdout();
    }
}

```

### Java

Update the `src/main/java/io/dagger/modules/mymodule/MyModule.java` file with the following code:

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns users matching the provided gender
   */
  @Function
  public String getUser() throws Exception {
    // NOTE: this uses the http service included in the container, see
    // https://docs.dagger.io/manuals/developer/services/
    try (Client client = Dagger.connect()) {
      Container ctr = client
          .container()
          .from("alpine")
          .withExec(List.of("apk", "add", "curl", "jq"))
          .withExec(
              List.of(
                  "curl",
                  "-s",
                  "https://randomuser.me/api/"
              )
          );

      // Use jq to extract the first name from the JSON response
      return ctr
          .withExec(List.of("jq", "-r", ".results[0].name | {title, first, last}"))
          .stdout()
          .get();
    }
  }
}

```

Dagger Functions must be public. The function must be decorated with the `@Function` annotation
and the class containing the functions must be decorated with the `@Object` annotation.

> **Caution:**
> You can try this Dagger Function by copying it into the default template generated by `dagger init`, but remember that you must update the module name in the code samples above to match the name used when your module was first initialized.

In simple terms, here is what this Dagger Function does:

- It initializes a new container from an `alpine` base image.
- It executes the `apk add ...`   command in the container to add the `curl` and `jq` utilities.
- It uses the `curl` utility to send an HTTP request to the URL `https://randomuser.me/api/` and parses the response using `jq`.
- It retrieves and returns the output stream of the last executed command as a string.

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'get-user'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
get-user
```

### Dagger CLI
```shell
dagger call get-user
```

Here's what you should see:

```shell
{
  "title": "Mrs",
  "first": "Beatrice",
  "last": "Lavigne"
}
```

> **Important:**
> Dagger Functions execute within containers spawned by the Dagger Engine. This "sandboxing" serves a few important purposes:
> 1. Reproducibility: Executing in a well-defined and well-controlled container ensures that a Dagger Function runs the same way every time it is invoked. It also guards against creating "hidden dependencies" on ambient properties of the execution environment that could change at any moment.
> 1. Caching: A reproducible containerized environment makes it possible to cache the result of Dagger Function execution, which in turn allows Dagger to automatically speed up operations.
> 1. Security: Even when running third-party Dagger Functions sourced from a Git repository, those Dagger Functions will not have default access to your host environment (host files, directories, environment variables, etc.). Access to these host resources can only be granted by explicitly passing them as argument values to the Dagger Function.

When implementing Dagger Functions, you are free to write arbitrary code that will execute inside the Dagger module's container. You have access to the Dagger API to make calls to the core Dagger API or other Dagger modules you depend on, but you are also free to just use the language's standard library and/or imported third-party libraries.

The process your code executes in will currently be with the `root` user, but without a full set of Linux capabilities and other standard container sandboxing provided by `runc`.

The current working directory of your code will be an initially empty directory. You can write and read files and directories in this directory if needed. This includes using the `Container.export()`, `Directory.export()` or `File.export()` APIs to write those artifacts to this local directory if needed.