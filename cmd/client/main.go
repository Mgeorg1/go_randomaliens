package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mgeorg1/go_randomaliens/internal/client"
	"github.com/Mgeorg1/go_randomaliens/internal/client/event_handler"
	"github.com/Mgeorg1/go_randomaliens/internal/client/repo"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, exiting...")
		cancel()
	}()

	addr := flag.String("addr", "127.0.0.1:50051", "The address of the gRPC server")
	k := flag.Float64("k", 3.0, "The coefficient for anomaly detection (diff > k * stddev)")
	dbHost := flag.String("db_host", "127.0.0.1", "The address of the PostgreSQL server")
	dbPort := flag.Int("p", 5432, "The port of the PostgreSQL server")
	dbUser := flag.String("u", "postgres", "The username for the PostgreSQL server")
	dbPassword := flag.String("pwd", "secret", "The password for the PostgreSQL server")
	dbName := flag.String("db_name", "randomaliens", "The name of the PostgreSQL database")

	flag.Parse()
	if *k <= 0 {
		log.Fatal("k must be greater than 0")
	}

	client, err := client.NewClient(*addr)
	if err != nil {
		log.Fatal(fmt.Errorf("can't create client: %w", err))
	}
	repo, err := repo.NewRepo(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	eventHandler := event_handler.NewWelfordEventHandler(*k, repo)
	if err != nil {
		log.Fatal(fmt.Errorf("can't create repo: %w", err))
	}
	err = client.Run(ctx, eventHandler.Handle)
	if err != nil {
		log.Fatal(fmt.Errorf("error running client: %w", err))
	}
}
