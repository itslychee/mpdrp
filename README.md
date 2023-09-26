# MPDRP, a MPD Rich Presence for Discord


<div style="display: inline-block; margin: 10px; padding: 10px;">
<span align=left>MPDRP is a <a href="https://discord.com">Discord Rich Preesence</a> that displays your MPD status<br/>via the Rich Presence accordingly.</span>
<img align=right src="https://media.discordapp.net/attachments/1148910978948407339/1156214030822805575/image.png" alt="showcase image" width=40% height=40%>
</div>

## Using
You will only need [Go](https://go.dev) and the dependencies listed in the `go.mod` file. As of right now, I do not provide
prebuilt binaries, but maybe if this project receives more recognition I will.

```bash
$ git clone https://github.com/itslychee/mpdrp && cd mpdrp
$ go build ./cmd/mpdrp # or alternatively go run ./cmd/mpdrp and use the arguments below
$ ./mpdrp -retry --retry-delay 1s --address 127.0.0.1:1234 --password "password!"
// 2021/08/22 02:33:31 ] attempting to connect to 1 address(es)
// ...
```
<br>
<br>
<br>
<br>

## Autostart

There is a batch file, a systemd service file, and a launchd file located in `config/` that you can use with your process manager.

## Copyright Notice
The image used for the front facing project icon uses assets from MPD (https://www.musicpd.org) and Discord (https://discord.com), 
<strong><u>all rights are reserved</u></strong> to these entities and any legitimate request made by either one for removal should be done 
through [email](mailto:itslychee@protonmail.com), or my discord `@itsalychee`. As such this excludes the image from the AGPL-3.0 license.

The rest of the project, including assets in `assets/`, is licensed under AGPLv3 unless explicitly stated otherwise. By contributing to this project
you also agree to license your code under the same license.