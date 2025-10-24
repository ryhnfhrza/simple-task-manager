package repository

import (
	"context"
	"database/sql"

	"github.com/ryhnfhrza/simple-task-manager/model/domain"
)

type usersRepositoryImpl struct {
}

func NewUserRepository() UsersRepository {
	return &usersRepositoryImpl{}
}

func (repo *usersRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user *domain.User) error {

	query := "insert into users(username,password_hash) values(?,?)"
	_, err := tx.ExecContext(ctx, query, user.Username, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (repo *usersRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (*domain.User, error) {
	query := "select id,username,password_hash from users where username = ? "
	row := tx.QueryRowContext(ctx, query, username)

	user := &domain.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
