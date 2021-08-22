package discord

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/google/uuid"
)

func (c *DiscordPresence) Close() error {
	if c.conn != nil {
		_, err := c.Send(Close, Payload{})
		if err != nil {
			return err
		}
		return c.conn.Close()
	}
	return nil
}

func (c *DiscordPresence) CreateHandshake() error {
	payload := Payload{
		Version:  1,
		ClientID: c.ClientID,
	}
	_, err := c.Send(Handshake, payload)
	return err
}

func (c *DiscordPresence) SetActivity(activity *Activity) (string, []byte, error) {
	nonce, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}
	payload := Payload{
		Cmd: "SET_ACTIVITY",
		Args: &Arguments{
			Pid:      os.Getpid(),
			Activity: activity,
		},
		Nonce: nonce.String(),
	}
	n, err := c.Send(Frame, payload)
	return nonce.String(), n, err
}

func (c *DiscordPresence) Send(opcode OpCode, payload Payload) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	if err := binary.Write(buffer, binary.LittleEndian, uint32(opcode)); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, uint32(len(data))); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(data); err != nil {
		return nil, err
	}
	_, err = c.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	prelude := make([]byte, 8)
	if _, err := c.Read(prelude); err != nil {
		return nil, err
	}
	_, length := binary.LittleEndian.Uint32(prelude[:4]), binary.LittleEndian.Uint32(prelude[4:])
	d := make([]byte, length)
	if _, err := c.Read(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *DiscordPresence) Write(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	if errors.Is(err, io.EOF) {
		c.Close()
	}
	return
}

func (c *DiscordPresence) Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	if errors.Is(err, io.EOF) {
		c.Close()
	}
	return
}

func NewDiscordPresence(ClientID string) *DiscordPresence {
	return &DiscordPresence{
		ClientID: ClientID,
	}
}
