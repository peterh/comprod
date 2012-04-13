package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"hash"
	"io"
	"log"
	"strings"
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

func doHash(thing, name string) string {
	digest.Reset()
	io.WriteString(digest, thing)
	io.WriteString(digest, name)
	sum := digest.Sum(nil)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(sum), "=")
}

func cookieHash(name string) string {
	return doHash("cookie:", name)
}

func inviteHash(name string) string {
	return doHash("invite:", name)
}

func pwdHash(name, password string) []byte {
	digest.Reset()
	io.WriteString(digest, "password/")
	io.WriteString(digest, name)
	io.WriteString(digest, ":")
	io.WriteString(digest, password)
	return digest.Sum(nil)
}
