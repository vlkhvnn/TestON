package store

import (
	"context"
	"database/sql"
)

type StatStore struct {
	db *sql.DB
}

func (s *StatStore) IncrementByLang(ctx context.Context, lang string, date string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO stats (lang, date, count)
	VALUES ($1, $2, 1)
	ON CONFLICT (lang, date) DO UPDATE
	SET count = stats.count + 1;
	`
	_, err := s.db.ExecContext(ctx, query, lang, date)
	return err
}

func (s *StatStore) Get(ctx context.Context, lang string, date string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT count FROM stats WHERE lang = $1 AND date = $2;`
	var count int
	err := s.db.QueryRowContext(ctx, query, lang, date).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return count, nil
}
