package p2p

const (
	IncommingMessage = 0x2
	IncommingStream  = 0x1
)

// holds any arbitrary data that is sent btw 2 nodes
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
