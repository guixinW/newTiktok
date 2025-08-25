package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"newTiktoken/internal/video_comment/application"
	"newTiktoken/internal/video_comment/infrastructure/persistence"
	grpcinterface "newTiktoken/internal/video_comment/interfaces/grpc"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	videoCommentPb "newTiktoken/pkg/pb/video_comment"
	"os"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 初始化日志
	appLogger := logger.New(cfg.LogLevel)
	appLogger.Info("starting user relation service...")

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		appLogger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	mysqlRepo := persistence.NewMySQLVideoCommentRepository(db)

	// 创建 Application
	videoCommentApp := application.NewCommentApplicationService(mysqlRepo, appLogger)

	// 创建 gRPC 服务器
	grpcServer := grpcinterface.NewCommentServer(videoCommentApp, appLogger)

	// 注册服务
	s := grpc.NewServer()
	videoCommentPb.RegisterCommentServiceServer(s, grpcServer)
	reflection.Register(s)

	// 开始服务
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
