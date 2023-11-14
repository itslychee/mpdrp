package mpd

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func (r Command) String() string {
	var builder strings.Builder
	var quoteReplacer *strings.Replacer = strings.NewReplacer("'", "\\'", "\"", "\\\"")

	builder.Write(append([]byte(r.Name), 0x20))

	for _, val := range r.Args {
		newVal := quoteReplacer.Replace(val)
		builder.WriteString(strings.Join([]string{"\"", "\""}, newVal))
		builder.WriteByte(0x20)
	}
	builder.WriteByte(0x0A)
	return builder.String()
}

type Response struct {
	Records map[string][]string
	Command []Command
	Binary []byte
	eol     []byte
}


func (r Response) Get(value string) string {
	return r.Records[value][0]
}

func (r Response) OK() []byte {
	return r.eol
}

// ResponseError is supposed to represent the line followed by ACK in MPD errors
type ResponseError struct {
	ErrorEnum uint64
	Offset    uint64
	Command   string
	Message   string
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("error: [%d] %s", r.ErrorEnum, r.Message)
}
