package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *domain.UserRegisterRequest) error {
	// 检查手机号是否已注册
	existingUser, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("该手机号已被注册")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &domain.User{
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Age:       req.Age,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.Create(ctx, user)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *domain.UserLoginRequest) (*domain.UserLoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("手机号或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("手机号或密码错误")
	}

	// 生成JWT
	tokenString, err := s.createUserToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.UserLoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

// createUserToken 创建用户JWT token
func (s *UserService) createUserToken(user *domain.User) (string, error) {
	secret := viper.GetString("general.jwt")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	claims := domain.UserClaims{
		UserID: user.ID,
		Phone:  user.Phone,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 300)), // 300天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GetUserById 根据ID获取用户信息
func (s *UserService) GetUserById(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return s.userRepo.FindById(ctx, id)
}

