package app

import (
	"context"
	"taskhub/internal/domains/user"
	"taskhub/internal/domains/user/repo"
	"taskhub/pkg/logger"

	"go.uber.org/fx"
)

var UserServiceModule = fx.Module(
	"user-service",
	fx.Provide(NewUserService),
)

type UserService struct {
	logger   *logger.Logger
	userRepo *repo.UserRepository
}

func NewUserService(logger *logger.Logger, userRepo *repo.UserRepository) *UserService {
	return &UserService{logger: logger, userRepo: userRepo}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
}

func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	_, err := s.userRepo.Create(ctx, &user.User{})
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{}, nil
}
