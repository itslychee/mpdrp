//go:build darwin || linux
// +build darwin linux

package discord

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
)

type DiscordPresence struct {
	ClientID string
	conn     net.Conn
}

var envKeys = [4]string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

func (c *DiscordPresence) Connect() error {
	var baseDirectory string
	for _, value := range envKeys {
		val, present := os.LookupEnv(value)
		if !present || val == "" {
			continue
		}
		baseDirectory = val
		break
	}
	if baseDirectory == "" {
		return ErrCannotConnect
	}
	for index := 0; index <= 9; index++ {
		pipePath := filepath.Join(baseDirectory, fmt.Sprintf("discord-ipc-%d", index))
		conn, err := net.Dial("unix", pipePath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}

		if err != nil {
			return err
		}

		c.conn = conn
		return nil

	}
	return ErrCannotConnect
}
