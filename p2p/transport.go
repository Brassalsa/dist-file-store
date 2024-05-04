package p2p

// represents the remote node
type Peer interface {
	Close() error
}

// handles commnunication btw nodes in network
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
