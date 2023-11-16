{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };
  outputs = { self, gomod2nix, nixpkgs}: let
    supportedSystems = [
      "x86_64-linux"
      "aarch64-linux" 
      "aarch64-darwin"
      "x86_64-darwin"
    ];
    withSystems = f: nixpkgs.lib.genAttrs supportedSystems
    (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ gomod2nix.overlays.default ];
      };
    in (f { inherit system pkgs; }));
  in {
    devShells = withSystems ({pkgs, system}: {
      default = pkgs.mkShell {
        packages = with pkgs; [ gomod2nix.packages.${system}.default go ];
      };
    });
    packages = withSystems ({ pkgs, system }: (import nix/packages.nix pkgs));
    nixosModules.default = import ./nix/module.nix;
  };
}
