{ 
    config, 
    pkgs,
    options,
    ...
}:
with pkgs.lib;
with pkgs.lib.types;
let
    cfg = config.programs.mpdrp;

in {
    options = {
        programs.mpdrp = {
            enable = mkEnableOption "mpdrp";
            package = mkPackageOption "mpdrp";
            settings = { 
                    reconnect = mkOption {
                      description = "Grace period between reconnections, in seconds";
                      type = int;
                      default = 5;
                    };
                    verbose = mkEnableOption "verbose";
                    password = mkOption { 
                        type = nullOr str; 
                        description = "Password to MPD server";
                        default = null;
                    };
                    albumCovers = mkOption {
                        type = bool;
                        description = ''
                            Whether to fetch from Cover Art Archive via MusicBrainz 
                            to utilize the album cover for the Rich Presence's LargeImage field.
                        '';
                        default = true;
                    };
                    clientID = mkOption {
                        type = nullOr int;
                        description = ''
                            Client ID for mpdrp to use, normally you shouldn't set this.

                            By default the program will use it's preconfigured ID 
                        '';
                        default = null;
                    };
                    timeout = mkOption {
                        type = int;
                        description = "TCP connection timeout, this isn't used for other connection protocols";
                        default = 30;
                    };
                    address = mkOption {
                        type = nullOr str;
                        description = "Address to use for connection, if unset mpdrp will choose from a list of defaults";
                        default = null;
                    };
                    withMpc = mkOption {
                        type = bool;
                        description = "Include cmd/mpc to path";
                        default = false;
                    };
                };
            };
    };

    config.home.packages = (mkIf cfg.settings.withMpc (with pkgs; [ pkgs.mpdrp-mpc ]));
    config.systemd.user.services.mpdrp = mkIf (cfg.enable) {
        Unit.Description = "A discord rich presence for MPD";
        Service = let 
           opts = [
              "--reconnect ${toString cfg.settings.reconnect}s"
              (if cfg.settings.password != null then "--password ${cfg.settings.password}" else "")
              (if (!cfg.settings.albumCovers) then "-no-album-covers" else "")
              (if (cfg.settings.clientID != null) then "--client-id ${toString cfg.settings.clientID}" else "")
              (if (cfg.settings.verbose != false) then "--verbose" else "")
              (if (cfg.settings.address != null) then "--address ${cfg.settings.address}" else "")
           ];
        in {
            Type = "exec";
            ExecStart = "${pkgs.mpdrp}/bin/mpdrp " + (concatStringsSep " " opts);
        };
    };
}
