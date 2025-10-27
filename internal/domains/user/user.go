package user

import (
	"taskhub/pkg/base/entity"

	"github.com/google/uuid"
)

type User struct {
	entity.BaseEntity
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) GetId() uuid.UUID {
	return u.Id
}
