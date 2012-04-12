package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"hash"
	"io"
	"log"
)

var digest hash.Hash

func (g *gameState) newKey() {
	g.Key = make([]byte, sha1.Size)
	_, err := io.ReadFull(rand.Reader, g.Key)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *gameState) initKey() {
	digest = hmac.New(sha1.New, g.Key)
}
