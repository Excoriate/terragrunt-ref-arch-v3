---
slug: /api/cli
---

# Dagger CLI

The Dagger CLI lets you call both the core and extended Dagger API (the core APIs plus the new APIs provided by external Dagger modules) directly from the command-line.

You can call the API interactively (`dagger`) or non-interactively (`dagger -c`, `dagger call`, or `dagger core`).

Here are a few examples:

- Create a simple pipeline that is fully satisfied by the core Dagger API, without needing to program a Dagger module:

    ### System shell
    ```shell
    dagger <<EOF
    container |
      from cgr.dev/chainguard/wolfi-base |
      with-exec apk add go |
      with-directory /src https://github.com/golang/example#master |
      with-workdir /src/hello |
      with-exec -- go build -o hello . |
      file ./hello |
      export ./hello-from-dagger
    EOF
    ```

    ### Dagger Shell
    ```shell title="First type 'dagger' for interactive mode."
    container |
      from cgr.dev/chainguard/wolfi-base |
      with-exec apk add go |
      with-directory /src https://github.com/golang/example#master |
      with-workdir /src/hello |
      with-exec -- go build -o hello . |
      file ./hello |
      export ./hello-from-dagger
    ```

    ### Dagger CLI
    ```shell
    dagger core container from --address="cgr.dev/chainguard/wolfi-base" \
      with-exec --args="apk","add","go" \
      with-directory --path="/src" --directory="https://github.com/golang/example#master" \
      with-workdir --path="/src/hello" \
      with-exec --args="go","build","-o","hello","." \
      file --path="./hello" \
      export --path="./hello-from-dagger"
    ```

- Use Dagger as an alternative to `docker run`.

    ### System shell
    ```shell
    dagger -c 'container | from cgr.dev/chainguard/wolfi-base | terminal'
    ```

    ### Dagger Shell
    ```shell title="First type 'dagger' for interactive mode."
    container | from cgr.dev/chainguard/wolfi-base | terminal
    ```

    ### Dagger CLI
    ```shell
    dagger core container \
      from --address=cgr.dev/chainguard/wolfi-base \
      terminal
    ```

    > **Tip:**
    > If only the core Dagger API is needed, the `-M` (`--no-mod`) flag can be provided. This results in quicker startup, because the Dagger CLI doesn't try to find and load a current module. This also makes `dagger -M` equivalent to `dagger core`.

- Call one of the auto-generated Dagger Functions:

    ### System shell
    ```shell
    dagger -c 'container-echo "Welcome to Dagger!" | stdout'
    ```

    ### Dagger Shell
    ```shell title="First type 'dagger' for interactive mode."
    container-echo "Welcome to Dagger!" | stdout
    ```

    ### Dagger CLI
    ```shell
    dagger call container-echo --string-arg="Welcome to Dagger!" stdout
    ```

    > **Tip:**
    > When using the Dagger CLI, all names (functions, arguments, struct fields, etc) are converted into a shell-friendly "kebab-case" style.

- Modules don't need to be installed locally. Dagger lets you consume modules from GitHub repositories and call their Dagger Functions as though you were calling them locally:

    ### System shell
    ```shell
    dagger <<EOF
    github.com/jpadams/daggerverse/trivy@v0.5.0 |
      scan-image ubuntu:latest
    EOF
    ```

    ### Dagger Shell
    ```shell title="First type 'dagger' for interactive mode."
    github.com/jpadams/daggerverse/trivy@v0.5.0 | scan-image ubuntu:latest
    ```

    ### Dagger CLI
    ```shell
    dagger -m github.com/jpadams/daggerverse/trivy@v0.5.0 call \
      scan-image --image-ref=ubuntu:latest
    ```

- List all the Dagger Functions available in a module using context-sensitive help:

    ### System shell
    ```shell
    dagger -c '.help github.com/jpadams/daggerverse/trivy@v0.5.0'
    ```

    ### Dagger Shell
    ```shell title="First type 'dagger' for interactive mode."
    .help github.com/jpadams/daggerverse/trivy@v0.5.0
    ```

    ### Dagger CLI
    ```shell
    dagger -m github.com/jpadams/daggerverse/trivy@v0.5.0 call --help
    ```

- List all the optional and required arguments for a Dagger Function using context-sensitive help:

    ### System shell
    ```shell
    dagger -c 'github.com/jpadams/daggerverse/trivy@v0.5.0 | scan-image | .help'
    ```

    ### Dagger Shell
    ```shell title="First type 'dagger' for interactive mode."
    github.com/jpadams/daggerverse/trivy@v0.5.0 | scan-image | .help
    ```

    ### Dagger CLI
    ```shell
    dagger -m github.com/jpadams/daggerverse/trivy@v0.5.0 call scan-image --help