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
	url := config.NatsUrl
	if url == "" {
		url = nats.DefaultURL
	}

	conn, err := nats.Connect(url)
	if err != nil {
		logger.Error("failed to connect to nats: %v", err)
		return nil
	}

	return &Nats{
		conn:   conn,
		logger: logger,
	}
}

func (n *Nats) Close() {
	if n == nil || n.conn == nil {
		return
	}

	n.conn.Close()
}

func (n *Nats) Publish(subject string, data []byte) error {
	if n == nil || n.conn == nil {
		return nil
	}

	return n.conn.Publish(subject, data)
}

func (n *Nats) Subscribe(subject string, handler func([]byte)) error {
	if n == nil || n.conn == nil {
		return nil
	}

	_, err := n.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})

	return err
}

func (n *Nats) QueueSubscribe(subject, queue string, handler func([]byte)) error {
	if n == nil || n.conn == nil {
		return nil
	}

	_, err := n.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		handler(msg.Data)
	})

	return err
}

func (n *Nats) IsConnected() bool {
	return n != nil && n.conn != nil && n.conn.IsConnected()
}
