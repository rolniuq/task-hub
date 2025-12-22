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

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	query := `INSERT INTO users (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id uuid.UUID
	if err := r.conn.QueryRowContext(ctx, query, u.Id, u.Name, u.Email, u.Password, u.CreatedAt).Scan(&id); err != nil {
		return nil, err
	}

	u.Id = id
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`

	var u user.User
	err := r.conn.QueryRowContext(ctx, query, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE id = $1`

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var u user.User
	err = r.conn.QueryRowContext(ctx, query, uid).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	query := `UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4`

	_, err := r.conn.ExecContext(ctx, query, u.Name, u.Email, u.UpdateAt, u.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.conn.ExecContext(ctx, query, id)
	return err
}
