package main

import (
	"fmt"
	"log"

	"github.com/Brassalsa/go-dist-file-store/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("hello world")
	select {}
}
