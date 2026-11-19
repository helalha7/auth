package tokens

import (
	"context"
	"database/sql"
)

type Store struct {
	DB *sql.DB
}

func (s *Store) Insert(ctx context.Context, token *Token) error {
	query := `
		INSERT INTO tokens (user_id, hash, expiry, scope)
		VALUES (?, ?, ?, ?)
	`
	args := []any{token.UserID, token.Hash, token.Expiry, token.Scope}
	_, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Get(ctx context.Context, hash string) (*Token, error) {
	query := `
		SELECT id, user_id, expiry, scope
		FROM tokens
		WHERE hash = ?
	`
	args := []any{hash}
	row := s.DB.QueryRowContext(ctx, query, args...)

	var token Token
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Expiry,
		&token.Scope,
	)
	if err != nil {
		return nil, err
	}

	token.Hash = []byte(hash)

	return &token, nil
}
