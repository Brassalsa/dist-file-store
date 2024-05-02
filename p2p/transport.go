package p2p

// represents the remote node
type Peer interface {
}

// handles commnunication btw nodes in network
type Transport interface {
	ListenAndAccept() error
}
