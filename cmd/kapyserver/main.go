package main

import (
	"log"

	"kapycluster.com/corp/kapyserver/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("kapyserver: %+v\n", err)
	}
}
