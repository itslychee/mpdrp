package main

import (
	"encoding/json"
	"flag"
	logging "log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ItsLychee/mpdrp/discord"
	"github.com/ItsLychee/mpdrp/mpd"
)

var log *logging.Logger

func init() {
	log = logging.New(os.Stderr, "]  ", logging.LstdFlags|logging.Lmsgprefix)
}

func main() {
	address := flag.String("address", "127.0.0.1:6600", "unix or tcp address to use to connect with MPD")
	password := flag.String("password", "", "password to use when authorizing with the server")
	flag.Parse()

	var mpdClient = new(mpd.MPDConnection)
	var discordSocket = ipc.NewDiscordPresence("856926322437521428")

	var networkType string
	switch {
	case strings.HasPrefix(*address, "/"):
		// Unix socket
		networkType = "unix"
	case strings.HasPrefix(*address, "@"):
		// Abstract unix socket
		networkType = "unixgram"
	default:
		// TCP, guessingly
		networkType = "tcp"
	}

	err := mpdClient.Connect(networkType, *address, nil)
	if err != nil {
		log.Fatalf("unable to connect to MPD server:\n%s", err.Error())
	}
	log.Printf("established connection to MPD server: %s\n", *address)
	if *password != "" {
		r, err := mpdClient.Exec(mpd.Command{Name: "password", Args: []string{*password}})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(r.Records)
	}

	log.Println("connecting to discord's IPC socket")
	if err := discordSocket.Connect(); err != nil {
		log.Fatalf("unable to connect to Discord's IPC socket: %s", err.Error())
	}
	log.Println("established connection to discord's ipc socket")

	log.Println("sending handshake to discord")
	if err := discordSocket.CreateHandshake(); err != nil {
		log.Fatalf("error while sending a handshake to Discord: %s", err)
	}
	log.Println("handshake established")

	updateRichPresence(discordSocket, mpdClient)
	updateDelay := time.Second * 15

	defer mpdClient.Close()
	defer discordSocket.Disconnect()

	for {
		// This is to ensure that Discord won't ratelimit us while ensuring
		// MPD won't timeout on us as well
		for nt := time.Now(); time.Now().Unix() < nt.Add(updateDelay).Unix(); {
			_, err := mpdClient.Exec(mpd.Command{Name: "ping"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second * 5)
		}
		// Thankfully, idle disables timeouts during its execution
		mpdClient.Exec(mpd.Command{Name: "idle", Args: []string{"player"}})
		updateRichPresence(discordSocket, mpdClient)

	}

}

func updateRichPresence(discordSocket *ipc.DiscordPresence, mpdClient *mpd.MPDConnection) {
	r, err := mpdClient.Exec(mpd.Command{Name: "status"}, mpd.Command{Name: "currentsong"})
	if err != nil {
		log.Panic(err)
	}

	var albumArtist = []string{r.Records["Album"], r.Records["Artist"]}
	var payload ipc.Activity = ipc.Activity{
		Assets: &ipc.Assets{
			LargeImage: "mpd_logo",
			LargeText:  "MPD",
		},
	}
	var state = strings.Join(albumArtist, " by ")
	var details = r.Records["Title"]

	switch r.Records["state"] {
	case "play":
		t, err := strconv.ParseInt(r.Records["Time"], 10, 64)
		if err != nil {
			panic(err)
		}
		songEnd := time.Now().Add(time.Second * time.Duration(t))

		// Payload
		payload.Details = &details
		payload.State = &state
		payload.Timestamps = &ipc.Timestamps{End: int(songEnd.Unix())}
		payload.Assets.SmallImage = "mpd_play"
		payload.Assets.SmallText = "Playing"

	case "pause":

		// Payload
		payload.Details = &details
		payload.State = &state
		payload.Timestamps = nil
		payload.Assets.SmallImage = "mpd_pause"
		payload.Assets.SmallText = "Paused"

	case "stop":
		details = "Stopped"
		// Payload
		payload.Details = &details
		payload.State = nil
		payload.Timestamps = nil
		payload.Assets.SmallText = "Stopped"
		payload.Assets.SmallImage = "mpd_stop"
	}

	if _, _, err = discordSocket.SetActivity(payload); err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}

	log.Printf("Sent presence update payload\n%s", data)
}
