package repository

import (
	"context"
	"database/sql"

	"github.com/ryhnfhrza/simple-task-manager/model/domain"
)

type UsersRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user *domain.User) error
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (*domain.User, error)
}
