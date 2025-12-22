package main

import (
	"taskhub/config"
	"taskhub/internal/app"
	"taskhub/internal/desktop"
	taskrepo "taskhub/internal/domains/task/repo"
	userrepo "taskhub/internal/domains/user/repo"
	"taskhub/pkg/logger"
	"taskhub/pkg/nats"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.ConfigModule,
		logger.LoggerModule,
		userrepo.UserRepositoryModule,
		taskrepo.TaskRepositoryModule,
		app.AuthServiceModule,
		app.TaskServiceModule,
		nats.NatsModule,
		fx.Provide(desktop.NewApp),
		fx.Invoke(desktop.RunDesktopApp),
	)

	app.Run()
}
