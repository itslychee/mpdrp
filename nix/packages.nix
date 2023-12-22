{ pkgs, lib, version ? "dev"}:
let
  # longest function name in existence!
  wrapGo = pkgs.buildGoApplication;
  pkg = lib.makeOverridable ({ withMpc }: {
    inner = {
      inherit version;
      pname = "mpdrp";
      doCheck = false; 
      modules = ./gomod2nix.toml;
      src = with lib.fileset; toSource {
          root = ../.;
          fileset = difference ../. (unions [
            ../config
            ../assets
            ../release.sh
          ]);
        };
      subPackages = lib.flatten [ 
        "cmd/mpdrp"
        (lib.optionals withMpc [ "cmd/mpc" ])
      ];
      meta.mainProgram = "mpdrp";
    };
  }) { withMpc = false; };
in {
  mpdrp = wrapGo (pkg.inner // {
    passthru.withMpc = wrapGo (pkg.override { withMpc = true; }).inner;
  });
}
