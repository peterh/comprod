package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var root, data *string
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

	var path string
	my, err := user.Current()
	if err != nil {
		log.Println(err)
	} else {
		path = filepath.Join(my.HomeDir, ".comprodState")
	}
	data = flag.String("data", path, "File where game data is stored")
}
