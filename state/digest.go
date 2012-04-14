package state

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"hash"
	"io"
	"log"
)

var digest hash.Hash

// newKey must be called under the lock (or in a context
// where the lock is unnecessary)
func (g *gameState) newKey() {
	g.Key = make([]byte, sha1.Size)
	_, err := io.ReadFull(rand.Reader, g.Key)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) GetHash() hash.Hash {
	return hmac.New(sha1.New, g.g.Key)
}
