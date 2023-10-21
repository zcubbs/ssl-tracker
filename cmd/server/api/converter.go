package api

import (
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	pb "github.com/zcubbs/tlz/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func convertDomain(domain db.Domain) *pb.Domain {
	return &pb.Domain{
		Name:              domain.Name,
		Status:            domain.Status.String,
		Issuer:            domain.Issuer.String,
		Owner:             domain.Owner.String(),
		CertificateExpiry: timestamppb.New(domain.CertificateExpiry.Time),
		CreatedAt:         timestamppb.New(domain.CreatedAt),
	}
}
