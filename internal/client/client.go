package client

import (
	"context"
	"fmt"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient pb.FrequencyGeneratorServiceClient
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

func (c *Client) Run(ctx context.Context, handle func(*pb.FrequencyEvent)) error {

	stream, err := c.getStream()
	if err != nil {
		return fmt.Errorf("getting stream: %w", err)
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("receiving event: %w", err)
		}
		handle(event)
	}
}
