package main

import (
	"encoding/json"
	"fmt"
	logging "log"
	"os"
)

var (
	// Default Discord Client ID
	ClientID = 1155715236167426089
	// Rich Presence asset keys
	NoAlbumAsset = "no_album"
	PlayAsset    = "mpd_play"
	PauseAsset   = "mpd_pause"
	StopAsset    = "mpd_stop"
	logger       = logging.New(os.Stderr, "", logging.Ltime)
)

type LogLevel uint8

const (
	Normal LogLevel = iota
	Debug
	Network
)

func log(level LogLevel, message string) {
	switch {
	case level == Normal:
		fmt.Println(message)
	case level == Debug && *debugLevel >= int(Debug):
		logger.Println("[  DEBUG]:", message)
	case level == Network && *debugLevel >= int(Network):
		logger.Println("[NETWORK] >", message)
	}
}

func logf(level LogLevel, message string, args ...any) {
	log(level, fmt.Sprintf(message, args...))
}

func logjson(level LogLevel, message string, jsond json.Marshaler) {
	d, err := json.MarshalIndent(jsond, "", " ")
	if err != nil {
		logf(level, "JSON encoding error: %s", err.Error())
	}
	log(level, fmt.Sprintf("%s:\n%s", message, d))

}
