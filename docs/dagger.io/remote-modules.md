---
slug: /api/remote-modules
---

# Using Modules from Remote Repositories

Dagger supports the use of HTTP and SSH protocols for accessing remote repositories as Dagger [modules](../features/modules.md), compatible with all major Git hosting platforms such as GitHub, GitLab, BitBucket, Azure DevOps, Codeberg, and Sourcehut. Dagger supports authentication via both HTTPS (using Git credential managers) and SSH (using a unified authentication approach).

Dagger supports various reference schemes for Dagger modules, as below:

| Protocol | Scheme            | Authentication | Example |
|----------|-------------------|----------------|---------|
| HTTP(S)  | Go-like ref style | Git credential manager | `github.com/username/repo[/subdir][@version]`  |
| HTTP(S)  | Git HTTP style    | Git credential manager | `https://github.com/username/repo.git[/subdir][@version]` |
| SSH      | SCP-like          | SSH keys | `git@github.com:username/repo.git[/subdir][@version]`     |
| SSH      | Explicit SSH      | SSH keys | `ssh://git@github.com/username/repo.git[/subdir][@version]` |

Dagger provides additional flexibility in referencing modules through the following options:

- The `.git` extension is optional for HTTP refs or explicit SSH refs, except for [GitLab, when referencing modules stored on a private repo or private subgroup](https://gitlab.com/gitlab-org/gitlab-foss/-/blob/master/lib/gitlab/middleware/go.rb#L229-237).
- Monorepo support: Append `/subpath` to access a specific subdirectory within the repository.
- Version specification: Add `@version` to target a particular version of the module. This can be a tag, branch name, or full commit hash. If omitted, the default branch is used.

Here is an example of using a Go builder Dagger module from a public repository over HTTPS:

```shell
dagger -m github.com/kpenfound/dagger-modules/golang@v0.2.0 call \
  build --source=https://github.com/dagger/dagger --args=./cmd/dagger \
  export --path=./build
```

Here is the same example using SSH authentication. Note that this requires [SSH authentication to be properly configured](#configuring-ssh-authentication) on your Dagger host).

```shell
dagger -m git@github.com:kpenfound/dagger-modules/golang@v0.2.0 call \
  build --source=https://github.com/dagger/dagger --args=./cmd/dagger \
  export --path=./build
```

## Authentication methods

Dagger supports both HTTPS and SSH authentication for accessing remote repositories.

### HTTPS authentication

For HTTPS authentication, Dagger uses your system's configured Git credential manager. This means if you're already authenticated with your Git provider, Dagger will automatically use these credentials when needed.

The following credential helpers are supported:
- [Git Credential Manager](https://github.com/git-ecosystem/git-credential-manager)
- macOS Keychain
- Windows Credential Manager
- Custom credential helpers configured in your `.gitconfig`

To verify if your credentials are properly configured, try cloning a private repository (replace the placeholders below with valid values):

```shell
git clone https://github.com/USER/PRIVATE_REPOSITORY.git
```

If this works, Dagger will be able to use the same credentials to access your private repositories.

#### Credential manager configuration

- GitHub: Use [`gh auth login`](https://docs.github.com/en/get-started/getting-started-with-git/caching-your-github-credentials-in-git#github-cli) or [configure credentials via Git Credential Manager](https://docs.github.com/en/get-started/getting-started-with-git/caching-your-github-credentials-in-git#git-credential-manager)
- GitLab: Use [`glab auth login`](https://gitlab.com/gitlab-org/cli/-/blob/main/docs/source/auth/login.md) or [configure credentials via Git Credential Manager](https://github.com/git-ecosystem/git-credential-manager)
- Azure DevOps: Use [Git Credential Manager](https://learn.microsoft.com/en-us/azure/devops/repos/git/set-up-credential-managers)
- BitBucket: Configure credentials using the [Git credential system](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage) (the widely-adopted implementation is [Git Credential Manager](https://github.com/git-ecosystem/git-credential-manager))

### SSH authentication

Dagger mounts the socket specified by your host's `SSH_AUTH_SOCK` environment variable to the Dagger Engine. This is essential for SSH refs, as most Git servers use your SSH key for authentication and tracking purposes, even when cloning public repositories.

This means that you must ensure that the `SSH_AUTH_SOCK` environment variable is properly set in your environment when using SSH refs with Dagger.

[Read detailed instructions on setting up SSH authentication](https://docs.github.com/en/authentication/connecting-to-github-with-ssh), including how to generate SSH keys, start the SSH agent, and add your keys.

## Best practices

For quick and easy referencing:
- Copy the repository ref from your preferred Git server's UI.
- To specify a particular version or commit, append `#version` (for directory arguments) or `@version` (for modules).
- To target a specific directory within the repository, use the format `#version:subpath` (for directory arguments) or add a `/subpath` (for modules). Remember that the version is mandatory when specifying a subpath.
- For private repositories:
  - HTTPS: Ensure your Git credentials are properly configured using your provider's recommended method.
  - SSH: Make sure your SSH keys are properly set up and added to the SSH agent.


## Known limitations and workarounds

This section outlines current limitations and provides workarounds for common issues. We're actively working on improvements for these areas.

### Windows is not supported

Currently, SSH refs are fully supported on UNIX-based systems (Linux and macOS). Windows support is under development. Track progress and contribute to the discussion in our [GitHub issue for Windows support](https://github.com/dagger/dagger/issues/8313).

### Multiple SSH keys may cause SSH forwarding to fail

SSH forwarding may fail when multiple keys are loaded in your SSH agent. This is under active investigation in our [GitHub issue](https://github.com/dagger/dagger/issues/8288). Until this is resolved, the following workaround may be used:

1. Clear all loaded keys: `ssh-add -D`
2. Add back only the required key: `ssh-add /path/to/key`