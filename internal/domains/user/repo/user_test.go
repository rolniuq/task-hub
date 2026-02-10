package repo

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
	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	for _, u := range m.users {
		if u.Id.String() == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	for email, existing := range m.users {
		if existing.Id == u.Id {
			m.users[email] = u
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for email, u := range m.users {
		if u.Id == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestUserRepositoryModule(t *testing.T) {
	assert.NotNil(t, UserRepositoryModule)
}

func TestMockUserRepository_Create(t *testing.T) {
	repo := NewMockUserRepository()
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

	created, err := repo.Create(ctx, newUser)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, userID, created.Id)
	assert.Equal(t, "test@example.com", created.Email)
}

func TestMockUserRepository_FindByEmail(t *testing.T) {
	repo := NewMockUserRepository()
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

	repo.Create(ctx, newUser)

	// Test find existing
	found, err := repo.FindByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "Test User", found.Name)

	// Test find non-existing
	notFound, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestMockUserRepository_FindById(t *testing.T) {
	repo := NewMockUserRepository()
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

	repo.Create(ctx, newUser)

	// Test find by ID
	found, err := repo.FindById(ctx, userID.String())
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "Test User", found.Name)

	// Test find non-existing
	notFound, err := repo.FindById(ctx, uuid.New().String())
	assert.NoError(t, err)
	assert.Nil(t, notFound)

	// Test find with invalid UUID
	invalid, err := repo.FindById(ctx, "invalid-uuid")
	assert.Error(t, err)
	assert.Nil(t, invalid)
}

func TestMockUserRepository_Update(t *testing.T) {
	repo := NewMockUserRepository()
	ctx := context.Background()

	userID := uuid.New()
	newUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        userID,
			CreatedAt: time.Now(),
		},
		Name:     "Original Name",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	repo.Create(ctx, newUser)

	// Update user
	updateAt := time.Now()
	updatedUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        userID,
			CreatedAt: newUser.CreatedAt,
			UpdateAt:  &updateAt,
		},
		Name:     "Updated Name",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	updated, err := repo.Update(ctx, updatedUser)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)

	// Test update non-existing
	nonExisting := &user.User{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Name:       "Non Existing",
	}
	_, err = repo.Update(ctx, nonExisting)
	assert.Error(t, err)
}

func TestMockUserRepository_Delete(t *testing.T) {
	repo := NewMockUserRepository()
	ctx := context.Background()

	userID := uuid.New()
	newUser := &user.User{
		BaseEntity: entity.BaseEntity{
			Id:        userID,
			CreatedAt: time.Now(),
		},
		Name:     "User to Delete",
		Email:    "delete@example.com",
		Password: "hashedpassword",
	}

	repo.Create(ctx, newUser)

	// Verify user exists
	found, _ := repo.FindById(ctx, userID.String())
	assert.NotNil(t, found)

	// Delete user
	err := repo.Delete(ctx, userID)
	assert.NoError(t, err)

	// Verify user is deleted
	notFound, _ := repo.FindById(ctx, userID.String())
	assert.Nil(t, notFound)

	// Test delete non-existing
	err = repo.Delete(ctx, uuid.New())
	assert.Error(t, err)
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
