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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	videopersistence "newTiktoken/internal/video/infrastructure/persistence"
	"newTiktoken/internal/video_favorite/application"
	"newTiktoken/internal/video_favorite/infrastructure/persistence"
	grpcfavorite "newTiktoken/internal/video_favorite/interfaces/grpc"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	videoFavoritePb "newTiktoken/pkg/pb/video_favorite"
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
	appLogger.Info("starting video_favorite service...")

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

	// 5. Create and assemble repositories
	mysqlRepo, err := persistence.NewMySQLVideoFavoriteRepository(db, appLogger)
	if err != nil {
		appLogger.Error("failed to initialize mysql video_favorite repository", "error", err)
		os.Exit(1)
	}
	appLogger.Info("mysql repository initialized")

	cachedRepo := persistence.NewRedisVideoFavoriteRepository(rdb, mysqlRepo)
	appLogger.Info("redis cache repository initialized")

	// Initialize video repository for dependency
	videoMySQLRepo, err := videopersistence.NewMySQLVideoRepository(db, appLogger)
	if err != nil {
		appLogger.Error("failed to initialize mysql video repository for dependency", "error", err)
		os.Exit(1)
	}
	videoCachedRepo := videopersistence.NewRedisVideoRepository(rdb, videoMySQLRepo)

	// 6. Initialize application service
	favoriteApp := application.NewVideoFavoriteService(cachedRepo, videoCachedRepo)
	appLogger.Info("video_favorite service application initialized")

	// 7. Initialize gRPC server
	favoriteServiceServer := grpcfavorite.NewServer(favoriteApp)
	appLogger.Info("gRPC server initialized")

	// 8. Create and register gRPC service
	s := grpc.NewServer()
	videoFavoritePb.RegisterFavoriteServiceServer(s, favoriteServiceServer)
	reflection.Register(s)

	// 9. Start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		appLogger.Error("failed to listen", "port", cfg.Port, "error", err)
		os.Exit(1)
	}

	appLogger.Info("video_favorite service listening", "port", cfg.Port)
	if err := s.Serve(lis); err != nil {
		appLogger.Error("failed to serve gRPC server", "error", err)
		os.Exit(1)
	}
}
