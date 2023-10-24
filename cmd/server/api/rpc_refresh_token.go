package api

import (
	"context"
	pb "github.com/zcubbs/ssl-tracker/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	refreshTokenPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %v", err)
	}

	session, err := s.store.GetSession(ctx, refreshTokenPayload.ID)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid session: %v", err)
	}

	if session.IsBlocked {
		return nil, status.Errorf(codes.Unauthenticated, "session is blocked")
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token has expired")
	}

	user, err := s.store.GetUser(ctx, session.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.ID,
		s.cfg.Auth.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessTokenPayload.ExpiredAt),
	}, nil
}
