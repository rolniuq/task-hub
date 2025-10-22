package repo

import (
	"context"
	"database/sql"
	"taskhub/config"
	"taskhub/internal/domains/user"
	"taskhub/pkg/db"
	"taskhub/pkg/logger"

	"go.uber.org/fx"
)

var UserRepositoryModule = fx.Module(
	"user-repo",
	fx.Provide(NewUserRepository),
)

type UserRepository struct {
	conn   *sql.DB
	logger *logger.Logger
}

func NewUserRepository(config *config.Config, logger *logger.Logger) *UserRepository {
	conn := db.NewDB(config).GetConnection()
	return &UserRepository{
		conn:   conn,
		logger: logger,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	return nil, nil
}
