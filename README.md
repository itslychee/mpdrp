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

## Screenshots

Find them in the `/assets` directory of the repository, their filenames will start with `showcase`. I mainly done this to make the
page faster to load.