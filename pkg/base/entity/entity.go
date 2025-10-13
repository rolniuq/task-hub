package entity

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdateAt  *time.Time
	UpdateBy  *uuid.UUID
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
}

func (e *BaseEntity) GetId() uuid.UUID {
	return e.Id
}
