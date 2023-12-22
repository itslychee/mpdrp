package mpd

import (
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var errParser = regexp.MustCompile(`^ACK\s\[(\d*)@(\d*)\]\s{([a-zA-Z _]*)}\s(.+)`)

type MPDConnInfo struct {
	Address  net.Addr
	Password string
	Timeout  time.Duration
}

func ResolveAddr(address string) (addr net.Addr) {
	var err error
	switch {
	case strings.HasPrefix(address, "@/"):
		addr, err = net.ResolveUnixAddr("unixgram", address)
	case strings.HasPrefix(address, "/"):
		addr, err = net.ResolveUnixAddr("unix", address)
	default:
		addr, err = net.ResolveTCPAddr("tcp", address)
	}
	if err != nil {
		panic(err)
	}
	return
}

type MPDConnection struct {
	Conn    *textproto.Conn
	RawConn net.Conn
}

func (mpd *MPDConnection) Connect(conn *MPDConnInfo) error {
	// We immediately default to connecting to the specified address only, otherwise
	// we'll try a list of defaults to connect to.
	if conn == nil {
		conn = &MPDConnInfo{}
	}
	if conn.Address != nil {
		return mpd.connect(
			conn.Address.Network(),
			conn.Address.String(),
			conn.Timeout,
		)
	}

	var timeout time.Duration
	val, ok := os.LookupEnv("MPD_TIMEOUT")
	if !ok {
		timeout = conn.Timeout
	} else {
		v, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("MPD_TIMEOUT parsing error: %w", err)
		}
		timeout = time.Second * time.Duration(v)
	}

	var mpdAddresses []MPDConnInfo
	if runtime.GOOS != "windows" {
		if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
			mpdAddresses = append(mpdAddresses, MPDConnInfo{
				Address:  ResolveAddr(filepath.Join(dir, "mpd/socket")),
				Timeout:  timeout,
				Password: conn.Password,
			})
		}
		mpdAddresses = append(mpdAddresses, MPDConnInfo{
			Address:  ResolveAddr("/run/mpd/socket"),
			Timeout:  timeout,
			Password: conn.Password,
		})
	}

	connInfo := MPDConnInfo{Timeout: timeout}
	mpdHost := os.Getenv("MPD_HOST")
	if mpdHost == "" {
		mpdHost = "127.0.0.1"
	}
	mpdPort, ok := os.LookupEnv("MPD_PORT")
	if !ok {
		mpdPort = "6600"
	}
	if strings.Contains(mpdHost, "@") {
		parts := strings.SplitN(mpdHost, "@", 2)
		mpdHost = parts[1]
		if parts[0] != "" {
			connInfo.Password = parts[0]
		}
		if strings.HasPrefix(mpdHost, "@/") {
			connInfo.Address = ResolveAddr(mpdHost)
			mpdAddresses = append(mpdAddresses, connInfo)
		}
	}

	addr := MPDConnInfo{
		Address: ResolveAddr(net.JoinHostPort(mpdHost, mpdPort)),
		Timeout: timeout,
	}
	mpdAddresses = append(mpdAddresses, addr)
	if addr.Address.String() != "127.0.0.1:6600" {
		mpdAddresses = append(mpdAddresses, MPDConnInfo{
			Timeout: timeout,
			Address: ResolveAddr("127.0.0.1:6600"),
		})
	}

	for _, address := range mpdAddresses {
		err := mpd.connect(address.Address.Network(), address.Address.String(), address.Timeout)
		if err == nil {
			return nil
		}

		// Code to distinguish connection refused from more serious
		// errors, as they shouldn't be silenced.
		switch v := err.(type) {
		case *net.OpError:
			if v.Op == "read" || v.Op == "dial" {
				continue
			}
		case syscall.Errno:
			if v == syscall.ECONNREFUSED {
				continue
			}
		}
		return err
	}
	return nil
}

func (mpd *MPDConnection) connect(network, address string, timeout time.Duration) error {
	// Support timeouts to allow more flexibility
	if mpd.Conn != nil {
		return errors.New("cannot connect to MPD, you're already connected")
	}
	c, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return err
	}
	mpd.RawConn = c 
	mpd.Conn = textproto.NewConn(c)
	s, err := mpd.Conn.R.ReadString(0x0A)
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
	sid := mpd.Conn.Next()

	mpd.Conn.StartRequest(sid)

	mpd.Conn.W.WriteString("command_list_begin\n")
	for _, val := range cmds {
		mpd.Conn.W.Write([]byte(val.String()))
	}
	mpd.Conn.W.WriteString("command_list_end\n")
	if err := mpd.Conn.W.Flush(); err != nil {
		return nil, err
	}

	mpd.Conn.EndRequest(sid)

	var response = &Response{
		Records: make(map[string][]string),
	}

	mpd.Conn.StartResponse(sid)
	defer mpd.Conn.EndResponse(sid)
	for {
		s, err := mpd.Conn.R.ReadString(0x0A)
		s = strings.TrimSpace(s)
		if err != nil {
			return response, err
		}

		switch {
		case strings.HasPrefix(s, "OK"):
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
			// TODO: Fix binary responses
			fields := strings.SplitN(string(s), ":", 2)
			fields[1] = strings.TrimSpace(fields[1])
			response.Records[fields[0]] = append(response.Records[fields[0]], fields[1])
			if fields[0] == "binary" {
				// After retrieving the binary data, we won't break out of this loop as
				// the "OK" should be present at the end
				// https://mpd.readthedocs.io/en/latest/protocol.html#binary-responses
				allocSize, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}
				response.Binary = make([]byte, allocSize)
				if _, err = mpd.Conn.R.Read(response.Binary); err != nil {
					panic(err)
				}
			}

		}

	}

}

func (mpd *MPDConnection) Close() error {
	if mpd.Conn != nil {
		return mpd.Conn.Close()
	}
	return nil
}
