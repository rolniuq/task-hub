package nats

import (
	"taskhub/config"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn *nats.Conn
}

func NewNats(config *config.Config) *Nats {
	conn, err := nats.Connect(config.NatsUrl)
	if err != nil {
		return nil
	}

	return &Nats{
		conn: conn,
	}
}

func (n *Nats) Close() {
	if n == nil || n.conn == nil {
		return
	}

	n.conn.Close()
}
