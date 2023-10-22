package api

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/zcubbs/tlz/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) LogoutUser(ctx context.Context, req *pb.LogoutUserRequest) (*pb.Empty, error) {
	authPayload, err := s.requireUser(ctx)
	if err != nil {
		return nil, unauthorizedError(err)
	}

	violations := validateLogoutUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	u, err := s.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid session id")
	}
	session, err := s.store.GetSession(ctx, sessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get session: %v", err)
	}

	if session.UserID != u.ID {
		return nil, status.Errorf(codes.PermissionDenied, "user does not own session")
	}

	_, err = s.store.BlockSession(ctx, session.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout session: %v", err)
	}

	return &pb.Empty{}, nil
}

func validateLogoutUserRequest(_ *pb.LogoutUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	return violations
}
