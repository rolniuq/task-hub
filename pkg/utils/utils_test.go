package utils

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetContextKey(t *testing.T) {
	ctx := context.Background()
	value := "test-value"
	ctx = context.WithValue(ctx, "test-key", &value)

	result := GetContextKey[string](ctx, "test-key")
	assert.NotNil(t, result)
	assert.Equal(t, value, *result)
}

func TestGetContextKey_NotFound(t *testing.T) {
	ctx := context.Background()

	result := GetContextKey[string](ctx, "non-existent-key")
	assert.Nil(t, result)
}

func TestNewPointer(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"string", "test"},
		{"int", 42},
		{"bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case string:
				ptr := NewPointer(v)
				assert.NotNil(t, ptr)
				assert.Equal(t, v, *ptr)
			case int:
				ptr := NewPointer(v)
				assert.NotNil(t, ptr)
				assert.Equal(t, v, *ptr)
			case bool:
				ptr := NewPointer(v)
				assert.NotNil(t, ptr)
				assert.Equal(t, v, *ptr)
			}
		})
	}
}

func TestGetPointerValue(t *testing.T) {
	value := "test"
	ptr := &value

	result := GetPointerValue(ptr)
	assert.Equal(t, value, result)
}

func TestGetPointerValue_Nil(t *testing.T) {
	var ptr *string

	result := GetPointerValue(ptr)
	assert.Equal(t, "", result)
}

func TestGetPointerValue_NilInt(t *testing.T) {
	var ptr *int

	result := GetPointerValue(ptr)
	assert.Equal(t, 0, result)
}

func TestNewUUID(t *testing.T) {
	id := NewUUID()
	assert.NotEqual(t, uuid.Nil, id)
}

func TestNewUUID_Unique(t *testing.T) {
	id1 := NewUUID()
	id2 := NewUUID()
	assert.NotEqual(t, id1, id2)
}
