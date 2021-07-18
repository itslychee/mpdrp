// +build windows

package ipc

import (
	"errors"
	"fmt"
	"io/fs"
	"net"

	windows "github.com/Microsoft/go-winio"
)

type DiscordPresence struct {
	ClientID string
	conn     net.Conn
}

func (c *DiscordPresence) Connect() error {
	for index := 0; index <= 9; index++ {
		pipePath := fmt.Sprintf("\\\\.\\pipe\\discord-ipc-%d", index)
		conn, err := windows.DialPipe(pipePath, nil)
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
