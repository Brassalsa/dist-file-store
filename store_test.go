package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

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

const key = "specialKey"

var data = []byte("some jpeg")

func TestCreate(t *testing.T) {
	s := newStore()
	if _, err := s.Write(s.Id, key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	s := newStore()
	if err := s.Delete(s.Id, key); err != nil {
		t.Error(err)
	}

	has := s.Has(s.Id, key)

	if has == true {
		t.Error("file is not deleted")
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("random_%v", i)
		// create file and write data
		if _, err := s.Write(s.Id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		// check file exists
		if has := s.Has(s.Id, key); !has {
			t.Error("file not found")
		}

		// read a file
		_, r, err := s.Read(s.Id, key)

		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		fmt.Printf("reading file: %s\n", string(b))

		if string(b) != string(data) {
			t.Errorf("want %s have %s", string(b), string(data))
		}

		// close file
		rc, ok := r.(io.ReadCloser)
		if ok {
			rc.Close()
		}
		// delete file
		if err := s.Delete(s.Id, key); err != nil {
			t.Error(err)
		}

		// check again file exists
		if has := s.Has(s.Id, key); has {
			t.Error("file is not deleted")
		}
	}

}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              ".test_store",
		Id:                generateId(),
	}

	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Errorf("error clearing store: %s\n", err)
	}
}
