package p2p

import "net"

// represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// handles commnunication btw nodes in network
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
