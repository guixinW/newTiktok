package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"newTiktoken/internal/video/application"
	"newTiktoken/internal/video/infrastructure/persistence"
	grpcvideo "newTiktoken/internal/video/interfaces/grpc"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	"newTiktoken/pkg/videopb"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Initialize logger
	appLogger := logger.New(cfg.LogLevel)
	appLogger.Info("logger initialized")
	appLogger.Info("starting video service...")

	// 3. Initialize database connection (using old method for video service)
	mysqlRepo, err := persistence.NewMySQLVideoRepository(cfg.Database.DSN)
	if err != nil {
		appLogger.Error("failed to connect to mysql", "error", err)
		os.Exit(1)
	}
	appLogger.Info("mysql repository initialized")

	// 4. Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		appLogger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	appLogger.Info("redis connected")

	// 5. Create and assemble repositories
	// Note: The video service repositories do not currently accept a logger.
	cachedRepo := persistence.NewRedisVideoRepository(rdb, mysqlRepo)
	appLogger.Info("redis cache repository initialized")

	// 6. Initialize application service
	// Note: The video application service does not currently accept a logger.
	videoApp := application.NewVideoService(cachedRepo)
	appLogger.Info("video service application initialized")

	// 7. Initialize gRPC server
	// Note: The video gRPC server does not currently accept a logger or have an interceptor.
	videoServiceServer := grpcvideo.NewServer(videoApp)
	appLogger.Info("gRPC server initialized")

	// 8. Create and register gRPC service
	s := grpc.NewServer()
	videopb.RegisterVideoServiceServer(s, videoServiceServer)
	reflection.Register(s)

	// 9. Start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		appLogger.Error("failed to listen", "port", cfg.Port, "error", err)
		os.Exit(1)
	}

	appLogger.Info("video service listening", "port", cfg.Port)
	if err := s.Serve(lis); err != nil {
		appLogger.Error("failed to serve gRPC server", "error", err)
		os.Exit(1)
	}
}