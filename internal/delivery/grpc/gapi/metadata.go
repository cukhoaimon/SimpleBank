package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	grpcGatewayClientIP  = "x-forwarded-host"
)

func (handler *Handler) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(grpcGatewayUserAgent); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
		if clientIp := md.Get(grpcGatewayClientIP); len(clientIp) > 0 {
			mtdt.ClientIP = clientIp[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
