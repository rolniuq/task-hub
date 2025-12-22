package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNats_Close_Nil(t *testing.T) {
	var n *Nats
	n.Close()
}

func TestNats_Close_NilConn(t *testing.T) {
	n := &Nats{conn: nil}
	n.Close()
}

func TestNats_Publish_Nil(t *testing.T) {
	var n *Nats
	err := n.Publish("test", []byte("data"))
	assert.NoError(t, err)
}

func TestNats_Publish_NilConn(t *testing.T) {
	n := &Nats{conn: nil}
	err := n.Publish("test", []byte("data"))
	assert.NoError(t, err)
}

func TestNats_Subscribe_Nil(t *testing.T) {
	var n *Nats
	err := n.Subscribe("test", func(data []byte) {})
	assert.NoError(t, err)
}

func TestNats_Subscribe_NilConn(t *testing.T) {
	n := &Nats{conn: nil}
	err := n.Subscribe("test", func(data []byte) {})
	assert.NoError(t, err)
}

func TestNats_QueueSubscribe_Nil(t *testing.T) {
	var n *Nats
	err := n.QueueSubscribe("test", "queue", func(data []byte) {})
	assert.NoError(t, err)
}

func TestNats_QueueSubscribe_NilConn(t *testing.T) {
	n := &Nats{conn: nil}
	err := n.QueueSubscribe("test", "queue", func(data []byte) {})
	assert.NoError(t, err)
}

func TestNats_IsConnected_Nil(t *testing.T) {
	var n *Nats
	assert.False(t, n.IsConnected())
}

func TestNats_IsConnected_NilConn(t *testing.T) {
	n := &Nats{conn: nil}
	assert.False(t, n.IsConnected())
}

func TestNatsModule(t *testing.T) {
	assert.NotNil(t, NatsModule)
}
