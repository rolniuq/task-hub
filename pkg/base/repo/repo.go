package repo

import "context"

type BaseRepository[T any] interface {
	FindById(ctx context.Context, id string) (*T, error)
}

type baseRepository[T any] struct {
}

func NewBaseRepository[T any]() BaseRepository[T] {
	return &baseRepository[T]{}
}

func (b *baseRepository[T]) FindById(ctx context.Context, id string) (*T, error) {
	return nil, nil
}
