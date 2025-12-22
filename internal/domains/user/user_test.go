package user

import (
	"testing"
	"time"

	"taskhub/pkg/base/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_GetId(t *testing.T) {
	id := uuid.New()
	user := &User{
		BaseEntity: entity.BaseEntity{
			Id: id,
		},
	}

	assert.Equal(t, id, user.GetId())
}

func TestUser_Fields(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	user := &User{
		BaseEntity: entity.BaseEntity{
			Id:        id,
			CreatedAt: now,
		},
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	assert.Equal(t, id, user.Id)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, now, user.CreatedAt)
}

func TestUser_EmptyFields(t *testing.T) {
	user := &User{}

	assert.Equal(t, uuid.Nil, user.Id)
	assert.Empty(t, user.Name)
	assert.Empty(t, user.Email)
	assert.Empty(t, user.Password)
}

func TestUser_WithBaseEntity(t *testing.T) {
	id := uuid.New()
	createdBy := uuid.New()
	now := time.Now()
	updateAt := now.Add(time.Hour)
	updateBy := uuid.New()

	user := &User{
		BaseEntity: entity.BaseEntity{
			Id:        id,
			CreatedAt: now,
			CreatedBy: createdBy,
			UpdateAt:  &updateAt,
			UpdateBy:  &updateBy,
		},
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	assert.Equal(t, id, user.Id)
	assert.Equal(t, createdBy, user.CreatedBy)
	assert.NotNil(t, user.UpdateAt)
	assert.Equal(t, updateAt, *user.UpdateAt)
	assert.NotNil(t, user.UpdateBy)
	assert.Equal(t, updateBy, *user.UpdateBy)
}

func TestUser_JSONTags(t *testing.T) {
	user := &User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "secret",
	}

	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "secret", user.Password)
}
