package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Mgeorg1/go_randomaliens/internal/client"
	"github.com/Mgeorg1/go_randomaliens/internal/client/event_handler"
)

func main() {
	addr := flag.String("address", "127.0.0.1:50051", "The address of the gRPC server")
	flag.Parse()

	client, err := client.NewClient(*addr)
	if err != nil {
		log.Fatal(fmt.Errorf("can't create client: %w", err))
	}
	eventHandler := event_handler.NewLogEventHandler()
	err = client.Run(context.Background(), eventHandler.Handle)
	if err != nil {
		log.Fatal(fmt.Errorf("error running client: %w", err))
	}
}
