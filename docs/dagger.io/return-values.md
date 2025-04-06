---
slug: /api/return-values
---

# Return Values

In addition to returning basic types (string, boolean, ...), Dagger Functions can also return any of Dagger's core types, such as `Directory`, `Container`, `Service`, `Secret`, and many more.

This opens powerful applications to Dagger Functions. For example, a Dagger Function that builds binaries could take a directory with the source code as argument and return another directory (a "just-in-time" directory) containing just binaries or a container image (a "just-in-time" container) with the binaries included.

> **Note:**
> If a function doesn't have a return type annotation, it'll be translated to the [dagger.Void][void-type] type in the API.
> 
> [void-type]: https://docs.dagger.io/api/reference/#definition-Void

## String return values

Here is an example of a Dagger Function that returns operating system information for the container as a string:

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns OS information
func (m *MyModule) OsInfo(ctx context.Context, ctr *dagger.Container) (string, error) {
	return ctr.WithExec([]string{"uname", "-a"}).Stdout(ctx)
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def os_info(self, ctr: dagger.Container) -> str:
        """Returns OS information"""
        return await ctr.with_exec(["uname", "-a"]).stdout()

```

### TypeScript
```typescript
import { dag, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns OS information
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
     * Returns OS information
     */
    #[DaggerFunction]
    public function osInfo(Container $ctr): string
    {
        return $ctr->withExec(['uname', '-a'])->stdout();
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.client.Container;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns OS information
   */
  @Function
  public String osInfo(Container ctr) throws Exception {
    return ctr.withExec(List.of("uname", "-a")).stdout().get();
  }
}

```

Here is an example call for this Dagger Function:

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

## Integer return values

Here is an example of a Dagger Function that returns the sum of two integers:

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

## Floating-point number return values

Here is an example of a Dagger Function that returns the sum of two floating-point numbers:

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
> To declare a `float` return type on the function signature, import `float` from `@dagger.io/dagger` and use it as return type.
> The imported `float` type is a `number` underneath, so you can return it as you would return a regular type `number`.

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

## Directory return values

Directory return values might be produced by a Dagger Function that:

- Builds language-specific binaries
- Downloads source code from git or another remote source
- Processes source code (for example, generating documentation or linting code)
- Downloads or processes datasets
- Downloads machine learning models

Here is an example of a Go builder Dagger Function that accepts a remote Git address as a `Directory` argument, builds a Go binary from the source code in that repository, and returns the build directory containing the compiled binary:

### Go
```go
package main

import (
	"fmt"

	"dagger.io/dagger"
)

type MyModule struct{}

// Build the application and return the build directory
func (m *MyModule) GoBuilder(
	// Source directory
	src *dagger.Directory,
	// Architecture to build for
	arch string,
	// OS to build for
	os string,
) *dagger.Directory {
	return dag.Container().
		From("golang:latest").
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", fmt.Sprintf("build/%s-%s", os, arch)}).
		Directory("build")
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def go_builder(
        self,
        src: dagger.Directory,
        arch: str,
        os: str,
    ) -> dagger.Directory:
        """Build the application and return the build directory"""
        return await (
            dag.container()
            .from_("golang:latest")
            .with_env_variable("GOOS", os)
            .with_env_variable("GOARCH", arch)
            .with_directory("/src", src)
            .with_workdir("/src")
            .with_exec(["go", "build", "-o", f"build/{os}-{arch}"])
            .directory("build")
        )

```

### TypeScript
```typescript
import { dag, Directory, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Build the application and return the build directory
   */
  @func()
  goBuilder(src: Directory, arch: string, os: string): Directory {
    return dag
      .container()
      .from("golang:latest")
      .withEnvVariable("GOOS", os)
      .withEnvVariable("GOARCH", arch)
      .withDirectory("/src", src)
      .withWorkdir("/src")
      .withExec(["go", "build", "-o", `build/${os}-${arch}`])
      .directory("build")
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
     * Build the application and return the build directory
     */
    #[DaggerFunction]
    public function goBuilder(
        Directory $src,
        string $arch,
        string $os
    ): Directory {
        return dag()
            ->container()
            ->from('golang:latest')
            ->withEnvVariable('GOOS', $os)
            ->withEnvVariable('GOARCH', $arch)
            ->withDirectory('/src', $src)
            ->withWorkdir('/src')
            ->withExec(['go', 'build', '-o', "build/{$os}-{$arch}"])
            ->directory('build');
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
   * Build the application and return the build directory
   */
  @Function
  public Directory goBuilder(
      @Description("Source directory") Directory src,
      @Description("Architecture to build for") String arch,
      @Description("OS to build for") String os) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("golang:latest")
          .withEnvVariable("GOOS", os)
          .withEnvVariable("GOARCH", arch)
          .withDirectory("/src", src)
          .withWorkdir("/src")
          .withExec(List.of("go", "build", "-o", String.format("build/%s-%s", os, arch)))
          .directory("build");
    }
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'go-builder https://github.com/golang/example#master:/hello amd64 linux'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
go-builder https://github.com/golang/example#master:/hello amd64 linux
```

### Dagger CLI
```shell
dagger call go-builder --src=https://github.com/golang/example#master:/hello --arch=amd64 --os=linux
```

Once the command completes, you should see this output:

```shell
_type: Directory
entries:
    - hello
```

This means that the build succeeded, and a `Directory` type representing the build directory was returned. This `Directory` is called a "just-in-time" directory: a dynamically-produced artifact of a Dagger pipeline.

## File return values

Similar to just-in-time directories, Dagger Functions can produce just-in-time files by returning the `File` type.

Just-in-time files might be produced by a Dagger Function that:

- Builds language-specific binaries
- Combines multiple input files into a single output file, such as a composite video or a compressed archive

Here is an example of a Dagger Function that accepts a filesystem path or remote Git address as a `Directory` argument and  returns a ZIP archive of that directory:

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a ZIP archive of the directory
func (m *MyModule) Archiver(ctx context.Context, src *dagger.Directory) *dagger.File {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "zip"}).
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"zip", "-r", "out.zip", "."}).
		File("out.zip")
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def archiver(self, src: dagger.Directory) -> dagger.File:
        """Returns a ZIP archive of the directory"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "zip"])
            .with_directory("/src", src)
            .with_workdir("/src")
            .with_exec(["zip", "-r", "out.zip", "."])
            .file("out.zip")
        )

```

### TypeScript
```typescript
import { dag, Directory, File, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a ZIP archive of the directory
   */
  @func()
  archiver(src: Directory): File {
    return dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "zip"])
      .withDirectory("/src", src)
      .withWorkdir("/src")
      .withExec(["zip", "-r", "out.zip", "."])
      .file("out.zip")
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
use Dagger\Client\File;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a ZIP archive of the directory
     */
    #[DaggerFunction]
    public function archiver(Directory $src): File
    {
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withExec(['apk', 'add', 'zip'])
            ->withDirectory('/src', $src)
            ->withWorkdir('/src')
            ->withExec(['zip', '-r', 'out.zip', '.'])
            ->file('out.zip');
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Directory;
import io.dagger.client.File;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a ZIP archive of the directory
   */
  @Function
  public File archiver(@Description("Source directory") Directory src) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("alpine:latest")
          .withExec(List.of("apk", "add", "zip"))
          .withDirectory("/src", src)
          .withWorkdir("/src")
          .withExec(List.of("zip", "-r", "out.zip", "."))
          .file("out.zip");
    }
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'archiver https://github.com/dagger/dagger#main:./docs/current_docs/quickstart'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
archiver https://github.com/dagger/dagger#main:./docs/current_docs/quickstart
```

### Dagger CLI
```shell
dagger call archiver --src=https://github.com/dagger/dagger#main:./docs/current_docs/quickstart
```

Once the command completes, you should see this output:

```shell
_type: File
name: out.zip
size: 13744
```

This means that the build succeeded, and a `File` type representing the ZIP archive was returned.

## Container return values

Similar to directories and files, just-in-time containers are produced by calling a Dagger Function that returns the `Container` type. This type provides a complete API for building, running and distributing containers.

Just-in-time containers might be produced by a Dagger Function that:

- Builds a container
- Minifies a container
- Downloads a container image from a running registry
- Exports a container from Docker or other container runtimes
- Snapshots the state of a running container

You can think of a just-in-time container, and the `Container` type that represents it, as a build stage in Dockerfile. Each operation produces a new immutable state, which can be further processed, or exported as an OCI image. Dagger Functions can accept, return and pass containers between themselves, just like regular variables.

Here's an example of a Dagger Function that returns a base `alpine` container image with a list of additional specified packages:

### Go
```go
package main

import (
	"dagger.io/dagger"
)

type MyModule struct{}

// Returns an Alpine container with the specified packages installed
func (m *MyModule) AlpineBuilder(
	// Packages to install
	packages []string,
) *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithExec(append([]string{"apk", "add"}, packages...))
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
    async def alpine_builder(self, packages: List[str]) -> dagger.Container:
        """Returns an Alpine container with the specified packages installed"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", *packages])
        )

```

### TypeScript
```typescript
import { dag, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns an Alpine container with the specified packages installed
   */
  @func()
  alpineBuilder(packages: string[]): Container {
    return dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", ...packages])
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

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns an Alpine container with the specified packages installed
     *
     * @param string[] $packages Packages to install
     */
    #[DaggerFunction]
    public function alpineBuilder(array $packages): Container
    {
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withExec(['apk', 'add', ...$packages]);
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

import java.util.ArrayList;
import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns an Alpine container with the specified packages installed
   */
  @Function
  public Container alpineBuilder(@Description("Packages to install") List<String> packages) throws Exception {
    try (Client client = Dagger.connect()) {
      List<String> args = new ArrayList<>();
      args.add("apk");
      args.add("add");
      args.addAll(packages);
      return client
          .container()
          .from("alpine:latest")
          .withExec(args);
    }
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'alpine-builder curl,openssh'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
alpine-builder curl,openssh
```

### Dagger CLI
```shell
dagger call alpine-builder --packages=curl,openssh
```

Once the command completes, you should see this output:

```shell
_type: Container
defaultArgs:
    - /bin/sh
entrypoint: []
mounts: []
platform: linux/amd64
user: ""
workdir: ""
```

This means that the build succeeded, and a `Container` type representing the built container image was returned.

> **Note:**
> When calling Dagger Functions that produce a just-in-time artifact, you can use the Dagger CLI to add more functions to the pipeline for further processing - for example, inspecting the contents of a directory artifact, exporting a file artifact to the local filesystem, publishing a container artifact to a registry, and so on. This is called ["function chaining"](./index.md#chaining), and it is one of Dagger's most powerful features.

## Chaining

So long as a Dagger Function returns an object that can be JSON-serialized, its state will be preserved and passed to the next function in the chain. This makes it possible to write custom Dagger Functions that support function chaining in the same style as the Dagger API.

Here is an example module with support for function chaining:

### Go
```go
package main

import (
	"fmt"
)

type MyModule struct {
	Greeting string
	Name     string
}

func New(
	// +optional
	// +default="Hello"
	greeting string,
	// +optional
	// +default="World"
	name string,
) *MyModule {
	return &MyModule{
		Greeting: greeting,
		Name:     name,
	}
}

// Return the greeting message
func (m *MyModule) Message() string {
	return fmt.Sprintf("%s, %s!", m.Greeting, m.Name)
}

// Update the greeting message
func (m *MyModule) WithGreeting(greeting string) *MyModule {
	m.Greeting = greeting
	return m
}

// Update the name
func (m *MyModule) WithName(name string) *MyModule {
	m.Name = name
	return m
}

```

### Python
```python
from typing import Annotated

import dagger
from dagger import Doc, field, function, object_type


@object_type
class MyModule:
    greeting: Annotated[str, Doc("The greeting to use")] = field(default="Hello")
    name: Annotated[str, Doc("Who to greet")] = field(default="World")

    @function
    def message(self) -> str:
        """Return the greeting message"""
        return f"{self.greeting}, {self.name}!"

    @function
    def with_greeting(self, greeting: str) -> "MyModule":
        """Update the greeting message"""
        self.greeting = greeting
        return self

    @function
    def with_name(self, name: str) -> "MyModule":
        """Update the name"""
        self.name = name
        return self

```

### TypeScript
```typescript
import { field, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * The greeting to use
   *
   * @default "Hello"
   */
  @field()
  greeting = "Hello"

  /**
   * Who to greet
   *
   * @default "World"
   */
  @field()
  name = "World"

  /**
   * Return the greeting message
   */
  @func()
  message(): string {
    return `${this.greeting}, ${this.name}!`
  }

  /**
   * Update the greeting message
   */
  @func()
  withGreeting(greeting: string): this {
    this.greeting = greeting
    return this
  }

  /**
   * Update the name
   */
  @func()
  withName(name: string): this {
    this.name = name
    return this
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
use Dagger\Attribute\DaggerField;

#[DaggerObject]
class MyModule
{
    /**
     * The greeting to use
     */
    #[DaggerField(default: 'Hello')]
    public string $greeting;

    /**
     * Who to greet
     */
    #[DaggerField(default: 'World')]
    public string $name;

    /**
     * Return the greeting message
     */
    #[DaggerFunction]
    public function message(): string
    {
        return sprintf('%s, %s!', $this->greeting, $this->name);
    }

    /**
     * Update the greeting message
     */
    #[DaggerFunction]
    public function withGreeting(string $greeting): self
    {
        $this->greeting = $greeting;
        return $this;
    }

    /**
     * Update the name
     */
    #[DaggerFunction]
    public function withName(string $name): self
    {
        $this->name = $name;
        return $this;
    }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;
import io.dagger.module.annotation.Default;

@Module
@Object
public class MyModule {

  @Description("The greeting to use")
  @Default("\"Hello\"")
  public String greeting;

  @Description("Who to greet")
  @Default("\"World\"")
  public String name;

  public MyModule() {}

  public MyModule(String greeting, String name) {
    this.greeting = greeting;
    this.name = name;
  }

  /**
   * Return the greeting message
   */
  @Function
  public String message() {
    return String.format("%s, %s!", this.greeting, this.name);
  }

  /**
   * Update the greeting message
   */
  @Function
  public MyModule withGreeting(String greeting) {
    this.greeting = greeting;
    return this;
  }

  /**
   * Update the name
   */
  @Function
  public MyModule withName(String name) {
    this.name = name;
    return this;
  }
}

```

And here is an example call for this module:

### System shell
```shell
dagger -c 'with-name Monde | with-greeting Bonjour | message'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
with-name Monde | with-greeting Bonjour | message
```

### Dagger CLI
```shell
dagger call with-name --name=Monde with-greeting --greeting=Bonjour message
```

The result will be:

```shell
Bonjour, Monde!