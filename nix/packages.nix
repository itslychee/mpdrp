pkgs:
{
    mpdrp = pkgs.buildGoApplication rec {
        name = "mpdrp";
        pname = name;
        go = pkgs.go;
        modules = ./gomod2nix.toml;
        src = ../.;
        doCheck = false;
        subPackages = [
            "cmd/mpdrp"
        ];
    };
    mpdrp-mpc = pkgs.buildGoApplication rec {
        name = "mpdrp-mpc";
        pname = "mpc";
        go = pkgs.go;
        modules = ./gomod2nix.toml;
        src = ../.;
        doCheck = false;
        subPackages = [
            "cmd/mpc"
        ];
    };
} 
