package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	userAgentHeader            = "user-agent"
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (s *Server) extractMetaData(ctx context.Context) *Metadata {
	meta := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(xForwardedForHeader); len(clientIps) > 0 {
			meta.ClientIp = clientIps[0]
		}
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		if meta.ClientIp == "" {
			meta.ClientIp = p.Addr.String()
		}
	}

	return meta
}
