package api

import (
	"github.com/google/uuid"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	pb "github.com/zcubbs/tlz/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUserToPb(user db.User) *pb.User {
	return &pb.User{
		Id:                user.ID.String(),
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		Role:              pb.Role(pb.Role_value[user.Role]),
	}
}

func convertPbToUser(user *pb.User) db.User {
	id, _ := uuid.Parse(user.Id)
	return db.User{
		ID:                id,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.AsTime(),
		CreatedAt:         user.CreatedAt.AsTime(),
		Role:              user.Role.String(),
	}
}

func convertDomain(domain db.Domain) *pb.Domain {
	return &pb.Domain{
		Id:                domain.ID.String(),
		Name:              domain.Name,
		Status:            domain.Status.String,
		Issuer:            domain.Issuer.String,
		Namespace:         domain.Namespace.String(),
		CertificateExpiry: timestamppb.New(domain.CertificateExpiry.Time),
		CreatedAt:         timestamppb.New(domain.CreatedAt),
	}
}

func convertNamespace(namespace db.Namespace) *pb.Namespace {
	return &pb.Namespace{
		Id:        namespace.ID.String(),
		Name:      namespace.Name,
		UserId:    namespace.UserID.String(),
		CreatedAt: timestamppb.New(namespace.CreatedAt),
	}
}
