---
slug: /api/enumerations
displayed_sidebar: "current"
toc_max_heading_level: 2
title: "Enumerations"
---

# Enumerations

> **Important:**
> The information on this page is only applicable to Go, Python and TypeScript SDKs. Enumerations are not currently supported in the PHP SDK.

Dagger supports custom enumeration (enum) types, which can be used to restrict possible values for a string argument. Enum values are strictly validated, preventing common mistakes like accidentally passing null, true, or false.

> **Note:**
> Following the [GraphQL specification](https://spec.graphql.org/October2021/#Name), enums are represented as strings in the Dagger API GraphQL schema and follow these rules:
> - Enum names cannot start with digits, and can only be composed of alphabets, digits or `_`.
> - Enum values are case-sensitive, and by convention should be upper-cased.

Here is an example of a Dagger Function that takes two arguments: an image reference and a severity filter. The latter is defined as an enum named `Severity`:

### Go
```go
package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type MyModule struct{}

type Severity string

const (
	UnknownSeverity  Severity = "UNKNOWN"
	LowSeverity      Severity = "LOW"
	MediumSeverity   Severity = "MEDIUM"
	HighSeverity     Severity = "HIGH"
	CriticalSeverity Severity = "CRITICAL"
)

// Scan an image and return a report
func (m *MyModule) Scan(
	ctx context.Context,
	// Image reference
	ref string,
	// Severity filter
	severity Severity,
) (string, error) {
	return dag.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			"image",
			"--quiet",
			"--format", "json",
			"--severity", string(severity),
			ref,
		}).
		Stdout(ctx)
}

```

### Python
```python
from enum import Enum
from typing import Annotated

import dagger
from dagger import Doc, dag, function, object_type


@dagger.enum
class Severity(Enum):
    """Severity choices"""

    UNKNOWN = "UNKNOWN"
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"
    CRITICAL = "CRITICAL"


@object_type
class MyModule:
    @function
    async def scan(
        self,
        ref: Annotated[str, Doc("Image reference")],
        severity: Annotated[Severity, Doc("Severity filter")],
    ) -> str:
        """Scan an image and return a report"""
        return await (
            dag.container()
            .from_("aquasec/trivy:latest")
            .with_exec(
                [
                    "image",
                    "--quiet",
                    "--format",
                    "json",
                    "--severity",
                    severity.value,
                    ref,
                ]
            )
            .stdout()
        )

```

> **Note:**
> `dagger.Enum` is a convenience base class for defining documentation, but you can also use `enum.Enum` directly.

### TypeScript
```typescript
import { dag, func, object, enumType } from "@dagger.io/dagger"

@enumType()
export enum Severity {
  UNKNOWN = "UNKNOWN",
  LOW = "LOW",
  MEDIUM = "MEDIUM",
  HIGH = "HIGH",
  CRITICAL = "CRITICAL",
}

@object()
class MyModule {
  /**
   * Scan an image and return a report
   *
   * @param ref Image reference
   * @param severity Severity filter
   */
  @func()
  async scan(ref: string, severity: Severity): Promise<string> {
    return await dag
      .container()
      .from("aquasec/trivy:latest")
      .withExec([
        "image",
        "--quiet",
        "--format",
        "json",
        "--severity",
        severity,
        ref,
      ])
      .stdout()
  }
}

```

### Java
```java
package io.dagger.modules.mymodule;

import io.dagger.module.annotation.Enum;

@Enum
public enum Severity {
  UNKNOWN,
  LOW,
  MEDIUM,
  HIGH,
  CRITICAL
}

```

> **Note:**
> Please note the `@Enum` annotation that is required for Dagger to recognize the enum.

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
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
   * Scan an image and return a report
   */
  @Function
  public String scan(
      @Description("Image reference") String ref,
      @Description("Severity filter") Severity severity) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("aquasec/trivy:latest")
          .withExec(
              List.of(
                  "image",
                  "--quiet",
                  "--format",
                  "json",
                  "--severity",
                  severity.name(),
                  ref))
          .stdout()
          .get();
    }
  }
}

```

Enumeration choices will be displayed when calling `--help` or `.help` on a Dagger Function:

### System shell
```shell
dagger -c '.help scan'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
.help scan
```

### Dagger CLI
```shell
dagger call scan --help
```

The result will be:

### System shell
```shell
USAGE
  scan <ref> <severity>

REQUIRED ARGUMENTS
  ref string
  severity MyModuleSeverity   (possible values: UNKNOWN, LOW, MEDIUM, HIGH, CRITICAL)

RETURNS
  string - Primitive type.
```

### Dagger Shell
```shell
USAGE
  scan <ref> <severity>

REQUIRED ARGUMENTS
  ref string
  severity MyModuleSeverity   (possible values: UNKNOWN, LOW, MEDIUM, HIGH, CRITICAL)

RETURNS
  string - Primitive type.
```

### Dagger CLI
```shell
USAGE
  dagger call scan [arguments]

ARGUMENTS
      --ref string                                  [required]
      --severity UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL   [required]
```

Here's an example of calling the Dagger Function with an invalid enum argument:

### System shell
```shell
dagger -c 'scan hello-world:latest FOO'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
scan hello-world:latest FOO
```

### Dagger CLI
```shell
dagger call scan --ref=hello-world:latest --severity=FOO
```

This will result in an error that displays possible values, as follows:

### System shell
```shell
! function "scan": invalid argument "FOO" for "--severity" flag: value should be one of UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL
! Usage: scan <ref> <severity>
```

### Dagger Shell
```shell
! function "scan": invalid argument "FOO" for "--severity" flag: value should be one of UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL
! Usage: scan <ref> <severity>
```

### Dagger CLI
```shell
Error: invalid argument "FOO" for "--severity" flag: value should be one of UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL
Run 'dagger call scan --help' for usage.