package ipc

import "errors"

var (
	ErrCannotConnect = errors.New("cannot connect to the socket/pipe")
)
