package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	id := NewUUID()
	assert.NotEqual(t, uuid.Nil, id)
}

func TestNewUUID_Unique(t *testing.T) {
	id1 := NewUUID()
	id2 := NewUUID()
	assert.NotEqual(t, id1, id2)
}
