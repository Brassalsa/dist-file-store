package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolder = "_file_store"

const (
	TypeGetFile    = 0x1
	TypeDeleteFile = 0x2
)

// content addressable storage function
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

func (p *PathKey) FullPathWRoot(root, id string) string {
	return fmt.Sprintf("%s/%s/%s/%s", root, id, p.PathName, p.FileName)
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
	Root string
	// Id of the owner of storage
	Id                string
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

// clears the whole root directory and files inside it
func (s *Store) Clear() error {
	defer log.Printf("deleted [%s]", s.Root)
	return os.RemoveAll(s.Root)
}

// check if file exists
func (s *Store) Has(id, key string) bool {
	pathkey := s.PathTransformFunc(key)

	_, err := os.Stat(pathkey.FullPathWRoot(s.Root, id))

	return err == nil
}

// delete a file
func (s *Store) Delete(id, key string) error {
	pathKey := s.PathTransformFunc(key)
	defer fmt.Printf("deleted from disk: [%s]\n", pathKey.FullPathWRoot(s.Root, id))

	path := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstName())
	return os.RemoveAll(path)
}

// read a file
func (s *Store) Read(id, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Store) readStream(id, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	file, err := os.Open(pathKey.FullPathWRoot(s.Root, id))

	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {

		return 0, nil, err
	}

	return fi.Size(), file, err
}

func (s *Store) openFileForWriting(id, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)

	path := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	fullPath := pathKey.FullPathWRoot(s.Root, id)

	return os.Create(fullPath)
}

// save file to disk
func (s *Store) Write(id, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

func (s *Store) WriteDecrypt(encKey []byte, id, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	m, err := copyDecrypt(encKey, r, f)
	return int64(m), err
}

func (s *Store) writeStream(id, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	return io.Copy(f, r)
}
