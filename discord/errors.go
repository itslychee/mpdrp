package ipc

import "errors"

var (
	ErrNotConnected  = errors.New("not connected to the socket/pipe")
	ErrCannotConnect = errors.New("cannot connect to the socket/pipe")
)
