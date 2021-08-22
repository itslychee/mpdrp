package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ItsLychee/mpdrp/discord"
	"github.com/ItsLychee/mpdrp/mpd"
	"golang.org/x/sys/windows/svc"
)

type windowsHandler struct {
	ipc chan svc.Status
}

func (w windowsHandler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	// s <- svc.Status{State: svc.StartPending}
	s <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptShutdown | svc.AcceptStop,
	}

loop:
	for {
		select {
		case cr := <-r:
			switch cr.Cmd {
			case svc.Interrogate:
				s <- cr.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				s <- cr.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending}
				break loop
			}
		case signal := <-w.ipc:
			s <- signal
		}
	}
	s <- svc.Status{State: svc.StopPending}
	return
}

func main() {
	address := flag.String("address", "", "MPD's address, if left unset then the program will try a list of defaults")
	password := flag.String("password", "", "Password to authenticate to MPD with")
	timeout := flag.Duration("timeout", time.Duration(0), "how long the program should wait for the connection to respond before quitting, unique to TCP-based MPD addresses")
	retry := flag.Bool("retry", false, "Always reconnect to MPD and/or Discord if one of their connections get dropped")
	retryDelay := flag.Duration("retry-delay", time.Duration(time.Second*5), "Grace period between reconnections, this flag is useless without -retry being passed")
	clientID := flag.Int64("client-id", 856926322437521428, "Client ID that MPDRP will use, it's best not to ignorantly pass this flag")
	flag.Parse()

	isManaged, err := svc.IsWindowsService()
	if err != nil {
		panic(err)
	}
	// This enables use of Windows Services, which most users
	// will most likely use
	if isManaged {
		win := windowsHandler{}
		go func() {
			if err = svc.Run("mpdrp", win); err != nil {
				panic(err)
			}
		}()
		defer func() {
			win.ipc <- svc.Status{State: svc.StopPending}
			// Let Windows' SCM know that we're exiting, just in case
		}()
	}
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
	for index := 0 ;; index++{
		if index > 0 && !*retry {
			os.Exit(-1)
		}

		log.Printf("attempting to connect to %d address(es)\n", len(addressPool))
		for index, val := range addressPool {
			err = mpc.Connect(val.address.Network(), val.address.String(), *timeout) 
			if err == nil {
				log.Printf("mpd: connected to %s/%s\n", val.address.Network(), val.address.String())
				break
			}
			log.Printf("mpd: could not connect to %s\n%s\n", val.address.String(), err)
			if index == len(addressPool) - 1 {
				log.Printf("mpd: could not connect to a viable address\n%s\n", err)
				time.Sleep(*retryDelay)
				continue
			}
		}

		if *password != "" {
			if err, _ := mpc.Exec(mpd.Command{Name: "password", Args: []string{*password}}); err != nil {
				log.Fatalf("mpd: incorrect password\n%s\n", err)
			}
			log.Println("mpd: successfully authenticated to mpd")
		}

		
		if err := ipc.Connect(); err != nil {
			log.Printf("discord: %s\n", err)
			time.Sleep(*retryDelay)
			continue
		} else {
			log.Println("discord: connected to ipc pipe")
		}

		if err := ipc.CreateHandshake(); err != nil {
			log.Printf("discord: %s\n", err)
			time.Sleep(*retryDelay)
			continue
		} else {
			log.Println("discord: successfully sent handshake")
		}		

		for {
			if err := updateRichPresence(mpc, ipc); err != nil {
				log.Printf("discord: error while updating activity:\n%s\n", err)
			}
			for nt := time.Now(); time.Now().Unix() < nt.Add(time.Second*15).Unix(); {
				if err, _ := mpc.Exec(mpd.Command{Name: "ping"}); err != nil {
					panic(err)
				}
				time.Sleep(time.Second * 5)
			}
			err, _ := mpc.Exec(mpd.Command{Name: "idle", Args: []string{"player"}})
			if err != nil {
				log.Fatalf("mpd: error while idling:\n%s\n", err)
			}

			
		}

	}


}
