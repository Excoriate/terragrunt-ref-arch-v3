---
slug: /api/secrets
---

# Secrets

Dagger has first-class support for "secrets", such as passwords, API keys, SSH keys and so on. These secrets can be securely used in Dagger functions without exposing them in plaintext logs, writing them into the filesystem of containers you're building, or inserting them into the cache.

Here is an example, which uses a secret in a Dagger function chain:

```shell
export API_TOKEN="guessme"
```

### System shell
```shell
dagger <<'EOF'
container |
  from alpine:latest |
  with-secret-variable MY_SECRET env://API_TOKEN |
  with-exec -- sh -c 'echo this is the secret: $MY_SECRET' |
  stdout
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from alpine:latest |
  with-secret-variable MY_SECRET env://API_TOKEN |
  with-exec -- sh -c 'echo this is the secret: $MY_SECRET' |
  stdout
```

### Dagger CLI
```shell
dagger core container \
  from --address=alpine:latest \
  with-secret-variable --name="MY_SECRET" --secret="env://API_TOKEN" \
  with-exec --args="sh","-c",'echo this is the secret: $MY_SECRET' \
  stdout
```

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Use a secret in a container
func (m *MyModule) UseSecret(ctx context.Context, secret *dagger.Secret) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithSecretVariable("MY_SECRET", secret).
		WithExec([]string{"sh", "-c", "echo this is the secret: $MY_SECRET"}).
		Stdout(ctx)
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def use_secret(self, secret: dagger.Secret) -> str:
        """Use a secret in a container"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_secret_variable("MY_SECRET", secret)
            .with_exec(["sh", "-c", "echo this is the secret: $MY_SECRET"])
            .stdout()
        )

```

### TypeScript
```typescript
import { dag, Secret, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Use a secret in a container
   */
  @func()
  async useSecret(secret: Secret): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withSecretVariable("MY_SECRET", secret)
      .withExec(["sh", "-c", "echo this is the secret: $MY_SECRET"])
      .stdout()
  }
}

```

[Secret arguments](./arguments.md#secret-arguments) can be sourced from multiple providers: the host environment, the host filesystem, the result of host command execution, and external secret managers [1Password](https://1password.com/) and [Vault](https://www.hashicorp.com/products/vault).

## Security considerations

- Dagger automatically scrubs secrets from its various logs and output streams. This ensures that sensitive data does not leak - for example, in the event of a crash.
- Secret plaintext should be handled securely within your Dagger pipeline. For example, you should not write secret plaintext to a file, as it could then be stored in the Dagger cache.