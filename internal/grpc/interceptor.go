package grpc

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				zap.S().Error("gRPC 서비스 Panic 발생",
					"panic", r,
					"stack", debug.Stack(),
					"method", info.FullMethod,
				)

				// 보안을 위해 패닉 오류의 원본 메시지를 숨긴다.
				err = status.Error(codes.Internal, "Internal Server Error")
			}
		}()
		return handler(ctx, req)
	}
}
