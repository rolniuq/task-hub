package app

import (
	"context"
	"errors"
	"taskhub/config"
	"taskhub/internal/domains/user"
	"taskhub/internal/domains/user/repo"
	"taskhub/pkg/base/entity"
	"taskhub/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

var AuthServiceModule = fx.Module(
	"auth-service",
	fx.Provide(NewAuthService),
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type AuthService struct {
	config   *config.Config
	userRepo *repo.UserRepository
}

func NewAuthService(config *config.Config, userRepo *repo.UserRepository) *AuthService {
	return &AuthService{
		config:   config,
		userRepo: userRepo,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User *user.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        utils.NewUUID(),
			CreatedAt: time.Now(),
		},
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	createdUser.Password = ""

	return &RegisterResponse{User: createdUser}, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User   *user.User `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || existingUser == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.GenerateTokenPair(existingUser)
	if err != nil {
		return nil, err
	}

	existingUser.Password = ""

	return &LoginResponse{
		User:   existingUser,
		Tokens: tokens,
	}, nil
}

func (s *AuthService) GenerateTokenPair(u *user.User) (*TokenPair, error) {
	accessExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := &Claims{
		UserID: u.Id.String(),
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taskhub",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &RefreshClaims{
		UserID: u.Id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taskhub",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenPair, error) {
	token, err := jwt.ParseWithClaims(req.RefreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	u, err := s.userRepo.FindById(ctx, claims.UserID)
	if err != nil || u == nil {
		return nil, ErrInvalidToken
	}

	return s.GenerateTokenPair(u)
}
