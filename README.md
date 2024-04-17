# mpdrp, a MPD Rich Presence for Discord

mpdrp is a MPD (Music Player Daemon) Rich Presence for Discord!

## features

- Fetches album covers from [Cover Art Archive](https://coverartarchive.org) utilizing Discord's dynamic assets feature
- Automatically clears your Rich Presence after 5 minutes of MPD being paused.
- Handles connection failures with Discord and MPD alike.

## Usage

### Nix

Nix has first class support, prebuilt closures can be substituted from [`cache.garnix.io`](https://garnix.io), details
on how to utilize Garnix's binary cache are provided [here](https://garnix.io/docs/caching).

See the available outputs under `packages` via `nix flake show github:itslychee/mpdrp`, or just
look at [`default.nix`](./default.nix).

### Other

You will need Go installed to compile the program

* `go install ./cmd/mpdrp ./cmd/mpc` (omit any unwanted sub packages)

And place it somewhere suitable to be added and/or used in `$PATH`

After that just run the binary and its defaults should be good enough for the average user, otherwise
pass `--help` to see available options.

