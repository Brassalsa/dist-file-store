package p2p

import "net"

// represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// handles commnunication btw nodes in network
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
