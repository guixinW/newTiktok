package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"newTiktoken/internal/gateway/handlers"
	"newTiktoken/pkg/config"
	"newTiktoken/pkg/logger"
	"newTiktoken/pkg/userpb"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	appLogger := logger.New(cfg.LogLevel)
	appLogger.Info("starting API gateway...")

	userConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.UserService.ServiceName, cfg.UserService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		appLogger.Error("failed to connect to user service", "error", err)
		os.Exit(1)
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)
	appLogger.Info("connected to user service", "port", cfg.UserService.Port)

	userHandler := handlers.NewUserHandler(userClient, appLogger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.GET("/info", userHandler.GetUserInfo)
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Gateway.Port),
		Handler: router,
	}

	go func() {
		appLogger.Info("gateway listening", "port", cfg.Gateway.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("shutting down gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("gateway forced to shutdown", "error", err)
	}

	appLogger.Info("gateway exited")
}
