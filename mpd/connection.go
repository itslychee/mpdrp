package mpd

import (
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var errParser = regexp.MustCompile(`^ACK\s\[(\d*)@(\d*)\]\s{([a-zA-Z _]*)}\s(.+)`)

type MPDConnection struct {
	conn *textproto.Conn
}

func (mpd *MPDConnection) Connect(network, address string, timeout time.Duration) error {
	// Support timeouts to allow more flexibility
	c, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return err
	}
	mpd.conn = textproto.NewConn(c)
	s, err := mpd.conn.R.ReadString(0x0A)
	if !strings.HasPrefix(s, "OK MPD") {
		return fmt.Errorf("the server did not answer correctly, got %s instead", s)
	}

	if err != nil {
		return err
	}
	return nil
}

func (mpd *MPDConnection) Exec(cmds ...Command) (*Response, error) {
	// Command execution
	sid := mpd.conn.Next()

	mpd.conn.StartRequest(sid)

	mpd.conn.W.WriteString("command_list_begin\n")
	for _, val := range cmds {
		mpd.conn.W.Write([]byte(val.String()))
	}
	mpd.conn.W.WriteString("command_list_end\n")
	if err := mpd.conn.W.Flush(); err != nil {
		return nil, err
	}

	mpd.conn.EndRequest(sid)

	var response = &Response{
		Records: make(map[string]string),
	}

	mpd.conn.StartResponse(sid)
	defer mpd.conn.EndResponse(sid)
	for {
		s, err := mpd.conn.R.ReadString(0x0A)
		s = strings.TrimSpace(s)
		if err != nil {
			return response, err 
		}

		switch {
		case strings.HasPrefix(s, "OK"):
			// As defined by the spec, this is where we should stop reading
			response.eol = []byte(s)
			return response, nil
		case strings.HasPrefix(s, "ACK"):
			// TODO: Add logging

			// This is suppose to denote an error from the server
			res := errParser.FindSubmatch([]byte(s))
			if res == nil {
				panic("nil slice from parsing error returned by mpd")
			}

			// Converts the resulting group to uint64, respecting
			// the size of C++'s enums as defined here:
			// https://github.com/MusicPlayerDaemon/MPD/blob/d39b11ba5d0f9e36e59f1fdf7321bcd64c3bfe26/src/protocol/Ack.hxx#L26L40
			//
			// It also seems if MPD doesn't support 32 bit platforms (but don't take my word for it), so
			// returning a 64 bit from calling ParseUint doesn't seem like a bad choice here.

			enum, err := strconv.ParseUint(string(res[1]), 10, 64)
			if err != nil {
				panic(err)
			}
			offset, err := strconv.ParseUint(string(res[2]), 10, 64)
			if err != nil {
				panic(err)
			}

			respErr := ResponseError{
				ErrorEnum: enum,
				Offset:    offset,
				Command:   string(res[3]),
				Message:   string(res[4]),
			}
			return nil, respErr
		default:
			fields := strings.SplitN(string(s), ":", 2)
			fields[1] = strings.TrimSpace(fields[1])
			response.Records[fields[0]] = fields[1]
			if fields[0] == "binary" {
				// After retrieving the binary data, we won't break out of this loop as
				// the "OK" should be present at the end
				// https://mpd.readthedocs.io/en/latest/protocol.html#binary-responses
				allocSize, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}
				response.Data = make([]byte, allocSize)
				if _, err = mpd.conn.R.Read(response.Data); err != nil {
					panic(err)
				}
			}

		}

	}

}

func (mpd *MPDConnection) Close() error {
	if mpd.conn != nil {
		return mpd.conn.Close()
	}
	return nil
}
