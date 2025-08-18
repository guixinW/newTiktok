package application

import (
	"context"
	"errors"
	"log/slog"
	"newTiktoken/internal/user/domain/model"
	"newTiktoken/internal/user/domain/password"
	"newTiktoken/internal/user/domain/repository"
)

// ErrUserAlreadyExists 是一个标准错误，表示用户已存在
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrPasswordNotMatch  = errors.New("password does not match")
)

// UserService 提供了用户相关的应用服务
type UserService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

// NewUserService 创建一个新的 UserService
func NewUserService(userRepo repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger.With("layer", "application"),
	}
}

// RegisterUser 处理用户注册的逻辑
func (s *UserService) RegisterUser(ctx context.Context, username, plainPassword string) (*model.User, error) {
	// 1. 检查用户是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.Error("failed to check for existing user", "username", username, "error", err)
		return nil, err
	}
	if existingUser != nil {
		s.logger.Warn("registration attempt for existing user", "username", username)
		return nil, ErrUserAlreadyExists
	}

	// 2. 使用 Argon2 对密码进行哈希处理
	hashedPassword, err := password.Hash(plainPassword)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, err
	}

	// 3. 创建新的用户领域模型
	newUser := model.NewUser(username, hashedPassword)
	s.logger.Debug("creating new user object", "username", username)

	// 4. 保存用户到仓库
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		s.logger.Error("failed to save new user", "username", username, "error", err)
		return nil, err
	}

	s.logger.Info("user created successfully", "user_id", newUser.ID, "username", username)
	return newUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, plainPassword string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.Error("failed to find user", "username", username, "error", err)
		return nil, ErrUserNotFound
	}
	ok, err := user.CheckPassword(plainPassword)
	if err != nil {
		s.logger.Error("hashed password err", "username", username, "error", err)
		return nil, ErrPasswordNotMatch
	}
	if !ok {
		s.logger.Warn("password incorrect", "username", username)
		return nil, ErrPasswordNotMatch
	}
	return user, nil
}

// GetUserInfo 处理获取用户信息的逻辑
func (s *UserService) GetUserInfo(ctx context.Context, userID uint64) (*model.User, error) {
	s.logger.Debug("getting user info", "user_id", userID)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user info from repository", "user_id", userID, "error", err)
		return nil, err
	}
	if user == nil {
		s.logger.Warn("user not found in repository", "user_id", userID)
	}
	return user, nil
}
