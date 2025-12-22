package main

import (
	"context"
	"taskhub/config"
	"taskhub/internal/app"
	taskrepo "taskhub/internal/domains/task/repo"
	userrepo "taskhub/internal/domains/user/repo"
	"taskhub/internal/gateway"
	"taskhub/pkg/logger"
	"taskhub/pkg/nats"

	"go.uber.org/fx"
)

func startApp(lc fx.Lifecycle, gw *gateway.Gateway) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go gw.Start()
			return nil
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
		userrepo.UserRepositoryModule,
		taskrepo.TaskRepositoryModule,
		app.AuthServiceModule,
		app.TaskServiceModule,
		gateway.GatewayModule,
		nats.NatsModule,
		fx.Invoke(startApp),
	)

	app.Run()
}
