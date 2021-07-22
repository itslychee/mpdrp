# mpdrp

mpdrp is a [Discord](https://discord.com) Rich Presence for the music application, [MPD](https://musicpd.org). It supports any platform that 
Discord supports, although I don't test on OSX, there shouldn't be any issues with that platform.

Additionally, I've created a small client that's mainly targeted for testing on Windows. It seems 
to work with each command with ease, but is not intended for scripting for now.

mpdrp is licensed under the Apache License 2.0, and you can find a full copy in the `LICENSE` file within this directory.

# Building

The process is rather simple
```bash
$ git clone https://github.com/ItsLychee/mpdrp
$ cd mpdrp
$ go build cmd/mpdrp/main.go
$ ./mpdrp
```

Replace `mpdrp` with `mpc` if you would like to instead build `cmd/mpc`

## Todo

- [ ] Scripts to aid in adding mpdrp to the user's process manager
- [ ] Support a reconnection logic to make better use of process managers for error handling
- [ ] More options for `cmd/mpdrp`
    - [ ] A configuration file
- [x] Honor MPD's connection defaults in an efficient manner (https://mpd.readthedocs.io/en/latest/client.html#connecting-to-mpd)


## Goals
- Maintain a nice, manageable codebase to encourage contributions
- Maybe move `discord` and make a library for all RPC functionalities?

## Screenshots 

![MPD Playing](https://raw.githubusercontent.com/ItsLychee/mpdrp/main/assets/showcase-playing.png)

![MPD Paused](https://raw.githubusercontent.com/ItsLychee/mpdrp/main/assets/showcase-paused.png)

![MPD Stopped](https://raw.githubusercontent.com/ItsLychee/mpdrp/main/assets/showcase-stopped.png)