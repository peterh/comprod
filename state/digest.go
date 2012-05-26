package state

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"hash"
	"io"
	"log"
)

// newKey must be called under the lock (or in a context
// where the lock is unnecessary)
func (g *GameState) newKey() {
	g.Key = make([]byte, sha1.Size)
	_, err := io.ReadFull(rand.Reader, g.Key)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) GetHash() hash.Hash {
	g.Lock()
	defer g.Unlock()
	return hmac.New(sha1.New, g.g.Key)
}

func GetSeed() (rv int64) {
	binary.Read(rand.Reader, binary.LittleEndian, &rv)
	return
}
