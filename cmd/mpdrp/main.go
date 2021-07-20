package main

import (
	"encoding/json"
	"flag"
	"path/filepath"
	"runtime"

	logging "log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ItsLychee/mpdrp/discord"
	"github.com/ItsLychee/mpdrp/mpd"
)

var log *logging.Logger

func init() {
	log = logging.New(os.Stderr, "] ", logging.LstdFlags|logging.Lmsgprefix)
}

func main() {
	address := flag.String("address", "", "address to connect to MPD with, if -address is not provided then mpdrp will try to connect by a predefined list of defaults")
	password := flag.String("password", "", "password to authorize with in order to use MPD")
	timeoutSeconds := flag.Uint("timeout", 0, "the connection timeout value used by mpdrp when connecting to mpd")
	forcePassword := flag.Bool("forcepassword", false, "use the provided -password if present, even if there's a password set in MPD_HOST")

	var mpdAddresses []net.Addr
	var mpc = new(mpd.MPDConnection)
	var discord = ipc.NewDiscordPresence("856926322437521428")

	var timeout *time.Duration
	if mpd_timeout := os.Getenv("MPD_TIMEOUT"); mpd_timeout != "" {
		val, err := strconv.ParseUint(mpd_timeout, 10, 64)
		if err != nil {
			log.Printf("could not parse environment variable MPD_TIMEOUT value: %s\n", mpd_timeout)
			panic(err)
		}
		*timeout = time.Duration(val * uint64(time.Second))
	}
	if *timeoutSeconds > 1 {
		*timeout = time.Duration(*timeoutSeconds * uint(time.Second))
	}

	// If *address is not an empty string, then we will only add *address to
	// mpdAddresses, otherwise add a minimal set of defaults
	if *address == "" {
		var formattedAddr string
		var mpd_host = "127.0.0.1"

		getenv := func(env string, d string) string {
			val, ok := os.LookupEnv(env)
			if !ok {
				val = d
			}
			return val
		}

		if runtime.GOOS != "windows" {
			// XDG_RUNTIME_DIR envvar takes precedence over /run, for obvious reasons
			socketPath := filepath.Join(getenv("XDG_RUNTIME_DIR", "/run"), "mpd/socket")
			if r, err := resolveAddr(socketPath); err != nil {
				log.Printf("could not resolve socket address: %s\n", socketPath)
				panic(err)

			} else {
				mpdAddresses = append(mpdAddresses, r)
			}
		}

		if val, ok := os.LookupEnv("MPD_HOST"); ok {
			hostSlice := strings.SplitN(val, "@", 2)
			mpd_host = hostSlice[len(hostSlice)-1]
			if q, err := strconv.Unquote(mpd_host); err == nil {
				mpd_host = q
			}
			if len(hostSlice) == 2 && !(*forcePassword) {
				*password = hostSlice[0]
			}
		}

		if strings.HasPrefix(mpd_host, "@/") || strings.HasPrefix(mpd_host, "/") {
			// We can assume that this is a unix (abstract)? socket
			formattedAddr = mpd_host
		} else {
			formattedAddr = net.JoinHostPort(mpd_host, getenv("MPD_PORT", "6600"))
		}

		if r, err := resolveAddr(formattedAddr); err != nil {
			log.Printf("could not resolve address: %s\n", formattedAddr)
			panic(err)
		} else {
			mpdAddresses = append(mpdAddresses, r)
		}

	} else {
		val, err := resolveAddr(*address)
		if err != nil {
			log.Printf("could not resolve address: %s\n", *address)
			panic(err)
		}
		mpdAddresses = append(mpdAddresses, val)
	}

	log.Printf("Attempting to establish a connection to MPD with %d address(es)\n", len(mpdAddresses))
	// Connect to MPD
	for index, val := range mpdAddresses {
		err := mpc.Connect(val.Network(), val.String(), timeout)
		if err == nil {
			break
		} else {
			log.Println("unable to connect to mpd address: ", val.String())
			log.Println(err)
		}
		if index == len(mpdAddresses) {
			log.Fatalln("mpdrp cannot find a suitable address to connect to MPD with")
		}
	}
	defer mpc.Close()

	if *password != "" {
		cmd := mpd.Command{
			Name: "password",
			Args: []string{
				*password,
			},
		}
		if _, err := mpc.Exec(cmd); err != nil {
			log.Fatal("password authentication failed: ", err)
		}
	}

	if err := discord.Connect(); err != nil {
		log.Println("error while trying to connect to Discord")
		panic(err)
	}
	defer discord.Disconnect()

	if err := discord.CreateHandshake(); err != nil {
		log.Println("discord handshake failed")
		panic(err)
	}
	updateDelay := time.Second * 15

	updateRichPresence(discord, mpc)

	for {
		// This is to ensure that Discord won't ratelimit us while ensuring
		// MPD won't timeout on us as well
		for nt := time.Now(); time.Now().Unix() < nt.Add(updateDelay).Unix(); {
			_, err := mpc.Exec(mpd.Command{Name: "ping"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second * 5)
		}
		// Thankfully, idle disables timeouts during its execution
		mpc.Exec(mpd.Command{Name: "idle", Args: []string{"player"}})
		updateRichPresence(discord, mpc)

	}

}

func resolveAddr(address string) (net.Addr, error) {
	var addr net.Addr
	var err error
	switch {
	case strings.HasPrefix(address, "@/"):
		addr, err = net.ResolveUnixAddr("unixgram", address)
	case strings.HasPrefix(address, "/"):
		addr, err = net.ResolveUnixAddr("unix", address)
	default:
		addr, err = net.ResolveTCPAddr("tcp", address)
	}
	return addr, err

}

func updateRichPresence(ipcSocket *ipc.DiscordPresence, mpc *mpd.MPDConnection) {
	r, err := mpc.Exec(mpd.Command{Name: "status"}, mpd.Command{Name: "currentsong"})
	if err != nil {
		log.Panic(err)
	}

	var albumArtist = []string{r.Records["Album"], r.Records["Artist"]}
	var state = strings.Join(albumArtist, " by ")
	var details = r.Records["Title"]

	var payload = ipc.Activity{
		State:      &state,
		Details:    &details,
		Timestamps: nil,
		Assets:     &ipc.Assets{LargeImage: "mpd_logo", LargeText: "MPD"},
	}

	switch r.Records["state"] {
	case "play":
		elapsed, err := strconv.ParseFloat(r.Records["elapsed"], 64)
		duration, err1 := strconv.ParseFloat(r.Records["duration"], 64)

		if err != nil {
			panic(err)
		}
		if err1 != nil {
			panic(err)
		}

		now := time.Now()
		secondsLeft := now.Add(time.Second * time.Duration(duration-elapsed))

		payload.Timestamps = &ipc.Timestamps{End: int(secondsLeft.Unix())}
		payload.Assets.SmallImage = "mpd_play"
		payload.Assets.SmallText = "Playing"

	case "pause":
		payload.Assets.SmallImage = "mpd_pause"
		payload.Assets.SmallText = "Paused"

	case "stop":
		details = "Stopped"

		payload.Details = &details
		payload.State = nil
		payload.Assets.SmallText = details
		payload.Assets.SmallImage = "mpd_stop"
	}

	if _, _, err = ipcSocket.SetActivity(payload); err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}

	log.Printf("Sent presence update payload\n%s", data)
}
