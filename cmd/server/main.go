package main

import (
	"context"
	"fmt"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/logger"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/configuration"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/db"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/db/repository"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/handler"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/interceptor"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/security"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	logger.InitLogger()

	config := configuration.Configure()

	// Init database
	dbConnection, err := db.NewDBStorage(config.DatabaseDsn)
	if err != nil {
		zap.S().Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	err = db.RunMigrations(dbConnection)
	if err != nil {
		zap.S().Fatal("Failed to run migrations", zap.Error(err))
	}

	// Check if port is free
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
	if err != nil {
		zap.L().Fatal("Failed to listen", zap.Error(err))
	}

	// Create dependencies
	jwtService := security.NewJwtService([]byte(config.JwtSecret), config.JwtLifetimeHours)

	// User
	userRepository := repository.NewUserRepository(dbConnection)
	userService := service.NewUserService(userRepository, jwtService)

	// Secret
	secretRepository := repository.NewSecretRepository(dbConnection)
	secretService := service.NewSecretService(secretRepository)

	// Create list of interceptors
	var interceptors []grpc.UnaryServerInterceptor
	interceptors = append(interceptors, interceptor.AuthorizationInterceptor(jwtService))

	// Create gRPCServer
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	// Register handlers
	pb.RegisterUserServiceServer(gRPCServer, handler.NewUserHandler(userService))
	pb.RegisterSecretServiceServer(gRPCServer, handler.NewSecretHandler(secretService))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	// Start gRPC server in separate goroutine
	go func() {
		zap.L().Info("gRPC server started", zap.Int("port", config.GRPCPort))
		if err := gRPCServer.Serve(listen); err != nil {
			zap.L().Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Waiting for signal
	<-ctx.Done()
	zap.L().Info("Shutting down gRPC server...")
	gRPCServer.GracefulStop()
}
