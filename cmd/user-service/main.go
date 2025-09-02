package main

import (
	"context"
	"google.golang.org/grpc"
	userpb "newTiktoken/internal/common/genproto/user"
	"newTiktoken/internal/common/server"
	"newTiktoken/internal/user/ports"
	"newTiktoken/internal/user/service"
)

func main() {
	ctx := context.Background()
	application := service.NewApplication(ctx)
	server.RunGRPCServer(func(srv *grpc.Server) {
		svc := ports.NewGrpcServer(application)
		userpb.RegisterUserServiceServer(srv, svc)
	})
}
