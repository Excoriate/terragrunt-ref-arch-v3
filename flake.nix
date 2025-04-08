{
  description = "Terragrunt Reference Architecture Development Environment";

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

        # Essential tools - minimal set for fast shell startup
        essentialTools = with pkgs; [
          # Core infrastructure tools
          terraform
          terragrunt

          # Basic utilities
          git
          bash
        ];

        # Extended tools - available in the full development shell
        extendedTools = with pkgs; [
          # Go toolchain
          go
          go-tools
          golangci-lint

          # Additional Terraform tools
          terraform-ls
          tflint
          opentofu
          terraform-docs

          # Development and utility tools
          just

          # YAML tools
          yamllint
          yamlfmt

          # Shell scripting
          shellcheck

          # Environment management
          direnv
        ];

        # Show versions of key tools - extracted to a function to keep shellHook minimal
        showVersions = pkgs.writeShellScriptBin "show-versions" ''
          echo "Tool versions:"
          echo "-------------"
          echo "Terraform: $(terraform version | head -n 1)"
          echo "Terragrunt: $(terragrunt --version)"
          [ -x "$(command -v tofu)" ] && echo "OpenTofu: $(tofu version | head -n 1)"
          [ -x "$(command -v go)" ] && echo "Go: $(go version)"
          echo ""
          echo "Run 'show-versions' anytime to see this information again."
        '';
      in
      {
        # Fast startup default shell with minimal tools
        devShells.default = pkgs.mkShell {
          buildInputs = essentialTools ++ [ showVersions ];

          shellHook = ''
            echo "🚀 Fast Terragrunt Ref Arch Shell (minimal) 🛠️"
            echo "Type 'show-versions' for tool version information"
            echo "For full development environment with all tools: nix develop .#full"
          '';
        };

        # Full development shell with all tools
        devShells.full = pkgs.mkShell {
          buildInputs = essentialTools ++ extendedTools ++ [ showVersions ];

          shellHook = ''
            echo "🚀 Complete Terragrunt Ref Arch Development Environment 🛠️"
            echo "Type 'show-versions' for tool version information"
          '';
        };
      }
    );
}
