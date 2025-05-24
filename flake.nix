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

        # --- Helper functions for architecture strings ---
        systemToArchCommon = sys:
          if pkgs.lib.strings.hasPrefix "x86_64-linux" sys then "linux_amd64"
          else if pkgs.lib.strings.hasPrefix "aarch64-linux" sys then "linux_arm64"
          else if pkgs.lib.strings.hasPrefix "x86_64-darwin" sys then "darwin_amd64"
          else if pkgs.lib.strings.hasPrefix "aarch64-darwin" sys then "darwin_arm64"
          else null;

        tgArchString = systemToArchCommon system;
        daggerArchString = systemToArchCommon system;
        tfArchString = systemToArchCommon system;

        # --- Pinned Terragrunt Derivation ---
        terragruntVersion = "0.80.2";
        terragrunt_pinned = if tgArchString == null
          then pkgs.terragrunt # Fallback or could throw error
          else pkgs.stdenvNoCC.mkDerivation {
            pname = "terragrunt-pinned";
            version = terragruntVersion;
            src = pkgs.fetchurl {
              url = "https://github.com/gruntwork-io/terragrunt/releases/download/v${terragruntVersion}/terragrunt_${tgArchString}";
              sha256 = "sha256-V6JLUqRk7SjE14aiOO57o6Rwsm9N6SFtPIzbKpNMM1A="; # Updated from placeholder
            };
            dontUnpack = true;
            installPhase = ''
              mkdir -p $out/bin
              cp $src $out/bin/terragrunt
              chmod +x $out/bin/terragrunt
            '';
            meta.platforms = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
          };

        # --- Pinned Dagger CLI Derivation ---
        daggerVersion = "0.18.8";
        dagger_cli_pinned = if daggerArchString == null
          then pkgs.dagger # Fallback, though dagger might not be in older nixpkgs by default
          else pkgs.stdenvNoCC.mkDerivation {
            pname = "dagger-cli-pinned";
            version = daggerVersion;
            src = pkgs.fetchurl {
              url = "https://dl.dagger.io/dagger/releases/v${daggerVersion}/dagger_v${daggerVersion}_${daggerArchString}.tar.gz";
              sha256 = "0000000000000000000000000000000000000000000000000000"; # Replace with actual hash
            };
            nativeBuildInputs = [ pkgs.gnutar ];
            installPhase = ''
              mkdir -p $out/bin
              tar -xzf $src -C $out/bin dagger
            '';
            meta.platforms = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
          };

        # --- Pinned Terraform Derivation (Example with 1.8.0) ---
        # TODO: User needs to confirm the Terraform version for "1.11.3" as it's not standard.
        # Using 1.8.0 as a placeholder for the derivation structure.
        terraformPinnedVersion = "1.8.0"; # Placeholder - replace with a valid, confirmed version
        terraform_pinned = if tfArchString == null
          then pkgs.terraform # Fallback
          else pkgs.stdenvNoCC.mkDerivation {
            pname = "terraform-pinned";
            version = terraformPinnedVersion;
            src = pkgs.fetchurl {
              url = "https://releases.hashicorp.com/terraform/${terraformPinnedVersion}/terraform_${terraformPinnedVersion}_${tfArchString}.zip";
              sha256 = "sha256-q/sG64DxrNGauKAfbSSkpfmbqbYow7AKOwyJhwnuo7M="; # Updated from placeholder
            };
            nativeBuildInputs = [ pkgs.unzip ];
            installPhase = ''
              mkdir -p $out/bin
              unzip -j $src terraform -d $out/bin # -j to junk paths, place directly in $out/bin
              chmod +x $out/bin/terraform
            '';
            meta.license = pkgs.lib.licenses.mpl20; # Terraform is MPL 2.0
            meta.platforms = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
          };

        # Essential tools - minimal set for fast shell startup
        essentialTools = [
          terraform_pinned # Pinned Terraform
          terragrunt_pinned # Pinned Terragrunt
          pkgs.git
          pkgs.bash
        ];

        # Extended tools - available in the full development shell
        extendedTools = with pkgs; [
          # Go toolchain
          go
          go-tools
          golangci-lint

          # Additional Terraform ecosystem tools (use from nixpkgs)
          terraform-ls
          tflint
          opentofu # User had this, keep it if desired alongside pinned Terraform
          terraform-docs

          # Development and utility tools
          just
          dagger_cli_pinned # Pinned Dagger CLI

          # YAML tools
          yamllint
          yamlfmt

          # Shell scripting
          shellcheck

          # Environment management
          direnv
        ];

        # Show versions of key tools
        showVersions = pkgs.writeShellScriptBin "show-versions" ''
          echo "Tool versions:"
          echo "-------------"
          echo "Terraform (pinned): $(terraform version | head -n 1)"
          echo "Terragrunt (pinned): $(terragrunt --version | head -n 1)" # terragrunt --version can be verbose
          echo "Dagger CLI (pinned): $(dagger version)"
          [ -x "$(command -v tofu)" ] && echo "OpenTofu: $(tofu version | head -n 1)"
          [ -x "$(command -v go)" ] && echo "Go: $(go version)"
          echo ""
          echo "NOTE: Terraform version 1.11.3 requested by user is non-standard."
          echo "The pinned Terraform is currently set to ${terraformPinnedVersion} as a placeholder."
          echo "Please update flake.nix with the correct version and sha256 hash."
          echo ""
          echo "Run 'show-versions' anytime to see this information again."
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = essentialTools ++ [ showVersions ];

          shellHook = ''
            echo "🚀 Fast Terragrunt Ref Arch Shell (minimal) 🛠️"
            echo "Pinned Tools: Terragrunt v${terragruntVersion}, Terraform v${terraformPinnedVersion} (Placeholder)"
            echo "Type 'show-versions' for tool version information"
            echo "For full development environment with all tools: nix develop .#full"
            echo ""
            echo "IMPORTANT: The requested Terraform version 1.11.3 is not standard."
            echo "This shell uses Terraform v${terraformPinnedVersion} as a placeholder."
            echo "Please verify the required Terraform version and update flake.nix."
            echo "You will also need to update sha256 hashes for downloaded binaries."
          '';
        };

        devShells.full = pkgs.mkShell {
          buildInputs = essentialTools ++ extendedTools ++ [ showVersions ];

          shellHook = ''
            echo "🚀 Complete Terragrunt Ref Arch Development Environment 🛠️"
            echo "Pinned Tools: Terragrunt v${terragruntVersion}, Terraform v${terraformPinnedVersion} (Placeholder), Dagger v${daggerVersion}"
            echo "Type 'show-versions' for tool version information"
            echo ""
            echo "IMPORTANT: The requested Terraform version 1.11.3 is not standard."
            echo "This shell uses Terraform v${terraformPinnedVersion} as a placeholder."
            echo "Please verify the required Terraform version and update flake.nix."
            echo "You will also need to update sha256 hashes for downloaded binaries."
          '';
        };
      }
    );
}
