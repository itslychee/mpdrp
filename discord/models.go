package ipc

type OpCode uint32

const (
	Handshake OpCode = iota
	Frame
	Close
	Ping
	Pong
)

type Payload struct {
	Cmd   string      `json:"cmd,omitempty"`
	Args  *Arguments  `json:"args,omitempty"`
	Nonce string      `json:"nonce,omitempty"`
	Data  interface{} `json:"data,omitempty"`

	// Handshake specific data
	ClientID string `json:"client_id,omitempty"`
	Version  int8   `json:"v,omitempty"`
}

type Arguments struct {
	Pid      int       `json:"pid,omitempty"`
	Activity *Activity `json:"activity,omitempty"`
}

type Activity struct {
	State      *string      `json:"state,omitempty"`
	Details    *string      `json:"details,omitempty"`
	Timestamps *Timestamps `json:"timestamps,omitempty"`
	Assets     *Assets     `json:"assets,omitempty"`
}

type Timestamps struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type Assets struct {
	LargeImage string `json:"large_image"`
	LargeText  string `json:"large_text"`
	SmallImage string `json:"small_image"`
	SmallText  string `json:"small_text"`
}
