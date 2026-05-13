package service

import (
	"context"
	"errors"
	"time"

	"cms/internal/model"
	"cms/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrWrongPassword     = errors.New("密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")
	ErrInvalidToken       = errors.New("无效的token")
	ErrTokenExpired       = errors.New("token已过期")
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Register 注册
func (s *UserService) Register(ctx context.Context, username, password, email string) (*model.User, error) {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}
	
	user := &model.User{
		Username: username,
		Email:    email,
		Nickname: username,
		Status:   1,
	}
	
	// 加密密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	
	return user, nil
}

// Login 登录，返回JWT token
func (s *UserService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, ErrUserNotFound
		}
		return "", nil, err
	}
	
	// 检查用户状态
	if user.Status != 1 {
		return "", nil, ErrUserDisabled
	}
	
	// 验证密码
	if !user.CheckPassword(password) {
		return "", nil, ErrWrongPassword
	}
	
	// 生成JWT token
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}
	
	// 缓存token
	cache.Set(ctx, cache.UserToken+token, user.ID, 24*time.Hour)
	
	return token, &user, nil
}

// generateToken 生成JWT token
func (s *UserService) generateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key-change-in-production"))
}

// ValidateToken 验证token
func (s *UserService) ValidateToken(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key-change-in-production"), nil
	})
	
	if err != nil {
		return nil, ErrInvalidToken
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	
	// 检查是否过期
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return nil, ErrTokenExpired
	}
	
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	var user model.User
	if err := s.db.First(&user, uint(userID)).Error; err != nil {
		return nil, ErrUserNotFound
	}
	
	return &user, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
