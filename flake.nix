{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };
  outputs = {
    self,
    gomod2nix,
    nixpkgs,
  }: let
    supportedSystems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];
    withSystems = f:
      nixpkgs.lib.genAttrs supportedSystems (system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [gomod2nix.overlays.default];
        };
      in
        f {inherit system pkgs;});
  in {
    packages = withSystems ({
      pkgs,
      system,
    }:
      import nix/packages.nix {
        inherit (nixpkgs) lib;
        inherit pkgs;
      });
    homeManagerModules.default = import ./nix/module.nix self;
  };
}
