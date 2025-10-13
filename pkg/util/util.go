package util

import "context"

func GetContextKey[T any](ctx context.Context, key string) *T {
	val := ctx.Value(key)
	if val == nil {
		return nil
	}

	return val.(*T)
}
