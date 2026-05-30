package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"github.com/Mgeorg1/go_randomaliens/internal/client/event_handler"
	"github.com/Mgeorg1/go_randomaliens/internal/client/repo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient pb.FrequencyGeneratorServiceClient
	repo       *repo.Repo
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("creating new client: %w", err)
	}
	client := pb.NewFrequencyGeneratorServiceClient(conn)
	return &Client{grpcClient: client}, nil
}

func (c *Client) getStream() (grpc.ServerStreamingClient[pb.FrequencyEvent], error) {
	stream, err := c.grpcClient.StreamRandom(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("getting stream: %w", err)
	}
	return stream, nil
}

func (c *Client) Run(ctx context.Context, handle func(*pb.FrequencyEvent) bool) error {
	stream, err := c.getStream()
	if err != nil {
		return fmt.Errorf("getting stream: %w", err)
	}
	eventCh := make(chan *pb.FrequencyEvent, 100)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go event_handler.HandlerWorker(wg, ctx, eventCh, handle)

	defer func() {
		close(eventCh)
		wg.Wait()
	}()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		event, err := stream.Recv()
		if err != nil {

			return fmt.Errorf("receiving event: %w", err)
		}
		eventCh <- event
	}
}
