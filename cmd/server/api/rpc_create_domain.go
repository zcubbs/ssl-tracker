package api

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zcubbs/ssl-tracker/cmd/server/db/sqlc"
	pb "github.com/zcubbs/ssl-tracker/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

func validateCreateDomainRequest(req *pb.CreateDomainRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateDomainName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}
