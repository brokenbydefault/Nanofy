package Nanofy

import (
	"golang.org/x/crypto/blake2b"
	"io"
)

// CreateFileHash computes the Blake2b-256 for given file reader.
func CreateFileHash(file io.Reader) (hash []byte, err error) {
	blake, _ := blake2b.New(32, nil)

	_, err = io.Copy(blake, file)
	if err != nil {
		return
	}

	return blake.Sum(nil), nil
}

// CreateFileHash computes the Blake2b-256 for given byte.
func CreateHash(text []byte) (hash []byte, err error) {
	blake, _ := blake2b.New(32, nil)

	_, err = blake.Write(text)
	if err != nil {
		return
	}

	return blake.Sum(nil), nil
}
