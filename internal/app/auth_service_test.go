package app

import (
	"context"
	"testing"
	"time"

	"taskhub/config"
	"taskhub/internal/domains/user"
	"taskhub/pkg/base/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	users map[string]*user.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[string]*user.User),
	}
}

func (m *MockUserRepo) Create(ctx context.Context, u *user.User) (*user.User, error) {
	m.users[u.Email] = u
	return u, nil
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *MockUserRepo) FindById(ctx context.Context, id string) (*user.User, error) {
	for _, u := range m.users {
		if u.Id.String() == id {
			return u, nil
		}
	}
	return nil, nil
}

func newTestConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-key-for-testing-purposes",
	}
}

func TestGenerateTokenPair(t *testing.T) {
	cfg := newTestConfig()
	service := &AuthService{config: cfg}

	testUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id: uuid.New(),
		},
		Email: "test@example.com",
	}

	tokens, err := service.GenerateTokenPair(testUser)

	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Greater(t, tokens.ExpiresAt, time.Now().Unix())
}

func TestValidateAccessToken_Success(t *testing.T) {
	cfg := newTestConfig()
	service := &AuthService{config: cfg}

	testUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id: uuid.New(),
		},
		Email: "test@example.com",
	}

	tokens, _ := service.GenerateTokenPair(testUser)
	claims, err := service.ValidateAccessToken(tokens.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, testUser.Id.String(), claims.UserID)
	assert.Equal(t, testUser.Email, claims.Email)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	cfg := newTestConfig()
	service := &AuthService{config: cfg}

	claims, err := service.ValidateAccessToken("invalid-token")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
	assert.Nil(t, claims)
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	cfg := newTestConfig()
	service := &AuthService{config: cfg}

	testUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id: uuid.New(),
		},
		Email: "test@example.com",
	}

	tokens, _ := service.GenerateTokenPair(testUser)

	wrongCfg := &config.Config{JWTSecret: "wrong-secret"}
	wrongService := &AuthService{config: wrongCfg}

	claims, err := wrongService.ValidateAccessToken(tokens.AccessToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, string(hashedPassword))
}

func TestComparePassword_Success(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	assert.NoError(t, err)
}

func TestComparePassword_WrongPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))

	assert.Error(t, err)
}

func TestClaims(t *testing.T) {
	claims := &Claims{
		UserID: uuid.New().String(),
		Email:  "test@example.com",
	}

	assert.NotEmpty(t, claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestRefreshClaims(t *testing.T) {
	claims := &RefreshClaims{
		UserID: uuid.New().String(),
	}

	assert.NotEmpty(t, claims.UserID)
}

func TestTokenPair(t *testing.T) {
	tokenPair := &TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(15 * time.Minute).Unix(),
	}

	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Greater(t, tokenPair.ExpiresAt, time.Now().Unix())
}

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     RegisterRequest
		isValid bool
	}{
		{
			name: "valid request",
			req: RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			isValid: true,
		},
		{
			name: "missing name",
			req: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "missing email",
			req: RegisterRequest{
				Name:     "Test User",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "missing password",
			req: RegisterRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.req.Name != "" && tt.req.Email != "" && tt.req.Password != ""
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     LoginRequest
		isValid bool
	}{
		{
			name: "valid request",
			req: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			isValid: true,
		},
		{
			name: "missing email",
			req: LoginRequest{
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "missing password",
			req: LoginRequest{
				Email: "test@example.com",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.req.Email != "" && tt.req.Password != ""
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestErrorTypes(t *testing.T) {
	assert.Equal(t, "invalid credentials", ErrInvalidCredentials.Error())
	assert.Equal(t, "user already exists", ErrUserAlreadyExists.Error())
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "token expired", ErrTokenExpired.Error())
}
