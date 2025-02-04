package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vlkhvnn/TestON/internal/models"
)

var (
	ErrNotFound          = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Event interface {
		Add(ctx context.Context, lang string, event *models.RecentChangeEvent) error
		GetRecent(ctx context.Context, lang string, limit int) ([]*models.RecentChangeEvent, error)
	}
	Stat interface {
		IncrementByLang(ctx context.Context, lang string, date string) error
		Get(ctx context.Context, lang string, date string) (int, error)
	}
	Lang interface {
		SetUserLang(ctx context.Context, userID, lang string) error
		GetUserLang(ctx context.Context, userID string) (string, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Event: &EventStore{db: db},
		Stat:  &StatStore{db: db},
		Lang:  &LangStore{db: db},
	}
}
