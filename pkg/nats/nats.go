package nats

import (
	"taskhub/config"
	"taskhub/pkg/logger"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(NewNats),
)

type Nats struct {
	conn   *nats.Conn
	logger *logger.Logger
}

func NewNats(config *config.Config, logger *logger.Logger) *Nats {
	conn, err := nats.Connect(config.NatsUrl)
	if err != nil {
		logger.Error("failed to connect to nats: %v", err)
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
