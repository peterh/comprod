package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var root, data, hostname *string
var port = flag.String("port", ":2012", "TCP port to listen on")

func init() {
	var r string
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		path := filepath.Join(p, "src", "comprod")
		test := filepath.Join(path, "static", "cp.css")
		_, err := os.Stat(test)
		if err == nil {
			r = path
			break
		}
	}
	root = flag.String("root", r, "Directory where static/* and templates/* are found")

	const defState = ".comprodState"
	var path string
	my, err := user.Current()
	if err != nil {
		path = filepath.Join(os.Getenv("HOME"), defState)
		log.Println(err)
	} else {
		path = filepath.Join(my.HomeDir, defState)
	}
	data = flag.String("data", path, "File where game data is stored")

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	hostname = flag.String("hostname", host, "Name used in invitation links")
}
