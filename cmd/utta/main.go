package main

import (
	"log"

	"github.com/julian7/utta/tunnel"
)

func main() {
	config, err := tunnel.GetConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = config.Tunnel()
	if err != nil {
		log.Fatal(err)
	}
}
