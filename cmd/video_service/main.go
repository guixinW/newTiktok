package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"newTiktoken/internal/video/application"
	"newTiktoken/internal/video/domain/repository"
	"newTiktoken/internal/video/infrastructure/persistence"
	grpcvideo "newTiktoken/internal/video/interfaces/grpc"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	"newTiktoken/pkg/videopb"
)

func main() {
	// Initialize logger
	log := logger.New("video-service", "dev") // Or load level from config

	// Load configuration
	cfg, err := config.Load("config.yaml") // Ensure config.yaml is present
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repository based on config
	repo, err := initRepo(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Start gRPC server
	listenAddress := fmt.Sprintf(":%d", cfg.VideoService.Port)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Dependency Injection
	videoAppService := application.NewVideoService(repo)
	grpcServer := grpc.NewServer()
	videoServiceServer := grpcvideo.NewServer(videoAppService)
	videopb.RegisterVideoServiceServer(grpcServer, videoServiceServer)

	log.Infof("Video gRPC service starting on %s", listenAddress)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

// initRepo initializes and returns the appropriate video repository based on the configuration.
func initRepo(cfg *config.Config, log logger.Logger) (repository.VideoRepository, error) {
	// Default to in-memory if no database is configured
	if cfg.Database.DSN == "" {
		log.Info("Using in-memory repository")
		return persistence.NewInMemoryVideoRepository(), nil
	}

	// Initialize MySQL repository
	log.Infof("Connecting to MySQL database...")
	mysqlRepo, err := persistence.NewMySQLVideoRepository(cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	log.Info("MySQL connection successful")

	// If Redis is not configured, return the MySQL repository directly
	if cfg.Redis.Addr == "" {
		log.Info("Using MySQL repository without Redis cache")
		return mysqlRepo, nil
	}

	// Initialize Redis client
	log.Infof("Connecting to Redis at %s...", cfg.Redis.Addr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
	if _, err := redisClient.Ping(redisClient.Context()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Info("Redis connection successful")

	// Wrap MySQL repository with Redis cache
	log.Info("Using MySQL repository with Redis cache")
	return persistence.NewRedisVideoRepository(redisClient, mysqlRepo), nil
}
