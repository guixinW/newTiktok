package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	"syscall"
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

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册服务
	// user_relation.RegisterUserRelationServiceServer(server, &userRelationService{})

	// 启动服务器
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.UserRelationService.Port))
		if err != nil {
			appLogger.Error("failed to listen", "error", err)
			os.Exit(1)
		}
		if err := server.Serve(lis); err != nil {
			appLogger.Error("failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	appLogger.Info("user relation service started", "port", cfg.UserRelationService.Port)

	// 优雅地关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("shutting down user relation service...")
	server.GracefulStop()
	appLogger.Info("user relation service stopped")
}
