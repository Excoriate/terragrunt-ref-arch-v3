---
slug: /api/cache-volumes
---

# Cache Volumes

Volume caching involves caching specific parts of the filesystem and reusing them on subsequent function calls if they are unchanged. This is especially useful when dealing with package managers such as `npm`, `maven`, `pip` and similar. Since these dependencies are usually locked to specific versions in the application's manifest, re-downloading them on every session is inefficient and time-consuming. By using a cache volume for these dependencies, Dagger can reuse the cached contents across Dagger Function runs and reduce execution time.

Here's an example:

### System shell
```shell
dagger <<EOF
container |
  from node:21 |
  with-directory /src https://github.com/dagger/hello-dagger |
  with-workdir /src |
  with-mounted-cache /root/.npm node-21 |
  with-exec npm install
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
container |
  from node:21 |
  with-directory /src https://github.com/dagger/hello-dagger |
  with-workdir /src |
  with-mounted-cache /root/.npm node-21 |
  with-exec npm install
```

### Dagger CLI
```shell
dagger core container \
  from --address=node:21 \
  with-directory --path=/src --directory=https://github.com/dagger/hello-dagger \
  with-workdir --path=/src \
  with-mounted-cache --path=/root/.npm --cache=node-21 \
  with-exec --args="npm","install"
```

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

func (m *MyModule) Build(ctx context.Context) (string, error) {
	// create cache volume for npm dependencies
	npmCache := dag.CacheVolume("node-21")

	// get reference to source code directory
	src := dag.Git("https://github.com/dagger/hello-dagger").Branch("main").Tree()

	// build application
	return dag.Container().
		From("node:21").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		// mount cache volume to /root/.npm
		WithMountedCache("/root/.npm", npmCache).
		WithExec([]string{"npm", "install"}).
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
    async def build(self) -> str:
        """Build the application"""
        # create cache volume for npm dependencies
        npm_cache = dag.cache_volume("node-21")

        # get reference to source code directory
        src = dag.git("https://github.com/dagger/hello-dagger").branch("main").tree()

        # build application
        return await (
            dag.container()
            .from_("node:21")
            .with_directory("/src", src)
            .with_workdir("/src")
            # mount cache volume to /root/.npm
            .with_mounted_cache("/root/.npm", npm_cache)
            .with_exec(["npm", "install"])
            .stdout()
        )

```

### TypeScript

```typescript
import { CacheVolume, dag, Directory, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async build(): Promise<string> {
    // create cache volume for npm dependencies
    const npmCache: CacheVolume = dag.cacheVolume("node-21")

    // get reference to source code directory
    const src: Directory = dag
      .git("https://github.com/dagger/hello-dagger")
      .branch("main")
      .tree()

    // build application
    return await dag
      .container()
      .from("node:21")
      .withDirectory("/src", src)
      .withWorkdir("/src")
      // mount cache volume to /root/.npm
      .withMountedCache("/root/.npm", npmCache)
      .withExec(["npm", "install"])
      .stdout()
  }
}

```

This example will take some time to complete on the first run, as the cache volumes will not exist at that point. Subsequent runs will be significantly faster (assuming there is no other change), since Dagger will simply use the dependencies from the cache volumes instead of downloading them again.