package main

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ItsLychee/mpdrp/discord"
	"github.com/ItsLychee/mpdrp/mpd"
)

func updateRichPresence(mpc *mpd.MPDConnection, ipc *discord.DiscordPresence) error {
	// status: Get the player's current positioning
	// currentsong: Get the metadata of the song
	err, r := mpc.Exec(mpd.Command{Name: "currentsong"}, mpd.Command{Name: "status"})
	if err != nil {
		return err
	}

	// Parse the records "elapsed" and "duration" to a native Go uint64
	duration, err := strconv.ParseUint(r.Records["duration"], 10, 64)
	if err != nil {
		panic(err)
	}
	elapsed, err := strconv.ParseUint(r.Records["elapsed"], 10, 64)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	songEnd := now.Add(time.Duration(duration - elapsed))
	f.WriteString(songEnd.String())

	switch r.Records["state"] {
	// set parameters for the following rich presence update
	}

	return nil	
}


func getDefaultAddresses() (addresses []Addr, err error) {
	// A handy function for env var defaults
	getenv := func(env string, d string) string {
		val, ok := os.LookupEnv(env)
		if !ok {
			val = d
		}
		return val
	}

	if runtime.GOOS != "windows" {
		if addr, er := resolveAddr(filepath.Join(getenv("XDG_RUNTIME_DIR", "/run"), "mpd/socket")); er != nil {
			err = er
			return
		} else {
			addresses = append(addresses, addr)
		}
	}

	var mpdHost, mpdPort string = "127.0.0.1", "6600"
	var mpdPassword string

	if val, ok := os.LookupEnv("MPD_HOST"); ok {
		if nv, err := strconv.Unquote(val); err == nil {
			val = nv
		}
		segments := strings.SplitN(val, "@", 2)
		if len(segments) == 2 {
			mpdHost = segments[0]
			mpdPassword = segments[1]
		}
	}
	if val, ok := os.LookupEnv("MPD_PORT"); ok {
		if nv, err := strconv.Unquote(val); err == nil {
			val = nv
		}
		mpdPort = val
	}

	var address Addr
	if strings.HasPrefix(mpdHost, "@/") || strings.HasPrefix(mpdHost, "/") {
		address, err = resolveAddr(mpdHost)
		if err != nil {
			return
		}

	} else {
		address, err = resolveAddr(net.JoinHostPort(mpdHost, mpdPort))
		if err != nil {
			return
		}
	}
	if mpdPassword != "" {
		address.password = mpdPassword
	}
	addresses = append(addresses, address)
	return
}

type Addr struct {
	address  net.Addr
	password string
}


func resolveAddr(address string) (addr Addr, err error) {
	switch {
	case strings.HasPrefix(address, "@/"):
		addr.address, err = net.ResolveUnixAddr("unixgram", address)
	case strings.HasPrefix(address, "/"):
		addr.address, err = net.ResolveUnixAddr("unix", address)
	default:
		addr.address, err = net.ResolveTCPAddr("tcp", address)
	}
	return

}