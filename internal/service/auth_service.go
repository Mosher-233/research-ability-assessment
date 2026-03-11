package service

import (
	"context"
	"errors"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/repository/postgres"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *postgres.UserRepo
}

func NewAuthService(userRepo *postgres.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, user *models.User) error {
	// 生成UUID
	user.ID = uuid.New().String()
	
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	
	// 创建用户
	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("用户不存在")
	}
	
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("密码错误")
	}
	
	// 生成JWT令牌
	token, err := s.generateToken(user)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("无效的令牌")
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}