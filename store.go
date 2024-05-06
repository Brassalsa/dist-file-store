package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolder = "_file_store"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

func DefaultPathTransformFunc(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type PathTransformFunc func(str string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func (p *PathKey) FullPathWithRoot(root string) string {
	return fmt.Sprintf("%s/%s/%s", root, p.PathName, p.FileName)
}

func (p *PathKey) FirstName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[0]
}

type StoreOpts struct {
	// Root is the root folder in which folders/files will be saved
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if opts.Root == "" {
		opts.Root = defaultRootFolder
	}
	return &Store{
		StoreOpts: opts,
	}
}

// clear the whole root directory and files inside it
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

// check if file exists
func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(key)

	_, err := os.Stat(pathkey.FullPathWithRoot(s.Root))

	return err == nil
}

// delete a file
func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	return os.RemoveAll(s.Root + "/" + pathKey.FirstName())
}

// read a file
func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(pathKey.FullPathWithRoot(s.Root))
}

// save file to disk
func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(s.Root+"/"+pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	fullPath := pathKey.FullPathWithRoot(s.Root)

	f, err := os.Create(fullPath)

	if err != nil {
		return err
	}

	defer f.Close()

	m, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", m, fullPath)

	return nil
}
