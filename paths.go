package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var data, hostname *string
var root = flag.String("root", "", "Directory where static/* and templates/* are found (default: use built-in)")
var port = flag.String("port", ":2012", "TCP port to listen on")

func init() {
	var path string
	my, err := user.Current()
	if err != nil {
		log.Println(err)
	} else {
		path = filepath.Join(my.HomeDir, ".comprodState")
	}
	data = flag.String("data", path, "File where game data is stored")

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	hostname = flag.String("hostname", host, "Name used in invitation links")
}
