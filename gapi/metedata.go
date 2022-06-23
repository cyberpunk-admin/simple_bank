package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grapcGatewayUserAgent = "grpcgateway-user-agent"
	userAgent = "user-agent"
	xForwardedFor = "x-forwarded-for"
)

type Metedata struct {
	UserAget string
	ClientIP string
}

func extractMetedata(ctx context.Context) *Metedata {
	mtdt := &Metedata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if UserAgents := md.Get(grapcGatewayUserAgent); len(UserAgents) > 0 {
			mtdt.UserAget = UserAgents[0]
		}

		if UserAgents := md.Get(userAgent); len(UserAgents) > 0 {
			mtdt.UserAget = UserAgents[0]
		}

		if ClientIPs := md.Get(xForwardedFor); len(ClientIPs) > 0 {
			mtdt.ClientIP = ClientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
