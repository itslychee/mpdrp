# MPDRP, a MPD Rich Presence for Discord

<img src="https://cdn.discordapp.com/app-icons/1155715236167426089/db147772ae4b4e494cb9ff61a0e2e9f1.png?size=256" alt="mpdrp logo" height=200 width=200>

MPDRP is a [Discord Rich Presence](https://discord.com) that displays your 
MPD status via the Rich Presence accordingly.

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

# Autostart
It's ideal that the program shouldn't be manually started by the user, so I've attempted to create decent scripts. For Linux and OSX users, you
will have to change it up to meet your needs.

## Windows
Open up the folder and execute `mpdrp.bat` while following its instructions.

## OSX
- Open up the folder
- Move `mpdrp.plist` to `/Library/LaunchAgents/`
- Open up your Terminal
- Enter `launchctl load /Library/LaunchAgents/mpdrp.plist`
- Finally, enter `launchctl start com.itslychee.mpdrp` and it should be running

As I don't have a Mac, I cannot assure you that the program or the plist file provided will work as expected. So, if something doesn't look right 
to you, please make a PR or an issue. I would highly be appreciative!

## Linux
For systemd users, you will need to copy `mpdrp.service` in the extracted tarball folder. I recommend using `systemd-user` for MPDRP as it is locally based
and as such, it doesn't require root. For other process managers, you're expected to know how to create one yourself, or feel free to create a PR that improves/adds on to
the current release configuration, I'd appreciate your contribution greatly.

## Screenshots
Find them in the `/assets` directory of the repository, their filenames will start with `showcase`. I mainly done this to make the
page faster to load.

## Copyright Notice
The image used for the front facing project icon uses assets from MPD (https://www.musicpd.org) and Discord (https://discord.com), 
<strong><u>all rights are reserved</u></strong> to these entities and any legitimate request made by either one for removal should be done 
through [email](mailto:itslychee@protonmail.com), or my discord `@itsalychee`. As such this excludes the image from the AGPL-3.0 license.

The rest of the project, including assets in `assets/`, is licensed under AGPLv3 unless explicitly stated otherwise. By contributing to this project
you also agree to license your code under the same license.