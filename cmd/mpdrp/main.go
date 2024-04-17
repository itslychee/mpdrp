package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/itslychee/mpdrp/discord"
	music "github.com/itslychee/mpdrp/mpd"
)

var (
	debugLevel    = flag.Int("debug", 1, "debug level for program, zero disables debugging output")
	clientID      = flag.Int("client-id", ClientID, "rich presence client id, this normally should not be changed")
	noAlbumCovers = flag.Bool("no-album-covers", false, "Do not set album covers and retrieve them via Cover Art Archive")
	reconnect     = flag.Duration("reconnect", time.Duration(time.Second*5), "grace period before reattempting to reconnect to MPD & Discord, must be above 5 seconds or zero to disable this")
)

func init() {
	flag.Parse()
}

func formatActivity(activity *discord.Activity, currentStatus *music.Response, cachedURLs map[string]string) {
	if currentStatus.Get("state") != "stop" {
		// Activity metadata
		album, ok := currentStatus.Records["Album"]
		if !ok {
			album = []string{"??"}
		}
		artist, ok := currentStatus.Records["Artist"]
		if !ok {
			artist = []string{"??"}
		}
		details, ok := currentStatus.Records["Title"]
		if !ok {
			details = []string{"??"}
		}

		state := fmt.Sprintf("%s by %s", album[0], artist[0])
		activity.Details = &details[0]
		activity.State = &state

		var err error

		// Reuse last image if available
		v, ok := cachedURLs[currentStatus.Get("songid")]
		if !ok {
			v, err = GetCoverArt(*currentStatus)
			if err != nil {
				logf(Normal, "[error] encountered error while fetching cover art: %s", err)
			} else {
				cachedURLs[currentStatus.Get("songid")] = v
			}
		}
		activity.Assets.LargeImage = v
		activity.Assets.SmallImage = PauseAsset
		activity.Assets.SmallText = "Paused"
		if currentStatus.Get("state") == "play" {
			activity.Assets.SmallImage = PlayAsset
			activity.Assets.SmallText = "Playing"
			var duration, elapsed float64
			if duration, err = strconv.ParseFloat(currentStatus.Get("duration"), 64); err != nil {
				panic(err)
			}
			if elapsed, err = strconv.ParseFloat(currentStatus.Get("elapsed"), 64); err != nil {
				panic(err)
			}
			activity.Timestamps = &discord.Timestamps{
				End: time.Now().Add(time.Second * time.Duration(duration-elapsed)).Unix(),
			}
		}
	} else {
		activity = nil
	}

}

// retryOrScream sucks but I intend to package mpd and discord
// into their own full fledge repos some day :tm:
//
// function should not be called if it is out of reconnection context
// and the program can recover from this error by reconnecting
func retryOrScream(err error) {
	handledErrors := []syscall.Errno{
		syscall.ECONNREFUSED,
		syscall.ECONNABORTED,
		syscall.ECONNRESET,
		syscall.EPIPE,
	}

	// Get syscall errno hex
	if v, ok := err.(*net.OpError); ok {
		if err, ok := v.Err.(*os.SyscallError); ok {
			errno := err.Err.(syscall.Errno)
			src := "https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/syscall/zerrors_linux_amd64.go;l=1183"
			logf(Normal, "panic: ERRNO (0x%x)\n%s", int(errno), src)
		}
	}

	handledErrs := []error{
		discord.ErrCannotConnect,
		io.EOF,
	}


        // handle generic errors
	for _, v := range handledErrs {
		if errors.Is(err, v) {
			log(Network, "reattempting connection")
			time.Sleep(5 * time.Second)
			return
		}
	}

        // handle specific socket errors
	for _, e := range handledErrors {
		if errors.Is(err, e) {
			log(Network, "reattempting connection")
			time.Sleep(5 * time.Second)
			return
		}
	}

	panic(err)
}

type Signal uint8

const (
	Play = iota
	Pause
	Stopped
	Retry
)

func main() {
conn:
	client := discord.NewDiscordPresence(strconv.Itoa(*clientID))
	mpd := new(music.MPDConnection)
	log(Normal, "connecting to discord")
	if err := client.Connect(); err != nil {
		logf(Normal, "failed: %s", err)
		retryOrScream(err)
		goto conn
	}
	defer client.Close()
	logf(Normal, "connected to discord: %s", client.Conn.RemoteAddr())
	log(Network, "starting discord handshake")
	b, err := client.CreateHandshake()
	logjson(Network, "handshake result", json.RawMessage(b))
	if err != nil {
		logf(Normal, "failed %s", err)
		retryOrScream(err)
		goto conn
	}
	log(Normal, "connecting to mpd")
	if err := mpd.Connect(nil); err != nil {
		logf(Normal, "failed %s", err)
		retryOrScream(err)
		goto conn
	}
	defer mpd.Close()
	logf(Normal, "connected to mpd instance: %s", mpd.RawConn.RemoteAddr())

	var cachedURLs = make(map[string]string)

	var lastSignal Signal
	ch := make(chan Signal, 1)

	go func() {
		for {
			// Get current status of song
			currentStatus, err := mpd.Exec(
				music.Command{Name: "currentsong"},
				music.Command{Name: "status"},
			)
			if err != nil {
				retryOrScream(err)
				ch <- Retry
				return
			}

			var activity *discord.Activity

			if currentStatus.Records["state"][0] != "stop" {
				activity = &discord.Activity{
					Assets: &discord.Assets{LargeText: "Music Player Daemon"},
				}
			}

			formatActivity(activity, currentStatus, cachedURLs)

			_, body, err := client.SetActivity(activity)
			logjson(Network, "set activity", json.RawMessage(body))
			if err != nil {
				retryOrScream(err)
				ch <- Retry
				return
			}

			switch currentStatus.Records["state"][0] {
			case "play":
				ch <- Play
			case "pause":
				ch <- Pause
			case "stop":
				ch <- Stopped
			}
			// Idle and wait
			_, err = mpd.Exec(music.Command{Name: "idle", Args: []string{"player"}})
			if err != nil {
				retryOrScream(err)
				ch <- Retry
				return
			}
		}

	}()
	for {

		select {
		case <-time.After(time.Minute * 5):

			if lastSignal != Pause {
				continue
			}

			_, body, err := client.SetActivity(nil)
			logjson(Network, "paused for too long, timing out", json.RawMessage(body))
			if err != nil {
				retryOrScream(err)

			}

		case sig := <-ch:
			lastSignal = sig
			if sig == Retry {
				goto conn
			}

		}
	}
}
