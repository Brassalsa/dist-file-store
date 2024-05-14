package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct {
}

func (dec GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return nil
	}

	// incase of stream, don't decode and return
	stream := peekBuf[0] == IncommingStream
	if stream {
		rpc.Stream = true
		return nil
	}

	buf := make([]byte, 1028)
	m, err := r.Read(buf)
	if err != nil {
		return err
	}
	rpc.Payload = buf[:m]
	return nil
}
