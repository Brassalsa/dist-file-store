package p2p

import (
	"fmt"
	"net"
	"sync"
)

// respresents remote node over tcp connection
type TCPPeer struct {
	// underlying connection of peer
	conn net.Conn

	// if dialed/send -> outbound = true
	// if accept/recieved -> inbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, oubound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: oubound,
	}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener

	mu   sync.RWMutex
	peer map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("New incomming connection %v\n", peer)
}
