package repository

import (
	"context"
	"github.com/amakmurr/dans-multi-pro-test/internal/core"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) core.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*core.User, error) {
	var user core.User
	err := u.db.GetContext(ctx, &user, `SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
