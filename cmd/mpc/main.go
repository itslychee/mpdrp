package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/itslychee/mpdrp/mpd"
)

func formatResponse(command mpd.Command, r *mpd.Response) {
	fmt.Println(command.String())
	for k, v := range r.Records {
		fmt.Printf("%s: %s\n", k, v)
	}
	if len(r.OK()) > 0 {
		fmt.Println(string(r.OK()))
	}
}

func main() {
	address := flag.String("address", "127.0.0.1:6600", "address to connect to, must be in HOST:PORT or /f/i/l/e format")
	password := flag.String("password", "", "password to authorize the connection with MPD")
	flag.Parse()

	var client mpd.MPDConnection

	err := client.Connect(&mpd.MPDConnInfo{
		Address: mpd.ResolveAddr(*address),
		Password: *password,
		Timeout: time.Duration(time.Second * 30),
	})
	if err != nil {
		fmt.Printf("error while trying to connect to MPD [%s]\n", *address)
		panic(err)
	}
	if *password != "" {
		cmd := mpd.Command{
			Name: "password",
			Args: []string{*password},
		}
		r, err := client.Exec(cmd)
		if err != nil {
			fmt.Println(err)
		}
		formatResponse(cmd, r)
	}
	if len(flag.Args()) == 0 {
		fmt.Println("no command supplied")
		os.Exit(-1)
	}


	cmd := mpd.Command{
		Name: flag.Arg(0),
		Args: flag.Args()[1:],
	}
	r, err := client.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	} else {
		formatResponse(cmd, r)
	}
}
