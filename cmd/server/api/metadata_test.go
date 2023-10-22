package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"testing"
)

func TestExtractMetaData_NoMetadata(t *testing.T) {
	server := Server{}
	ctx := context.Background()

	meta := server.extractMetaData(ctx)

	assert.Empty(t, meta.UserAgent)
	assert.Empty(t, meta.ClientIp)
}

func TestExtractMetaData_GrpcGatewayUserAgent(t *testing.T) {
	server := Server{}
	md := metadata.Pairs(grpcGatewayUserAgentHeader, "grpc-gateway-agent")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	meta := server.extractMetaData(ctx)

	assert.Equal(t, "grpc-gateway-agent", meta.UserAgent)
}

func TestExtractMetaData_GeneralUserAgent(t *testing.T) {
	server := Server{}
	md := metadata.Pairs(userAgentHeader, "general-agent")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	meta := server.extractMetaData(ctx)

	assert.Equal(t, "general-agent", meta.UserAgent)
}

func TestExtractMetaData_XForwardedFor(t *testing.T) {
	server := Server{}
	md := metadata.Pairs(xForwardedForHeader, "192.168.0.1")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	meta := server.extractMetaData(ctx)

	assert.Equal(t, "192.168.0.1", meta.ClientIp)
}

func TestExtractMetaData_PeerContext(t *testing.T) {
	server := Server{}
	peerCtx := &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("192.168.0.2")}}
	ctx := peer.NewContext(context.Background(), peerCtx)

	meta := server.extractMetaData(ctx)

	assert.Equal(t, "192.168.0.2:0", meta.ClientIp)
}
