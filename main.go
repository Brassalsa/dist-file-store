package main

import (
	"bytes"
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
		EncKey:            newEncryptionKey(),
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
	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")
	s3 := makeServer(":5000", ":3000", ":4000")

	go s1.Start()
	go s2.Start()

	time.Sleep(time.Second * 1)
	go s3.Start()

	time.Sleep(time.Second * 1)

	for i := range 10 {

		key := fmt.Sprintf("coolPicture_%d.jpg", i)
		fileData := fmt.Sprintf("very long file please help {%d}", i)

		data := bytes.NewReader([]byte(fileData))
		if err := s3.Store(key, data); err != nil {
			log.Println(err)
		}
		if err := s3.store.Delete(key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)

		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("file is found: ", string(b))
	}

}
