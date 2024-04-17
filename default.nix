{ pkgs ? import <nixpkgs>}:
rec {
    mpdrp = pkgs.callPackage ./nix/package.nix {
        name = "mpdrp";
        vendorHash = "sha256-LcYXXlHPHze9zWpJ6wB5o0py/wVzYW0r2m7liJF0uWg=";
    };
    mpdrp-mpc = pkgs.callPackage ./nix/package.nix {
        name = "mpc";
        vendorHash = "sha256-LcYXXlHPHze9zWpJ6wB5o0py/wVzYW0r2m7liJF0uWg=";
    };

    default = mpdrp;
}
