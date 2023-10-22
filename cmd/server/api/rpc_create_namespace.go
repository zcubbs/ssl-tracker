package api

import (
    "context"
    "github.com/google/uuid"
    db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
    pb "github.com/zcubbs/tlz/pb"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *Server) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {

    userId, err := uuid.Parse(req.GetUserId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %v", err)
    }

    namespace, err := s.store.InsertNamespace(ctx, db.InsertNamespaceParams{
        Name:   req.GetName(),
        UserID: userId,
    })
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create namespace: %v", err)
    }

    return &pb.CreateNamespaceResponse{
        Namespace: convertNamespace(namespace),
    }, nil
}
