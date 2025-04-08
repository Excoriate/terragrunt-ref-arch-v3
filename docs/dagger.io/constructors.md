---
slug: /api/constructor
---

# Constructors

Every Dagger module has a constructor. The default one is generated automatically and has no arguments.

It's possible to write a custom constructor. The mechanism to do this is SDK-specific.

This is a simple way to accept module-wide configuration, or just to set a few attributes without having to create setter functions for them.

## Simple constructor

The default constructor for a module can be overridden by registering a custom constructor. Its parameters are available as flags in the `dagger` command directly.

> **Important:**
> Dagger modules have only one constructors. Constructors of [custom types](./custom-types.md) are not registered; they are constructed by the function that [chains](./index.md#chaining) them.

Here is an example module with a custom constructor:

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

// Returns the greeting message
func (m *MyModule) Message() string {
	return fmt.Sprintf("%s, %s!", m.Greeting, m.Name)
}

```

### Python

```python
from typing import Annotated, Optional

import dagger
from dagger import Doc, field, function, object_type


@object_type
class MyModule:
    greeting: Annotated[
        str,
        Doc("The greeting to use"),
    ] = field(default="Hello")

    name: Annotated[
        str,
        Doc("Who to greet"),
    ] = field(default="World")

    @function
    def message(self) -> str:
        """Returns the greeting message"""
        return f"{self.greeting}, {self.name}!"

```

> **Info:**
> In the Python SDK, the [`@dagger.object_type`](https://dagger-io.readthedocs.io/en/latest/module.html#dagger.object_type) decorator wraps [`@dataclasses.dataclass`](https://docs.python.org/3/library/dataclasses.html), which means that an `__init__()` method is automatically generated, with parameters that match the declared class attributes.

The code listing above is an example of an object that has typed attributes.

If a constructor argument needs an asynchronous call to set the default value, it's
possible to replace the default constructor function from `__init__()` to
a factory class method named `create`, as in the following code listing:

> **Warning:**
> This factory class method must be named `create`.

```python
import dagger
from dagger import dag, field, function, object_type


@object_type
class MyModule:
    version: str = field()

    @classmethod
    async def create(cls) -> "MyModule":
        """Create a new instance of MyModule."""
        version = await dag.container().from_("alpine").with_exec(["cat", "/etc/alpine-release"]).stdout()
        return cls(version=version.strip())

    @function
    def version(self) -> str:
        """Return the version."""
        return self.version

```

### TypeScript

```typescript
import { field, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  @field()
  greeting: string = "Hello"

  @field()
  name: string = "World"

  constructor(greeting?: string, name?: string) {
    if (greeting) {
      this.greeting = greeting
    }
    if (name) {
      this.name = name
    }
  }

  /**
   * Returns the greeting message
   */
  @func()
  message(): string {
    return `${this.greeting}, ${this.name}!`
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
use Dagger\Attribute\DaggerArgument;

#[DaggerObject]
class MyModule
{
    public function __construct(
        #[DaggerArgument("The greeting to use")]
        public string $greeting = 'Hello',

        #[DaggerArgument("Who to greet")]
        public string $name = 'World',
    ) {
    }

    /**
     * Returns the greeting message
     */
    #[DaggerFunction]
    public function message(): string
    {
        return sprintf('%s, %s!', $this->greeting, $this->name);
    }
}

```

> **Info:**
> In the PHP SDK the constructor must be the [magic method `__construct`](https://www.php.net/manual/en/language.oop5.decon.php#object.construct).
> As with any method, only public methods with the `#[DaggerFunction]` attribute will be registered with Dagger.

### Java

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

  private String greeting;
  private String name;

  public MyModule() {
    this("Hello", "World");
  }

  public MyModule(
      @Description("The greeting to use") @Default("Hello") String greeting,
      @Description("Who to greet") @Default("World") String name) {
    this.greeting = greeting;
    this.name = name;
  }

  /**
   * Returns the greeting message
   */
  @Function
  public String message() {
    return String.format("%s, %s!", this.greeting, this.name);
  }
}

```

> **Info:**
> In the Java SDK, the constructor must be public. A public **empty** constructor is also required in order to create the object from the serialized data.

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c '. --name=Foo | message'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
. --name=Foo | message
```

### Dagger CLI
```shell
dagger call --name=Foo message
```

The result will be:

```shell
Hello, Foo!
```

> **Important:**
> If you plan to use constructor fields in other module functions, ensure that they are declared as public (in Go and TypeScript). This is because Dagger stores fields using serialization and private fields are omitted during the serialization process. As a result, if a field is not declared as public, calling methods that use it will produce unexpected results.

## Default values for complex types

Constructors can be passed both simple and complex types (such as `Container`, `Directory`, `Service` etc.) as arguments. Default values can be assigned in both cases.

Here is an example of a Dagger module with a default constructor argument of type `Container`:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct {
	Ctr *dagger.Container
}

func New(
	// +optional
	ctr *dagger.Container,
) *MyModule {
	if ctr == nil {
		ctr = dag.Container().From("alpine:3.14.0")
	}
	return &MyModule{
		Ctr: ctr,
	}
}

// Returns the container's alpine version
func (m *MyModule) Version(ctx context.Context) (string, error) {
	return m.Ctr.
		WithExec([]string{"cat", "/etc/alpine-release"}).
		Stdout(ctx)
}

```

### Python

```python
from typing import Annotated, Optional

import dagger
from dagger import Doc, dag, field, function, object_type


@object_type
class MyModule:
    ctr: dagger.Container = field(default=lambda: dag.container().from_("alpine:3.14.0"))

    @function
    async def version(self) -> str:
        """Returns the container's alpine version"""
        return await self.ctr.with_exec(["cat", "/etc/alpine-release"]).stdout()

```

For default values that are more complex, dynamic or just [mutable](https://docs.python.org/3/library/dataclasses.html#mutable-default-values),
use a [factory function](https://docs.python.org/3/library/dataclasses.html#default-factory-functions) without arguments in
[dataclasses.field(default_factory=...)](https://docs.python.org/3/library/dataclasses.html#dataclasses.field):

```python
import random
from dataclasses import dataclass, field
from typing import List

import dagger
from dagger import function, object_type


def random_words() -> List[str]:
    return random.sample(["foo", "bar", "baz"], 2)


@object_type
class MyModule:
    words: List[str] = field(default_factory=random_words)

    @function
    def words(self) -> List[str]:
        return self.words

```

### TypeScript

```typescript
import { Container, dag, field, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  @field()
  ctr: Container

  constructor(ctr?: Container) {
    this.ctr = ctr ?? dag.container().from("alpine:3.14.0")
  }

  /**
   * Returns the container's alpine version
   */
  @func()
  async version(): Promise<string> {
    return await this.ctr.withExec(["cat", "/etc/alpine-release"]).stdout()
  }
}

```

This default value can also be assigned directly in the field:

```typescript
import { Container, dag, field, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  @field()
  ctr: Container = dag.container().from("alpine:3.14.0")

  constructor(ctr?: Container) {
    if (ctr) {
      this.ctr = ctr
    }
  }

  /**
   * Returns the container's alpine version
   */
  @func()
  async version(): Promise<string> {
    return await this.ctr.withExec(["cat", "/etc/alpine-release"]).stdout()
  }
}

```

> **Important:**
> When assigning default values to complex types in TypeScript, it is necessary to use the `??` notation for this assignment. It is not possible to use the classic TypeScript notation for default arguments because the argument in this case is not a TypeScript primitive.

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
    public Container $ctr;

    public function __construct(?Container $ctr = null)
    {
        $this->ctr = $ctr ?? dag()->container()->from('alpine:3.14.0');
    }

    /**
     * Returns the container's alpine version
     */
    #[DaggerFunction]
    public function version(): string
    {
        return $this->ctr
            ->withExec(['cat', '/etc/alpine-release'])
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
import java.util.Optional;

@Module
@Object
public class MyModule {

  private Container ctr;

  public MyModule() {
    this(Optional.empty());
  }

  public MyModule(Optional<Container> ctr) {
    try (Client client = Dagger.connect()) {
      this.ctr = ctr.orElse(client.container().from("alpine:3.14.0"));
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

  /**
   * Returns the container's alpine version
   */
  @Function
  public String version() throws Exception {
    return this.ctr.withExec(List.of("cat", "/etc/alpine-release")).stdout().get();
  }
}

```

It is necessary to explicitly declare the type even when a default value is assigned, so that the Dagger SDK can extend the GraphQL schema correctly.

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c version
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
version
```

### Dagger CLI
```shell
dagger call version
```

The result will be:

```shell
VERSION_ID=3.14.0
