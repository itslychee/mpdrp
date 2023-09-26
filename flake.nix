{
    inputs = {
        utils.url = "github:numtide/flake-utils";
    };
    outputs = { utils, nixpkgs, ...}@inputs: utils.lib.eachDefaultSystem (system: let
        pkgs = import nixpkgs { inherit system; }; 
    in with pkgs; {
        devShells.default = mkShell {
            packages = [ go ];
        };
    });
}