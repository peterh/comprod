package main

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/peterh/comprod/state"
)

func doHash(g *state.Game, thing, name string) string {
	digest := g.GetHash()
	io.WriteString(digest, thing)
	io.WriteString(digest, name)
	sum := digest.Sum(nil)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(sum), "=")
}

func cookieHash(g *state.Game, name string) string {
	return doHash(g, "cookie:", name)
}

func inviteHash(g *state.Game, name string) string {
	return doHash(g, "invite:", name)
}

func pwdHash(g *state.Game, name, password string) []byte {
	digest := g.GetHash()
	io.WriteString(digest, "password/")
	io.WriteString(digest, name)
	io.WriteString(digest, ":")
	io.WriteString(digest, password)
	return digest.Sum(nil)
}
