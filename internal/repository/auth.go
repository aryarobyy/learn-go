package repository

import (
	"auth/internal/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthRepo interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	VerifyEmail(ctx context.Context, email string, id int) error
	VerifyUsername(ctx context.Context, username string, id int) error
}

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *authRepo {
	return &authRepo{db: db}
}

func (s *authRepo) Create(ctx context.Context, user model.User) (*model.User, error) {
	data := model.User{
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Role:     model.RoleUser,
	}

	query := `
		INSERT INTO users (name, username, password, role)
		VALUES (:name, :username, :password, :role)
		RETURNING id, name, username, role, password`

	rows, err := s.db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&data); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("insert succeeded but returned no rows")
	}

	return &data, nil
}

func (s *authRepo) VerifyEmail(ctx context.Context, email string, id int) error {
	var dest int

	query := `SELECT 1 
		FROM users
		WHERE
		email = $1 AND 
		id != $2`

	if err := s.db.GetContext(ctx,
		&dest,
		query,
		email, id); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	return nil
}

func (s *authRepo) VerifyUsername(ctx context.Context, username string, id int) error {
	var dest int

	query := `SELECT 1
		FROM users
		WHERE
		username = $1 AND 
		id != $2`

	if err := s.db.GetContext(ctx,
		&dest,
		query,
		username, id); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	return nil
}
