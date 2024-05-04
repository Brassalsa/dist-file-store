package p2p

import "net"

// holds any arbitrary data that is sent btw 2 nodes
type RPC struct {
	From    net.Addr
	Payload []byte
}
