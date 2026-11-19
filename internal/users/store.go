package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type Store struct {
	DB *sql.DB
}

func (s *Store) Insert(ctx context.Context, newUser *User) error {
	query := `
		INSERT INTO users (full_name, email, username, password_hash)
		VALUES (?, ?, ?, ?)
	`
	args := []any{newUser.FullName, newUser.Email, newUser.Username, newUser.PasswordHash}
	res, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			switch {
			case mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users.username"):
				return ErrDuplicateUsername
			case mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users.email"):
				return ErrDuplicateEmail
			}
		}
		return err
	}

	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	lastInsertedUser, err := s.Get(ctx, int(lastInsertedID))
	if err != nil {
		return err
	}

	newUser.ID = lastInsertedUser.ID
	newUser.CreatedAt = lastInsertedUser.CreatedAt
	newUser.Version = lastInsertedUser.Version

	return nil
}

func (s *Store) Get(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, created_at, version, full_name, email, username, password_hash
		FROM users
		WHERE id = ?
	`
	args := []any{id}
	row := s.DB.QueryRowContext(ctx, query, args...)

	var user User
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
		&user.FullName,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, created_at, version, full_name, email, username, password_hash
		FROM users
		WHERE username = ?
	`
	args := []any{username}
	row := s.DB.QueryRowContext(ctx, query, args...)

	var user User
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
		&user.FullName,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
