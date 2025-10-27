package utils

import (
	"context"

	"github.com/google/uuid"
)

func GetContextKey[T any](ctx context.Context, key string) *T {
	val := ctx.Value(key)
	if val == nil {
		return nil
	}

	return val.(*T)
}

func NewPointer[V any](v V) *V {
	return &v
}

func GetPointerValue[V any](v *V) V {
	var val V
	if v != nil {
		val = *v
	}

	return val
}

func NewUUID() uuid.UUID {
	return uuid.New()
}
