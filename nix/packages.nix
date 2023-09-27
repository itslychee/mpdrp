pkgs:
with pkgs;
let 
    drvTemplate = name: output: buildGoApplication {
        inherit name;
        pname = name;
        go = pkgs.go;
        modules = ./gomod2nix.toml;
        src = ../.;
        subPackages = [
            "cmd/${name}"
        ];
    };
in rec {
    mpdrp = drvTemplate "mpdrp" "mpdrp";
} 