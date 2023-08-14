package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	sti := status.New(codes.InvalidArgument, "invalid parameters")
	br := &errdetails.BadRequest{
		FieldViolations: violations,
	}
	std, err := sti.WithDetails(br)
	if err != nil {
		return sti.Err()
	}
	return std.Err()
}
