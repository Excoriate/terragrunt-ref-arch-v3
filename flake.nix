{
  description = "Terragrunt Reference Architecture Development Shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11-small";
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
          config.allowUnfree = true;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = [ pkgs.bashInteractive ];
          buildInputs = with pkgs; [
            terraform_1
            terragrunt
            just
            direnv
          ];

          shellHook = ''
            echo "🚀 Terragrunt Dev Shell | Run 'just --list' for commands"
          '';
        };
      }
    );
}
