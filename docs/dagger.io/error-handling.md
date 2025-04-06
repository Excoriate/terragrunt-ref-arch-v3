---
slug: /api/error-handling
---

# Error Handling

Dagger modules handle errors in the same way as the language they are written in. This allows you to support any kind of error handling that your application requires. You can also use error handling to verify user input.

Here is an example Dagger Function that performs division and throws an error if the denominator is zero:

### Go

```go
package main

import (
	"errors"
)

type MyModule struct{}

// Divide two numbers
func (m *MyModule) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

```

Error handling in Go modules follows typical Go error patterns with explicit `error` return values and `if err != nil` checks. You can also use error handling to verify user input.

### Python

```python
import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def divide(self, a: int, b: int) -> int:
        """Divide two numbers"""
        if b == 0:
            raise ValueError("cannot divide by zero")
        return a // b

```

### TypeScript

```typescript
import { func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Divide two numbers
   */
  @func()
  divide(a: number, b: number): number {
    if (b === 0) {
      throw new Error("cannot divide by zero")
    }
    return Math.floor(a / b)
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
     * Divide two numbers
     */
    #[DaggerFunction]
    public function divide(int $a, int $b): int
    {
        if ($b === 0) {
            throw new \InvalidArgumentException('cannot divide by zero');
        }
        return intdiv($a, $b);
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

@Module
@Object
public class MyModule {

  /**
   * Divide two numbers
   */
  @Function
  public int divide(int a, int b) {
    if (b == 0) {
      throw new IllegalArgumentException("cannot divide by zero");
    }
    return a / b;
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'divide 4 2'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
divide 4 2
```

### Dagger CLI
```shell
dagger call divide --a=4 --b=2
```

The result will be:

```shell
2
```

Here is another example call for this Dagger Function, this time dividing by zero:

### System shell
```shell
dagger -c 'divide 4 0'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
divide 4 0
```

### Dagger CLI
```shell
dagger call divide --a=4 --b=0
```

The result will be:

```
cannot divide by zero