{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };
  outputs = { self, nixpkgs, }: let
    supportedSystems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];
    withSystems = f:
      nixpkgs.lib.genAttrs supportedSystems (system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in
        f {inherit system pkgs;});
  in {
    packages = withSystems ({
      pkgs,
      system,
    }: import ./. { inherit pkgs; });
    formatter = withSystems ({ pkgs, system, }: pkgs.alejandra);
    homeManagerModules.default = import ./nix/module.nix self;
  };
}
