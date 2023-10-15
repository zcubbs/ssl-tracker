package api

import (
	"context"
	"github.com/google/uuid"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	"github.com/zcubbs/tlz/internal/validator"
	"github.com/zcubbs/tlz/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	id, err := uuid.Parse(req.GetEmailId())
	if err != nil {
		return nil, invalidArgumentError([]*errdetails.BadRequest_FieldViolation{
			fieldViolation("email_id", err),
		})
	}

	txResult, err := s.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    id,
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email : %v", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmailID(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := validator.ValidatePassword(req.SecretCode); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
