package main

import (
	"bytes"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {

	var (
		payload = "foo not bar"
		src     = bytes.NewReader([]byte(payload))
		dst     = new(bytes.Buffer)
		key     = newEncryptionKey()
	)

	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}
	if nw != (16 + len(payload)) {
		t.Error("bytes length must be payload length + 16")
	}

	if out.String() != payload {
		t.Errorf("have %s want %s", out.String(), payload)
	}
}
