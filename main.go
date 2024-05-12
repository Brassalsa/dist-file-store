package main

import (
	"bytes"
	"log"
	"time"

	"github.com/Brassalsa/dist-file-store/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTranportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTranport := p2p.NewTCPTransport(tcpTranportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr[1:] + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTranport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTranport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(time.Second * 1)
	go s2.Start()

	time.Sleep(time.Second * 2)
	data := bytes.NewReader([]byte("very big file please help"))
	if err := s2.StoreData("store_key", data); err != nil {
		log.Println(err)
	}

	select {}
}
