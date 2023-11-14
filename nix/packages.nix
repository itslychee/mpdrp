pkgs:
with pkgs;
let 

in rec {
    mpdrp = buildGoApplication rec {
        name = "mpdrp";
        pname = name;
        go = pkgs.go;
        modules = ./gomod2nix.toml;
        src = ../.;
        subPackages = [
            "cmd/mpdrp"
        ];
    };
    mpdrp-mpc = buildGoApplication rec {
        name = "mpdrp-mpc";
        pname = name;
        go = pkgs.go;
        modules = ./gomod2nix.toml;
        src = ../.;
        subPackages = [
            "cmd/mpc"
        ];
    };
} 
