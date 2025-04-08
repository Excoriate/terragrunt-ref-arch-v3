---
slug: /api/state
---

# State and Getters

Object state can be exposed as a Dagger Function, without having to create a getter function explicitly. Depending on the language used, this state is exposed using struct fields (Go), object attributes (Python) or object properties (TypeScript).

### Go
Dagger only exposes a struct's public fields; private fields will not be exposed.

Here's an example where one struct field is exposed as a Dagger Function, while the other is not:

```go
package main

import (
	"fmt"
)

type MyModule struct {
	// The greeting to use
	Greeting string
	// Who to greet
	name string
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
		name:     name,
	}
}

// Return the greeting message
func (m *MyModule) Message() string {
	return fmt.Sprintf("%s, %s!", m.Greeting, m.name)
}

```

### Python
The [`dagger.field`](https://dagger-io.readthedocs.io/en/latest/module.html#dagger.field) descriptor is a wrapper of
[`dataclasses.field`](https://docs.python.org/3/library/dataclasses.html#mutable-default-values). It creates a getter function for the attribute as well so that it's accessible from the Dagger API.

Here's an example where one attribute is exposed as a Dagger Function, while the other is not:

```python
from typing import Annotated

import dagger
from dagger import Doc, field, function, object_type


@object_type
class MyModule:
    greeting: Annotated[str, Doc("The greeting to use")] = field(default="Hello")
    name: str = field(default="World", init=False)

    def __init__(self, name: str = "World"):
        self.name = name

    @function
    def message(self) -> str:
        """Return the greeting message"""
        return f"{self.greeting}, {self.name}!"

```

Notice that compared to [`dataclasses.field`](https://docs.python.org/3/library/dataclasses.html#mutable-default-values), the [`dagger.field`](https://dagger-io.readthedocs.io/en/latest/module.html#dagger.field) wrapper only supports setting `init: bool`, and both `default` and `default_factory` in the same `default` parameter.

> **Note:**
> In a future version of the Python SDK, the `dagger.function` decorator will be used as a descriptor in place of `dagger.field` to make the distinction clearer.

### TypeScript
TypeScript already offers `private`, `protected` and `public` keywords to handle member visibility in a class. However, Dagger will only expose those members of a Dagger module that are explicitly decorated with the `@func()` decorator. Others will remain private.

Here's an example where one field is exposed as a Dagger Function, while the other is not:

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
  name = "World"

  constructor(name?: string, greeting?: string) {
    if (name) {
      this.name = name
    }
    if (greeting) {
      this.greeting = greeting
    }
  }

  /**
   * Return the greeting message
   */
  @func()
  message(): string {
    return `${this.greeting}, ${this.name}!`
  }
}

```

### Java
Dagger will automatically expose all public fields of a class as Dagger Functions. It's also possible to expose a package, `protected` or `private` field by annotating it with the `@Function` annotation.

In case of a field that shouldn't be serialized at all, this can be achieved by marking it as `transient` in Java.

Here's an example where one field is exposed as a Dagger Function, while the other is not:

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
  private String name;

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
}

```

Confirm with `dagger call --help` or `.help my-module` that only the `greeting` function was created, with `name` remaining only a constructor argument:

```
FUNCTIONS
  greeting      The greeting to use
  message       Return the greeting message

ARGUMENTS
      --greeting string   The greeting to use (default "Hello")
      --name string       Who to greet (default "World")