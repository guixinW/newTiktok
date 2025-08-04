package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"newTiktoken/pkg/userpb"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
	logger     *slog.Logger
}

func NewUserHandler(userClient userpb.UserServiceClient, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		logger:     logger,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,max=32"`
	Password string `json:"password" binding:"required,max=32"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,max=32"`
	Password string `json:"password" binding:"required,max=32"`
	DeviceId string `json:"device_id" binding:"required,max=32"`
}

type RefreshTokenRequest struct {
	UserId       uint64 `json:"user_id" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceId     string `json:"device_id" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "Invalid request parameters",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userClient.UserRegister(ctx, &userpb.UserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Error("user register failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_id":     resp.UserId,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "Invalid request parameters",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userClient.Login(ctx, &userpb.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
		DeviceId: req.DeviceId,
	})
	if err != nil {
		h.logger.Error("user login failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code":   resp.StatusCode,
		"status_msg":    resp.StatusMsg,
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	tokenUserIdStr := c.Query("token")

	if userIdStr == "" || tokenUserIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "Missing required parameters",
		})
		return
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "Invalid user_id format",
		})
		return
	}

	tokenUserId, err := strconv.ParseUint(tokenUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "Invalid token format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userClient.UserInfo(ctx, &userpb.UserInfoRequest{
		QueryUserId: userId,
		TokenUserId: tokenUserId,
	})
	if err != nil {
		h.logger.Error("get user info failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "Internal server error",
		})
		return
	}

	var user gin.H
	if resp.User != nil {
		user = gin.H{
			"id":              resp.User.Id,
			"name":            resp.User.Name,
			"following_count": resp.User.FollowingCount,
			"follower_count":  resp.User.FollowerCount,
			"is_follow":       resp.User.IsFollow,
			"total_favorite":  resp.User.TotalFavorite,
			"work_count":      resp.User.WorkCount,
			"favorite_count":  resp.User.FavoriteCount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user":        user,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "Invalid request parameters",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userClient.RefreshAccessToken(ctx, &userpb.AccessTokenRequest{
		UserId:       req.UserId,
		RefreshToken: req.RefreshToken,
		DeviceId:     req.DeviceId,
	})
	if err != nil {
		h.logger.Error("refresh token failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code":  resp.StatusCode,
		"status_msg":   resp.StatusMsg,
		"access_token": resp.AccessToken,
	})
}
