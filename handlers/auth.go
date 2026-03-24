package handlers

import (
	"errors"
	"sync"
	"time"

	"shorturl/models"

	"gorm.io/gorm"
)

var authService *AuthService

type AuthService struct {
	mu sync.RWMutex
}

// UserCredentials 用户登录请求
type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token       string `json:"token"`
	FirstLogin  bool   `json:"first_login"`
	Username    string `json:"username"`
}

func NewAuthService() *AuthService {
	if authService == nil {
		authService = &AuthService{}
	}
	return authService
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*LoginResponse, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if user.Password != password {
		return nil, errors.New("用户名或密码错误")
	}

	if !user.IsActive {
		return nil, errors.New("账户已被禁用")
	}

	// 判断是否首次登录（密码是初始密码）
	firstLogin := (password == "password")

	return &LoginResponse{
		Token:      "shorturl-secret-token-2026",
		FirstLogin: firstLogin,
		Username:   user.Username,
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(username, oldPassword, newPassword string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证旧密码
	if user.Password != oldPassword {
		return errors.New("原密码错误")
	}

	// 更新密码
	if err := models.DB.Model(&user).Update("password", newPassword).Error; err != nil {
		return err
	}

	return nil
}

// InitDefaultUser 初始化默认用户
func InitDefaultUser() error {
	var count int64
	if err := models.DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		user := models.User{
			Username:  "admin",
			Password:  "password",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return models.DB.Create(&user).Error
	}

	return nil
}
