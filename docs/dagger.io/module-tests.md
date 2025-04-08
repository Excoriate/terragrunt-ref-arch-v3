---
slug: /api/module-tests
---

# Module Tests

Like any other piece of software, Dagger Functions and modules should be thoroughly tested. This section documents proven patterns and best practices for effectively testing reusable modules.

The following examples rely on a module called `greeter` that provides a function to greet a person:


### Go

```go
package main

import (
	"fmt"
)

type Greeter struct{}

// Returns a greeting message
func (m *Greeter) Hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

```

### Python

```python
import dagger
from dagger import function, object_type


@object_type
class Greeter:
    @function
    def hello(self, name: str) -> str:
        """Returns a greeting message"""
        return f"Hello, {name}!"

```

### TypeScript

```typescript
import { func, object } from "@dagger.io/dagger"

@object()
export class Greeter {
  /**
   * Returns a greeting message
   */
  @func()
  hello(name: string): string {
    return `Hello, ${name}!`
  }
}

```

## Test module

Well-written tests often provide the best documentation for your software, and this holds true for Dagger modules as well. It's considered a best practice to keep your tests close to the module's code so they can serve as both verification and reference. Additionally, tests that rely on the module's public API act as functional examples, clearly illustrating how to use the module.

Following these principles leads to writing tests for your Dagger modules using Dagger modules themselves. In practice, this means creating a test module in the same directory as your main module and writing your tests as Dagger Functions, as shown below:

### Go

```bash
mkdir tests
cd tests
dagger init --name tests --sdk go --source .
```

Then add the following to `main.go`:

```go
func (m *Tests) Hello(ctx context.Context) error {
	greeting, err := dag.Greeter().Hello(ctx, "World")
	if err != nil {
		return err
	}

	if greeting != "Hello, World!" {
		return errors.New("unexpected greeting")
	}

	return nil
}
```

### Python

```bash
mkdir tests
cd tests
dagger init --name tests --sdk python --source .
```

Then add the following to `src/tests/main.py`:

```python
@object_type
class Tests:
    @function
    async def hello(self):
        greeting = await dag.greeter().hello("World")

        if greeting != "Hello, World!":
            raise Exception("unexpected greeting")
```

### TypeScript

```bash
mkdir tests
cd tests
dagger init --name tests --sdk typescript --source .
```

Then add the following to `src/index.ts`:

```typescript
@object()
export class Tests {
  @func()
  hello(): Promise<void> {
    return dag
      .greeter()
      .hello("World")
      .then((value: string) => {
        if (value != "Hello, World!") {
          throw new Error("unexpected greeting");
        }

        return;
      });
  }
}
```

> **Tip:**
> `tests` is a logical name to use for the test module, but this is not mandatory. Some people call it `dev` to indicate it contains other, development related functions, not just tests.

## Testable examples

In the Daggerverse, [example modules](https://docs.dagger.io/api/daggerverse#examples) are special modules designed to showcase your own modules, offering better demonstrations than the automatically generated ones.

You can combine example modules with the test module pattern to turn those examples into executable tests. Often, this approach provides enough coverage to eliminate the need for a separate test module.

### Go

```bash
mkdir -p examples/go
cd examples/go
dagger init --name examples/go --sdk go --source .
```

Then add the following to `main.go`:

```go
func (m *Examples) GreeterHello(ctx context.Context) error {
	greeting, err := dag.Greeter().Hello(ctx, "World")
	if err != nil {
		return err
	}

	// Do something with the greeting
	_ = greeting

	return nil
}
```

### Python

```bash
mkdir -p examples/python
cd examples/python
dagger init --name examples/python --sdk python --source .
```

Then add the following to `src/examples/main.py`:

```python
@object_type
class Examples:
    @function
    async def greeter_hello(self):
        greeting = await dag.greeter().hello("World")

       	# Do something with the greeting
```

### TypeScript

```bash
mkdir -p examples/typescript
cd examples/typescript
dagger init --name examples/typescript --sdk typescript --source .
```

Then add the following to `src/index.ts`:

```typescript
@object()
export class Examples {
  @func()
  greeterHello(): Promise<void> {
    return dag
      .greeter()
      .hello("World")
      .then((_: string) => {
        // Do something with the greeting

        return;
      });
  }
}
```

If you require more in-depth testing, you can still create a dedicated test module as demonstrated earlier.

> **Tip:**
> Make sure to check out the documentation o [example function naming](https://docs.dagger.io/api/daggerverse#examples).

## Test function signature

Since test functions are ordinary Dagger functions, you can return any value that's allowed. While this approach works fine when running a test with `dagger call`, there are scenarios where a single return value isn't sufficient. For example, you might need to handle multiple output objects, wait for asynchronous operations, manage errors, or (in some languages) provide additional context.

Another challenge arises when you have multiple test functions with parameters, as you must remember to call each test function with the correct arguments.

In such cases, it can be helpful to standardize your test function signature. Consider generating inputs from within the function, synchronizing any asynchronous tasks there, and returning a single value or an error. This approach keeps your tests consistent and easier to maintain.

### Go

```go
func (m *Tests) YourTest(ctx context.Context) error {
	// Your test here

	if false { // Your error condition here
		return errors.New("test failed")
	}

	return nil
}
```

### Python

```python
@object_type
class Tests:
    @function
    async def your_test(self):
        # Your test here

        if false: # Your error condition here
            raise Exception("test failed")
```

### TypeScript

```typescript
@object()
export class Tests {
  @func()
  hello(): Promise<void> {
    return dag
      .yourModule()
      .yourFunction()
      .then(() => {
        if (false) { // Your error condition here
          throw new Error("test failed");
        }

        return;
      });
  }
}
```

In some situations, you may need to provide specific values to your test modules, such as when authenticating against an external service. In these cases, you can rely on [module constructors](./constructors.md) to inject any required inputs.

## "All" function pattern

Regardless of whether you employ the test or the example module pattern, you probably want the ability to run all tests at once (for example, in CI or just to verify everything works locally), while reserving the capability to run individual tests for debugging purposes.

This is where the `all` function comes into the picture. It's basically a single function that executes all your tests or examples.

Depending on the SDK you use, this may be as simple as calling each test function after the other:

### Go

```go
func (m *Tests) All(ctx context.Context) error {
	var err error

	err = m.FirstTest(ctx)
	if err != nil {
		return err
	}

	err = m.SecondTest(ctx)
	if err != nil {
		return err
	}

	return nil
}
```

### Python

```python
@object_type
class Tests:
    @function
    async def all(self):
        await self.hello()
        await self.custom_greeting()
```

### TypeScript

```typescript
@object()
export class Tests {
  @func()
  async all(): Promise<void> {
    await this.hello();
    await this.customGreeting();
  }
}
```

Alternatively, if the SDK/language you use supports this, you can run tests in parallel:

### Go

```go
import "github.com/sourcegraph/conc/pool"

type Tests struct{}

func (m *Tests) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Hello)
	p.Go(m.CustomGreeting)

	return p.Wait()
}
```

### Python

```python
import anyio

@object_type
class Tests:
    @function
    async def all(self):
        async with anyio.create_task_group() as tg:
            tg.start_soon(self.first_test)
            tg.start_soon(self.second_test)
```

### TypeScript

```typescript
@object()
export class Tests {
  @func()
  async all(): Promise<void> {
    await Promise.all([this.firstTest(), this.secondTest()]);
  }
}
```

You can now run all tests for the module using `dagger call -m tests all`.

> **Tip:**
> Adopting a standard test function signature greatly simplifies both kinds of `all` functions.