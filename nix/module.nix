self: {
  config,
  pkgs,
  lib,
  options,
  ...
}: let
  cfg = config.programs.mpdrp;
  inherit (self.packages.${pkgs.system}) mpdrp;
  inherit (lib) mkOption mkEnableOption mkPackageOption types;
in {
  options = {
    programs.mpdrp = {
      enable = mkEnableOption "mpdrp";
      package = mkOption {
        description = "Package to use for MDPRP";
        type = types.package;
        default =
          if cfg.withMpc
          then mpdrp.withMpc
          else mpdrp;
      };
      withMpc = mkOption {
        description = "include mpc";
        default = false;
        type = types.bool;
      };
      settings = {
        reconnect = mkOption {
          description = "Grace period between reconnections, in seconds";
          type = types.int;
          default = 5;
        };
        verbose = mkEnableOption "verbose";
        password = mkOption {
          type = types.nullOr types.str;
          description = "Password to MPD server";
          default = null;
        };
        albumCovers = mkOption {
          type = types.bool;
          description = ''
            Whether to fetch from Cover Art Archive via MusicBrainz
            to utilize the album cover for the Rich Presence's LargeImage field.
          '';
          default = true;
        };
        clientID = mkOption {
          type = types.nullOr types.int;
          description = ''
            Client ID for mpdrp to use, normally you shouldn't set this.

            By default the program will use it's preconfigured ID
          '';
          default = null;
        };
        timeout = mkOption {
          type = types.int;
          description = "TCP connection timeout, this isn't used for other connection protocols";
          default = 30;
        };
        address = mkOption {
          type = types.nullOr types.str;
          description = "Address to use for connection, if unset mpdrp will choose from a list of defaults";
          default = null;
        };
      };
    };
  };
  config = {
    home.packages = [cfg.package];
    systemd.user.services.mpdrp = lib.mkIf cfg.enable {
      Unit.Description = "A discord rich presence for MPD";
      Unit.After = ["mpd.socket"];
      Install.WantedBy = ["default.target"];
      Service = let
        opts = [
          "--reconnect ${toString cfg.settings.reconnect}s"
          (
            if cfg.settings.password != null
            then "--password ${cfg.settings.password}"
            else ""
          )
          (
            if (!cfg.settings.albumCovers)
            then "-no-album-covers"
            else ""
          )
          (
            if (cfg.settings.clientID != null)
            then "--client-id ${toString cfg.settings.clientID}"
            else ""
          )
          (
            if cfg.settings.verbose
            then "--verbose"
            else ""
          )
          (
            if (cfg.settings.address != null)
            then "--address ${cfg.settings.address}"
            else ""
          )
        ];
      in {
        Type = "exec";
        Restart = "on-success";
        ExecStart =
          (lib.getExe cfg.package)
          + " "
          + (lib.concatStringsSep " " opts);
      };
    };
  };
}
