package discord

import "errors"

var (
	ErrCannotConnect = errors.New("cannot connect to the socket/pipe")
)
