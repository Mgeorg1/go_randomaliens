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
	k := flag.Float64("k", 3.0, "The coefficient for anomaly detection (diff > k * stddev)")
	flag.Parse()
	if *k <= 0 {
		log.Fatal("k must be greater than 0")
	}

	client, err := client.NewClient(*addr)
	if err != nil {
		log.Fatal(fmt.Errorf("can't create client: %w", err))
	}
	eventHandler := event_handler.NewWelfordEventHandler(*k)
	err = client.Run(context.Background(), eventHandler.Handle)
	if err != nil {
		log.Fatal(fmt.Errorf("error running client: %w", err))
	}
}
