# DEPRECATION NOTICE!

I have stopped using Discord so the default client ID that was provided in mpdrp's code will be 
rendered invalid. However this can still be remedied without code modifications so you can
simply pass `-client-id xxxxx...` and use your own rich presence application, but creating an
application is not enough as you must upload assets, naming is provided below.

- `assets/MPDPlay.png` must be named `mpd_play`
- `assets/MPDPause.png` must be named `mpd_pause`
- `assets/MPDStop.png` may be uploaded and named `mpd_stop`, but this is optional

After you do all of that, mpdrp should work as intended, granted that Discord hasn't changed the API up at the
time of writing this. I may start this project back up, perhaps if a few variables such as interest develops up on
both my and your side. 

This was actually my first Go project, and I really enjoyed working on this and I hope 
that I'll give it the love it deserves in the near future, but for right now I am just not in the
current position to think about this at all, and I wish I was.


# MPDRP

MPDRP is a [Discord Rich Presence](https://discord.com) that accordingly displays your 
MPD status via the Rich Presence.

# Using

I have made it to the best of my ability (or to the extent of my productivity) to make MPDRP as user friendly as possible. Grab a release
from the [Releases page](/releases) and follow the instructions that are stated there. On a bit of an unrelated note, I haven't cared to tag commits during
the under-development era of the project since it was not only in a fluctuating state of change but there wasn't any reason to.

While i test on Linux and Windows, there are no OSX-specific features or
issues that I am aware of, so building/using MPDRP on OSX ideally shouldn't be a problem.

## Building

You will only need [Go](https://go.dev) and the dependencies listed in the `go.mod` file. 

```bash
$ go build ./cmd/mpdrp
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
the current release configuration.

## Screenshots

Find them in the `/assets` directory of the repository, their filenames will start with `showcase`. I mainly done this to make the
page faster to load.