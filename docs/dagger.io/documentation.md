---
slug: /api/documentation
---

# Inline Documentation

Dagger modules and Dagger Functions should be documented so that descriptions are shown in the API and the CLI - for example, when calling `dagger functions`, `dagger call ... --help`, or `.help`.

### Go

The following code snippet shows how to add documentation for:

- The whole module
- Function methods
- Function arguments

```go
package main

import (
	"context"
	"fmt"
)

// A simple example module to say hello.
//
// Further documentation for the module here.
type MyModule struct{}

// Return a greeting.
func (m *MyModule) Hello(
	ctx context.Context,
	// The greeting to display.
	greeting string,
	// Who to greet.
	name string,
) string {
	return fmt.Sprintf("%s, %s!", greeting, name)
}

// Return a loud greeting.
func (m *MyModule) LoudHello(
	ctx context.Context,
	// The greeting to display.
	greeting string,
	// Who to greet.
	name string,
) string {
	return fmt.Sprintf("%s, %s!", greeting, name)
}

```

### Python

The following code snippet shows how to use Python's [documentation string conventions](https://peps.python.org/pep-0008/#documentation-strings) for adding descriptions to:
- The whole [module](https://docs.python.org/3/tutorial/modules.html) or [package](https://docs.python.org/3/tutorial/modules.html#packages)
- Object type classes (group of functions)
- Function methods
- Function arguments

For function arguments, [annotate](https://peps.python.org/pep-0727/#specification) with the [`dagger.Doc`](https://dagger-io.readthedocs.io/en/latest/module.html#dagger.Doc) metadata.

> **Note:**
> [`dagger.Doc`](https://dagger-io.readthedocs.io/en/latest/module.html#dagger.Doc) is just an alias for [`typing_extensions.Doc`](https://typing-extensions.readthedocs.io/en/latest/#Doc).

<!-- loud_hello has a multi-line docstring on purpose -->
```python
"""A simple example module to say hello.

Further documentation for the module here.
"""

from typing import Annotated

import dagger
from dagger import Doc, function, object_type


@object_type
class MyModule:
    """MyModule functions"""

    @function
    def hello(
        self,
        greeting: Annotated[str, Doc("The greeting to display")],
        name: Annotated[str, Doc("Who to greet")],
    ) -> str:
        """Return a greeting."""
        return f"{greeting}, {name}!"

    @function
    def loud_hello(
        self,
        greeting: Annotated[str, Doc("The greeting to display")],
        name: Annotated[str, Doc("Who to greet")],
    ) -> str:
        """Return a loud greeting.

        This is a multi-line docstring.
        """
        return f"{greeting.upper()}, {name.upper()}!"

```

### TypeScript

The following code snippet shows how to add documentation for:
- The whole module
- Function methods
- Function arguments

```typescript
import { func, object, field } from "@dagger.io/dagger"

/**
 * A simple example module to say hello.
 *
 * Further documentation for the module here.
 */
@object()
class MyModule {
  /**
   * Return a greeting.
   *
   * @param greeting The greeting to display.
   * @param name Who to greet.
   */
  @func()
  hello(greeting: string, name: string): string {
    return `${greeting}, ${name}!`
  }

  /**
   * Return a loud greeting.
   *
   * @param greeting The greeting to display.
   * @param name Who to greet.
   */
  @func()
  loudHello(greeting: string, name: string): string {
    return `${greeting.toUpperCase()}, ${name.toUpperCase()}!`
  }
}

```

### PHP

The following code snippet shows how to add documentation for:

- The whole module
- Function methods
- Function arguments

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Attribute\DaggerArgument;

/**
 * A simple example module to say hello.
 *
 * Further documentation for the module here.
 */
#[DaggerObject]
class MyModule
{
    /**
     * Return a greeting.
     */
    #[DaggerFunction]
    public function hello(
        #[DaggerArgument("The greeting to display.")]
        string $greeting,

        #[DaggerArgument("Who to greet.")]
        string $name
    ): string {
        return sprintf('%s, %s!', $greeting, $name);
    }

    /**
     * Return a loud greeting.
     */
    #[DaggerFunction]
    public function loudHello(
        #[DaggerArgument("The greeting to display.")]
        string $greeting,

        #[DaggerArgument("Who to greet.")]
        string $name
    ): string {
        return sprintf('%s, %s!', strtoupper($greeting), strtoupper($name));
    }
}

```

### Java

The following code snippet shows how to add documentation for:

- Function methods
- Function arguments

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

@Module
@Object
public class MyModule {

  /**
   * Return a greeting.
   */
  @Function
  public String hello(
      @Description("The greeting to display.") String greeting,
      @Description("Who to greet.") String name) {
    return String.format("%s, %s!", greeting, name);
  }

  /**
   * Return a loud greeting.
   */
  @Function
  public String loudHello(
      @Description("The greeting to display.") String greeting,
      @Description("Who to greet.") String name) {
    return String.format("%s, %s!", greeting.toUpperCase(), name.toUpperCase());
  }
}

```

The documentation of the module is added in the `package-info.java` file.

```java
/**
 * A simple example module to say hello.
 *
 * Further documentation for the module here.
 */
@Module
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;

```

Here is an example of the result from `dagger functions`:

```
Name         Description
hello        Return a greeting.
loud-hello   Return a loud greeting.
```

Here is an example of the result from `dagger call hello --help`:

```
Return a greeting.

USAGE
  dagger call hello [arguments]

ARGUMENTS
      --greeting string   The greeting to display [required]
      --name string       Who to greet [required]
```


The following code snippet shows how to add documentation for an object and its fields in your Dagger module:

### Go

```go
package main

// User represents a user with a name and age.
type User struct {
	// The name of the user.
	Name string
	// The age of the user.
	Age int
}

```

### Python

```python
from typing import Annotated

import dagger
from dagger import Doc, field, object_type


@object_type
class User:
    """User represents a user with a name and age."""

    name: Annotated[str, Doc("The name of the user.")] = field()
    age: Annotated[int, Doc("The age of the user.")] = field()

```

### TypeScript

```typescript
import { field, object } from "@dagger.io/dagger"

/**
 * User represents a user with a name and age.
 */
@object()
export class User {
  /**
   * The name of the user.
   */
  @field()
  name: string

  /**
   * The age of the user.
   */
  @field()
  age: number

  constructor(name: string, age: number) {
    this.name = name
    this.age = age
  }
}

```

### PHP

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerObject;
use Dagger\Attribute\DaggerField;
use Dagger\Attribute\DaggerArgument;

/**
 * User represents a user with a name and age.
 */
#[DaggerObject]
class User
{
    public function __construct(
        /**
         * The name of the user.
         */
        #[DaggerField]
        public string $name,

        /**
         * The age of the user.
         */
        #[DaggerField]
        public int $age,
    ) {
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Description;

/**
 * User represents a user with a name and age.
 */
@Object
public class User {

  /**
   * The name of the user.
   */
  @Description("The name of the user.")
  public String name;

  /**
   * The age of the user.
   */
  @Description("The age of the user.")
  public int age;

  public User(String name, int age) {
    this.name = name;
    this.age = age;
  }
}

```

Here is an example of the result from `dagger call --help`:

```
ARGUMENTS
      --age int       The age of the user. [required]
      --name string   The name of the user. [required]
