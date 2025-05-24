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
              sha256 = "sha256-V6okpSpGTSjE13hqoje3o6R4LJb4Ga3xwU8t08yaTNE="; # Terragrunt 0.80.2 darwin_arm64
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
              sha256 = "sha256-qCS8dznieOG0DqD4xOT0xveW/Dtjr/pumsf+ngHo2Tg="; # Dagger CLI 0.18.8 darwin_arm64
            };
            nativeBuildInputs = [ pkgs.gnutar ];
            installPhase = ''
              mkdir -p $out/bin
              tar -xzf $src -C $out/bin dagger
            '';
            meta.platforms = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
          };

        # --- Pinned Terraform Derivation ---
        terraformPinnedVersion = "1.11.3";
        terraform_pinned = if tfArchString == null
          then pkgs.terraform # Fallback
          else pkgs.stdenvNoCC.mkDerivation {
            pname = "terraform-pinned";
            version = terraformPinnedVersion;
            src = pkgs.fetchurl {
              url = "https://releases.hashicorp.com/terraform/${terraformPinnedVersion}/terraform_${terraformPinnedVersion}_${tfArchString}.zip";
              sha256 = "sha256-wMZPp7hZ9QX9zv2riTF+mLJo9o1AHah98LACHoJ88Zc="; # Terraform 1.11.3 darwin_arm64 - CORRECTED HASH
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
          echo "All tools pinned with verified SHA256 hashes:"
          echo "- Terraform v${terraformPinnedVersion}"
          echo "- Terragrunt v${terragruntVersion}"
          echo "- Dagger CLI v${daggerVersion}"
          echo ""
          echo "Run 'show-versions' anytime to see this information again."
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = essentialTools ++ [ showVersions ];

          shellHook = ''
            echo "🚀 Fast Terragrunt Ref Arch Shell (minimal) 🛠️"
            echo "Pinned Tools: Terragrunt v${terragruntVersion}, Terraform v${terraformPinnedVersion}"
            echo "Type 'show-versions' for tool version information"
            echo "For full development environment with all tools: nix develop .#full"
          '';
        };

        devShells.full = pkgs.mkShell {
          buildInputs = essentialTools ++ extendedTools ++ [ showVersions ];

          shellHook = ''
            echo "🚀 Complete Terragrunt Ref Arch Development Environment 🛠️"
            echo "Pinned Tools: Terragrunt v${terragruntVersion}, Terraform v${terraformPinnedVersion}, Dagger v${daggerVersion}"
            echo "Type 'show-versions' for tool version information"
          '';
        };
      }
    );
}
