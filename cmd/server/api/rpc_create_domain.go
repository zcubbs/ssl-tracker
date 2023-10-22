package api

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	pb "github.com/zcubbs/tlz/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateDomain(ctx context.Context, req *pb.CreateDomainRequest) (*pb.CreateDomainResponse, error) {
	namespaceId, err := uuid.Parse(req.GetNamespace())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid namespace id: %v", err)
	}

	domain, err := s.store.InsertDomain(ctx, db.InsertDomainParams{
		Name:      req.GetName(),
		Status:    pgtype.Text{},
		Namespace: namespaceId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create domain: %v", err)
	}

	return &pb.CreateDomainResponse{
		Domain: convertDomain(domain),
	}, nil
}
