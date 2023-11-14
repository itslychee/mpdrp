//go:build windows
// +build windows

package discord

import (
	"errors"
	"fmt"
	"io/fs"
	"net"

	"github.com/Microsoft/go-winio"
)

type DiscordPresence struct {
	ClientID string
	Conn     net.Conn
}

func (c *DiscordPresence) Connect() error {
	for index := 0; index <= 9; index++ {
		conn, err := winio.DialPipe(fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, index), nil)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}

		if err != nil {
			return err
		}

		c.Conn = conn
		return nil

	}
	return ErrCannotConnect

}
