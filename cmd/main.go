package main

import (
	"context"
	"flag"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/GDH-Project/auth/cmd/config"
	projectGrpc "github.com/GDH-Project/auth/internal/grpc"
	"github.com/GDH-Project/auth/internal/grpc/authpb"
	"github.com/GDH-Project/auth/internal/grpc/userpb"
	"github.com/GDH-Project/auth/internal/repository"
	"github.com/GDH-Project/auth/internal/resource"
	"github.com/GDH-Project/auth/internal/service"
	"github.com/GDH-Project/auth/internal/use_case"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	if !*debug {
		debug = flag.Bool("d", false, "enable debug mode")
	}
	flag.Parse()

	config.InitLogger(*debug)
	cfg := config.GetConfig()

	db := resource.InitDB(cfg)
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userUseCase := usecase.NewUserUseCase(userService)

	authRepository := repository.NewAuthRepository(cfg, db)
	authService := service.NewAuthService(authRepository)
	authUseCase := usecase.NewAuthUseCase(cfg, authService, userService)

	// grpc server 세팅
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			projectGrpc.RecoveryInterceptor(),
		),
	)
	authpb.RegisterAuthServiceServer(grpcServer, projectGrpc.NewGrpcAuthService(authUseCase))
	userpb.RegisterUserServiceServer(grpcServer, projectGrpc.NewGrpcUserService(userUseCase))

	if *debug {
		// 개발모드 일 경우 grpc reflection 활성화
		reflection.Register(grpcServer)
	}

	lis, err := net.Listen("tcp", ":50501")
	if err != nil {
		zap.S().Fatal("gRPC 서버 초기화 오류",
			zap.Error(err),
		)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		zap.S().Info("gRPC server is running on port :50501")
		if err := grpcServer.Serve(lis); err != nil {
			zap.S().Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// 프로세스 종료 시그널 발생 대기
	<-ctx.Done()
	zap.S().Info("프로세스 종료 시그널이 감지되었습니다.")

	// 5초 타임 아웃 설정
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// gRPC 서비스 다운 상태를 확인하기 위한 채널 생성
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		// gRPC 서버가 정상적으로 종료되면 채널을 닫음
		close(stopped)
	}()

	select {
	case <-stopped:
		zap.S().Info("gRPC 서버가 정상적으로 종료되었습니다.")
	case <-shutdownCtx.Done():
		zap.S().Warn("서버 종료 제한시간에 도달했습니다. 강제 종료를 진행합니다.")
		grpcServer.Stop()
	}

	zap.S().Info("서버가 종료되었습니다.")
}
