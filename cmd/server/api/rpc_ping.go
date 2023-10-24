package api

import (
	"context"
	"github.com/zcubbs/ssl-tracker/cmd/server/config"
	pb "github.com/zcubbs/ssl-tracker/pb"
)

func (s *Server) Ping(_ context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message:   "Pong",
		Version:   config.Version,
		Commit:    config.Commit,
		BuildTime: config.Date,
	}, nil
}
