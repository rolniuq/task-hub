package repo

import (
	"context"
	"database/sql"
	"taskhub/config"
	"taskhub/internal/domains/user"
	"taskhub/pkg/db"
	"taskhub/pkg/logger"

	"github.com/google/uuid"
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
	query := `INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3) RETURNING id`

	var id uuid.UUID
	if err := r.conn.QueryRowContext(ctx, query, user.Id, user.Name, user.Email, user.Password).Scan(&id); err != nil {
		return nil, err
	}

	return user, nil
}
