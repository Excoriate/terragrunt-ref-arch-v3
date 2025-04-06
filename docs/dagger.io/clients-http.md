---
slug: /api/http
---

# Raw HTTP

The Dagger API is an HTTP API that uses GraphQL as its low-level language-agnostic framework. Therefore, it's possible to call the Dagger API using raw HTTP queries, from [any language that supports GraphQL](https://graphql.org/code/). GraphQL has a large and growing list of client implementations in over 20 languages.

> **Note:**
> In practice, calling the API using HTTP or GraphQL is optional. Typically, you will instead use a custom Dagger function created with a type-safe Dagger SDK, or from the command line using the Dagger CLI.

Dagger creates a unique local API endpoint for GraphQL HTTP queries for every Dagger session. This API endpoint is served by the local host at the port specified by the `DAGGER_SESSION_PORT` environment variable, and can be directly read from the environment in your client code. For example, if `DAGGER_SESSION_PORT` is set to `12345`, the API endpoint can be reached at `http://127.0.0.1:$DAGGER_SESSION_PORT/query`

> **Warning:**
> Dagger protects the exposed API with an HTTP Basic authentication token which can be retrieved from the `DAGGER_SESSION_TOKEN` variable. Treat the `DAGGER_SESSION_TOKEN` value as you would any other sensitive credential. Store it securely and avoid passing it to, or over, insecure applications and networks.

## Command-line HTTP clients

This example demonstrates how to connect to the Dagger API and run a simple pipeline using `curl`:

```shell
echo '{"query":"{
  container {
    from(address:\"alpine:latest\") {
      file(path:\"/etc/os-release\") {
        contents
      }
    }
  }
}"}'|   dagger run sh -c 'curl -s \
    -u $DAGGER_SESSION_TOKEN: \
    -H "content-type:application/json" \
    -d @- \
    http://127.0.0.1:$DAGGER_SESSION_PORT/query'
```

## Language-native HTTP clients

This example demonstrates how to connect to the Dagger API and run a simple pipeline in the following languages:

- Rust, using the [gql_client library](https://github.com/arthurkhlghatyan/gql-client-rs) (MIT License)
- PHP, using the [php-graphql-client library](https://github.com/mghoneimy/php-graphql-client) (MIT License)

Create a new directory for the project and install the GraphQL client.

### Rust

```shell
mkdir my-project
cd my-project
cargo init
cargo add gql_client@1.0.7
cargo add serde_json@1.0.125
cargo add tokio@1.39.3 -F full
cargo add base64@0.22.1
```

### PHP

```shell
mkdir my-project
cd my-project
composer require gmostafa/php-graphql-client
```

Once the client library is installed, create a Dagger API client.

### Rust

Add the following code to `src/main.rs`:

```rust
use base64::{engine::general_purpose, Engine as _};
use gql_client::Client;
use serde_json::Value;
use std::collections::HashMap;
use std::env;

#[tokio::main]
async fn main() {
    let port = env::var("DAGGER_SESSION_PORT").expect("DAGGER_SESSION_PORT is not set");
    let token = env::var("DAGGER_SESSION_TOKEN").expect("DAGGER_SESSION_TOKEN is not set");
    let endpoint = format!("http://127.0.0.1:{}/query", port);

    let auth_header = format!(
        "Basic {}",
        general_purpose::STANDARD.encode(format!("{}:", token))
    );

    let mut headers = HashMap::new();
    headers.insert("Authorization", auth_header);

    let client = Client::new_with_headers(endpoint, headers);

    let query = r#"
        query {
            container {
                from(address: "alpine:latest") {
                    withExec(args: ["uname", "-a"]) {
                        stdout
                    }
                }
            }
        }
    "#;

    let response = client
        .query_unwrap::<Value>(query)
        .await
        .expect("GraphQL query failed");

    println!(
        "{}",
        response["container"]["from"]["withExec"]["stdout"]
            .as_str()
            .unwrap()
    );
}

```

### PHP

Create a new file named `client.php` and add the following code to it:

```php
<?php

require_once(__DIR__ . '/vendor/autoload.php');

use GraphQL\Client;
use GraphQL\Exception\QueryError;
use GraphQL\Query;

$port = getenv('DAGGER_SESSION_PORT');
$token = getenv('DAGGER_SESSION_TOKEN');
$endpoint = "http://127.0.0.1:$port/query";

$client = new Client(
    $endpoint,
    ['Authorization' => 'Basic ' . base64_encode($token . ':')]
);

$gql = (new Query('container'))
    ->setSelectionSet([
        (new Query('from'))
            ->setArguments(['address' => 'alpine:latest'])
            ->setSelectionSet([
                (new Query('withExec'))
                    ->setArguments(['args' => ['uname', '-a']])
                    ->setSelectionSet(['stdout'])
            ])
    ]);

try {
    $results = $client->runQuery($gql);
    print_r($results->getData()['container']['from']['withExec']['stdout']);
}
catch (QueryError $exception) {
    print_r($exception->getErrorDetails());
    exit;
}

?>

```

This code listing initializes the GraphQL client library and defines the Dagger pipeline to be executed as a Dagger API query. The `dagger run` command takes care of initializing a new local instance (or reusing a running instance) of the Dagger Engine on the host system and executing a specified command against it.

Run the Dagger API client using the Dagger CLI as follows:

### Rust
```shell
dagger run cargo run
```

### PHP
```shell
dagger run php client.php
```

Here is an example of the output:

```shell
dagger 6.1.0-23-cloud-amd64 unknown Linux
```

## Dagger CLI

The Dagger CLI offers a `dagger query` sub-command, which provides an easy way to send raw GraphQL queries to the Dagger API from the command line.

This example demonstrates how to build a Go application by cloning the [canonical Git repository for Go](https://go.googlesource.com/example/+/HEAD/hello) and building the "Hello, world" example program from it by calling the Dagger API via `dagger query`.

Create a new shell script named `build.sh` and add the following code to it:

```shell
#!/bin/sh

set -e

# get source code directory ID
src=$(dagger query <<EOF | jq -re '.git."https://go.googlesource.com/example".branch("master").tree("hello").id'
{
  git(url: "https://go.googlesource.com/example") {
    branch(name: "master") {
      tree(path: "hello") {
        id
      }
    }
  }
}
EOF
)

# build application using source code directory ID
dagger query <<EOF | jq -re '.container.build.file("hello").export("./dagger-builds-hello")'
{
  container {
    from(address: "golang:latest") {
      withDirectory(path: "/src", directory: "$src") {
        withWorkdir(path: "/src") {
          withExec(args: ["go", "build", "-o", "hello"]) {
            build: directory(path: ".") {
              file(path: "hello") {
                export(path: "./dagger-builds-hello")
              }
            }
          }
        }
      }
    }
  }
}
EOF

```

This script uses `dagger query` to send two GraphQL queries to the Dagger API. The first query returns a content-addressed identifier of the source code directory from the remote Git repository. This is interpolated into the second query, which initializes a new container, mounts the source code directory, compiles the source code, and writes the compiled binary back to the host filesystem.

Add the executable bit to the shell script and then run it by executing the commands below:

```shell
chmod +x ./build.sh
./build.sh
```

On completion, the built Go application will be available in the working directory on the host, as shown below:

```shell
tree
.
├── build.sh
└── dagger-builds-hello

1 directory, 2 files