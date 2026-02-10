package repo

import (
	"context"
	"database/sql"
	"taskhub/config"
	"taskhub/internal/db"
	"taskhub/internal/domains/user"
	db2 "taskhub/pkg/db"
	"taskhub/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var UserRepositoryModule = fx.Module(
	"user-repo",
	fx.Provide(NewUserRepository),
)

type UserRepository struct {
	queries *db.Queries
	logger  *logger.Logger
}

func NewUserRepository(config *config.Config, logger *logger.Logger) *UserRepository {
	conn := db2.NewDB(config).GetConnection()
	queries := db.New(conn)
	return &UserRepository{
		queries: queries,
		logger:  logger,
	}
}

// Helper function to convert domain User to sqlc CreateUserParams
func toCreateUserParams(u *user.User) db.CreateUserParams {
	return db.CreateUserParams{
		ID:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

// Helper function to convert sqlc User to domain User
func toDomainUser(u db.User) *user.User {
	result := &user.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}

	// Set BaseEntity fields
	result.Id = u.ID
	result.CreatedAt = u.CreatedAt

	if u.UpdatedAt.Valid {
		result.UpdateAt = &u.UpdatedAt.Time
	}
	if u.UpdatedBy.Valid {
		result.UpdateBy = &u.UpdatedBy.UUID
	}

	return result
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	params := toCreateUserParams(u)

	created, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(created), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	found, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toDomainUser(found), nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	found, err := r.queries.GetUserByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toDomainUser(found), nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	params := db.UpdateUserParams{
		ID:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}

	if u.UpdateAt != nil {
		params.UpdatedAt = sql.NullTime{Time: *u.UpdateAt, Valid: true}
	}

	updated, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(updated), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
