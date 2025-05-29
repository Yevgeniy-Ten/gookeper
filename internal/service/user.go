package service

import (
	"context"
	"errors"

	"gokeeper/internal/models"
	"gokeeper/internal/repository/postgres"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repo *postgres.Postgres
}

func NewUserService(repo *postgres.Postgres) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, input *models.CreateUserDTO) (*models.User, error) {
	existingUser, err := s.repo.GetUserByLogin(ctx, input.Login)
	if err != nil && !errors.Is(err, postgres.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err := s.repo.CreateUser(ctx, input)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, input *models.LoginUserDTO) (*models.User, error) {
	user, err := s.repo.VerifyUser(ctx, input.Login, input.Password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
