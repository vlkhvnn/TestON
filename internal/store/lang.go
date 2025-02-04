package store

import (
	"context"
	"database/sql"
)

type LangStore struct {
	db *sql.DB
}

func (s *LangStore) SetUserLang(ctx context.Context, userID, lang string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO user_languages (user_id, lang)
	VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET lang = $2;
	`
	_, err := s.db.ExecContext(ctx, query, userID, lang)
	return err
}

func (s *LangStore) GetUserLang(ctx context.Context, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT lang FROM user_languages WHERE user_id = $1;`
	var lang string
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&lang)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return lang, nil
}
