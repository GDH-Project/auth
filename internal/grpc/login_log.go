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

	if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
		userIP = xff[0]
	} else if xri := md.Get("x-real-ip"); len(xri) > 0 {
		userIP = xri[0]
	} else if xci := md.Get("x-client-ip"); len(xci) > 0 {
		userIP = xci[0]
	}

	if ua := md.Get("user-agent"); len(ua) > 0 {
		userAgent = ua[0]
	}

	return &domain.LoginLog{
		UserIP:    net.ParseIP(userIP),
		UserAgent: userAgent,
	}
}
