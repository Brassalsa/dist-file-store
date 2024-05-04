package p2p

import "errors"

// inavlid handshake return when conn btw 2 nodes couldn't be established
var ErrInvalidHandshake = errors.New("invalid handshake")

// handshake
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {
	return nil
}
