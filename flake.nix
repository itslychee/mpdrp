{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };
  outputs = {
    self,
    nixpkgs,
  }: let
    supportedSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "aarch64-darwin"
      "x86_64-darwin"
    ];
    withSystems = f:
      nixpkgs.lib.genAttrs
      supportedSystems
      (system: f nixpkgs.legacyPackages.${system});
  in {
    packages = withSystems (pkgs: import ./. {inherit pkgs;});
    formatter = withSystems (pkgs: pkgs.alejandra);
    homeManagerModules.default = import ./nix/module.nix self;
    devShell = withSystems (pkgs:
      pkgs.mkShell {
        packages = builtins.attrValues {
          inherit (pkgs) gopls go;
        };
      });
  };
}
