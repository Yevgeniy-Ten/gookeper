package postgres

import (
	"context"
	"errors"

	"gokeeper/internal/models"
	pwd "gokeeper/internal/pkg/password"

	"database/sql"

	"github.com/lib/pq"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrNotFound = errors.New("not found")

func (p *Postgres) CreateUser(ctx context.Context, user *models.CreateUserDTO) (*models.User, error) {
	hashedPassword, err := pwd.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id, login`

	var createdUser models.User

	err = p.db.QueryRowContext(ctx, query,
		user.Login,
		hashedPassword,
	).Scan(
		&createdUser.ID,
		&createdUser.Login,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return &createdUser, nil
}

func (p *Postgres) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT id, login, password, created_at, updated_at
		FROM users
		WHERE login = $1`

	var user models.User
	err := p.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) VerifyUser(ctx context.Context, login, plainPassword string) (*models.User, error) {
	user, err := p.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !pwd.CheckPassword(plainPassword, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
