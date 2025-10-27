package main

import (
	"context"
	"taskhub/config"
	"taskhub/internal/gateway"
	"taskhub/pkg/logger"
	"taskhub/pkg/nats"

	"go.uber.org/fx"
)

func startApp(lc fx.Lifecycle, gw *gateway.Gateway) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return gw.Start()
		},
		OnStop: func(ctx context.Context) error {
			return gw.Shutdown(ctx)
		},
	})
}

func main() {
	app := fx.New(
		config.ConfigModule,
		logger.LoggerModule,
		gateway.GatewayModule,
		nats.NatsModule,
		fx.Invoke(startApp),
	)

	app.Run()
}
