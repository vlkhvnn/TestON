package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/vlkhvnn/TestON/internal/models"
)

type EventStore struct {
	db *sql.DB
}

func (s *EventStore) Add(ctx context.Context, lang string, event *models.RecentChangeEvent) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	eventID := event.ID.String()

	query := `
	INSERT INTO events (event_id, lang, title, username, comment, timestamp, wiki, server_name)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := s.db.ExecContext(ctx, query, eventID, lang, event.Title, event.User, event.Comment, event.Timestamp, event.Wiki, event.ServerName)
	if err != nil {
		return err
	}

	cleanupQuery := `
	DELETE FROM events WHERE event_id IN (
		SELECT event_id FROM events WHERE lang = $1
		ORDER BY timestamp DESC OFFSET 100
	);
	`
	_, err = s.db.ExecContext(ctx, cleanupQuery, lang)

	return err
}

func (s *EventStore) GetRecent(ctx context.Context, lang string, limit int) ([]*models.RecentChangeEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT event_id, title, username, comment, timestamp, wiki, server_name
	FROM events WHERE lang = $1
	ORDER BY timestamp DESC
	LIMIT $2;
	`
	rows, err := s.db.QueryContext(ctx, query, lang, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.RecentChangeEvent
	for rows.Next() {
		var e models.RecentChangeEvent
		var eventID string

		err := rows.Scan(&eventID, &e.Title, &e.User, &e.Comment, &e.Timestamp, &e.Wiki, &e.ServerName)
		if err != nil {
			return nil, err
		}

		e.ID = json.Number(eventID)
		events = append(events, &e)
	}

	if len(events) == 0 {
		return nil, ErrNotFound
	}

	return events, nil
}
