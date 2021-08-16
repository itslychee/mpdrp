package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"time"
)

func init() {
	logging := log.New(os.Stderr, "] ", log.LstdFlags|log.Lmsgprefix)
}

func main() {
	address := flag.String("address", "", "MPD's address, if left unset then the program will try a list of defaults")
	password := flag.String("password", "", "Password to authenticate to MPD with")
	timeout := flag.Duration("timeout", time.Duration(0), "how long the program should wait for the connection to respond before quitting, unique to TCP-based MPD addresses")
	retry := flag.Bool("retry", false, "Always reconnect to MPD and/or Discord if one of their connections get dropped")
	retryDelay := flag.Duration("retry-delay", time.Duration(time.Second*5), "Grace period between reconnections, this flag is useless without -retry being passed")

	if runtime.GOOS == "windows" {
		scmMode := flag.Bool("scm-mode", false, "Program is installed under a Windows Service and should redirect all logs to their respective place, you shouldn't pass this flag yourself")
	}
	flag.Parse()

}
