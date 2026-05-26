package server

import (
	"net"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"github.com/Mgeorg1/go_randomaliens/internal/server/handler"
	"google.golang.org/grpc"
)

type Server struct {
	grpc *grpc.Server
}

func NewServer() *Server {
	s := grpc.NewServer()
	handler := handler.NewHandler()
	pb.RegisterFrequencyGeneratorServiceServer(s, handler)

	return &Server{grpc: s}
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.grpc.Serve(lis)
}
