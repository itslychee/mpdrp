# mpdrp

mpdrp is a [Discord](https://discord.com) Rich Presence for the music application, [MPD](https://musicpd.org). It supports any platform that 
Discord supports, although I don't test on OSX, there shouldn't be any issues with that platform specifically.

Additionally, I've created a small client that's mainly targeted for testing on Windows. It seems 
to work with each command with ease, but is not intended for scripting for now.

mpdrp is licensed under the Apache License 2.0, and you can find a full copy in the `LICENSE` file within this directory.

## Autostart

For right now, things are a bit wacky and I promise to create a few scripts dedicated to each platform for this. Windows isn't making 
this easy due to UAC and the execution policy of PowerShell scripts being locked off, so here we are.

### Windows

You may need administrative privileges to execute these commands

```powershell
# Add program to autostart, optionally add the desired program flags
$ reg add HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run /v MPDRP /d "C:\An\Absolute\Path\Here.exe"

# Remove it
$ reg delete HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run /f /v MPDRP
```
If you would simply like to disable mpdrp from being autostarted, but not delete the registry key, go
to `Apps > Settings` in the Windows Settings (Windows 10)

### Linux

I cannot cover all process managers in Linux, so I will be covering `systemd` here. Systemd makes it easier to autostart local programs by 
providing a utility called systemd-user, which is pretty much like `systemd` but based locally therefore not requiring you to use root every time you want to
control it.

I have provided a service file at `config/mpdrp.service`, and despite https://wiki.archlinux.org/title/Systemd/User being the Arch Wiki, it should cover
most if not all Linux distros, using `systemd` of course. If you have trouble getting this set up, carefully re-read the aforementioned page again.


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