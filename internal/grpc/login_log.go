package grpc

import (
	"context"
	"net"

	"github.com/GDH-Project/auth/internal/domain"
	"google.golang.org/grpc/metadata"
)

func makeLoginLog(ctx context.Context) *domain.LoginLog {
	var userIP, userAgent string
	md, _ := metadata.FromIncomingContext(ctx)

	if xci := md.Get("x-client-ip"); len(xci) > 0 {
		userIP = xci[0]
	}

	if ua := md.Get("x-user-agent"); len(ua) > 0 {
		userAgent = ua[0]
	}

	return &domain.LoginLog{
		UserIP:    net.ParseIP(userIP),
		UserAgent: userAgent,
	}
}
