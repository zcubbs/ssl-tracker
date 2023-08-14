package gapi

import (
	"fmt"
	"github.com/charmbracelet/log"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pb"
	"github.com/zcubbs/tlz/pkg/token"
	"github.com/zcubbs/tlz/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	pb.UnimplementedTlzServer
	store      db.Store
	tokenMaker token.Maker
	cfg        util.Config
}

func NewServer(store db.Store, cfg util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.Auth.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create new tokenMaker: %w", err)
	}

	s := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		cfg:        cfg,
	}

	return s, nil
}

func (s *Server) Start() error {
	grpcServer := grpc.NewServer()
	pb.RegisterTlzServer(grpcServer, s)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	log.Info("starting gRPC server", "port", s.cfg.GrpcServer.Port)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	return nil
}
