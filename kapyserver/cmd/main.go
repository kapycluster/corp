package main

import (
	"log"

	"github.com/kapycluster/corpy/kapyserver/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start kapy-server: %+v\n", err)
	}
}
