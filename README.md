# MPDRP, a MPD Rich Presence for Discord


<div style="display: inline-block; margin: 10px; padding: 10px;">
<span align=left>MPDRP is a <a href="https://discord.com">Discord Rich Preesence</a> that displays your MPD status<br/>via the Rich Presence accordingly.</span>
<img align=right src="https://media.discordapp.net/attachments/1148910978948407339/1156214030822805575/image.png" alt="showcase image" width=40% height=40%>
</div>

## Usage

### Nix

* `nix build .#mpdrp` for the base `mpdrp` package
* `nix build .#mpdrp.withMpc` for the base `mpdrp` package that also includes `cmd/mpc`
* `nix run .#mpdrp` runs mpdrp 

The home manager module can be found at `homeManagerModules.default`, I currently do not
have a NixOS module as I use mpd on a user level, but other than that, there's nothing
preventing support for one so feel free to PR.

### Other (including Windows)

You will need Go installed to compile the program

* `go install ./cmd/mpdrp ./cmd/mpc` (omit any unwanted sub packages)

And place it somewhere suitable to be added and/or used in `$PATH`
