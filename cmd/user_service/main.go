package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"newTiktoken/internal/user/application"
	"newTiktoken/internal/user/infrastructure/persistence"
	grpcinterface "newTiktoken/internal/user/interfaces/grpc"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	"newTiktoken/pkg/pb/user"
	"os"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Use standard log for fatal errors before logger is initialized
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Initialize logger
	appLogger := logger.New(cfg.LogLevel)
	appLogger.Info("logger initialized")
	appLogger.Info("starting user service...")

	// 3. Initialize database connection
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		appLogger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	appLogger.Info("database connected")

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

	// 5. Create and assemble repositories (Infrastructure)
	mysqlRepo, err := persistence.NewMySQLUserRepository(db, appLogger)
	if err != nil {
		appLogger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}
	appLogger.Info("mysql repository initialized")

	cachedRepo := persistence.NewRedisUserRepository(rdb, mysqlRepo, appLogger)
	appLogger.Info("redis cache repository initialized")

	// 6. Initialize application service (Application)
	userApp := application.NewUserService(cachedRepo, appLogger)
	appLogger.Info("user service application initialized")

	// 7. Initialize gRPC server (Interfaces)
	grpcServer := grpcinterface.NewUserServer(userApp, appLogger)
	appLogger.Info("gRPC server initialized")

	// 8. Create and register gRPC service with interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcServer.LogInterceptor),
	)
	user.RegisterUserServiceServer(s, grpcServer)
	reflection.Register(s)

	// 9. Start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		appLogger.Error("failed to listen", "port", cfg.Port, "error", err)
		os.Exit(1)
	}

	appLogger.Info("user service listening", "port", cfg.Port)
	if err := s.Serve(lis); err != nil {
		appLogger.Error("failed to serve gRPC server", "error", err)
		os.Exit(1)
	}
}
