package main

import (
	"flag"
	"log"

	"github.com/Mgeorg1/go_randomaliens/internal/server"
)

func main() {
	addr := flag.String("address", "127.0.0.1:50051", "address string e.g. 127.0.0.1:50051")
	server := server.NewServer()
	err := server.Run(*addr)
	if err != nil {
		log.Fatal(err)
	}
}
