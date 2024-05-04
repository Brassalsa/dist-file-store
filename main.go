package main

import (
	"fmt"
	"log"

	"github.com/Brassalsa/go-dist-file-store/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// peer.Close()
	fmt.Println("TODO")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		rpc := tr.Consume()
		for msg := range rpc {
			fmt.Printf("%v\n", msg)
		}

	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("hello world")
	select {}
}
