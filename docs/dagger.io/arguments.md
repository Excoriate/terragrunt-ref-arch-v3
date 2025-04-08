---
slug: /api/arguments
---

# Arguments

Dagger Functions, just like regular functions, can accept arguments. In addition to basic types (string, boolean, integer, arrays...), Dagger also defines powerful core types which Dagger Functions can use for their arguments, such as `Directory`, `Container`, `Service`, `Secret`, and many more.

When calling a Dagger Function from the CLI, its arguments are exposed as command-line flags. How the flag is interpreted depends on the argument type.

> **Important:**
> Dagger Functions execute in containers and thus do not have default access to your host environment (host files, directories, environment variables, etc.). Access to these host resources can only be granted by explicitly passing them as argument values to the Dagger Function.
> - Files and directories: Dagger Functions can accept arguments of type `File` or `Directory`. Pass files and directories on your host by specifying their path as the value of the argument.
> - Environment variables: Pass environment variable values as argument values when invoking a function by just using the standard shell convention of using `$ENV_VAR_NAME.
> - Local network services: Dagger Functions that accept an argument of type `Service` can be passed local network services in the form `tcp://HOST:PORT`.

> **Note:**
> When passing values to Dagger Functions within Dagger Shell, required arguments are positional, while flags can be placed anywhere.

## String arguments

Here is an example of a Dagger Function that accepts a string argument:

### Go
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type MyModule struct{}

type User struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

// Returns users matching the provided gender
func (m *MyModule) GetUser(ctx context.Context, gender string) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("https://randomuser.me/api/?gender=%s", gender))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Name User `json:"name"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	return &result.Results[0].Name, nil
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def get_user(self, gender: str) -> str:
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
                    f"https://randomuser.me/api/?gender={gender}",
                ]
            )
        )

        # Use jq to extract the first name from the JSON response
        return await (
            ctr.with_exec(["jq", "-r", ".results[0].name | {title, first, last}"])
            .stdout()
        )

```

Even though the Python runtime doesn't enforce [type annotations][typing] at runtime,
it's important to define them with Dagger Functions. The Python SDK needs the
typing information at runtime to correctly report to the API. It can't rely on
[type inference][inference], which is only possible for external [static type
checkers][type-checker].

If a function doesn't have a return type annotation, it'll be declared as `None`,
which translates to the [dagger.Void][void] type in the API:

```python
@function
def hello(self):
    return "Hello world!"

# Error: cannot convert string to Void
```

It's fine however, when no value actually needs to be returned:

```python
@function
def hello(self):
    ...
    # no return
```

[@function]: https://dagger-io.readthedocs.io/en/latest/module.html#dagger.function
[@object_type]: https://dagger-io.readthedocs.io/en/latest/module.html#dagger.object_type
[typing]: https://docs.python.org/3/library/typing.html
[inference]: https://mypy.readthedocs.io/en/stable/type_inference_and_annotations.html
[type-checker]: https://realpython.com/python-type-checking/#static-type-checking
[void]: https://dagger-io.readthedocs.io/en/latest/client.html#dagger.Void


### TypeScript
```typescript
import { dag, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns users matching the provided gender
   */
  @func()
  async getUser(gender: string): Promise<string> {
    // NOTE: this uses the http service included in the container, see
    // https://docs.dagger.io/manuals/developer/services/
    const ctr = dag
      .container()
      .from("alpine")
      .withExec(["apk", "add", "curl", "jq"])
      .withExec([
        "curl",
        "-s",
        `https://randomuser.me/api/?gender=${gender}`,
      ])

    // Use jq to extract the first name from the JSON response
    return await ctr
      .withExec(["jq", "-r", ".results[0].name | {title, first, last}"])
      .stdout()
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
    /**
     * Returns users matching the provided gender
     */
    #[DaggerFunction]
    public function getUser(string $gender): string
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
                sprintf('https://randomuser.me/api/?gender=%s', $gender),
            ]);

        // Use jq to extract the first name from the JSON response
        return $ctr
            ->withExec(['jq', '-r', '.results[0].name | {title, first, last}'])
            ->stdout();
    }
}

```

Even though PHP doesn't enforce [type annotations][typing] at runtime,
it's important to define them with Dagger Functions. The PHP SDK needs the
typing information at runtime to correctly report to the API.

[typing]: https://www.php.net/manual/en/language.types.type-system.php


### Java
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
  public String getUser(String gender) throws Exception {
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
                  String.format("https://randomuser.me/api/?gender=%s", gender)
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

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'get-user male'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
get-user male
```

### Dagger CLI
```shell
dagger call get-user --gender=male
```

The result will look something like this:

```shell
{
  "title": "Mr",
  "first": "Hans-Werner",
  "last": "Thielen"
}
```

To pass host environment variables as arguments when invoking a Dagger Function, use the standard shell convention of `$ENV_VAR_NAME`.

Here is an example of passing a host environment variable containing a string value to a Dagger Function:

```shell
export GREETING=bonjour
```

### System shell
```shell
dagger -c 'github.com/shykes/daggerverse/hello@v0.3.0 | hello --greeting=$GREETING'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/shykes/daggerverse/hello@v0.3.0 | hello --greeting=$GREETING
```

### Dagger CLI
```shell
dagger -m github.com/shykes/daggerverse/hello@v0.3.0 call hello --greeting=$GREETING
```

## Boolean arguments

Here is an example of a Dagger Function that accepts a Boolean argument:

### Go
```go
package main

import (
	"strings"
)

type MyModule struct{}

// Returns a greeting message
func (m *MyModule) Hello(shout bool) string {
	msg := "Hello, world"
	if shout {
		msg = strings.ToUpper(msg)
	}
	return msg
}

```

### Python
```python
import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def hello(self, shout: bool) -> str:
        """Returns a greeting message"""
        msg = "Hello, world"
        if shout:
            msg = msg.upper()
        return msg

```

### TypeScript
```typescript
import { func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a greeting message
   */
  @func()
  hello(shout: boolean): string {
    let msg = "Hello, world"
    if (shout) {
      msg = msg.toUpperCase()
    }
    return msg
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

#[DaggerObject]
class MyModule
{
    /**
     * Returns a greeting message
     */
    #[DaggerFunction]
    public function hello(bool $shout): string
    {
        $msg = 'Hello, world';
        if ($shout) {
            $msg = strtoupper($msg);
        }
        return $msg;
    }
}

```

### Java
> **Note:**
> You can either use the primitive `boolean` type or the boxed `java.lang.Boolean` type.

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

@Module
@Object
public class MyModule {

  /**
   * Returns a greeting message
   */
  @Function
  public String hello(boolean shout) {
    String msg = "Hello, world";
    if (shout) {
      msg = msg.toUpperCase();
    }
    return msg;
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'hello true'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
hello true
```

### Dagger CLI
```shell
dagger call hello --shout=true
```

The result will look like this:

```shell
HELLO, WORLD
```

> **Note:**
> When passing optional boolean flags:
> - To set the argument to true: `--foo=true` or `--foo`
> - To set the argument to false: `--foo=false`

## Integer arguments

Here is an example of a Dagger function that accepts an integer argument:

### Go
```go
package main

type MyModule struct{}

// Returns the sum of two integers
func (m *MyModule) AddInteger(a, b int) int {
	return a + b
}

```

### Python
```python
import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def add_integer(self, a: int, b: int) -> int:
        """Returns the sum of two integers"""
        return a + b

```

### TypeScript
```typescript
import { func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns the sum of two integers
   */
  @func()
  addInteger(a: number, b: number): number {
    return a + b
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

#[DaggerObject]
class MyModule
{
    /**
     * Returns the sum of two integers
     */
    #[DaggerFunction]
    public function addInteger(int $a, int $b): int
    {
        return $a + $b;
    }
}

```

### Java
> **Note:**
> You can either use the primitive `int` type or the boxed `java.lang.Integer` type.

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

@Module
@Object
public class MyModule {

  /**
   * Returns the sum of two integers
   */
  @Function
  public int addInteger(int a, int b) {
    return a + b;
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'add-integer 1 2'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
add-integer 1 2
```

### Dagger CLI
```shell
dagger call add-integer --a=1 --b=2
```

The result will look like this:

```shell
3
```

## Floating-point number arguments

Here is an example of a Dagger function that accepts a floating-point number as argument:

### Go
```go
package main

type MyModule struct{}

// Returns the sum of two floats
func (m *MyModule) AddFloat(a, b float64) float64 {
	return a + b
}

```

### Python
```python
import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def add_float(self, a: float, b: float) -> float:
        """Returns the sum of two floats"""
        return a + b

```

### TypeScript

> **Note:**
> There's no `float` type keyword in TypeScript because the type keyword `number` already supports floating point numbers.
>
> To declare a `float` argument on the function signature, import `float` from `@dagger.io/dagger` and use it as an argument's type.
> The imported `float` type is a `number` underneath, so you can use it as you would use a `number` inside your function.

```typescript
import { func, object, float } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns the sum of two floats
   */
  @func()
  addFloat(a: float, b: float): float {
    return a + b
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

#[DaggerObject]
class MyModule
{
    /**
     * Returns the sum of two floats
     */
    #[DaggerFunction]
    public function addFloat(float $a, float $b): float
    {
        return $a + $b;
    }
}

```

### Java
> **Note:**
> You can either use the primitive `float` type or the boxed `java.lang.Float` type.

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

@Module
@Object
public class MyModule {

  /**
   * Returns the sum of two floats
   */
  @Function
  public float addFloat(float a, float b) {
    return a + b;
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'add-float 1.4 2.7'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
add-float 1.4 2.7
```

### Dagger CLI
```shell
dagger call add-float --a=1.4 --b=2.7
```

The result will look like this:

```shell
4.1
```

## Array arguments

To pass an array argument to a Dagger Function, use a comma-separated list of values.

### Go
```go
package main

import (
	"strings"
)

type MyModule struct{}

// Returns a greeting message for a list of names
func (m *MyModule) Hello(names []string) string {
	return "Hello " + strings.Join(names, ", ")
}

```

### Python
```python
from typing import Annotated, List

import dagger
from dagger import Doc, function, object_type


@object_type
class MyModule:
    @function
    def hello(
        self,
        names: Annotated[
            List[str],
            Doc("List of names to greet"),
        ],
    ) -> str:
        """Returns a greeting message for a list of names"""
        return f"Hello {', '.join(names)}"

```

### TypeScript
```typescript
import { func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a greeting message for a list of names
   */
  @func()
  hello(names: string[]): string {
    return `Hello ${names.join(", ")}`
  }
}

```

### PHP
> **Note:**
> Lists must have their subtype specified by adding the `#[ListOfType]` attribute to the relevant function argument.
>
> The PHP SDK needs the typing information at runtime to correctly report to the API.

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Attribute\ListOfType;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a greeting message for a list of names
     *
     * @param string[] $names List of names to greet
     */
    #[DaggerFunction]
    public function hello(
        #[ListOfType('string')]
        array $names
    ): string {
        return 'Hello ' . implode(', ', $names);
    }
}

```

### Java
> **Note:**
> You can also use the `java.util.List` interface to represent a list of values.
> For instance instead of the `String[] names` argument in the example, you can have `List<String> names`.

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

@Module
@Object
public class MyModule {

  /**
   * Returns a greeting message for a list of names
   */
  @Function
  public String hello(String[] names) {
    return "Hello " + String.join(", ", names);
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'hello John,Jane'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
hello John,Jane
```

### Dagger CLI
```shell
dagger call hello --names=John,Jane
```

The result will look like this:

```shell
Hello John, Jane
```

## Directory arguments

You can also pass a directory argument from the command-line. To do so, add the corresponding flag, followed by a local filesystem path or a remote Git reference. In both cases, the CLI will convert it to an object referencing the contents of that filesystem path or Git repository location, and pass the resulting `Directory` object as argument to the Dagger Function.

Dagger Functions do not have access to the filesystem of the host you invoke the Dagger Function from (i.e. the host you execute a CLI command like `dagger call` from). Instead, host directories need to be explicitly passed as arguments to Dagger Functions.

Here's an example of a Dagger Function that accepts a `Directory` as argument. The Dagger Function returns a tree representation of the files and directories at that path.

### Go
```go
package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a tree representation of the directory's contents
func (m *MyModule) Tree(
	ctx context.Context,
	// Source directory
	src *dagger.Directory,
	// Depth of the tree
	depth int,
) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "tree"}).
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"tree", "-L", fmt.Sprintf("%d", depth)}).
		Stdout(ctx)
}

```

### Python
```python
from typing import Annotated

import dagger
from dagger import Doc, dag, function, object_type


@object_type
class MyModule:
    @function
    async def tree(
        self,
        src: Annotated[
            dagger.Directory,
            Doc("Source directory"),
        ],
        depth: Annotated[
            int,
            Doc("Depth of the tree"),
        ],
    ) -> str:
        """Returns a tree representation of the directory's contents"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "tree"])
            .with_directory("/src", src)
            .with_workdir("/src")
            .with_exec(["tree", "-L", str(depth)])
            .stdout()
        )

```

### TypeScript
```typescript
import { dag, Directory, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a tree representation of the directory's contents
   *
   * @param src Source directory
   * @param depth Depth of the tree
   */
  @func()
  async tree(src: Directory, depth: number): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "tree"])
      .withDirectory("/src", src)
      .withWorkdir("/src")
      .withExec(["tree", "-L", depth.toString()])
      .stdout()
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
use Dagger\Client\Directory;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a tree representation of the directory's contents
     */
    #[DaggerFunction]
    public function tree(
        /**
         * Source directory
         */
        Directory $src,

        /**
         * Depth of the tree
         */
        int $depth
    ): string {
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withExec(['apk', 'add', 'tree'])
            ->withDirectory('/src', $src)
            ->withWorkdir('/src')
            ->withExec(['tree', '-L', (string) $depth])
            ->stdout();
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
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
   * Returns a tree representation of the directory's contents
   */
  @Function
  public String tree(
      @Description("Source directory") Directory src,
      @Description("Depth of the tree") int depth) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("alpine:latest")
          .withExec(List.of("apk", "add", "tree"))
          .withDirectory("/src", src)
          .withWorkdir("/src")
          .withExec(List.of("tree", "-L", String.valueOf(depth)))
          .stdout()
          .get();
    }
  }
}

```

Here is an example of passing a local directory to this Dagger Function as argument:

```shell
mkdir -p mydir/mysubdir
touch mydir/a mydir/b mydir/c mydir/mysubdir/y mydir/mysubdir/z
```

### System shell
```shell
dagger -c 'tree mydir 2'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
tree mydir 2
```

### Dagger CLI
```shell
dagger call tree --src=mydir --depth=2
```

The result will look like this:

```shell
.
├── a
├── b
├── c
└── mysubdir
    ├── y
    └── z

2 directories, 5 files
```

Here is an example of passing a remote repository (Dagger's open-source repository) over HTTPS as a `Directory` argument:

### System shell
```shell
dagger <<EOF
container |
  from alpine:latest |
  with-directory /src https://github.com/dagger/dagger |
  with-exec ls /src |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from alpine:latest |
  with-directory /src https://github.com/dagger/dagger |
  with-exec ls /src |
  stdout
```

### Dagger CLI
```shell
dagger core \
  container \
  from --address=alpine:latest \
  with-directory --path=/src --directory=https://github.com/dagger/dagger \
  with-exec --args="ls","/src" \
  stdout
```

The same repository can also be accessed using SSH. Note that this requires [SSH authentication to be properly configured](./remote-modules.md#ssh-authentication) on your Dagger host. Here is the same example, this time using SSH:

### System shell
```shell
dagger <<EOF
container |
  from alpine:latest |
  with-directory /src ssh://git@github.com/dagger/dagger |
  with-exec ls /src |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from alpine:latest |
  with-directory /src ssh://git@github.com/dagger/dagger |
  with-exec ls /src |
  stdout
```

### Dagger CLI
```shell
dagger core \
  container \
  from --address=alpine:latest \
  with-directory --path=/src --directory=ssh://git@github.com/dagger/dagger \
  with-exec --args="ls","/src" \
  stdout
```

For more information about remote repository access, refer to the documentation on [reference schemes](#reference-schemes-for-remote-repositories) and [authentication methods](./remote-modules.md#authentication-methods).

> **Note:**
> Dagger offers two important features for working with `Directory` arguments:
> - [Default paths](./default-paths.md): Set a default directory path to use no value is specified for the argument.
> - [Filters](./fs-filters.md): Control which files and directories are uploaded to a Dagger Function.

## File arguments

File arguments work in the same way as [directory arguments](#directory-arguments). To pass a file to a Dagger Function as an argument, add the corresponding flag, followed by a local filesystem path or a remote Git reference. In both cases, the CLI will convert it to an object referencing that filesystem path or Git repository location, and pass the resulting `File` object as argument to the Dagger Function.

Here's an example of a Dagger Function that accepts a `File` as argument, reads it, and returns its contents:

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns the contents of a file
func (m *MyModule) ReadFile(ctx context.Context, source *dagger.File) (string, error) {
	return source.Contents(ctx)
}

```

### Python
```python
from typing import Annotated

import dagger
from dagger import Doc, function, object_type


@object_type
class MyModule:
    @function
    async def read_file(
        self,
        source: Annotated[
            dagger.File,
            Doc("Source file"),
        ],
    ) -> str:
        """Returns the contents of a file"""
        return await source.contents()

```

### TypeScript
```typescript
import { File, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns the contents of a file
   *
   * @param source Source file
   */
  @func()
  async readFile(source: File): Promise<string> {
    return await source.contents()
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
use Dagger\Client\File;

#[DaggerObject]
class MyModule
{
    /**
     * Returns the contents of a file
     */
    #[DaggerFunction]
    public function readFile(
        /**
         * Source file
         */
        File $source
    ): string {
        return $source->contents();
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.File;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

@Module
@Object
public class MyModule {

  /**
   * Returns the contents of a file
   */
  @Function
  public String readFile(@Description("Source file") File source) throws Exception {
    return source.contents().get();
  }
}

```

Here is an example of passing a local file to this Dagger Function as argument:

### System shell
```shell
dagger -c 'read-file /my/file/path/README.md'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
read-file /my/file/path/README.md
```

### Dagger CLI
```shell
dagger call read-file --source=/my/file/path/README.md
```

And here is an example of passing a file from a remote Git repository as argument:

### System shell
```shell
dagger -c 'read-file https://github.com/dagger/dagger.git#main:README.md'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
read-file https://github.com/dagger/dagger.git#main:README.md
```

### Dagger CLI
```shell
dagger call read-file --source=https://github.com/dagger/dagger.git#main:README.md
```

For more information about remote repository access, refer to the documentation on [reference schemes](#reference-schemes-for-remote-repositories) and [authentication methods](./remote-modules.md#authentication-methods).

> **Note:**
> Dagger offers two important features for working with `File` arguments:
> - [Default paths](./default-paths.md): Set a default file path to use no value is specified for the argument.
> - [Filters](./fs-filters.md): Control which files are uploaded to a Dagger Function.

## Container arguments

Just like directories, you can pass a container to a Dagger Function from the command-line. To do so, add the corresponding flag, followed by the address of an OCI image. The CLI will dynamically pull the image, and pass the resulting `Container` object as argument to the Dagger Function.

Here is an example of a Dagger Function that accepts a container image reference as an argument. The Dagger Function returns operating system information for the container.

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns OS information for the container
func (m *MyModule) OsInfo(ctx context.Context, ctr *dagger.Container) (string, error) {
	return ctr.WithExec([]string{"uname", "-a"}).Stdout(ctx)
}

```

### Python
```python
from typing import Annotated

import dagger
from dagger import Doc, function, object_type


@object_type
class MyModule:
    @function
    async def os_info(
        self,
        ctr: Annotated[
            dagger.Container,
            Doc("Container to get OS information from"),
        ],
    ) -> str:
        """Returns OS information for the container"""
        return await ctr.with_exec(["uname", "-a"]).stdout()

```

### TypeScript
```typescript
import { Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns OS information for the container
   *
   * @param ctr Container to get OS information from
   */
  @func()
  async osInfo(ctr: Container): Promise<string> {
    return await ctr.withExec(["uname", "-a"]).stdout()
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

#[DaggerObject]
class MyModule
{
    /**
     * Returns OS information for the container
     */
    #[DaggerFunction]
    public function osInfo(
        /**
         * Container to get OS information from
         */
        Container $ctr
    ): string {
        return $ctr->withExec(['uname', '-a'])->stdout();
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns OS information for the container
   */
  @Function
  public String osInfo(@Description("Container to get OS information from") Container ctr) throws Exception {
    return ctr.withExec(List.of("uname", "-a")).stdout().get();
  }
}

```

Here is an example of passing a container image reference to this Dagger Function as an argument.

### System shell
```shell
dagger -c 'os-info ubuntu:latest'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
os-info ubuntu:latest
```

### Dagger CLI
```shell
dagger call os-info --ctr=ubuntu:latest
```

The result will look like this:

```shell
Linux dagger 6.1.0-22-cloud-amd64 #1 SMP PREEMPT_DYNAMIC Debian 6.1.94-1 (2024-06-21) x86_64 x86_64 x86_64 GNU/Linux
```

Here is another example of passing a container image reference to a Dagger Function as an argument. The Dagger Function scans the container using Trivy and reports any vulnerabilities found.

### System shell
```shell
dagger -c 'github.com/jpadams/daggerverse/trivy@v0.3.0 | scan-container index.docker.io/alpine:latest'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/jpadams/daggerverse/trivy@v0.3.0 | scan-container index.docker.io/alpine:latest
```

### Dagger CLI
```shell
dagger -m github.com/jpadams/daggerverse/trivy@v0.3.0 call scan-container --ctr=index.docker.io/alpine:latest
```

## Secret arguments

Dagger allows you to utilize confidential information, such as passwords, API keys, SSH keys and so on, in your Dagger [modules](../features/modules.md) and Dagger Functions, without exposing those secrets in plaintext logs, writing them into the filesystem of containers you're building, or inserting them into the cache.

Secrets can be passed to Dagger Functions as arguments using the `Secret` core type. Here is an example of a Dagger Function which accepts a GitHub personal access token as a secret, and uses the token to authorize a request to the GitHub API:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a list of issues from the Dagger repository
func (m *MyModule) GithubApi(ctx context.Context, token *dagger.Secret) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithSecretVariable("GITHUB_API_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			`curl "https://api.github.com/repos/dagger/dagger/issues" -H "Authorization: Bearer $(cat /run/secrets/GITHUB_API_TOKEN)"`,
		}).
		Stdout(ctx)
}

```

### Python

```python
from typing import Annotated

import dagger
from dagger import Doc, Secret, dag, function, object_type


@object_type
class MyModule:
    @function
    async def github_api(
        self,
        token: Annotated[
            Secret,
            Doc("GitHub API token"),
        ],
    ) -> str:
        """Returns a list of issues from the Dagger repository"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "curl"])
            .with_secret_variable("GITHUB_API_TOKEN", token)
            .with_exec(
                [
                    "sh",
                    "-c",
                    'curl "https://api.github.com/repos/dagger/dagger/issues" -H "Authorization: Bearer $(cat /run/secrets/GITHUB_API_TOKEN)"',
                ]
            )
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, func, object, Secret } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a list of issues from the Dagger repository
   *
   * @param token GitHub API token
   */
  @func()
  async githubApi(token: Secret): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "curl"])
      .withSecretVariable("GITHUB_API_TOKEN", token)
      .withExec([
        "sh",
        "-c",
        'curl "https://api.github.com/repos/dagger/dagger/issues" -H "Authorization: Bearer $(cat /run/secrets/GITHUB_API_TOKEN)"',
      ])
      .stdout()
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
use Dagger\Client\Secret;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a list of issues from the Dagger repository
     */
    #[DaggerFunction]
    public function githubApi(
        /**
         * GitHub API token
         */
        Secret $token
    ): string {
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withExec(['apk', 'add', 'curl'])
            ->withSecretVariable('GITHUB_API_TOKEN', $token)
            ->withExec([
                'sh',
                '-c',
                'curl "https://api.github.com/repos/dagger/dagger/issues" -H "Authorization: Bearer $(cat /run/secrets/GITHUB_API_TOKEN)"',
            ])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Secret;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a list of issues from the Dagger repository
   */
  @Function
  public String githubApi(@Description("GitHub API token") Secret token) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("alpine:latest")
          .withExec(List.of("apk", "add", "curl"))
          .withSecretVariable("GITHUB_API_TOKEN", token)
          .withExec(
              List.of(
                  "sh",
                  "-c",
                  "curl \"https://api.github.com/repos/dagger/dagger/issues\" -H \"Authorization: Bearer $(cat /run/secrets/GITHUB_API_TOKEN)\""))
          .stdout()
          .get();
    }
  }
}

```

The result will be a JSON-formatted list of issues from Dagger's repository.

When invoking the Dagger Function using the Dagger CLI, secrets can be sourced from multiple providers. Dagger can read secrets from the host environment, the host filesystem, and the result of host command execution, as well as from external secret managers [1Password](https://1password.com/) and [Vault](https://www.hashicorp.com/products/vault).

### Host secret providers

Here is an example call for this Dagger Function, with the secret sourced from a host environment variable named `GITHUB_API_TOKEN` via the `env` provider:

### System shell
```shell
dagger -c 'github-api env://GITHUB_API_TOKEN'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github-api env://GITHUB_API_TOKEN
```

### Dagger CLI
```shell
dagger call github-api --token=env://GITHUB_API_TOKEN
```

Secrets can also be passed from a host file using the `file` provider:

### System shell
```shell
dagger -c 'github-api file://./github.txt'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github-api file://./github.txt
```

### Dagger CLI
```shell
dagger call github-api --token=file://./github.txt
```

...or as the result of executing a command on the host using the `cmd` provider:

### System shell
```shell
dagger -c 'github-api cmd://"gh auth token"'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github-api cmd://"gh auth token"
```

### Dagger CLI
```shell
dagger call github-api --token=cmd://"gh auth token"
```

### External secret providers

Secrets can also be sourced from external secret managers. Currently, Dagger supports 1Password and Vault.

1Password requires creating a [service account](https://developer.1password.com/docs/service-accounts/get-started) and then setting the `OP_SERVICE_ACCOUNT_TOKEN` environment variable. Alternatively, if no `OP_SERVICE_ACCOUNT_TOKEN` is provided, the integration will attempt to execute the (official) `op` CLI if installed in the system.

1Password [secret references](https://developer.1password.com/docs/cli/secret-references/), in the format `op://VAULT-NAME/ITEM-NAME/[SECTION-NAME/]FIELD-NAME` are supported. Here is an example:

```shell
export OP_SERVICE_ACCOUNT_TOKEN="mytoken"
```

### System shell
```shell
dagger -c 'github-api op://infra/github/credential'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github-api op://infra/github/credential
```

### Dagger CLI
```shell
dagger call github-api --token=op://infra/github/credential
```

Vault can be authenticated with either token or AppRole methods. The Vault host can be specified by setting the environment variable `VAULT_ADDR`. For token authentication, set the environment variable `VAULT_TOKEN`. For AppRole authentication, set the environment variables `VAULT_APPROLE_ROLE_ID` and `VAULT_APPROLE_SECRET_ID`. Additional client configuration can be specified by the default environment variables accepted by Vault.

Vault KvV2 secrets are accessed with the scheme `vault://PATH/TO/SECRET.ITEM`. If your KvV2 is not mounted at `/secret`, specify the mount location with the environment variable `VAULT_PATH_PREFIX`. Here is an example:

```shell
export VAULT_ADDR='https://example.com:8200'
export VAULT_TOKEN=abcd_1234
```

### System shell
```shell
dagger -c 'github-api vault://infra/github.credential'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github-api vault://infra/github.credential
```

### Dagger CLI
```shell
dagger call github-api --token=vault://infra/github.credential
```

## Service arguments

Host network services or sockets can be passed to Dagger Functions as arguments. To do so, add the corresponding flag, followed by a service or socket reference.

### TCP and UDP services

To pass host TCP or UDP network services as arguments when invoking a Dagger Function, specify them in the form `tcp://HOST:PORT` or `udp://HOST:PORT`.

Assume that you have a PostgresQL database running locally on port 5432, as with:

```shell
docker run --rm -d -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
```

Here is an example of passing this host service as argument to a PostgreSQL client Dagger Function, which drops you to a prompt from where you can execute SQL queries:

### System shell
```shell
dagger <<EOF
github.com/kpenfound/dagger-modules/postgres@v0.1.0 |
  client tcp://localhost:5432 postgres postgres postgres
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/kpenfound/dagger-modules/postgres@v0.1.0 | client tcp://localhost:5432 postgres postgres postgres
```

### Dagger CLI
```shell
dagger -m github.com/kpenfound/dagger-modules/postgres@v0.1.0 call \
  client --db=postgres --user=postgres --password=postgres --server=tcp://localhost:5432
```

### Unix sockets

Similar to host TCP/UDP services, Dagger Functions can also be granted access to host Unix sockets when the client is running on Linux or MacOS.

To pass host Unix sockets as arguments when invoking a Dagger Function, specify them by their path on the host.

For example, assuming you have Docker on your host with the Docker daemon listening on a Unix socket at `/var/run/docker.sock`, you can pass this socket to a Docker client Dagger Function as follows:

### System shell
```shell
dagger -c 'github.com/sipsma/daggerverse/docker-client@v0.0.1 /var/run/docker.sock | version'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
github.com/sipsma/daggerverse/docker-client@v0.0.1 /var/run/docker.sock | version
```

### Dagger CLI
```shell
dagger -m github.com/sipsma/daggerverse/docker-client@v0.0.1 call \
  --sock=/var/run/docker.sock version
```

## Optional arguments

Function arguments can be marked as optional. In this case, the Dagger CLI will not display an error if the argument is omitted in the function call.

Here's an example of a Dagger Function with an optional argument:

### Go
```go
package main

import (
	"fmt"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a greeting message
// +optional name
func (m *MyModule) Hello(name dagger.Optional[string]) string {
	val, ok := name.Get()
	if !ok {
		val = "world"
	}
	return fmt.Sprintf("Hello, %s", val)
}

```

### Python
```python
from typing import Annotated, Optional

import dagger
from dagger import Doc, function, object_type


@object_type
class MyModule:
    @function
    def hello(
        self,
        name: Annotated[
            Optional[str],
            Doc("Who to greet"),
        ] = None,
    ) -> str:
        """Returns a greeting message"""
        if name is None:
            name = "world"
        return f"Hello, {name}"

```

### TypeScript
```typescript
import { func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a greeting message
   *
   * @param name Who to greet
   */
  @func()
  hello(name?: string): string {
    if (!name) {
      name = "world"
    }
    return `Hello, ${name}`
  }
}

```

### PHP
> **Note:**
> The definition of optional varies between Dagger and PHP.
>
> An optional argument to PHP is one that has a default value.
>
> An optional argument to Dagger can be omitted entirely. It is truly optional.
>
> To specify a function argument as optional, simply make it nullable. When using the Dagger CLI, if the argument is omitted; the PHP SDK will treat this as receiving the value `null`.

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a greeting message
     */
    #[DaggerFunction]
    public function hello(
        /**
         * Who to greet
         */
        ?string $name = null
    ): string {
        if ($name === null) {
            $name = 'world';
        }
        return sprintf('Hello, %s', $name);
    }
}

```

### Java
> **Note:**
> Because of the usage of `Optional`, primitive types can not be marked as optional. You have to use the boxed types like `Integer` or `Boolean`.
>
> When an argument is not set as optional, Dagger will ensure the value is not `null` by adding a call to `Objects.requireNonNull` against the argument.

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

import java.util.Optional;

@Module
@Object
public class MyModule {

  /**
   * Returns a greeting message
   */
  @Function
  public String hello(@Description("Who to greet") Optional<String> name) {
    return String.format("Hello, %s", name.orElse("world"));
  }
}

```

Here is an example call for this Dagger Function, with the optional argument:

### System shell
```shell
dagger -c 'hello --name=John'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
hello --name=John
```

### Dagger CLI
```shell
dagger call hello --name=John
```

The result will look like this:

```shell
Hello, John
```

Here is an example call for this Dagger Function, without the optional argument:

### System shell
```shell
dagger -c hello
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
hello
```

### Dagger CLI
```shell
dagger call hello
```

The result will look like this:

```shell
Hello, world
```

## Default values

Function arguments can define a default value if no value is supplied for them.

Here's an example of a Dagger Function with a default value for a string argument:

### Go
```go
package main

import (
	"fmt"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a greeting message
// +default="world"
func (m *MyModule) Hello(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}

```

### Python
```python
from typing import Annotated

import dagger
from dagger import Doc, function, object_type


@object_type
class MyModule:
    @function
    def hello(
        self,
        name: Annotated[
            str,
            Doc("Who to greet"),
        ] = "world",
    ) -> str:
        """Returns a greeting message"""
        return f"Hello, {name}"

```

### TypeScript
```typescript
import { func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a greeting message
   *
   * @param name Who to greet (default: "world")
   */
  @func()
  hello(name = "world"): string {
    return `Hello, ${name}`
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

#[DaggerObject]
class MyModule
{
    /**
     * Returns a greeting message
     */
    #[DaggerFunction]
    public function hello(
        /**
         * Who to greet
         */
        string $name = 'world'
    ): string {
        return sprintf('Hello, %s', $name);
    }
}

```

### Java
> **Note:**
> The default value provided must be a valid JSON string representation of the type.
>
> For instance if the argument is of type `Integer` and the default value is `123`, then the annotation must be `@Default("123")`.
> If the argument is of type `String` and the default value is `world`, then the annotation should be `@Default("\"world\"")`.
> In order to simplify this very specific case, if the argument is of type `String` and the value doesn't start with an escaped quote,
> then the SDK will automatically add the escaped quotes for you. That way you can simply write `@Default("world")`.

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Default;
import io.dagger.module.annotation.Description;

@Module
@Object
public class MyModule {

  /**
   * Returns a greeting message
   */
  @Function
  public String hello(@Description("Who to greet") @Default("world") String name) {
    return String.format("Hello, %s", name);
  }
}

```

Here is an example call for this Dagger Function, without the required argument:

### System shell
```shell
dagger -c hello
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
hello
```

### Dagger CLI
```shell
dagger call hello
```

The result will look like this:

```shell
Hello, world
```

Passing null to an optional argument signals that no default value should be used.

> **Note:**
> Dagger supports [default paths](./default-paths.md) for `Directory` or `File` arguments. Dagger will automatically use this default path when no value is specified for the corresponding argument.

## Reference schemes for remote repositories

Dagger supports the use of HTTP and SSH protocols for accessing files and directories in remote repositories, compatible with all major Git hosting platforms such as GitHub, GitLab, BitBucket, Azure DevOps, Codeberg, and Sourcehut. Dagger supports authentication via both HTTPS (using Git credential managers) and SSH (using a unified authentication approach).

Dagger supports the following reference schemes for file and directory arguments:

| Protocol | Scheme     | Authentication | Example |
|----------|------------|----------------|---------|
| HTTP(S)  | Git HTTP   | Git credential manager | `https://github.com/username/repo.git[#version[:subpath]]` |
| SSH      | Explicit   | SSH keys | `ssh://git@github.com/username/repo.git[#version[:subpath]]` |
| SSH      | SCP-like   | SSH keys | `git@github.com:username/repo.git[#version[:subpath]]`     |

> **Note:**
> The reference scheme for directory arguments is currently under discussion [here](https://github.com/dagger/dagger/issues/6957) and [here](https://github.com/dagger/dagger/issues/6944) and may change in future.

Dagger provides additional flexibility in referencing file and directory arguments through the following options:

- Version specification: Add `#version` to target a particular version of the repository. This can be a tag, branch name, or full commit hash. If omitted, the default branch is used.
- Monorepo support: Append `:subpath` after the version specification to access a specific subdirectory within the repository. Note that specifying a version is mandatory when including a subpath.

> **Important:**
> When referencing a specific subdirectory (subpath) within a repository, you must always include a version specification. The format is always `#version:subpath`.