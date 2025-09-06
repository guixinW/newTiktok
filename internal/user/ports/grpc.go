package ports

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"newTiktoken/internal/common/auth"
	userPb "newTiktoken/internal/common/genproto/user"
	"newTiktoken/internal/user/app"
	"newTiktoken/internal/user/app/command"
	"newTiktoken/internal/user/app/query"
)

type GrpcServer struct {
	userPb.UnimplementedUserServiceServer
	app app.Application
}

func NewGrpcServer(application app.Application) *GrpcServer {
	return &GrpcServer{app: application}
}

func (g *GrpcServer) CreateUser(ctx context.Context, req *userPb.CreateUserRequest) (*emptypb.Empty, error) {
	if err := g.app.Commands.CreateUser.Handle(ctx, command.CreateUser{
		UUID:   req.GetUuid(),
		Name:   req.GetName(),
		Age:    uint16(req.GetAge()),
		Gender: uint16(req.GetGender()),
	}); err != nil {
		log.Printf("ERROR: failed to handle command: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) UpdateUser(ctx context.Context, req *userPb.UpdateUserRequest) (*emptypb.Empty, error) {
	if err := g.app.Commands.UpdateUser.Handle(ctx, command.UpdateUser{
		UUID:   req.GetUuid(),
		Name:   req.GetName(),
		Age:    uint16(req.GetAge()),
		Gender: uint16(req.GetGender()),
	}); err != nil {
		log.Printf("ERROR: failed to handle command: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) GetUserInformation(ctx context.Context, request *userPb.GetUserInformationRequest) (*userPb.User, error) {
	usr, err := g.app.Queries.InformationOfUser.Handle(ctx, query.InformationOfUser{
		User: auth.User{UUID: request.GetUuid()},
	})
	pbUser := queryUserToProtoUser(usr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return pbUser, nil
}

func queryUserToProtoUser(user *query.User) *userPb.User {
	return &userPb.User{
		Uuid:           user.UUID,
		Name:           user.Name,
		Age:            uint32(user.Age),
		Gender:         uint32(user.Gender),
		FollowingCount: user.FollowingCount,
		FollowerCount:  user.FollowerCount,
		TotalFavorite:  user.TotalFavorite,
		WorkCount:      user.WorkCount,
		FavoriteCount:  user.FavoriteCount,
	}
}
