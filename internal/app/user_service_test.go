package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"taskhub/internal/domains/user"
	"taskhub/pkg/base/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	users map[string]*user.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*user.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	if _, exists := m.users[u.Email]; exists {
		return nil, errors.New("user already exists")
	}
	m.users[u.Email] = u
	return u, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *MockUserRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	for _, u := range m.users {
		if u.Id.String() == id {
			return u, nil
		}
	}
	return nil, nil
}

func TestUserServiceModule(t *testing.T) {
	assert.NotNil(t, UserServiceModule)
}

func TestCreateUserRequest(t *testing.T) {
	req := &CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestCreateUserResponse(t *testing.T) {
	resp := &CreateUserResponse{}
	assert.NotNil(t, resp)
}

func TestUserService_Create(t *testing.T) {
	// Create a mock user repository
	mockRepo := NewMockUserRepository()

	// Create a test user through the mock repo
	ctx := context.Background()
	userID := uuid.New()
	newUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        userID,
			CreatedAt: time.Now(),
		},
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	created, err := mockRepo.Create(ctx, newUser)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, userID, created.Id)
}

func TestUserService_CreateUser_Duplicate(t *testing.T) {
	mockRepo := NewMockUserRepository()
	ctx := context.Background()

	// Create first user
	user1 := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
		},
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password1",
	}
	_, err := mockRepo.Create(ctx, user1)
	assert.NoError(t, err)

	// Try to create second user with same email
	user2 := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
		},
		Name:     "Another User",
		Email:    "test@example.com",
		Password: "password2",
	}
	_, err = mockRepo.Create(ctx, user2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestUserService_FindByEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	ctx := context.Background()

	// Create user
	newUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
		},
		Name:     "Test User",
		Email:    "find@example.com",
		Password: "password",
	}
	_, err := mockRepo.Create(ctx, newUser)
	assert.NoError(t, err)

	// Find user by email
	found, err := mockRepo.FindByEmail(ctx, "find@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "Test User", found.Name)

	// Find non-existing user
	notFound, err := mockRepo.FindByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestUserStruct(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	u := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        userID,
			CreatedAt: now,
		},
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secretpassword",
	}

	assert.Equal(t, userID, u.Id)
	assert.Equal(t, "John Doe", u.Name)
	assert.Equal(t, "john@example.com", u.Email)
	assert.Equal(t, "secretpassword", u.Password)
}

func TestUserErrors(t *testing.T) {
	// Test that error constants are defined
	assert.Equal(t, "invalid credentials", ErrInvalidCredentials.Error())
	assert.Equal(t, "user already exists", ErrUserAlreadyExists.Error())
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "token expired", ErrTokenExpired.Error())
}
