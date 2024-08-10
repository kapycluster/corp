package main

import (
	"log"

	"github.com/kapycluster/corpy/kapyserver/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("kapyserver: %+v\n", err)
	}
}
