package util

import (
	"crypto/rand"
	"io"
	"log"
)

func GetRandomBytesOrDie(len int) []byte {
	b := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
