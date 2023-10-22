package api

import (
	"context"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	dbUtil "github.com/zcubbs/tlz/cmd/server/db/util"
	pb "github.com/zcubbs/tlz/pb"
	"github.com/zcubbs/x/password"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser handles the creation of a new user via gRPC.
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, err := s.requireAdmin(ctx)
	if err != nil {
		return nil, unauthorizedError(err)
	}

	violations := validateCreateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	hashedPass, err := password.Hash(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	role := req.GetRole().String()
	if role == pb.Role_ROLE_UNSPECIFIED.String() {
		role = pb.Role_ROLE_USER.String()
	}

	// Prepare parameters for the database function.
	params := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPass,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		Role:           role,
	}

	// Use the Store interface to save the user details to the database.
	user, err := s.store.CreateUser(ctx, params)
	if err != nil {
		if err == dbUtil.ErrUniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username already taken: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Convert the db.User to pb.User for the response.
	respUser := convertUserToPb(user)

	return &pb.CreateUserResponse{User: respUser}, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
