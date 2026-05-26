// Package handler is uses for stream random frequencies
package handler

import (
	"time"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"github.com/Mgeorg1/go_randomaliens/internal/server/service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	pb.UnimplementedFrequencyGeneratorServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) StreamRandom(req *pb.Empty, stream pb.FrequencyGeneratorService_StreamRandomServer) error {
	freqGenerator := service.NewGenerator()
	sessionID := uuid.NewString()
	for {
		timestamp := timestamppb.New(time.Now())
		event := &pb.FrequencyEvent{
			SessionId: sessionID,
			Frequency: freqGenerator.Generate(),
			Timestamp: timestamp,
		}
		if err := stream.Send(event); err != nil {
			return err
		}
	}
}
