---
slug: /api/sdk
---

# Dagger SDKs

Dagger SDKs make it easy to call the Dagger API from your favorite programming language, by developing Dagger Functions or custom applications.

A Dagger SDK provides two components:

- A client library to call the Dagger API from your code
- Tooling to extend the Dagger API with your own Dagger Functions (bundled in a Dagger module)

The Dagger API uses GraphQL as its low-level language-agnostic framework, and can also be accessed using any standard GraphQL client. However, you do not need to know GraphQL to call the Dagger API; the translation to underlying GraphQL API calls is handled internally by the Dagger SDKs.

Official Dagger SDKs are currently available for Go, TypeScript and Python. There are also [experimental and community SDKs contributed by the Dagger community](https://github.com/dagger/dagger/tree/main/sdk).

## Dagger Functions

The recommended, and most common way, to interact with the Dagger API is through Dagger Functions. Dagger Functions are just regular code, written in your usual language using a type-safe Dagger SDK.

Dagger Functions are packaged, shared and reused using Dagger modules. A new Dagger module is initialized by calling `dagger init`. This creates a new `dagger.json` configuration file in the current working directory, together with sample Dagger Function source code. The configuration file will default the name of the module to the current directory name, unless an alternative is specified with the `--name` argument.

Once a module is initialized, `dagger develop --sdk=...` sets up or updates all the resources needed to develop the module locally using a Dagger SDK. By default, the module source code will be stored in the current working directory, unless an alternative is specified with the `--source` argument.

Here is an example of initializing a Dagger module:

### Go
```shell
dagger init --name=my-module
dagger develop --sdk=go
```

### Python
```shell
dagger init --name=my-module
dagger develop --sdk=python
```

### TypeScript
```shell
dagger init --name=my-module
dagger develop --sdk=typescript
```

### PHP
```shell
dagger init --name=my-module
dagger develop --sdk=php
```

### Java
```shell
dagger init --name=my-module
dagger develop --sdk=java
```

> **Warning:**
> Running `dagger develop` regenerates the module's code based on dependencies, the current state of the module, and the current Dagger API version. This can result in unexpected results if there are significant changes between the previous and latest installed Dagger API versions. Always refer to the [changelog](https://github.com/dagger/dagger/blob/main/CHANGELOG.md) for a complete list of changes (including breaking changes) in each Dagger release before running `dagger develop`, or use the `--compat=skip` option to bypass updating the Dagger API version.

The default template from `dagger develop` creates the following structure:

### Go

```
.
├── LICENSE
├── dagger.gen.go
├── dagger.json
├── go.mod
├── go.sum
├── internal
│   ├── dagger
│   ├── querybuilder
│   └── telemetry
└── main.go
```

In this structure:

- `dagger.json` is the [Dagger module configuration file](../configuration/modules.md).
- `go.mod`/`go.sum` manage the Go module and its dependencies.
- `main.go` is where your Dagger module code goes. It contains sample code to help you get started.
- `internal` contains automatically-generated types and helpers needed to configure and run the module:
    - `dagger` contains definitions for the Dagger API that's tied to the currently running Dagger Engine container.
    - `querybuilder` has utilities for building GraphQL queries (used internally by the `dagger` package).
    - `telemetry` has utilities for sending Dagger Engine telemetry.

### Python

```
.
├── LICENSE
├── pyproject.toml
├── uv.lock
├── sdk
├── src
│   └── my_module
│       ├── __init__.py
│       └── main.py
└── dagger.json
```

In this structure:

- `dagger.json` is the [Dagger module configuration file](../configuration/modules.md).
- `pyproject.toml` manages the Python project configuration.
- `uv.lock` manages the module's pinned dependencies.
- `src/my_module/` is where your Dagger module code goes. It contains sample code to help you get started.
- `sdk/` contains the vendored Python SDK [client library](https://pypi.org/project/dagger-io/).

This structure hosts a Python import package, with a name derived from the project name (in `pyproject.toml`), inside a `src` directory. This follows a [Python convention](https://packaging.python.org/en/latest/discussions/src-layout-vs-flat-layout/) that requires a project to be installed in order to run its code. This convention prevents accidental usage of development code since the Python interpreter includes the current working directory as the first item on the import path (more information is available in this [blog post on Python packaging](https://blog.ionelmc.ro/2014/05/25/python-packaging/)).

### TypeScript

```
.
├── LICENSE
├── dagger.json
├── package.json
├── sdk
├── src
│   └── index.ts
├── tsconfig.json
└── yarn.lock
```

In this structure:

- `dagger.json` is the [Dagger module configuration file](../configuration/modules.md).
- `package.json` manages the module dependencies.
- `src/` is where your Dagger module code goes. It contains sample code to help you get started.
- `sdk/` contains the TypeScript SDK.

### PHP

```
.
├── LICENSE
├── README.md
├── composer.json
├── composer.lock
├── entrypoint.php
├── sdk
├── src
│   └── MyModule.php
└── vendor
```

In this structure:

- `dagger.json` is the [Dagger module configuration file](../configuration/modules.md).
- `composer.json` manages the module dependencies.
- `src/` is where your Dagger module code goes. It contains sample code to help you get started.
- `sdk/` contains the PHP SDK.

### Java
```
.
├── LICENSE
├── dagger.json
├── pom.xml
├── src
│   └── main
│       └── java
│           └── io
│               └── dagger
│                   └── modules
│                       └── mymodule
│                           ├── MyModule.java
│                           └── package-info.java
└── target

9 directories, 5 files
```

In this structure:

- `dagger.json` is the [Dagger module configuration file](../configuration/modules.md).
- `pom.xml` manages the module dependencies.
- `src/main/java/io/dagger/modules/mymodule/` is where your Dagger module code goes. It contains sample code to help you get started.
- `target` contains the generated Java source classes.

> **Note:**
> While you can use the utilities defined in the automatically-generated code above, you *cannot* edit these files. Even if you edit them locally, any changes will not be persisted when you run the module.

You can now write Dagger Functions using the selected Dagger SDK. Here is an example, which calls a remote API method and returns the result:

### Go
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

> **Caution:**
> You can try this Dagger Function by copying it into the default template generated by `dagger init`, but remember that you must update the module name in the code samples above to match the name used when your module was first initialized.

In simple terms, this Dagger Function:

- initializes a new container from an `alpine` base image.
- executes the `apk add ...` command in the container to add the `curl` and `jq` utilities.
- uses the `curl` utility to send an HTTP request to the URL `https://randomuser.me/api/` and parses the response using `jq`.
- retrieves and returns the output stream of the last executed command as a string.

> **Important:**
> Every Dagger Function has access to the `dag` client, which is a pre-initialized Dagger API client. This client contains all the core types (like `Container`, `Directory`, etc.), as well as bindings to any dependencies your Dagger module has declared.

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

## Custom applications

An alternative approach is to develop a custom application using a Dagger SDK. This involves:

- installing the SDK for your selected language in your development environment
- initializing a Dagger API client in your application code
- calling and combining Dagger API methods from your application to achieve the required result
- executing your application using `dagger run`

### Go

> **Note:**
> The Dagger Go SDK requires [Go 1.22 or later](https://go.dev/doc/install).

From an existing Go module, install the Dagger Go SDK using the commands below:

```shell
go get dagger.io/dagger@latest
```

After importing `dagger.io/dagger` in your Go module code, run the following command to update `go.sum`:

```shell
go mod tidy
```

This example demonstrates how to build a Go application for multiple architectures and Go versions using the Go SDK.

Clone an example project and create a new Go module in the project directory:

```shell
git clone https://go.googlesource.com/example
cd example/hello
mkdir multibuild && cd multibuild
go mod init multibuild
```

Create a new file in the `multibuild` directory named `main.go` and add the following code to it:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"dagger.io/dagger"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := build(context.Background()); err != nil {
		panic(err)
	}
}

func build(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory("../..")

	// create a matrix of Go versions and architectures
	goVersions := []string{"1.22", "1.23"}
	oses := []string{"linux", "darwin"}
	arches := []string{"amd64", "arm64"}

	eg, gctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())

	for _, version := range goVersions {
		for _, goos := range oses {
			for _, goarch := range arches {
				// create a closure to capture the loop variables
				version, goos, goarch := version, goos, goarch

				// create a new container for each Go version, OS, and architecture
				image := fmt.Sprintf("golang:%s", version)
				builder := client.Container().From(image).
					WithDirectory("/src", src).
					WithWorkdir("/src/hello").
					WithEnvVariable("GOOS", goos).
					WithEnvVariable("GOARCH", goarch)

				// run the build command
				path := fmt.Sprintf("build/%s/%s/%s/", version, goos, goarch)
				eg.Go(func() error {
					outpath := fmt.Sprintf("%s/hello", path)
					build := builder.WithExec([]string{"go", "build", "-o", outpath})

					// get the build output directory
					output := build.Directory(path)

					// write the build output directory to the host
					_, err = output.Export(gctx, path)
					if err != nil {
						return err
					}
					return nil
				})
			}
		}
	}

	// wait for all builds to complete
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

```

This Go program imports the Dagger SDK and defines two functions. The `build()` function represents the pipeline and creates a Dagger client, which provides an interface to the Dagger API. It also defines the build matrix, consisting of two OSs (`darwin` and `linux`) and two architectures (`amd64` and `arm64`), and builds the Go application for each combination. The Go build process is instructed via the `GOOS` and `GOARCH` build variables, which are reset for each case.

Try the Go program by executing the command below from the project directory:

```shell
dagger run go run multibuild/main.go
```

The `dagger run` command executes the specified command in a Dagger session and displays live progress. The Go program builds the application for each OS/architecture combination and writes the build results to the host. You will see the build process run four times, once for each combination. Note that the builds are happening concurrently, because the builds do not depend on eachother.

Use the `tree` command to see the build artifacts on the host, as shown below:

```shell
tree build
build
├── 1.22
│   ├── darwin
│   │   ├── amd64
│   │   │   └── hello
│   │   └── arm64
│   │       └── hello
│   └── linux
│       ├── amd64
│       │   └── hello
│       └── arm64
│           └── hello
└── 1.23
    ├── darwin
    │   ├── amd64
    │   │   └── hello
    │   └── arm64
    │       └── hello
    └── linux
        ├── amd64
        │   └── hello
        └── arm64
            └── hello
```

### Python

> **Note:**
> The Dagger Python SDK requires [Python 3.10 or later](https://docs.python.org/3/using/index.html).

Install the Dagger Python SDK in your project:

```shell
uv add dagger-io
```

If you prefer, you can alternatively add the Dagger Python SDK in your Python program. This is useful in case of dependency conflicts, or to keep your Dagger code self-contained.

```shell
uv add --script myscript.py dagger-io
```

This example demonstrates how to test a Python application against multiple Python versions using the Python SDK.

Clone an example project:

```shell
git clone --branch 0.101.0 https://github.com/tiangolo/fastapi
cd fastapi
```

Create a new file named `test.py` in the project directory and add the following code to it.

```python
import sys
from typing import List

import anyio
import dagger


async def test(versions: List[str]):
    async with dagger.Connection(dagger.Config(log_output=sys.stderr)) as client:
        # get reference to the local project
        src = client.host().directory(".")

        async def test_version(version: str):
            print(f"Starting tests for Python {version}")

            python = (
                client.container().from_(f"python:{version}-slim-buster")
                # mount cloned repository
                .with_directory("/src", src)
                # set current working directory
                .with_workdir("/src")
                # install test dependencies
                .with_exec(["pip", "install", "-r", "requirements.txt"])
                # run tests
                .with_exec(["pytest", "tests"])
            )

            # execute
            await python.sync()

            print(f"Tests for Python {version} succeeded!")

        # create task group to run tests in parallel
        async with anyio.create_task_group() as tg:
            for version in versions:
                tg.start_soon(test_version, version)

    print("All tasks have finished")


async def main():
    versions = ["3.8", "3.9", "3.10", "3.11"]
    await test(versions)


if __name__ == "__main__":
    anyio.run(main)

```

This Python program imports the Dagger SDK and defines an asynchronous function named `test()`. This `test()` function creates a Dagger client, which provides an interface to the Dagger API. It also defines the test matrix, consisting of Python versions `3.8` to `3.11` and iterates over this matrix, downloading a Python container image for each specified version and testing the source application in that version.

Add the dependency:

```shell
uv add --script test.py dagger-io
```

Run the Python program by executing the command below from the project directory:

```shell
dagger run uv run test.py
```

The `dagger run` command executes the specified command in a Dagger session and displays live progress. The tool tests the application against each version concurrently and displays the following final output:

```shell
Starting tests for Python 3.8
Starting tests for Python 3.9
Starting tests for Python 3.10
Starting tests for Python 3.11
Tests for Python 3.8 succeeded!
Tests for Python 3.9 succeeded!
Tests for Python 3.11 succeeded!
Tests for Python 3.10 succeeded!
All tasks have finished
```

### TypeScript
> **Note:**
> The Dagger TypeScript SDK requires [TypeScript 5.0 or later](https://www.typescriptlang.org/download/). This SDK currently only [supports Node.js (stable) and Bun (experimental)](../configuration/modules.md). To execute the TypeScript program, you must also have an TypeScript executor like `ts-node` or `tsx`.

Install the Dagger TypeScript SDK in your project using `npm` or `yarn`:

```shell
// using npm
npm install @dagger.io/dagger@latest --save-dev

// using yarn
yarn add @dagger.io/dagger --dev
```

This example demonstrates how to test a Node.js application against multiple Node.js versions using the TypeScript SDK.

Create an example React project (or use an existing one) in TypeScript:

```shell
npx create-react-app my-app --template typescript
cd my-app
```

In the project directory, create a new file named `build.mts` and add the following code to it:

```typescript
import { connect, Directory, Container } from "@dagger.io/dagger"

connect(
  async (client) => {
    // get reference to the local project
    const src: Directory = client.host().directory(".")

    // define the Node versions to test against
    const nodeVersions = ["16", "18", "20"]

    // initialize a matrix of Node containers
    const nodeMatrix: Container[] = nodeVersions.map((version) => {
      return client
        .container()
        .from(`node:${version}`)
        .withDirectory("/src", src)
        .withWorkdir("/src")
        .withExec(["npm", "install"])
    })

    // run tests and build the application for each Node version
    await Promise.all(
      nodeMatrix.map(async (node) => {
        const version = await node.withExec(["node", "-v"]).stdout()

        // run tests
        await node.withExec(["npm", "test", "--", "--watchAll=false"]).sync()

        // build the application
        const buildDir = await node.withExec(["npm", "run", "build"]).directory("build")

        // write the build output to the host
        await buildDir.export(`./build-node-${version.trim()}`)
      }),
    )
  },
  { LogOutput: process.stderr },
)

```

This TypeScript program imports the Dagger SDK and defines an asynchronous function. This function creates a Dagger client, which provides an interface to the Dagger API. It also defines the test/build matrix, consisting of Node.js versions `16`, `18` and `20`, and iterates over this matrix, downloading a Node.js container image for each specified version and testing and building the source application against that version.

Run the program with a Typescript executor like `ts-node`, as shown below:

```shell
dagger run node --loader ts-node/esm ./build.mts
```

The `dagger run` command executes the specified command in a Dagger session and displays live progress. The program tests and builds the application against each version in sequence. At the end of the process, a built application is available for each Node.js version in a `build-node-XX` folder in the project directory, as shown below:

```shell
tree -L 2 -d build-*
build-node-16
└── static
    ├── css
    ├── js
    └── media
build-node-18
└── static
    ├── css
    ├── js
    └── media
build-node-20
└── static
    ├── css
    ├── js
    └── media
```

### PHP

> **Note:**
> The Dagger PHP SDK requires [PHP 8.2 or later](https://www.php.net/downloads.php).

Install the Dagger PHP SDK in your project using `composer`:

```shell
composer require dagger/dagger
```

This example demonstrates how to test a PHP application against multiple PHP versions using the PHP SDK.

Clone an example project:

```shell
git clone https://github.com/slimphp/Slim-Skeleton.git
cd Slim-Skeleton
```

Create a new file named `test.php` in the project directory and add the following code to it.

```php
<?php

require_once(__DIR__ . '/vendor/autoload.php');

use Dagger\Client;
use Dagger\Connection;
use Dagger\Container;
use Dagger\Directory;
use Dagger\Dagger;

use function Amp\async;
use function Amp\Future\await;

function test(array $versions): void
{
    $config = new \Dagger\Config();
    $client = Dagger::connect($config);

    // get reference to the local project
    $src = $client->host()->directory('.');

    $futures = [];
    foreach ($versions as $version) {
        $futures[] = async(function () use ($client, $src, $version) {
            echo "Starting tests for PHP $version...\n";

            $php = $client->container()->from("php:$version-alpine")
                // mount cloned repository
                ->withDirectory('/src', $src)
                // set current working directory
                ->withWorkdir('/src')
                // install dependencies
                ->withExec(['composer', 'install'])
                // run tests
                ->withExec(['vendor/bin/phpunit', 'tests']);

            // execute
            $php->sync();

            echo "Completed tests for PHP $version\n";
            echo "**********\n";
        });
    }

    await($futures);
}

$versions = ['8.2', '8.3'];
test($versions);

```

This PHP program imports the Dagger SDK and defines a function named `test()`. This `test()` function creates a Dagger client, which provides an interface to the Dagger API. It also defines the test matrix, consisting of PHP versions `8.2` to `8.4` and iterates over this matrix, downloading a PHP container image for each specified version and testing the source application in that version.

Run the PHP program by executing the command below from the project directory:

```shell
dagger run php test.php
```

The `dagger run` command executes the specified command in a Dagger session and displays live progress. The program tests the application against each version concurrently and displays the following final output:

```shell
Starting tests for PHP 8.2...
PHPUnit 9.6.22 by Sebastian Bergmann and contributors.

Warning:       Your XML configuration validates against a deprecated schema.
Suggestion:    Migrate your XML configuration using "--migrate-configuration"!

...................                                               19 / 19 (100%)

Time: 00:00.038, Memory: 12.00 MB

OK (19 tests, 37 assertions)
Completed tests for PHP 8.2
**********
Starting tests for PHP 8.3...
PHPUnit 9.6.22 by Sebastian Bergmann and contributors.

Warning:       Your XML configuration validates against a deprecated schema.
Suggestion:    Migrate your XML configuration using "--migrate-configuration"!

...................                                               19 / 19 (100%)

Time: 00:00.039, Memory: 12.00 MB

OK (19 tests, 37 assertions)
Completed tests for PHP 8.3
**********
```

## Differences

Here is a quick summary of differences between these two approaches.

|  | Dagger Functions | Custom applications |
|:---|:---|:---|
| Pre-initialized Dagger API client | Y | N |
| Direct host access | N | Y |
| Direct third-party module access | Y | N |
| Cross-language interoperability | Y | N |
