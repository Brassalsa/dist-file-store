package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              "_test_store",
	}

	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Errorf("error clearing store: %s\n", err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "test picture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "8706061afc327915b34f73db0dfc57dce7ff0520"
	expectedPathName := "87060/61afc/32791/5b34f/73db0/dfc57/dce7f/f0520"
	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.FileName != expectedOriginalKey {
		t.Errorf("have %s want %s", pathKey.FileName, expectedOriginalKey)
	}
}

var key = "specialKey"

var data = []byte("some jpeg")

func TestCreate(t *testing.T) {
	s := newStore()
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	s := newStore()
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	has := s.Has(key)

	if has == true {
		t.Error("file is not deleted")
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)
	// create file and write data
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	// check file exists
	if has := s.Has(key); has == false {
		t.Error("file not found")
	}

	// read a file
	r, err := s.Read(key)

	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Printf("reading file: %s\n", string(b))

	if string(b) != string(data) {
		t.Errorf("want %s have %s", string(b), string(data))
	}

	// delete file
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
