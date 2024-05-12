package p2p

// holds any arbitrary data that is sent btw 2 nodes
type RPC struct {
	From    string
	Payload []byte
}
