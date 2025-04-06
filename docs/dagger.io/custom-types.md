---
slug: /api/custom-types
---

# Custom Types

A Dagger module can have multiple object types defined. It's important to understand that they are only accessible through [chaining](./index.md#chaining), starting from a function in the main object.

### Go

Here is an example of a `github` Dagger module, with a function named `DaggerOrganization`
that returns a custom `Organization` type, itself containing a collection of
`Account` types:

```go
package main

import (
	"fmt"
)

// Github module
type Github struct{}

// Github organization
type Organization struct {
	Name    string
	Members []*Account
}

// Github account
type Account struct {
	Login string
	URL   string
}

// Returns the Dagger organization
func (m *Github) DaggerOrganization() *Organization {
	return &Organization{
		Name: "Dagger",
		Members: []*Account{
			{Login: "jane", URL: "https://github.com/jane"},
			{Login: "john", URL: "https://github.com/john"},
		},
	}
}

// Returns the organization's members
func (o *Organization) Members() []*Account {
	return o.Members
}

// Returns the account's login
func (a *Account) Login() string {
	return a.Login
}

// Returns the account's URL
func (a *Account) URL() string {
	return a.URL
}

```

### Python

Here is an example of a `github` Dagger module, with a function named `dagger_organization`
that returns a custom `Organization` type, itself containing a collection of
`Account` types:

```python
from typing import List

import dagger
from dagger import field, function, object_type


@object_type
class Account:
    login: str = field()
    url: str = field()

    @function
    def login(self) -> str:
        """Returns the account's login"""
        return self.login

    @function
    def url(self) -> str:
        """Returns the account's URL"""
        return self.url


@object_type
class Organization:
    name: str = field()
    members: List[Account] = field()

    @function
    def members(self) -> List[Account]:
        """Returns the organization's members"""
        return self.members


@object_type
class Github:
    @function
    def dagger_organization(self) -> Organization:
        """Returns the Dagger organization"""
        return Organization(
            name="Dagger",
            members=[
                Account(login="jane", url="https://github.com/jane"),
                Account(login="john", url="https://github.com/john"),
            ],
        )

```

The [`dagger.field`](https://dagger-io.readthedocs.io/en/latest/module.html#dagger.field) descriptors expose getter functions without arguments, for their [attributes](./state.md).

### TypeScript

Here is an example of a `github` Dagger module, with a function named `daggerOrganization`
that returns a custom `Organization` type, itself containing a collection of
`Account` types:

```typescript
import { field, func, object } from "@dagger.io/dagger"

@object()
class Account {
  @field()
  login: string

  @field()
  url: string

  constructor(login: string, url: string) {
    this.login = login
    this.url = url
  }

  /**
   * Returns the account's login
   */
  @func()
  getLogin(): string {
    return this.login
  }

  /**
   * Returns the account's URL
   */
  @func()
  getUrl(): string {
    return this.url
  }
}

@object()
class Organization {
  @field()
  name: string

  @field()
  members: Account[]

  constructor(name: string, members: Account[]) {
    this.name = name
    this.members = members
  }

  /**
   * Returns the organization's members
   */
  @func()
  getMembers(): Account[] {
    return this.members
  }
}

@object()
class Github {
  /**
   * Returns the Dagger organization
   */
  @func()
  daggerOrganization(): Organization {
    return new Organization("Dagger", [
      new Account("jane", "https://github.com/jane"),
      new Account("john", "https://github.com/john"),
    ])
  }
}

```

TypeScript has multiple ways to support complex data types. Use a `class` when you need methods and privacy, use `type` for plain data objects with only public fields.

> **Note:**
> When the Dagger Engine extends the Dagger API schema with these types, it prefixes
> their names with the name of the main object:
> - Github
> - GithubAccount
> - GithubOrganization
> 
> This is to prevent possible naming conflicts when loading multiple modules,
> which is reflected in code generation (for example, when using this module in
> another one as a dependency).

Here's an example of calling a Dagger Function from this module to get all member URLs:

### System shell
```shell
dagger -c 'dagger-organization | members | url'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
dagger-organization | members | url
```

### Dagger CLI
```shell
dagger call dagger-organization members url
```


```shell
dagger call dagger-organization members url
```

The result will be:

```
https://github.com/jane
https://github.com/john