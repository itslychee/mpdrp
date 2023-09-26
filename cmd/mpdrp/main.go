package main

import (
	"errors"
	"flag"
	"fmt"
	logging "log"
	"os"
	"strconv"
	"time"

	"github.com/itslychee/mpdrp/discord"
	"github.com/itslychee/mpdrp/mpd"
)

var (
	// This variable will be reassigned by the linker
	Version = "n/a"
	log     = logging.New(os.Stderr, "] ", logging.Lmsgprefix|logging.LstdFlags)
	// Flags
	address       = flag.String("address", "", "MPD's address, if left unset then the program will try a list of defaults")
	password      = flag.String("password", "", "Password to authenticate to MPD with")
	timeout       = flag.Duration("timeout", time.Duration(0), "how long the program should wait for the connection to respond before quitting, unique to TCP-based MPD addresses")
	retry         = flag.Bool("retry", false, "Always reconnect to MPD and/or Discord if one of their connections get dropped")
	retryDelay    = flag.Duration("retry-delay", time.Duration(time.Second*5), "Grace period between reconnections, this flag is useless without -retry being passed")
	clientID      = flag.Int64("client-id", 1155715236167426089, "Client ID that MPDRP will use, it's best not to ignorantly pass this flag")
	version       = flag.Bool("version", false, "Display MPDRP's version and exit")
	verbose       = flag.Bool("verbose", false, "MPDRP will be more transparent with what it does internally, know this will produce a lot of output")
	noAlbumCovers = flag.Bool("no-album-covers", false, "Disables MPDRP from making HTTP requests to retrieve an album cover URL which is set to your Rich Presence's LargeImage field")
)

func debug(v ...interface{}) {
	if verbose != nil && *verbose {
		log.Println(v...)
	}
}

func main() {
	flag.Parse()
	debug("verbose logging enabled")

	if *version {
		fmt.Println("mpdrp version:", Version)
		return
	}

	if *timeout == 0 {
		if val, ok := os.LookupEnv("MPD_TIMEOUT"); ok {
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				panic(err)
			}
			*timeout = time.Duration(v) * time.Second
		}
	}

	// Account for no playing songs
	var addressPool []Addr
	if *address != "" {
		if addr, err := resolveAddr(*address); err != nil {
			panic(err)

		} else {
			addressPool = append(addressPool, addr)
		}
	} else {
		if addrs, err := getDefaultAddresses(); err != nil {
			panic(err)
		} else {
			addressPool = append(addressPool, addrs...)
		}
	}
	// Connection structs
	var mpc = new(mpd.MPDConnection)
	var ipc = discord.NewDiscordPresence(strconv.FormatInt(*clientID, 10))

connection_loop:
	for index := 0; index < 1 || *retry; index++ {
		if index >= 1 {
			mpc.Close()
			ipc.Close()
			time.Sleep(*retryDelay)
		}

		log.Printf("attempting to connect to %d address(es)\n", len(addressPool))
		for index, val := range addressPool {
			err := mpc.Connect(val.address.Network(), val.address.String(), *timeout)
			if err == nil {
				log.Printf("mpd: connected to %s/%s\n", val.address.Network(), val.address.String())
				break
			}
			log.Printf("mpd: could not connect to %s\n%s\n", val.address.String(), err)
			if index == len(addressPool)-1 {
				log.Printf("mpd: could not connect to a viable address")
				continue connection_loop
			}
		}

		if *password != "" {
			if err, _ := mpc.Exec(mpd.Command{Name: "password", Args: []string{*password}}); err != nil {
				log.Fatalf("mpd: incorrect password: %s\n", err)
			}
			log.Println("mpd: successfully authenticated to mpd")
		}

		if err := ipc.Connect(); err != nil {
			log.Printf("discord: %s\n", err)
			continue
		} else {
			log.Println("discord: connected to ipc pipe")
		}

		if err := ipc.CreateHandshake(); err != nil {
			log.Printf("discord: %s\n", err)
			continue
		} else {
			log.Println("discord: successfully sent handshake")
		}

		for {
			if err := updateRichPresence(mpc, ipc); err != nil {
				log.Println("discord: error while updating activity: ", err)
				continue connection_loop
			}
			for nt := time.Now(); time.Now().Unix() < nt.Add(time.Second*15).Unix(); {
				if _, err := mpc.Exec(mpd.Command{Name: "ping"}); err != nil {
					if errors.Is(err, mpd.ResponseError{}) {
						log.Fatalln(err)
					}
					continue connection_loop
				}
				time.Sleep(time.Second * 5)
			}

			if _, err := mpc.Exec(mpd.Command{Name: "idle", Args: []string{"player"}}); err != nil {
				if errors.Is(err, mpd.ResponseError{}) {
					log.Fatalln(err)
				}
				continue connection_loop
			}

		}

	}
}
