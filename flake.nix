{
  description = "Terraform Registry Module Template Devshell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
            allowUnfreePredicate =
              pkg:
              builtins.elem (pkgs.lib.getName pkg) [
                "terraform"
                "opentofu"
              ];
          };
        };

        # Consolidated tools list
        devTools = with pkgs; [
          # Go toolchain
          go
          go-tools
          golangci-lint

          # Terraform and related
          terraform
          terraform-ls
          tflint
          opentofu
          terraform-docs
          terragrunt

          # Development and utility tools
          just
          git
          bash

          # YAML tools
          yamllint
          yamlfmt

          # Shell scripting
          shellcheck

          # direnv for environment management
          direnv
        ];
      in
      {
        # Development shell configuration
        devShells.default = pkgs.mkShell {
          buildInputs = devTools;

          shellHook = ''
            echo "🚀 Devshell Terragrunt Ref Arch 🛠️"
            echo "Go version: $(go version)"
            echo "Terraform version: $(terraform version)"
            echo "OpenTofu version: $(tofu version)"
            echo "Terragrunt version: $(terragrunt --version)"
          '';
        };
      }
    );
}
