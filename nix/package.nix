{
    buildGoModule,
    lib,
    name,
    vendorHash
}:
let
    inherit (lib.fileset)
        toSource
        gitTracked
        difference
    ;
in
buildGoModule {
    inherit name vendorHash;
    src = toSource {
        root = ../.;
        fileset = difference
            (gitTracked ../.)
            ../README.md
        ;
    };
    doCheck = false;
    subPackages = [ "cmd/${name}" ];
    meta.mainProgram = name;
}
