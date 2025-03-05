{
  description = "Terragrunt Reference Architecture Development Shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            terraform_1_9
            terragrunt
            just
            direnv
          ];

          shellHook = ''
            echo "🚀 Welcome to Terragrunt Reference Architecture Development Shell"
            echo "📦 Terraform v1.9 and Terragrunt are available"
            echo "🔧 Use 'just --list' to see available commands"
            eval "$(direnv hook bash)"
          '';
        };
      }
    );
}