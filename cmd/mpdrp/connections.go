package main

import (
	"fmt"
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
	r, err := mpc.Exec(mpd.Command{Name: "currentsong"}, mpd.Command{Name: "status"})
	if err != nil {
		return err
	}

	if len(r.Records) == 0 {
		fmt.Println("empty")
		return nil
	}

	artistAlbum := []string{"??", "??"}
	if album := r.Records["Album"]; album != "" {
		artistAlbum[0] = album
	}
	if artist := r.Records["Artist"]; artist != "" {
		artistAlbum[1] = artist
	}

	details := "??"
	state := strings.Join(artistAlbum, " - ")
	if r.Records["Artist"] != "" {
		details = r.Records["Artist"]
	}

	var payload = discord.Activity{
		State:   &state,
		Details: &details,
		Assets: &discord.Assets{
			LargeImage: "mpd_logo",
			LargeText:  "Music Player Daemon",
		},
		Timestamps: nil,
	}

	switch r.Records["state"] {
	case "play":
		var duration, elapsed float64
		if duration, err = strconv.ParseFloat(r.Records["duration"], 64); err != nil {
			panic(err)
		}
		if elapsed, err = strconv.ParseFloat(r.Records["elapsed"], 64); err != nil {
			panic(err)
		}

		now := time.Now()
		songEnd := now.Add(time.Second * time.Duration(duration-elapsed))

		payload.Details = &details
		payload.Timestamps = &discord.Timestamps{
			End: int(songEnd.Unix()),
		}
		payload.Assets.SmallImage = "mpd_play"
		payload.Assets.SmallText = "Playing"
	case "pause":
		payload.Assets.SmallImage = "mpd_pause"
		payload.Assets.SmallText = "Paused"
	case "stop":
		details = "Stopped"
		payload.Details = &details
		payload.State = nil
		payload.Assets.SmallText = *payload.Details
		payload.Assets.SmallImage = "mpd_stop"
	}
	_, _, err = ipc.SetActivity(payload)
	return err
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
