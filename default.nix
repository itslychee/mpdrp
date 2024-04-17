{ pkgs ? import <nixpkgs>}:
rec {
    mpdrp = pkgs.callPackage ./nix/package.nix {
        name = "mpdrp";
        vendorHash = "sha256-LcYXXlHPHze9zWpJ6wB5o0py/wVzYW0r2m7liJF0uWg=";
    };
    # this is mostly for personal usage, I do not guarantee
    # stability with this.
    mpdrp-mpc = pkgs.callPackage ./nix/package.nix {
        name = "mpc";
        vendorHash = "sha256-LcYXXlHPHze9zWpJ6wB5o0py/wVzYW0r2m7liJF0uWg=";
    };

    default = mpdrp;
}
