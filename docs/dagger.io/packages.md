---
slug: /api/packages
---

# Third-Party Packages

Dagger Functions are just regular code, written in your usual programming language. One of the key advantages of this approach is that it opens up access to your language's existing ecosystem of packages or modules. You can easily import these packages/modules in your Dagger module via your language's package manager.

### Go

To add a Go module, add it to your `go.mod` file using `go get`. For example:

```shell
go get github.com/spf13/cobra
```

### Python

To add a Python package, add it to your `pyproject.toml` file using your chosen package manager. For example:

#### uv

```sh
uv add requests
```

#### poetry

```sh
poetry add requests
```

#### uv pip

Add the dependency manually to [`pyproject.toml`](https://packaging.python.org/en/latest/guides/writing-pyproject-toml/#dependencies-and-requirements):

```toml
[project]
dependencies = [
    "requirements>=2.32.3",
]
```

Then install into your virtual environment:

```sh
uv pip install -e ./sdk -e .
```

> **Note:**
> There's no need to activate the virtual environment before `uv pip install`, but it does need to exist.

#### pip

Add the dependency manually to [`pyproject.toml`](https://packaging.python.org/en/latest/guides/writing-pyproject-toml/#dependencies-and-requirements):

```toml
[project]
dependencies = [
    "requirements>=2.32.3",
]
```

Then install into your virtual environment:

```sh
python -m pip install -e ./sdk -e .
```

> **Tip:**
> If you haven't setup your local environment yet, see [IDE Integration](./ide-integration.md).

> **Note:**
> Third-party dependencies are managed in the same way as any normal Python project. The only limitation is in "pinning" the dependencies. Currently, Dagger can install directly from a `uv.lock` file, or a [pip-tools compatible](https://docs.astral.sh/uv/pip/compile/#locking-requirements) `requirements.lock` file (notice `.lock` extension, not `.txt`). See [Language-native packaging](./custom-functions.md#language-native-packaging) for more information.

### TypeScript

To add a TypeScript package, add it to the `package.json` file using your favorite package manager. For example:

```shell
npm install pm2
```

Pinning a specific dependency version or adding local dependencies are supported, in the same way as any Node.js project.

### PHP

To add a PHP package, add it to the `composer.json` file, the same way as any PHP project. For example:

```shell
composer require phpunit/phpunit
```

> **Note:**
> Dagger modules installed as packages via Composer are not registered with Dagger.
>
> You can access its code, like any other PHP package, but this is not the indended use-case of a Dagger module.
> This may lead to unexpected behaviour.
>
> Use Composer for standard third-party packages.
>
> Use Dagger to [install Dagger modules](./module-dependencies.md)

### Java

To add a Java package, add it to your `pom.xml` file using Maven. For example:

```xml
<dependency>
    <groupId>org.slf4j</groupId>
    <artifactId>slf4j-simple</artifactId>
    <scope>runtime</scope>
    <version>2.0.16</version>
</dependency>