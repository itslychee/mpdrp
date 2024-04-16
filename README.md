# MPDRP, a MPD Rich Presence for Discord


<div style="display: inline-block; margin: 10px; padding: 10px;">
<span align=left>MPDRP is a <a href="https://discord.com">Discord Rich Presence</a> that displays your MPD status<br/>via the Rich Presence accordingly.</span>
<img align=right src="https://private-user-images.githubusercontent.com/82718618/292594271-65820cf1-e1d8-4985-8d54-5819105f6ba5.png?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTEiLCJleHAiOjE3MDMyODQ0MjksIm5iZiI6MTcwMzI4NDEyOSwicGF0aCI6Ii84MjcxODYxOC8yOTI1OTQyNzEtNjU4MjBjZjEtZTFkOC00OTg1LThkNTQtNTgxOTEwNWY2YmE1LnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFJV05KWUFYNENTVkVINTNBJTJGMjAyMzEyMjIlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjMxMjIyVDIyMjg0OVomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPTVhMTVkNzhmYTQ2MDNjMWU3ZWVkYzg0ZTg3OTI3MDg2NzYyMDZmMzg5MzFlZDAyNjk4MjIxMzJlMWJmMzU2MzMmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JmFjdG9yX2lkPTAma2V5X2lkPTAmcmVwb19pZD0wIn0.8a4CT68CTjFeNWkLUMKtkQUQlMUt-zJ7i0TM5qrhKpI" alt="showcase image" width=40% height=40%>
</div>


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
