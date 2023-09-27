{
    inputs.utils.url = "github:numtide/flake-utils";
    inputs.gomod2nix.url = "github:nix-community/gomod2nix";
    outputs = { self, nixpkgs, utils, gomod2nix, ...}@inputs: utils.lib.eachDefaultSystem (system: 
    let
        pkgs = import nixpkgs {
            inherit system;
            overlays = [ gomod2nix.overlays.default ];
        }; 
    in with pkgs; {
        devShells.default = mkShell {
            packages = [ go gomod2nix.packages.${system}.default ];
        };
        # Making it accessible to other Nix users
        packages.default =  import nix/packages.nix pkgs;
        overlays.default = (_: _: { mpdrp = self.packages.${system}.default.mpdrp; });
    });
}