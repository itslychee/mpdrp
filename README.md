# mpdrp

mpdrp is a [Discord](https://discord.com) Rich Presence for the music application, [MPD](https://musicpd.org). It supports any platform that 
Discord supports, although I don't test on OSX, there shouldn't be any issues with that platform specifically.

Additionally, I've created a small client that's mainly targeted for testing on Windows. It seems 
to work with each command with ease, but is not intended for scripting for now.

mpdrp is licensed under the Apache License 2.0, and you can find a full copy in the `LICENSE` file within this directory.

## Building

Usually, you should check [Releases](/releases) as they contain stable and precompiled binaries. Otherwise, you will need to install [Go](golang.org)
```bash
$ git clone https://github.com/ItsLychee/mpdrp
$ cd mpdrp

# For cmd/mpdrp
$ go build -o mpdrp cmd/mpdrp/main.go

# For cmd/mpc
$ go build -o mpc cmd/mpc/main.go
```

## Screenshots 

![MPD Playing](https://raw.githubusercontent.com/ItsLychee/mpdrp/main/assets/showcase-playing.png)

![MPD Paused](https://raw.githubusercontent.com/ItsLychee/mpdrp/main/assets/showcase-paused.png)

![MPD Stopped](https://raw.githubusercontent.com/ItsLychee/mpdrp/main/assets/showcase-stopped.png)