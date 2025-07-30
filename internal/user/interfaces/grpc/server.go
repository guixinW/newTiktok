package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"newTiktoken/internal/user/application"
	"newTiktoken/internal/user/domain/model"
	pb "newTiktoken/pkg/userpb" // 引入生成的 pb 文件
	"time"
)

// Server 是 gRPC 服务器，它实现了 userpb.UserServiceServer 接口
type Server struct {
	pb.UnimplementedUserServiceServer // 必须嵌入，以保证向前兼容
	app    *application.UserService
	logger *slog.Logger
}

// NewUserServer 创建一个新的 gRPC 用户服务服务器
func NewUserServer(app *application.UserService, logger *slog.Logger) *Server {
	return &Server{
		app:    app,
		logger: logger.With("layer", "gRPC"),
	}
}

// LogInterceptor 是一个 gRPC 拦截器，用于记录请求信息
func (s *Server) LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	s.logger.Info("gRPC request received", "method", info.FullMethod)

	resp, err := handler(ctx, req)

	duration := time.Since(startTime)
	logger := s.logger.With("method", info.FullMethod, "duration", duration.String())

	if err != nil {
		logger.Error("gRPC request failed", "error", err)
	} else {
		logger.Info("gRPC request completed")
	}

	return resp, err
}

// UserRegister 实现了注册的 gRPC 方法
func (s *Server) UserRegister(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	user, err := s.app.RegisterUser(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		s.logger.Warn("user registration failed", "username", req.GetUsername(), "error", err)
		return &pb.UserRegisterResponse{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		}, nil
	}

	s.logger.Info("user registered successfully", "user_id", user.ID)
	return &pb.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     user.ID,
	}, nil
}

// UserInfo 实现了获取用户信息的 gRPC 方法
func (s *Server) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	user, err := s.app.GetUserInfo(ctx, req.GetQueryUserId())
	if err != nil {
		s.logger.Error("failed to get user info", "query_user_id", req.GetQueryUserId(), "error", err)
		return &pb.UserInfoResponse{
			StatusCode: 500,
			StatusMsg:  "internal server error",
		}, nil
	}

	if user == nil {
		s.logger.Warn("user not found", "query_user_id", req.GetQueryUserId())
		return &pb.UserInfoResponse{
			StatusCode: 404,
			StatusMsg:  "user not found",
		}, nil
	}

	return &pb.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User:       convertUserToPB(user),
	}, nil
}

// convertUserToPB 是一个辅助函数，用于将领域模型 User 转换为 Protobuf 的 User
func convertUserToPB(user *model.User) *pb.User {
	if user == nil {
		return nil
	}
	return &pb.User{
		Id:             user.ID,
		Name:           user.Username,
		FollowingCount: user.FollowingCount,
		FollowerCount:  user.FollowerCount,
		IsFollow:       false, // is_follow 的逻辑需要额外实现
		TotalFavorite:  user.TotalFavorited,
		WorkCount:      user.WorkCount,
		FavoriteCount:  user.FavoriteCount,
	}
}