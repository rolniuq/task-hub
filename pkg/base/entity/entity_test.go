package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBaseEntity_GetId(t *testing.T) {
	id := uuid.New()
	entity := &BaseEntity{
		Id: id,
	}

	assert.Equal(t, id, entity.GetId())
}

func TestBaseEntity_Fields(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	createdBy := uuid.New()
	updateAt := time.Now().Add(time.Hour)
	updateBy := uuid.New()
	deletedAt := time.Now().Add(2 * time.Hour)
	deletedBy := uuid.New()

	entity := &BaseEntity{
		Id:        id,
		CreatedAt: createdAt,
		CreatedBy: createdBy,
		UpdateAt:  &updateAt,
		UpdateBy:  &updateBy,
		DeletedAt: &deletedAt,
		DeletedBy: &deletedBy,
	}

	assert.Equal(t, id, entity.Id)
	assert.Equal(t, createdAt, entity.CreatedAt)
	assert.Equal(t, createdBy, entity.CreatedBy)
	assert.NotNil(t, entity.UpdateAt)
	assert.Equal(t, updateAt, *entity.UpdateAt)
	assert.NotNil(t, entity.UpdateBy)
	assert.Equal(t, updateBy, *entity.UpdateBy)
	assert.NotNil(t, entity.DeletedAt)
	assert.Equal(t, deletedAt, *entity.DeletedAt)
	assert.NotNil(t, entity.DeletedBy)
	assert.Equal(t, deletedBy, *entity.DeletedBy)
}

func TestBaseEntity_NilOptionalFields(t *testing.T) {
	entity := &BaseEntity{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}

	assert.Nil(t, entity.UpdateAt)
	assert.Nil(t, entity.UpdateBy)
	assert.Nil(t, entity.DeletedAt)
	assert.Nil(t, entity.DeletedBy)
}

func TestBaseEntity_EmptyEntity(t *testing.T) {
	entity := &BaseEntity{}

	assert.Equal(t, uuid.Nil, entity.Id)
	assert.True(t, entity.CreatedAt.IsZero())
	assert.Equal(t, uuid.Nil, entity.CreatedBy)
}
