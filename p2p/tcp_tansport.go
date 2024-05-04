package p2p

import (
	"fmt"
	"net"
)

// tcp peer --------------------->
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

// implements Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// tcp tarnsport------------------->
type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcCh    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

// implements Transport interface, returns read-only channel
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

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

		fmt.Printf("New incomming connection %v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	var err error
	defer func() {
		fmt.Printf("droping peer connection: %s\n", err)
		conn.Close()
	}()

	if err = t.HandShakeFunc(peer); err != nil {
		fmt.Printf("TCP handshake error: %s\n", err)
		conn.Close()
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	//read loop
	rpc := RPC{}
	for {
		if err = t.Decoder.Decode(conn, &rpc); err != nil {
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
	}

}
