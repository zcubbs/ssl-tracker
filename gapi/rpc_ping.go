package gapi

import (
	"context"
	"github.com/zcubbs/tlz/pb"
)

func (s *Server) Ping(_ context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "Pong",
	}, nil
}
