package main

import (
	"errors"
	"flag"
	logging "log"
	"os"
	"strconv"
	"time"

	"github.com/ItsLychee/mpdrp/discord"
	"github.com/ItsLychee/mpdrp/mpd"
)

var log *logging.Logger = logging.New(os.Stderr, "] ", logging.Lmsgprefix|logging.LstdFlags)

func main() {
	address := flag.String("address", "", "MPD's address, if left unset then the program will try a list of defaults")
	password := flag.String("password", "", "Password to authenticate to MPD with")
	timeout := flag.Duration("timeout", time.Duration(0), "how long the program should wait for the connection to respond before quitting, unique to TCP-based MPD addresses")
	retry := flag.Bool("retry", false, "Always reconnect to MPD and/or Discord if one of their connections get dropped")
	retryDelay := flag.Duration("retry-delay", time.Duration(time.Second*5), "Grace period between reconnections, this flag is useless without -retry being passed")
	clientID := flag.Int64("client-id", 856926322437521428, "Client ID that MPDRP will use, it's best not to ignorantly pass this flag")
	flag.Parse()

	// Account for no playing songs
	var addressPool []Addr
	if *address != "" {
		if addr, err := resolveAddr(*address); err != nil {
			addressPool = append(addressPool, addr)
		} else {
			panic(err)
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
			// For the sake of simplicity, I'm just going to 
			// close both of these connections (assuming if one of them aren't)
			// 
			// Perhaps in the near future I'll make this process smarter and only reconnect 
			// as needed?
			mpc.Close()
			ipc.Close()

			// Sleep so mpdrp won't consume the entirety of your CPU!
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
	os.Exit(1)
}
