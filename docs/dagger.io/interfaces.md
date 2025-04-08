---
slug: /api/interfaces
---

# Interfaces

> **Important:**
> The information on this page is only applicable to Go and TypeScript SDKs. Interfaces are not currently supported in the Python SDK.

## Declaration

### Go

The Go SDK supports interfaces, which allow you to define Go-style interface
definitions so that your module can accept arbitrary values from other modules
without being tightly coupled to the concrete type underneath.

To use an interface, define a Go interface that embeds `DaggerObject` and use
it in a function argument:

Here is an example of the definition of an interface `Fooer` with a single function `foo`:

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

type Fooer interface {
	dagger.Object
	Foo(ctx context.Context, msg string) (string, error)
}

// Takes a Fooer interface as argument
func (m *MyModule) UseFooer(ctx context.Context, f Fooer) (string, error) {
	return f.Foo(ctx, "hello from module")
}

```

Functions defined in interface definitions must match the client-side API
signature style. If they return a scalar value or an array, they must accept a
`context.Context` argument and return an `error` return value. If they return a
chainable object value, they must not return an `error` value, and they do not need
to include a `context.Context` argument.

Note that you must also provide argument names, since they directly translate
to the GraphQL field argument names.

### TypeScript

The TypeScript SDK supports interfaces, which allow you to define a set of functions
that an object must implement to be considered a valid instance of that interface.
That way, your module can accept arbitrary values from other modules without being
tightly coupled to the concrete type underneath.

To use an interface, use the TypeScript `interface` keyword and use it
as a function argument:

Here is an example of the definition of an interface `Fooer` with a single function `foo`:

```ts
import { func, object } from "@dagger.io/dagger"

export interface Fooer {
  foo(msg: string): Promise<string>
}

@object()
class MyModule {
  /**
   * Takes a Fooer interface as argument
   */
  @func()
  async useFooer(f: Fooer): Promise<string> {
    return await f.foo("hello from module")
  }
}

```

Functions defined in interface definitions must match the client-side API
signature style:
- Always define `async` functions in interfaces (wrap the return type in a `Promise<T>`).
- Declare it as a method signature (e.g., `foo(): Promise<string>`) or a property signature (e.g., `foo: () => Promise<string>`).
- Parameters must be properly named since they directly translate to the GraphQL field argument names.

## Implementation

Here is an example of a module `Example` that implements the `Fooer` interface:

### Go

```go
package main

import (
	"context"
	"fmt"
)

type Example struct{}

func (m *Example) Foo(ctx context.Context, msg string) (string, error) {
	return fmt.Sprintf("foo: %s", msg), nil
}

```

### TypeScript

```ts
import { func, object } from "@dagger.io/dagger"

@object()
class Example {
  @func()
  async foo(msg: string): Promise<string> {
    return `foo: ${msg}`
  }
}

```

## Usage

Any object module that implements the interface method can be passed as an argument to the function
that uses the interface.

Dagger automatically detects if an object coming from the module itself or one of its dependencies implements
an interface defined in the module or its dependencies.
If so, it will add new conversion functions to the object that implement the interface
so it can be passed as argument.

Here is an example of a module that uses the `Example` module defined above and
passes it as argument to the `foo` function of the `MyModule` object:

### Go

```go
package main

import (
	"context"
)

type Usage struct{}

func (m *Usage) CallFooer(ctx context.Context) (string, error) {
	// Call the function that uses the interface, passing the Example object
	// Dagger automatically detects that Example implements Fooer and adds the AsFooer function
	return dag.MyModule().UseFooer(ctx, dag.Example().AsFooer())
}

```

### TypeScript

```ts
import { func, object } from "@dagger.io/dagger"

@object()
class Usage {
  @func()
  async callFooer(): Promise<string> {
    // Call the function that uses the interface, passing the Example object
    // Dagger automatically detects that Example implements Fooer and adds the asFooer function
    return await dag.myModule().useFooer(dag.example().asFooer())
  }
}