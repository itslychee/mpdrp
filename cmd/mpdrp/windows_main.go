package main

import (
	"flag"
	"fmt"

	// "log"
	"os"
	"time"

	"golang.org/x/sys/windows/svc"
)

// var logger *log.Logger
var f *os.File

type windowsHandler struct{}

func (w windowsHandler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	// timer := time.After(time.Second * 10)

	s <- svc.Status{State: svc.StartPending}
	time.Sleep(100 * time.Millisecond)
	s <- svc.Status{
		State: svc.Running,
		Accepts: svc.AcceptShutdown | svc.AcceptStop,
	}

loop:
	for cr := range r {
		f.WriteString(fmt.Sprintf("%+v", cr))
		switch cr.Cmd {
		case svc.Interrogate:
			s <- cr.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			s <- cr.CurrentStatus
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			break loop
		default:
			f.WriteString(fmt.Sprintf("unknown request command: %d ", cr.Cmd))
		}
	}
	s <- svc.Status{State: svc.StopPending}
	return
}

// func init() {
// 	logger = log.New(os.Stderr, "]", log.Lmsgprefix|log.LstdFlags)
// }

func main() {
	address := flag.String("address", "", "MPD's address, if left unset then the program will try a list of defaults")
	password := flag.String("password", "", "Password to authenticate to MPD with")
	timeout := flag.Duration("timeout", time.Duration(0), "how long the program should wait for the connection to respond before quitting, unique to TCP-based MPD addresses")
	retry := flag.Bool("retry", false, "Always reconnect to MPD and/or Discord if one of their connections get dropped")
	retryDelay := flag.Duration("retry-delay", time.Duration(time.Second*5), "Grace period between reconnections, this flag is useless without -retry being passed")

	
	file, err := os.OpenFile("C:\\Users\\Logan\\Projects\\mpdrp\\log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 000)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	f = file
	win, err := svc.IsWindowsService()
	if err != nil {
		panic(err)
	}
	// Account for no playing songs
	if win {
		go func() {
			if err = svc.Run("mpdrp", windowsHandler{}); err != nil {
				panic(err)
			}
		}()
	}

	for {
		f.WriteString(fmt.Sprintf("%s %s %s %v %s\n", *address, *password, *timeout, *retry, *retryDelay))
		time.Sleep(time.Second * 5)
	}

	
}
