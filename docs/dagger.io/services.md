---
slug: /api/services
title: "Services"
---

# Services

Dagger Functions support service containers, enabling users to spin up additional long-running services (as containers) and communicate with those services from Dagger Functions.

This makes it possible to:
- Instantiate and return services from a Dagger Function, and then:
  - Use those services in other Dagger Functions (container-to-container networking)
  - Use those services from the calling host (container-to-host networking)
- Expose host services for use in a Dagger Function (host-to-container networking).

Some common scenarios for using services with Dagger Functions are:

- Running a database service for local storage or testing
- Running end-to-end integration tests against a service
- Running sidecar services

## Service containers

Services instantiated by a Dagger Function run in service containers, which have the following characteristics:

- Each service container has a canonical, content-addressed hostname and an optional set of exposed ports.
- Service containers are started just-in-time, de-duplicated, and stopped when no longer needed.
- Service containers are health checked prior to running clients.

## Bind services in functions

A Dagger Function can create and return a service, which can then be used from another Dagger Function or from the calling host. Services in Dagger Functions are returned using the `Service` core type.

Here is an example of a Dagger Function that returns an HTTP service. This service is used by another Dagger Function, which creates a service binding using the alias `www` and then accesses the HTTP service using this alias.

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns an HTTP service
func (m *MyModule) HttpService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("nginx:1.25-alpine").
		WithNewFile("/usr/share/nginx/html/index.html", dagger.ContainerWithNewFileOpts{
			Contents: "Hello, world!",
		}).
		WithExposedPort(80).
		AsService()
}

// Accesses the HTTP service
func (m *MyModule) Get(ctx context.Context) (string, error) {
	// bind HTTP service to container
	// access HTTP service using service binding
	// return response
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithServiceBinding("www", m.HttpService(ctx)).
		WithExec([]string{"curl", "http://www:80"}).
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
    def http_service(self) -> dagger.Service:
        """Returns an HTTP service"""
        return (
            dag.container()
            .from_("nginx:1.25-alpine")
            .with_new_file(
                "/usr/share/nginx/html/index.html", contents="Hello, world!"
            )
            .with_exposed_port(80)
            .as_service()
        )

    @function
    async def get(self) -> str:
        """Accesses the HTTP service"""
        # bind HTTP service to container
        # access HTTP service using service binding
        # return response
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "curl"])
            .with_service_binding("www", self.http_service())
            .with_exec(["curl", "http://www:80"])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Container, Directory, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns an HTTP service
   */
  @func()
  httpService(): Service {
    return dag
      .container()
      .from("nginx:1.25-alpine")
      .withNewFile("/usr/share/nginx/html/index.html", {
        contents: "Hello, world!",
      })
      .withExposedPort(80)
      .asService()
  }

  /**
   * Accesses the HTTP service
   */
  @func()
  async get(): Promise<string> {
    // bind HTTP service to container
    // access HTTP service using service binding
    // return response
    return await dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "curl"])
      .withServiceBinding("www", this.httpService())
      .withExec(["curl", "http://www:80"])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns an HTTP service
     */
    #[DaggerFunction]
    public function httpService(): Service
    {
        return dag()
            ->container()
            ->from('nginx:1.25-alpine')
            ->withNewFile('/usr/share/nginx/html/index.html', contents: 'Hello, world!')
            ->withExposedPort(80)
            ->asService();
    }

    /**
     * Accesses the HTTP service
     */
    #[DaggerFunction]
    public function get(): string
    {
        // bind HTTP service to container
        // access HTTP service using service binding
        // return response
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withExec(['apk', 'add', 'curl'])
            ->withServiceBinding('www', $this->httpService())
            ->withExec(['curl', 'http://www:80'])
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
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns an HTTP service
   */
  @Function
  public Service httpService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("nginx:1.25-alpine")
          .withNewFile(
              "/usr/share/nginx/html/index.html",
              new Container.WithNewFileArguments().withContents("Hello, world!"))
          .withExposedPort(80)
          .asService();
    }
  }

  /**
   * Accesses the HTTP service
   */
  @Function
  public String get() throws Exception {
    try (Client client = Dagger.connect()) {
      // bind HTTP service to container
      // access HTTP service using service binding
      // return response
      return client
          .container()
          .from("alpine:latest")
          .withExec(List.of("apk", "add", "curl"))
          .withServiceBinding("www", this.httpService())
          .withExec(List.of("curl", "http://www:80"))
          .stdout()
          .get();
    }
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c get
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
get
```

### Dagger CLI
```shell
dagger call get
```

The result will be:

```shell
Hello, world!
```

## Expose services returned by functions to the host

Services returned by Dagger Functions can also be exposed directly to the host. This enables clients on the host to communicate with services running in Dagger.

One use case is for testing, where you need to be able to spin up ephemeral databases against which to run tests. You might also use this to access a web UI in a browser on your desktop.

Here is another example call for the Dagger Function shown previously, this time exposing the HTTP service on the host

### System shell
```shell
dagger -c 'http-service | up'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
http-service | up
```

### Dagger CLI
```shell
dagger call http-service up
```

By default, each service port maps to the same port on the host - in this case, port 8080. The service can then be accessed by clients on the host. Here's an example:

```shell
curl localhost:8080
```

The result will be:

```shell
Hello, world!
```

To specify a different mapping, use the additional `--ports` argument with a list of host/service port mappings. Here's an example, which exposes the service on host port 9000:

### System shell
```shell
dagger -c 'http-service | up --ports 9000:8080'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
http-service | up --ports 9000:8080
```

### Dagger CLI
```shell
dagger call http-service up --ports 9000:8080
```

> **Note:**
> To bind ports randomly, use the `--random` argument.

## Expose host services to functions

Dagger Functions can also receive host services as function arguments of type `Service`, in the form `tcp://<host>:<port>`. This enables client containers in Dagger Functions to communicate with services running on the host.

> **Note:**
> This implies that a service is already listening on a port on the host, out-of-band of Dagger.

Here is an example of how a container running in a Dagger Function can access and query a MariaDB database service (bound using the alias `db`) running on the host.

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a list of users from the database
func (m *MyModule) UserList(ctx context.Context, svc *dagger.Service) (string, error) {
	return dag.Container().
		From("mariadb:10.11.2").
		WithServiceBinding("db", svc).
		WithEnvVariable("MARIADB_HOST", "db").
		WithEnvVariable("MARIADB_PASSWORD", "secret").
		WithExec([]string{"mariadb", "-u", "root", "-e", "SELECT Host, User FROM mysql.user;"}).
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
    async def user_list(self, svc: dagger.Service) -> str:
        """Returns a list of users from the database"""
        return await (
            dag.container()
            .from_("mariadb:10.11.2")
            .with_service_binding("db", svc)
            .with_env_variable("MARIADB_HOST", "db")
            .with_env_variable("MARIADB_PASSWORD", "secret")
            .with_exec(
                ["mariadb", "-u", "root", "-e", "SELECT Host, User FROM mysql.user;"]
            )
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a list of users from the database
   */
  @func()
  async userList(svc: Service): Promise<string> {
    return await dag
      .container()
      .from("mariadb:10.11.2")
      .withServiceBinding("db", svc)
      .withEnvVariable("MARIADB_HOST", "db")
      .withEnvVariable("MARIADB_PASSWORD", "secret")
      .withExec([
        "mariadb",
        "-u",
        "root",
        "-e",
        "SELECT Host, User FROM mysql.user;",
      ])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a list of users from the database
     */
    #[DaggerFunction]
    public function userList(Service $svc): string
    {
        return dag()
            ->container()
            ->from('mariadb:10.11.2')
            ->withServiceBinding('db', $svc)
            ->withEnvVariable('MARIADB_HOST', 'db')
            ->withEnvVariable('MARIADB_PASSWORD', 'secret')
            ->withExec(['mariadb', '-u', 'root', '-e', 'SELECT Host, User FROM mysql.user;'])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Description;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a list of users from the database
   */
  @Function
  public String userList(@Description("Database service") Service svc) throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("mariadb:10.11.2")
          .withServiceBinding("db", svc)
          .withEnvVariable("MARIADB_HOST", "db")
          .withEnvVariable("MARIADB_PASSWORD", "secret")
          .withExec(
              List.of("mariadb", "-u", "root", "-e", "SELECT Host, User FROM mysql.user;"))
          .stdout()
          .get();
    }
  }
}

```

Before calling this Dagger Function, use the following command to start a MariaDB database service on the host:

```shell
docker run --rm --detach -p 3306:3306 --name my-mariadb --env MARIADB_ROOT_PASSWORD=secret  mariadb:10.11.2
```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'user-list tcp://localhost:3306'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
user-list tcp://localhost:3306
```

### Dagger CLI
```shell
dagger call user-list --svc=tcp://localhost:3306
```

The result will be:

```shell
Host    User
%       root
localhost       mariadb.sys
localhost       root
```

## Create interdependent services

Global hostnames can be assigned to services. This feature is especially valuable for complex networking configurations, such as circular dependencies between services, by allowing services to reference each other by predefined hostnames, without requiring an explicit service binding.

Custom hostnames follow a structured format (`<host>.<module id>.<session id>.dagger.local`), ensuring unique identifiers across modules and sessions.

For example, you can now run two services that depend on each other, each using a hostname to refer to the other by name:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Create two interdependent services
func (m *MyModule) Services(ctx context.Context) *dagger.Container {
	// create service A
	serviceA := dag.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "flask"}).
		WithNewFile("/srv/app.py", dagger.ContainerWithNewFileOpts{
			Contents: `
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service B
    response = requests.get('http://svcb:8081')
    return f'Service A received response from Service B: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
`,
		}).
		WithExposedPort(8080).
		WithExec([]string{"python", "/srv/app.py"})

	// create service B
	serviceB := dag.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "flask"}).
		WithNewFile("/srv/app.py", dagger.ContainerWithNewFileOpts{
			Contents: `
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service A
    response = requests.get('http://svca:8080')
    return f'Service B received response from Service A: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)
`,
		}).
		WithExposedPort(8081).
		WithExec([]string{"python", "/srv/app.py"})

	// create client container with service bindings
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithServiceBinding("svca", serviceA.AsService()).
		WithServiceBinding("svcb", serviceB.AsService()).
		WithExec([]string{"curl", "http://svca:8080"})
}

```

### Python

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def services(self) -> dagger.Container:
        """Create two interdependent services"""
        # create service A
        service_a = (
            dag.container()
            .from_("python:3.11-slim")
            .with_exec(["pip", "install", "flask"])
            .with_new_file(
                "/srv/app.py",
                contents="""
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service B
    response = requests.get('http://svcb:8081')
    return f'Service A received response from Service B: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
""",
            )
            .with_exposed_port(8080)
            .with_exec(["python", "/srv/app.py"])
        )

        # create service B
        service_b = (
            dag.container()
            .from_("python:3.11-slim")
            .with_exec(["pip", "install", "flask"])
            .with_new_file(
                "/srv/app.py",
                contents="""
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service A
    response = requests.get('http://svca:8080')
    return f'Service B received response from Service A: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)
""",
            )
            .with_exposed_port(8081)
            .with_exec(["python", "/srv/app.py"])
        )

        # create client container with service bindings
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "curl"])
            .with_service_binding("svca", service_a.as_service())
            .with_service_binding("svcb", service_b.as_service())
            .with_exec(["curl", "http://svca:8080"])
        )

```

### TypeScript

```typescript
import { dag, Container, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Create two interdependent services
   */
  @func()
  async services(): Promise<Container> {
    // create service A
    const serviceA = dag
      .container()
      .from("python:3.11-slim")
      .withExec(["pip", "install", "flask"])
      .withNewFile("/srv/app.py", {
        contents: `
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service B
    response = requests.get('http://svcb:8081')
    return f'Service A received response from Service B: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
`,
      })
      .withExposedPort(8080)
      .withExec(["python", "/srv/app.py"])

    // create service B
    const serviceB = dag
      .container()
      .from("python:3.11-slim")
      .withExec(["pip", "install", "flask"])
      .withNewFile("/srv/app.py", {
        contents: `
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service A
    response = requests.get('http://svca:8080')
    return f'Service B received response from Service A: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)
`,
      })
      .withExposedPort(8081)
      .withExec(["python", "/srv/app.py"])

    // create client container with service bindings
    return await dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "curl"])
      .withServiceBinding("svca", serviceA.asService())
      .withServiceBinding("svcb", serviceB.asService())
      .withExec(["curl", "http://svca:8080"])
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
use Dagger\Client\Container;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Create two interdependent services
     */
    #[DaggerFunction]
    public function services(): Container
    {
        // create service A
        $serviceA = dag()
            ->container()
            ->from('python:3.11-slim')
            ->withExec(['pip', 'install', 'flask'])
            ->withNewFile('/srv/app.py', contents: <<<PYTHON
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service B
    response = requests.get('http://svcb:8081')
    return f'Service A received response from Service B: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
PYTHON)
            ->withExposedPort(8080)
            ->withExec(['python', '/srv/app.py']);

        // create service B
        $serviceB = dag()
            ->container()
            ->from('python:3.11-slim')
            ->withExec(['pip', 'install', 'flask'])
            ->withNewFile('/srv/app.py', contents: <<<PYTHON
from flask import Flask, request
import requests

app = Flask(__name__)

@app.route('/')
def index():
    # Make a request to service A
    response = requests.get('http://svca:8080')
    return f'Service B received response from Service A: {response.text}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)
PYTHON)
            ->withExposedPort(8081)
            ->withExec(['python', '/srv/app.py']);

        // create client container with service bindings
        return dag()
            ->container()
            ->from('alpine:latest')
            ->withExec(['apk', 'add', 'curl'])
            ->withServiceBinding('svca', $serviceA->asService())
            ->withServiceBinding('svcb', $serviceB->asService())
            ->withExec(['curl', 'http://svca:8080']);
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

@Module
@Object
public class MyModule {

  /**
   * Create two interdependent services
   */
  @Function
  public Container services() throws Exception {
    try (Client client = Dagger.connect()) {
      // create service A
      Container serviceA = client
          .container()
          .from("python:3.11-slim")
          .withExec(List.of("pip", "install", "flask"))
          .withNewFile(
              "/srv/app.py",
              new Container.WithNewFileArguments().withContents("""
                  from flask import Flask, request
                  import requests

                  app = Flask(__name__)

                  @app.route('/')
                  def index():
                      # Make a request to service B
                      response = requests.get('http://svcb:8081')
                      return f'Service A received response from Service B: {response.text}'

                  if __name__ == '__main__':
                      app.run(host='0.0.0.0', port=8080)
                  """))
          .withExposedPort(8080)
          .withExec(List.of("python", "/srv/app.py"));

      // create service B
      Container serviceB = client
          .container()
          .from("python:3.11-slim")
          .withExec(List.of("pip", "install", "flask"))
          .withNewFile(
              "/srv/app.py",
              new Container.WithNewFileArguments().withContents("""
                  from flask import Flask, request
                  import requests

                  app = Flask(__name__)

                  @app.route('/')
                  def index():
                      # Make a request to service A
                      response = requests.get('http://svca:8080')
                      return f'Service B received response from Service A: {response.text}'

                  if __name__ == '__main__':
                      app.run(host='0.0.0.0', port=8081)
                  """))
          .withExposedPort(8081)
          .withExec(List.of("python", "/srv/app.py"));

      // create client container with service bindings
      return client
          .container()
          .from("alpine:latest")
          .withExec(List.of("apk", "add", "curl"))
          .withServiceBinding("svca", serviceA.asService())
          .withServiceBinding("svcb", serviceB.asService())
          .withExec(List.of("curl", "http://svca:8080"));
    }
  }
}

```

In this example, service A and service B are set up with custom hostnames `svca` and `svcb`, allowing each service to communicate with the other by hostname. This capability provides enhanced flexibility for managing service dependencies and interconnections within modular workflows, making it easier to handle complex setups in Dagger.

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c 'services | up --ports 8080:80'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
services | up --ports 8080:80
```

### Dagger CLI
```shell
dagger call services up --ports 8080:80
```


## Persist service state

Dagger cancels each service run after a 10 second grace period to avoid frequent restarts. To avoid relying on the grace period, use a cache volume to persist a service's data, as in the following example:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	// create cache volume for redis data
	redisData := dag.CacheVolume("redis-data")

	// create redis container
	// mount cache volume to /data
	// expose redis port 6379
	// start service
	return dag.Container().
		From("redis:7.2-alpine").
		WithMountedCache("/data", redisData).
		WithExposedPort(6379).
		AsService()
}

// Sets a key in the Redis service
func (m *MyModule) Set(ctx context.Context, key string, value string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "save"}).
		Stdout(ctx)
}

// Gets a key from the Redis service
func (m *MyModule) Get(ctx context.Context, key string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key}).
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
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        # create cache volume for redis data
        redis_data = dag.cache_volume("redis-data")

        # create redis container
        # mount cache volume to /data
        # expose redis port 6379
        # start service
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_mounted_cache("/data", redis_data)
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set(self, key: str, value: str) -> str:
        """Sets a key in the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "save"])
            .stdout()
        )

    @function
    async def get(self, key: str) -> str:
        """Gets a key from the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, CacheVolume, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    // create cache volume for redis data
    const redisData: CacheVolume = dag.cacheVolume("redis-data")

    // create redis container
    // mount cache volume to /data
    // expose redis port 6379
    // start service
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withMountedCache("/data", redisData)
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets a key in the Redis service
   */
  @func()
  async set(key: string, value: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "save"])
      .stdout()
  }

  /**
   * Gets a key from the Redis service
   */
  @func()
  async get(key: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        // create cache volume for redis data
        $redisData = dag()->cacheVolume('redis-data');

        // create redis container
        // mount cache volume to /data
        // expose redis port 6379
        // start service
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withMountedCache('/data', $redisData)
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets a key in the Redis service
     */
    #[DaggerFunction]
    public function set(string $key, string $value): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'save'])
            ->stdout();
    }

    /**
     * Gets a key from the Redis service
     */
    #[DaggerFunction]
    public function get(string $key): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.CacheVolume;
import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      // create cache volume for redis data
      CacheVolume redisData = client.cacheVolume("redis-data");

      // create redis container
      // mount cache volume to /data
      // expose redis port 6379
      // start service
      return client
          .container()
          .from("redis:7.2-alpine")
          .withMountedCache("/data", redisData)
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets a key in the Redis service
   */
  @Function
  public String set(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "save"))
          .stdout()
          .get();
    }
  }

  /**
   * Gets a key from the Redis service
   */
  @Function
  public String get(String key) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key))
          .stdout()
          .get();
    }
  }
}

```


This example uses Redis's `SAVE` command to save the service's data to a cache volume. When a new instance of the service is created, it uses the same cache volume to recreate the original state.

Here is an example of using these Dagger Functions:

### System shell
```shell
dagger -c 'set foo 123'
dagger -c 'get foo'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
set foo 123
get foo
```

### Dagger CLI
```shell
dagger call set --key=foo --value=123
dagger call get --key=foo
```

The result will be:

```shell
123
```

## Start and stop services

Services are designed to be expressed as a Directed Acyclic Graph (DAG) with explicit bindings allowing services to be started lazily, just like every other DAG node. But sometimes, you may need to explicitly manage the lifecycle in a Dagger Function.

For example, this may be needed if the application in the service has certain behavior on shutdown (such as flushing data) that needs careful coordination with the rest of your logic.

The following example explicitly starts the Redis service and stops it at the end, ensuring the 10 second grace period doesn't get in the way, without the need for a persistent cache volume:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Sets and gets a key in the Redis service
func (m *MyModule) SetGet(ctx context.Context, key string, value string) (string, error) {
	// start redis service
	redisSrv, err := m.RedisService(ctx).Start(ctx)
	if err != nil {
		return "", err
	}

	// create redis client container
	// bind redis service
	// execute redis-cli command
	redisCLI := dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", redisSrv).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key})

	// get result
	val, err := redisCLI.Stdout(ctx)
	if err != nil {
		return "", err
	}

	// stop redis service
	_, err = redisSrv.Stop(ctx)
	if err != nil {
		return "", err
	}

	// return result
	return val, nil
}

```

### Python

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set_get(self, key: str, value: str) -> str:
        """Sets and gets a key in the Redis service"""
        # start redis service
        redis_srv = await self.redis_service().start()

        # create redis client container
        # bind redis service
        # execute redis-cli command
        redis_cli = (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", redis_srv)
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
        )

        # get result
        val = await redis_cli.stdout()

        # stop redis service
        await redis_srv.stop()

        # return result
        return val

```

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @func()
  async setGet(key: string, value: string): Promise<string> {
    // start redis service
    const redisSrv = await this.redisService().start()

    // create redis client container
    // bind redis service
    // execute redis-cli command
    const redisCLI: Container = dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", redisSrv)
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])

    // get result
    const val = await redisCLI.stdout()

    // stop redis service
    await redisSrv.stop()

    // return result
    return val
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets and gets a key in the Redis service
     */
    #[DaggerFunction]
    public function setGet(string $key, string $value): string
    {
        // start redis service
        $redisSrv = $this->redisService()->start();

        // create redis client container
        // bind redis service
        // execute redis-cli command
        $redisCLI = dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $redisSrv)
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key]);

        // get result
        $val = $redisCLI->stdout();

        // stop redis service
        $redisSrv->stop();

        // return result
        return $val;
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @Function
  public String setGet(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // start redis service
      Service redisSrv = this.redisService().start().get();

      // create redis client container
      // bind redis service
      // execute redis-cli command
      Container redisCLI = client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", redisSrv)
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key));

      // get result
      String val = redisCLI.stdout().get();

      // stop redis service
      redisSrv.stop().get();

      // return result
      return val;
    }
  }
}

```

## Example: MariaDB database service for application tests

The following example demonstrates how services can be used in Dagger Functions, by creating a Dagger Function for application unit/integration testing against a bound MariaDB database service.

The application used in this example is [Drupal](https://www.drupal.org/), a popular open-source PHP CMS. Drupal includes a large number of unit tests, including tests which require an active database connection. All Drupal 10.x tests are written and executed using the [PHPUnit](https://phpunit.de/) testing framework. Read more about [running PHPUnit tests in Drupal](https://phpunit.de/).

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a MariaDB service
func (m *MyModule) MariaDBService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("mariadb:10.11.2").
		WithEnvVariable("MARIADB_ROOT_PASSWORD", "secret").
		WithEnvVariable("MARIADB_DATABASE", "drupal").
		WithExposedPort(3306).
		AsService()
}

// Tests a Drupal application using a MariaDB service
func (m *MyModule) Test(ctx context.Context) (string, error) {
	// get drupal source code
	drupalDir := dag.Git("https://git.drupalcode.org/project/drupal.git").
		Branch("10.1.x").
		Tree()

	// get php container
	// mount drupal source code
	// mount composer cache
	php := dag.Container().
		From("php:8.2-cli").
		WithDirectory("/opt/drupal", drupalDir).
		WithWorkdir("/opt/drupal/web").
		WithMountedCache("/root/.composer/cache", dag.CacheVolume("composer-cache"))

	// install php dependencies
	// install drupal dependencies
	php = php.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "libsqlite3-dev", "libxml2-dev", "zip"}).
		WithExec([]string{"docker-php-ext-install", "gd", "pdo_mysql", "pdo_sqlite", "xml"}).
		WithExec([]string{"pecl", "install", "xdebug"}).
		WithExec([]string{"docker-php-ext-enable", "xdebug"}).
		WithExec([]string{"php", "-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"}).
		WithExec([]string{"php", "composer-setup.php"}).
		WithExec([]string{"php", "-r", "unlink('composer-setup.php');"}).
		WithExec([]string{"mv", "composer.phar", "/usr/local/bin/composer"}).
		WithExec([]string{"composer", "install"})

	// bind mariadb service
	// set database url env var
	// execute tests
	// return test output
	return php.
		WithServiceBinding("db", m.MariaDBService(ctx)).
		WithEnvVariable("SIMPLETEST_DB", "mysql://root:secret@db/drupal").
		WithExec([]string{"../../vendor/bin/phpunit", "-c", "core/phpunit.xml.dist", "core/modules/user/tests/src/Kernel"}).
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
    def mariadb_service(self) -> dagger.Service:
        """Returns a MariaDB service"""
        return (
            dag.container()
            .from_("mariadb:10.11.2")
            .with_env_variable("MARIADB_ROOT_PASSWORD", "secret")
            .with_env_variable("MARIADB_DATABASE", "drupal")
            .with_exposed_port(3306)
            .as_service()
        )

    @function
    async def test(self) -> str:
        """Tests a Drupal application using a MariaDB service"""
        # get drupal source code
        drupal_dir = dag.git("https://git.drupalcode.org/project/drupal.git").branch(
            "10.1.x"
        ).tree()

        # get php container
        # mount drupal source code
        # mount composer cache
        php = (
            dag.container()
            .from_("php:8.2-cli")
            .with_directory("/opt/drupal", drupal_dir)
            .with_workdir("/opt/drupal/web")
            .with_mounted_cache(
                "/root/.composer/cache", dag.cache_volume("composer-cache")
            )
        )

        # install php dependencies
        # install drupal dependencies
        php = (
            php.with_exec(["apt-get", "update"])
            .with_exec(
                ["apt-get", "install", "-y", "git", "libsqlite3-dev", "libxml2-dev", "zip"]
            )
            .with_exec(["docker-php-ext-install", "gd", "pdo_mysql", "pdo_sqlite", "xml"])
            .with_exec(["pecl", "install", "xdebug"])
            .with_exec(["docker-php-ext-enable", "xdebug"])
            .with_exec(["php", "-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"])
            .with_exec(["php", "composer-setup.php"])
            .with_exec(["php", "-r", "unlink('composer-setup.php');"])
            .with_exec(["mv", "composer.phar", "/usr/local/bin/composer"])
            .with_exec(["composer", "install"])
        )

        # bind mariadb service
        # set database url env var
        # execute tests
        # return test output
        return await (
            php.with_service_binding("db", self.mariadb_service())
            .with_env_variable("SIMPLETEST_DB", "mysql://root:secret@db/drupal")
            .with_exec(
                [
                    "../../vendor/bin/phpunit",
                    "-c",
                    "core/phpunit.xml.dist",
                    "core/modules/user/tests/src/Kernel",
                ]
            )
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Directory, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a MariaDB service
   */
  @func()
  mariadbService(): Service {
    return dag
      .container()
      .from("mariadb:10.11.2")
      .withEnvVariable("MARIADB_ROOT_PASSWORD", "secret")
      .withEnvVariable("MARIADB_DATABASE", "drupal")
      .withExposedPort(3306)
      .asService()
  }

  /**
   * Tests a Drupal application using a MariaDB service
   */
  @func()
  async test(): Promise<string> {
    // get drupal source code
    const drupalDir: Directory = dag
      .git("https://git.drupalcode.org/project/drupal.git")
      .branch("10.1.x")
      .tree()

    // get php container
    // mount drupal source code
    // mount composer cache
    let php: Container = dag
      .container()
      .from("php:8.2-cli")
      .withDirectory("/opt/drupal", drupalDir)
      .withWorkdir("/opt/drupal/web")
      .withMountedCache("/root/.composer/cache", dag.cacheVolume("composer-cache"))

    // install php dependencies
    // install drupal dependencies
    php = php
      .withExec(["apt-get", "update"])
      .withExec([
        "apt-get",
        "install",
        "-y",
        "git",
        "libsqlite3-dev",
        "libxml2-dev",
        "zip",
      ])
      .withExec([
        "docker-php-ext-install",
        "gd",
        "pdo_mysql",
        "pdo_sqlite",
        "xml",
      ])
      .withExec(["pecl", "install", "xdebug"])
      .withExec(["docker-php-ext-enable", "xdebug"])
      .withExec([
        "php",
        "-r",
        "copy('https://getcomposer.org/installer', 'composer-setup.php');",
      ])
      .withExec(["php", "composer-setup.php"])
      .withExec(["php", "-r", "unlink('composer-setup.php');"])
      .withExec(["mv", "composer.phar", "/usr/local/bin/composer"])
      .withExec(["composer", "install"])

    // bind mariadb service
    // set database url env var
    // execute tests
    // return test output
    return await php
      .withServiceBinding("db", this.mariadbService())
      .withEnvVariable("SIMPLETEST_DB", "mysql://root:secret@db/drupal")
      .withExec([
        "../../vendor/bin/phpunit",
        "-c",
        "core/phpunit.xml.dist",
        "core/modules/user/tests/src/Kernel",
      ])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a MariaDB service
     */
    #[DaggerFunction]
    public function mariadbService(): Service
    {
        return dag()
            ->container()
            ->from('mariadb:10.11.2')
            ->withEnvVariable('MARIADB_ROOT_PASSWORD', 'secret')
            ->withEnvVariable('MARIADB_DATABASE', 'drupal')
            ->withExposedPort(3306)
            ->asService();
    }

    /**
     * Tests a Drupal application using a MariaDB service
     */
    #[DaggerFunction]
    public function test(): string
    {
        // get drupal source code
        $drupalDir = dag()
            ->git('https://git.drupalcode.org/project/drupal.git')
            ->branch('10.1.x')
            ->tree();

        // get php container
        // mount drupal source code
        // mount composer cache
        $php = dag()
            ->container()
            ->from('php:8.2-cli')
            ->withDirectory('/opt/drupal', $drupalDir)
            ->withWorkdir('/opt/drupal/web')
            ->withMountedCache('/root/.composer/cache', dag()->cacheVolume('composer-cache'));

        // install php dependencies
        // install drupal dependencies
        $php = $php
            ->withExec(['apt-get', 'update'])
            ->withExec(['apt-get', 'install', '-y', 'git', 'libsqlite3-dev', 'libxml2-dev', 'zip'])
            ->withExec(['docker-php-ext-install', 'gd', 'pdo_mysql', 'pdo_sqlite', 'xml'])
            ->withExec(['pecl', 'install', 'xdebug'])
            ->withExec(['docker-php-ext-enable', 'xdebug'])
            ->withExec(['php', '-r', "copy('https://getcomposer.org/installer', 'composer-setup.php');"])
            ->withExec(['php', 'composer-setup.php'])
            ->withExec(['php', '-r', "unlink('composer-setup.php');"])
            ->withExec(['mv', 'composer.phar', '/usr/local/bin/composer'])
            ->withExec(['composer', 'install']);

        // bind mariadb service
        // set database url env var
        // execute tests
        // return test output
        return $php
            ->withServiceBinding('db', $this->mariadbService())
            ->withEnvVariable('SIMPLETEST_DB', 'mysql://root:secret@db/drupal')
            ->withExec(['../../vendor/bin/phpunit', '-c', 'core/phpunit.xml.dist', 'core/modules/user/tests/src/Kernel'])
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
import io.dagger.client.Directory;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a MariaDB service
   */
  @Function
  public Service mariadbService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("mariadb:10.11.2")
          .withEnvVariable("MARIADB_ROOT_PASSWORD", "secret")
          .withEnvVariable("MARIADB_DATABASE", "drupal")
          .withExposedPort(3306)
          .asService();
    }
  }

  /**
   * Tests a Drupal application using a MariaDB service
   */
  @Function
  public String test() throws Exception {
    try (Client client = Dagger.connect()) {
      // get drupal source code
      Directory drupalDir = client
          .git("https://git.drupalcode.org/project/drupal.git")
          .branch("10.1.x")
          .tree();

      // get php container
      // mount drupal source code
      // mount composer cache
      Container php = client
          .container()
          .from("php:8.2-cli")
          .withDirectory("/opt/drupal", drupalDir)
          .withWorkdir("/opt/drupal/web")
          .withMountedCache("/root/.composer/cache", client.cacheVolume("composer-cache"));

      // install php dependencies
      // install drupal dependencies
      php = php
          .withExec(List.of("apt-get", "update"))
          .withExec(
              List.of(
                  "apt-get",
                  "install",
                  "-y",
                  "git",
                  "libsqlite3-dev",
                  "libxml2-dev",
                  "zip"))
          .withExec(
              List.of("docker-php-ext-install", "gd", "pdo_mysql", "pdo_sqlite", "xml"))
          .withExec(List.of("pecl", "install", "xdebug"))
          .withExec(List.of("docker-php-ext-enable", "xdebug"))
          .withExec(
              List.of(
                  "php",
                  "-r",
                  "copy('https://getcomposer.org/installer', 'composer-setup.php');"))
          .withExec(List.of("php", "composer-setup.php"))
          .withExec(List.of("php", "-r", "unlink('composer-setup.php');"))
          .withExec(List.of("mv", "composer.phar", "/usr/local/bin/composer"))
          .withExec(List.of("composer", "install"));

      // bind mariadb service
      // set database url env var
      // execute tests
      // return test output
      return php
          .withServiceBinding("db", this.mariadbService())
          .withEnvVariable("SIMPLETEST_DB", "mysql://root:secret@db/drupal")
          .withExec(
              List.of(
                  "../../vendor/bin/phpunit",
                  "-c",
                  "core/phpunit.xml.dist",
                  "core/modules/user/tests/src/Kernel"))
          .stdout()
          .get();
    }
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c test
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
test
```

### Dagger CLI
```shell
dagger call test
```

The result will be:

```shell
PHPUnit 9.6.17 by Sebastian Bergmann and contributors.
Runtime:       PHP 8.2.5
Configuration: /opt/drupal/web/core/phpunit.xml.dist
Testing
.....................S                                            22 / 22 (100%)
Time: 00:15.806, Memory: 315.00 MB
There was 1 skipped test:

1) Drupal\Tests\pgsql\Kernel\pgsql\KernelTestBaseTest::testSetUp

This test only runs for the database driver 'pgsql'. Current database driver is 'mysql'.
/opt/drupal/web/core/tests/Drupal/KernelTests/Core/Database/DriverSpecificKernelTestBase.php:44
/opt/drupal/vendor/phpunit/phpunit/src/Framework/TestResult.php:728

OK, but incomplete, skipped, or risky tests!
Tests: 22, Assertions: 72, Skipped: 1.
```

## Reference: How service binding works in Dagger Functions

If you're not interested in what's happening in the background, you can skip this section and just trust that services are running when they need to be. If you're interested in the theory, keep reading.

Consider this example:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Pings the Redis service
func (m *MyModule) Ping(ctx context.Context) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "ping"}).
		Stdout(ctx)
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redisSrv`.
1. Dagger starts the `redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### Python

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def ping(self) -> str:
        """Pings the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "ping"])
            .stdout()
        )

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redis_srv`.
1. Dagger starts the `redis_srv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redis_srv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Pings the Redis service
   */
  @func()
  async ping(): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "ping"])
      .stdout()
  }
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redisSrv`.
1. Dagger starts the `redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### PHP

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Pings the Redis service
     */
    #[DaggerFunction]
    public function ping(): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'ping'])
            ->stdout();
    }
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `$redisSrv`.
1. Dagger starts the `$redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `$redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Pings the Redis service
   */
  @Function
  public String ping() throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "ping"))
          .stdout()
          .get();
    }
  }
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redisSrv`.
1. Dagger starts the `redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

> **Note:**
> Dagger cancels each service run after a 10 second grace period to avoid frequent restarts, unless the explicit `Start` and `Stop` APIs are used.

Services are based on containers, but they run a little differently. Whereas regular containers in Dagger are de-duplicated across the entire Dagger Engine, service containers are only de-duplicated within a Dagger client session. This means that if you run separate Dagger sessions that use the exact same services, they will each get their own "instance" of the service. This process is carefully tuned to preserve caching at each client call-site, while prohibiting "cross-talk" from one Dagger session's client to another Dagger session's service.

Content-addressed services are very convenient. You don't have to come up with names and maintain instances of services; just use them by value. You also don't have to manage the state of the service; you can just trust that it will be running when needed and stopped when not.

> **Tip:**
> If you need multiple instances of a service, just attach something unique to each one, such as an instance ID.

Here's a more detailed client-server example of running commands against a Redis service:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Sets and gets a key in the Redis service
func (m *MyModule) SetGet(ctx context.Context, key string, value string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key}).
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
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set_get(self, key: str, value: str) -> str:
        """Sets and gets a key in the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @func()
  async setGet(key: string, value: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets and gets a key in the Redis service
     */
    #[DaggerFunction]
    public function setGet(string $key, string $value): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @Function
  public String setGet(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key))
          .stdout()
          .get();
    }
  }
}

```

This example relies on the 10-second grace period, which you should try to avoid. Depending on the 10-second grace period is risky because there are many factors which could cause a 10-second delay between calls to Dagger, such as excessive CPU load, high network latency between the client and Dagger, or Dagger operations that require a variable amount of time to process.

It would be better to chain both commands together, which ensures that the service stays running for both, as in the revision below:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Sets and gets a key in the Redis service
func (m *MyModule) SetGet(ctx context.Context, key string, value string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key}).
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
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set_get(self, key: str, value: str) -> str:
        """Sets and gets a key in the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @func()
  async setGet(key: string, value: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets and gets a key in the Redis service
     */
    #[DaggerFunction]
    public function setGet(string $key, string $value): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @Function
  public String setGet(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key))
          .stdout()
          .get();
    }
  }
}		Stdout(ctx)
}

// Gets a key from the Redis service
func (m *MyModule) Get(ctx context.Context, key string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key}).
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
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        # create cache volume for redis data
        redis_data = dag.cache_volume("redis-data")

        # create redis container
        # mount cache volume to /data
        # expose redis port 6379
        # start service
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_mounted_cache("/data", redis_data)
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set(self, key: str, value: str) -> str:
        """Sets a key in the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "save"])
            .stdout()
        )

    @function
    async def get(self, key: str) -> str:
        """Gets a key from the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, CacheVolume, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    // create cache volume for redis data
    const redisData: CacheVolume = dag.cacheVolume("redis-data")

    // create redis container
    // mount cache volume to /data
    // expose redis port 6379
    // start service
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withMountedCache("/data", redisData)
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets a key in the Redis service
   */
  @func()
  async set(key: string, value: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "save"])
      .stdout()
  }

  /**
   * Gets a key from the Redis service
   */
  @func()
  async get(key: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        // create cache volume for redis data
        $redisData = dag()->cacheVolume('redis-data');

        // create redis container
        // mount cache volume to /data
        // expose redis port 6379
        // start service
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withMountedCache('/data', $redisData)
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets a key in the Redis service
     */
    #[DaggerFunction]
    public function set(string $key, string $value): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'save'])
            ->stdout();
    }

    /**
     * Gets a key from the Redis service
     */
    #[DaggerFunction]
    public function get(string $key): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.CacheVolume;
import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      // create cache volume for redis data
      CacheVolume redisData = client.cacheVolume("redis-data");

      // create redis container
      // mount cache volume to /data
      // expose redis port 6379
      // start service
      return client
          .container()
          .from("redis:7.2-alpine")
          .withMountedCache("/data", redisData)
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets a key in the Redis service
   */
  @Function
  public String set(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "save"))
          .stdout()
          .get();
    }
  }

  /**
   * Gets a key from the Redis service
   */
  @Function
  public String get(String key) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key))
          .stdout()
          .get();
    }
  }
}

```


This example uses Redis's `SAVE` command to save the service's data to a cache volume. When a new instance of the service is created, it uses the same cache volume to recreate the original state.

Here is an example of using these Dagger Functions:

### System shell
```shell
dagger -c 'set foo 123'
dagger -c 'get foo'
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
set foo 123
get foo
```

### Dagger CLI
```shell
dagger call set --key=foo --value=123
dagger call get --key=foo
```

The result will be:

```shell
123
```

## Start and stop services

Services are designed to be expressed as a Directed Acyclic Graph (DAG) with explicit bindings allowing services to be started lazily, just like every other DAG node. But sometimes, you may need to explicitly manage the lifecycle in a Dagger Function.

For example, this may be needed if the application in the service has certain behavior on shutdown (such as flushing data) that needs careful coordination with the rest of your logic.

The following example explicitly starts the Redis service and stops it at the end, ensuring the 10 second grace period doesn't get in the way, without the need for a persistent cache volume:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Sets and gets a key in the Redis service
func (m *MyModule) SetGet(ctx context.Context, key string, value string) (string, error) {
	// start redis service
	redisSrv, err := m.RedisService(ctx).Start(ctx)
	if err != nil {
		return "", err
	}

	// create redis client container
	// bind redis service
	// execute redis-cli command
	redisCLI := dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", redisSrv).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key})

	// get result
	val, err := redisCLI.Stdout(ctx)
	if err != nil {
		return "", err
	}

	// stop redis service
	_, err = redisSrv.Stop(ctx)
	if err != nil {
		return "", err
	}

	// return result
	return val, nil
}

```

### Python

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set_get(self, key: str, value: str) -> str:
        """Sets and gets a key in the Redis service"""
        # start redis service
        redis_srv = await self.redis_service().start()

        # create redis client container
        # bind redis service
        # execute redis-cli command
        redis_cli = (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", redis_srv)
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
        )

        # get result
        val = await redis_cli.stdout()

        # stop redis service
        await redis_srv.stop()

        # return result
        return val

```

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @func()
  async setGet(key: string, value: string): Promise<string> {
    // start redis service
    const redisSrv = await this.redisService().start()

    // create redis client container
    // bind redis service
    // execute redis-cli command
    const redisCLI: Container = dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", redisSrv)
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])

    // get result
    const val = await redisCLI.stdout()

    // stop redis service
    await redisSrv.stop()

    // return result
    return val
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets and gets a key in the Redis service
     */
    #[DaggerFunction]
    public function setGet(string $key, string $value): string
    {
        // start redis service
        $redisSrv = $this->redisService()->start();

        // create redis client container
        // bind redis service
        // execute redis-cli command
        $redisCLI = dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $redisSrv)
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key]);

        // get result
        $val = $redisCLI->stdout();

        // stop redis service
        $redisSrv->stop();

        // return result
        return $val;
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Container;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @Function
  public String setGet(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // start redis service
      Service redisSrv = this.redisService().start().get();

      // create redis client container
      // bind redis service
      // execute redis-cli command
      Container redisCLI = client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", redisSrv)
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key));

      // get result
      String val = redisCLI.stdout().get();

      // stop redis service
      redisSrv.stop().get();

      // return result
      return val;
    }
  }
}

```

## Example: MariaDB database service for application tests

The following example demonstrates how services can be used in Dagger Functions, by creating a Dagger Function for application unit/integration testing against a bound MariaDB database service.

The application used in this example is [Drupal](https://www.drupal.org/), a popular open-source PHP CMS. Drupal includes a large number of unit tests, including tests which require an active database connection. All Drupal 10.x tests are written and executed using the [PHPUnit](https://phpunit.de/) testing framework. Read more about [running PHPUnit tests in Drupal](https://phpunit.de/).

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a MariaDB service
func (m *MyModule) MariaDBService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("mariadb:10.11.2").
		WithEnvVariable("MARIADB_ROOT_PASSWORD", "secret").
		WithEnvVariable("MARIADB_DATABASE", "drupal").
		WithExposedPort(3306).
		AsService()
}

// Tests a Drupal application using a MariaDB service
func (m *MyModule) Test(ctx context.Context) (string, error) {
	// get drupal source code
	drupalDir := dag.Git("https://git.drupalcode.org/project/drupal.git").
		Branch("10.1.x").
		Tree()

	// get php container
	// mount drupal source code
	// mount composer cache
	php := dag.Container().
		From("php:8.2-cli").
		WithDirectory("/opt/drupal", drupalDir).
		WithWorkdir("/opt/drupal/web").
		WithMountedCache("/root/.composer/cache", dag.CacheVolume("composer-cache"))

	// install php dependencies
	// install drupal dependencies
	php = php.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "libsqlite3-dev", "libxml2-dev", "zip"}).
		WithExec([]string{"docker-php-ext-install", "gd", "pdo_mysql", "pdo_sqlite", "xml"}).
		WithExec([]string{"pecl", "install", "xdebug"}).
		WithExec([]string{"docker-php-ext-enable", "xdebug"}).
		WithExec([]string{"php", "-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"}).
		WithExec([]string{"php", "composer-setup.php"}).
		WithExec([]string{"php", "-r", "unlink('composer-setup.php');"}).
		WithExec([]string{"mv", "composer.phar", "/usr/local/bin/composer"}).
		WithExec([]string{"composer", "install"})

	// bind mariadb service
	// set database url env var
	// execute tests
	// return test output
	return php.
		WithServiceBinding("db", m.MariaDBService(ctx)).
		WithEnvVariable("SIMPLETEST_DB", "mysql://root:secret@db/drupal").
		WithExec([]string{"../../vendor/bin/phpunit", "-c", "core/phpunit.xml.dist", "core/modules/user/tests/src/Kernel"}).
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
    def mariadb_service(self) -> dagger.Service:
        """Returns a MariaDB service"""
        return (
            dag.container()
            .from_("mariadb:10.11.2")
            .with_env_variable("MARIADB_ROOT_PASSWORD", "secret")
            .with_env_variable("MARIADB_DATABASE", "drupal")
            .with_exposed_port(3306)
            .as_service()
        )

    @function
    async def test(self) -> str:
        """Tests a Drupal application using a MariaDB service"""
        # get drupal source code
        drupal_dir = dag.git("https://git.drupalcode.org/project/drupal.git").branch(
            "10.1.x"
        ).tree()

        # get php container
        # mount drupal source code
        # mount composer cache
        php = (
            dag.container()
            .from_("php:8.2-cli")
            .with_directory("/opt/drupal", drupal_dir)
            .with_workdir("/opt/drupal/web")
            .with_mounted_cache(
                "/root/.composer/cache", dag.cache_volume("composer-cache")
            )
        )

        # install php dependencies
        # install drupal dependencies
        php = (
            php.with_exec(["apt-get", "update"])
            .with_exec(
                ["apt-get", "install", "-y", "git", "libsqlite3-dev", "libxml2-dev", "zip"]
            )
            .with_exec(["docker-php-ext-install", "gd", "pdo_mysql", "pdo_sqlite", "xml"])
            .with_exec(["pecl", "install", "xdebug"])
            .with_exec(["docker-php-ext-enable", "xdebug"])
            .with_exec(["php", "-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"])
            .with_exec(["php", "composer-setup.php"])
            .with_exec(["php", "-r", "unlink('composer-setup.php');"])
            .with_exec(["mv", "composer.phar", "/usr/local/bin/composer"])
            .with_exec(["composer", "install"])
        )

        # bind mariadb service
        # set database url env var
        # execute tests
        # return test output
        return await (
            php.with_service_binding("db", self.mariadb_service())
            .with_env_variable("SIMPLETEST_DB", "mysql://root:secret@db/drupal")
            .with_exec(
                [
                    "../../vendor/bin/phpunit",
                    "-c",
                    "core/phpunit.xml.dist",
                    "core/modules/user/tests/src/Kernel",
                ]
            )
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Directory, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a MariaDB service
   */
  @func()
  mariadbService(): Service {
    return dag
      .container()
      .from("mariadb:10.11.2")
      .withEnvVariable("MARIADB_ROOT_PASSWORD", "secret")
      .withEnvVariable("MARIADB_DATABASE", "drupal")
      .withExposedPort(3306)
      .asService()
  }

  /**
   * Tests a Drupal application using a MariaDB service
   */
  @func()
  async test(): Promise<string> {
    // get drupal source code
    const drupalDir: Directory = dag
      .git("https://git.drupalcode.org/project/drupal.git")
      .branch("10.1.x")
      .tree()

    // get php container
    // mount drupal source code
    // mount composer cache
    let php: Container = dag
      .container()
      .from("php:8.2-cli")
      .withDirectory("/opt/drupal", drupalDir)
      .withWorkdir("/opt/drupal/web")
      .withMountedCache("/root/.composer/cache", dag.cacheVolume("composer-cache"))

    // install php dependencies
    // install drupal dependencies
    php = php
      .withExec(["apt-get", "update"])
      .withExec([
        "apt-get",
        "install",
        "-y",
        "git",
        "libsqlite3-dev",
        "libxml2-dev",
        "zip",
      ])
      .withExec([
        "docker-php-ext-install",
        "gd",
        "pdo_mysql",
        "pdo_sqlite",
        "xml",
      ])
      .withExec(["pecl", "install", "xdebug"])
      .withExec(["docker-php-ext-enable", "xdebug"])
      .withExec([
        "php",
        "-r",
        "copy('https://getcomposer.org/installer', 'composer-setup.php');",
      ])
      .withExec(["php", "composer-setup.php"])
      .withExec(["php", "-r", "unlink('composer-setup.php');"])
      .withExec(["mv", "composer.phar", "/usr/local/bin/composer"])
      .withExec(["composer", "install"])

    // bind mariadb service
    // set database url env var
    // execute tests
    // return test output
    return await php
      .withServiceBinding("db", this.mariadbService())
      .withEnvVariable("SIMPLETEST_DB", "mysql://root:secret@db/drupal")
      .withExec([
        "../../vendor/bin/phpunit",
        "-c",
        "core/phpunit.xml.dist",
        "core/modules/user/tests/src/Kernel",
      ])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a MariaDB service
     */
    #[DaggerFunction]
    public function mariadbService(): Service
    {
        return dag()
            ->container()
            ->from('mariadb:10.11.2')
            ->withEnvVariable('MARIADB_ROOT_PASSWORD', 'secret')
            ->withEnvVariable('MARIADB_DATABASE', 'drupal')
            ->withExposedPort(3306)
            ->asService();
    }

    /**
     * Tests a Drupal application using a MariaDB service
     */
    #[DaggerFunction]
    public function test(): string
    {
        // get drupal source code
        $drupalDir = dag()
            ->git('https://git.drupalcode.org/project/drupal.git')
            ->branch('10.1.x')
            ->tree();

        // get php container
        // mount drupal source code
        // mount composer cache
        $php = dag()
            ->container()
            ->from('php:8.2-cli')
            ->withDirectory('/opt/drupal', $drupalDir)
            ->withWorkdir('/opt/drupal/web')
            ->withMountedCache('/root/.composer/cache', dag()->cacheVolume('composer-cache'));

        // install php dependencies
        // install drupal dependencies
        $php = $php
            ->withExec(['apt-get', 'update'])
            ->withExec(['apt-get', 'install', '-y', 'git', 'libsqlite3-dev', 'libxml2-dev', 'zip'])
            ->withExec(['docker-php-ext-install', 'gd', 'pdo_mysql', 'pdo_sqlite', 'xml'])
            ->withExec(['pecl', 'install', 'xdebug'])
            ->withExec(['docker-php-ext-enable', 'xdebug'])
            ->withExec(['php', '-r', "copy('https://getcomposer.org/installer', 'composer-setup.php');"])
            ->withExec(['php', 'composer-setup.php'])
            ->withExec(['php', '-r', "unlink('composer-setup.php');"])
            ->withExec(['mv', 'composer.phar', '/usr/local/bin/composer'])
            ->withExec(['composer', 'install']);

        // bind mariadb service
        // set database url env var
        // execute tests
        // return test output
        return $php
            ->withServiceBinding('db', $this->mariadbService())
            ->withEnvVariable('SIMPLETEST_DB', 'mysql://root:secret@db/drupal')
            ->withExec(['../../vendor/bin/phpunit', '-c', 'core/phpunit.xml.dist', 'core/modules/user/tests/src/Kernel'])
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
import io.dagger.client.Directory;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a MariaDB service
   */
  @Function
  public Service mariadbService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("mariadb:10.11.2")
          .withEnvVariable("MARIADB_ROOT_PASSWORD", "secret")
          .withEnvVariable("MARIADB_DATABASE", "drupal")
          .withExposedPort(3306)
          .asService();
    }
  }

  /**
   * Tests a Drupal application using a MariaDB service
   */
  @Function
  public String test() throws Exception {
    try (Client client = Dagger.connect()) {
      // get drupal source code
      Directory drupalDir = client
          .git("https://git.drupalcode.org/project/drupal.git")
          .branch("10.1.x")
          .tree();

      // get php container
      // mount drupal source code
      // mount composer cache
      Container php = client
          .container()
          .from("php:8.2-cli")
          .withDirectory("/opt/drupal", drupalDir)
          .withWorkdir("/opt/drupal/web")
          .withMountedCache("/root/.composer/cache", client.cacheVolume("composer-cache"));

      // install php dependencies
      // install drupal dependencies
      php = php
          .withExec(List.of("apt-get", "update"))
          .withExec(
              List.of(
                  "apt-get",
                  "install",
                  "-y",
                  "git",
                  "libsqlite3-dev",
                  "libxml2-dev",
                  "zip"))
          .withExec(
              List.of("docker-php-ext-install", "gd", "pdo_mysql", "pdo_sqlite", "xml"))
          .withExec(List.of("pecl", "install", "xdebug"))
          .withExec(List.of("docker-php-ext-enable", "xdebug"))
          .withExec(
              List.of(
                  "php",
                  "-r",
                  "copy('https://getcomposer.org/installer', 'composer-setup.php');"))
          .withExec(List.of("php", "composer-setup.php"))
          .withExec(List.of("php", "-r", "unlink('composer-setup.php');"))
          .withExec(List.of("mv", "composer.phar", "/usr/local/bin/composer"))
          .withExec(List.of("composer", "install"));

      // bind mariadb service
      // set database url env var
      // execute tests
      // return test output
      return php
          .withServiceBinding("db", this.mariadbService())
          .withEnvVariable("SIMPLETEST_DB", "mysql://root:secret@db/drupal")
          .withExec(
              List.of(
                  "../../vendor/bin/phpunit",
                  "-c",
                  "core/phpunit.xml.dist",
                  "core/modules/user/tests/src/Kernel"))
          .stdout()
          .get();
    }
  }
}

```

Here is an example call for this Dagger Function:

### System shell
```shell
dagger -c test
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
test
```

### Dagger CLI
```shell
dagger call test
```

The result will be:

```shell
PHPUnit 9.6.17 by Sebastian Bergmann and contributors.
Runtime:       PHP 8.2.5
Configuration: /opt/drupal/web/core/phpunit.xml.dist
Testing
.....................S                                            22 / 22 (100%)
Time: 00:15.806, Memory: 315.00 MB
There was 1 skipped test:

1) Drupal\Tests\pgsql\Kernel\pgsql\KernelTestBaseTest::testSetUp

This test only runs for the database driver 'pgsql'. Current database driver is 'mysql'.
/opt/drupal/web/core/tests/Drupal/KernelTests/Core/Database/DriverSpecificKernelTestBase.php:44
/opt/drupal/vendor/phpunit/phpunit/src/Framework/TestResult.php:728

OK, but incomplete, skipped, or risky tests!
Tests: 22, Assertions: 72, Skipped: 1.
```

## Reference: How service binding works in Dagger Functions

If you're not interested in what's happening in the background, you can skip this section and just trust that services are running when they need to be. If you're interested in the theory, keep reading.

Consider this example:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Pings the Redis service
func (m *MyModule) Ping(ctx context.Context) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "ping"}).
		Stdout(ctx)
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redisSrv`.
1. Dagger starts the `redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### Python

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def ping(self) -> str:
        """Pings the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "ping"])
            .stdout()
        )

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redis_srv`.
1. Dagger starts the `redis_srv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redis_srv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Pings the Redis service
   */
  @func()
  async ping(): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "ping"])
      .stdout()
  }
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redisSrv`.
1. Dagger starts the `redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### PHP

```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Pings the Redis service
     */
    #[DaggerFunction]
    public function ping(): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'ping'])
            ->stdout();
    }
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `$redisSrv`.
1. Dagger starts the `$redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `$redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Pings the Redis service
   */
  @Function
  public String ping() throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "ping"))
          .stdout()
          .get();
    }
  }
}

```

Here's what happens on the last line:

1. The client requests the `ping` container's stdout, which requires the container to run.
1. Dagger sees that the `ping` container has a service binding, `redisSrv`.
1. Dagger starts the `redisSrv` container, which recurses into this same process.
1. Dagger waits for health checks to pass against `redisSrv`.
1. Dagger runs the `ping` container with the `redis-srv` alias magically added to `/etc/hosts`.

> **Note:**
> Dagger cancels each service run after a 10 second grace period to avoid frequent restarts, unless the explicit `Start` and `Stop` APIs are used.

Services are based on containers, but they run a little differently. Whereas regular containers in Dagger are de-duplicated across the entire Dagger Engine, service containers are only de-duplicated within a Dagger client session. This means that if you run separate Dagger sessions that use the exact same services, they will each get their own "instance" of the service. This process is carefully tuned to preserve caching at each client call-site, while prohibiting "cross-talk" from one Dagger session's client to another Dagger session's service.

Content-addressed services are very convenient. You don't have to come up with names and maintain instances of services; just use them by value. You also don't have to manage the state of the service; you can just trust that it will be running when needed and stopped when not.

> **Tip:**
> If you need multiple instances of a service, just attach something unique to each one, such as an instance ID.

Here's a more detailed client-server example of running commands against a Redis service:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Sets and gets a key in the Redis service
func (m *MyModule) SetGet(ctx context.Context, key string, value string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key}).
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
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set_get(self, key: str, value: str) -> str:
        """Sets and gets a key in the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @func()
  async setGet(key: string, value: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets and gets a key in the Redis service
     */
    #[DaggerFunction]
    public function setGet(string $key, string $value): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @Function
  public String setGet(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key))
          .stdout()
          .get();
    }
  }
}

```

This example relies on the 10-second grace period, which you should try to avoid. Depending on the 10-second grace period is risky because there are many factors which could cause a 10-second delay between calls to Dagger, such as excessive CPU load, high network latency between the client and Dagger, or Dagger operations that require a variable amount of time to process.

It would be better to chain both commands together, which ensures that the service stays running for both, as in the revision below:

### Go

```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type MyModule struct{}

// Returns a Redis service
func (m *MyModule) RedisService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:7.2-alpine").
		WithExposedPort(6379).
		AsService()
}

// Sets and gets a key in the Redis service
func (m *MyModule) SetGet(ctx context.Context, key string, value string) (string, error) {
	// bind redis service to container
	// execute redis-cli command
	// return response
	return dag.Container().
		From("redis:7.2-alpine").
		WithServiceBinding("redis-srv", m.RedisService(ctx)).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "set", key, value}).
		WithExec([]string{"redis-cli", "-h", "redis-srv", "get", key}).
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
    def redis_service(self) -> dagger.Service:
        """Returns a Redis service"""
        return (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_exposed_port(6379)
            .as_service()
        )

    @function
    async def set_get(self, key: str, value: str) -> str:
        """Sets and gets a key in the Redis service"""
        # bind redis service to container
        # execute redis-cli command
        # return response
        return await (
            dag.container()
            .from_("redis:7.2-alpine")
            .with_service_binding("redis-srv", self.redis_service())
            .with_exec(["redis-cli", "-h", "redis-srv", "set", key, value])
            .with_exec(["redis-cli", "-h", "redis-srv", "get", key])
            .stdout()
        )

```

### TypeScript

```typescript
import { dag, Container, Service, func, object } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Returns a Redis service
   */
  @func()
  redisService(): Service {
    return dag
      .container()
      .from("redis:7.2-alpine")
      .withExposedPort(6379)
      .asService()
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @func()
  async setGet(key: string, value: string): Promise<string> {
    // bind redis service to container
    // execute redis-cli command
    // return response
    return await dag
      .container()
      .from("redis:7.2-alpine")
      .withServiceBinding("redis-srv", this.redisService())
      .withExec(["redis-cli", "-h", "redis-srv", "set", key, value])
      .withExec(["redis-cli", "-h", "redis-srv", "get", key])
      .stdout()
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
use Dagger\Client\Service;

use function Dagger\dag;

#[DaggerObject]
class MyModule
{
    /**
     * Returns a Redis service
     */
    #[DaggerFunction]
    public function redisService(): Service
    {
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withExposedPort(6379)
            ->asService();
    }

    /**
     * Sets and gets a key in the Redis service
     */
    #[DaggerFunction]
    public function setGet(string $key, string $value): string
    {
        // bind redis service to container
        // execute redis-cli command
        // return response
        return dag()
            ->container()
            ->from('redis:7.2-alpine')
            ->withServiceBinding('redis-srv', $this->redisService())
            ->withExec(['redis-cli', '-h', 'redis-srv', 'set', $key, $value])
            ->withExec(['redis-cli', '-h', 'redis-srv', 'get', $key])
            ->stdout();
    }
}

```

### Java

```java
package io.dagger.modules.mymodule;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Service;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

import java.util.List;

@Module
@Object
public class MyModule {

  /**
   * Returns a Redis service
   */
  @Function
  public Service redisService() throws Exception {
    try (Client client = Dagger.connect()) {
      return client
          .container()
          .from("redis:7.2-alpine")
          .withExposedPort(6379)
          .asService();
    }
  }

  /**
   * Sets and gets a key in the Redis service
   */
  @Function
  public String setGet(String key, String value) throws Exception {
    try (Client client = Dagger.connect()) {
      // bind redis service to container
      // execute redis-cli command
      // return response
      return client
          .container()
          .from("redis:7.2-alpine")
          .withServiceBinding("redis-srv", this.redisService())
          .withExec(List.of("redis-cli", "-h", "redis-srv", "set", key, value))
          .withExec(List.of("redis-cli", "-h", "redis-srv", "get", key))
          .stdout()
          .get();
    }
  }
}
