package main

import (
	"fmt"
	"io"
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

	time.Sleep(time.Second * 1)
	key := "coolPicture.jpg"

	// data := bytes.NewReader([]byte("very big file please help part"))
	// if err := s2.Store(key, data); err != nil {
	// 	log.Println(err)
	// }

	// time.Sleep(time.Second)
	r, err := s2.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("file is found: ", string(b))

}
