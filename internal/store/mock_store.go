package store

import (
	"context"

	"github.com/vlkhvnn/TestON/internal/models"
)

type MockEventStore struct {
	RecentEvents []*models.RecentChangeEvent
}

func (m *MockEventStore) Add(ctx context.Context, lang string, event *models.RecentChangeEvent) error {
	m.RecentEvents = append(m.RecentEvents, event)
	return nil
}

func (m *MockEventStore) GetRecent(ctx context.Context, lang string, limit int) ([]*models.RecentChangeEvent, error) {
	if len(m.RecentEvents) == 0 {
		return nil, ErrNotFound
	}
	if limit > len(m.RecentEvents) {
		limit = len(m.RecentEvents)
	}
	return m.RecentEvents[:limit], nil
}

type MockLangStore struct {
	Langs map[string]string
}

func (m *MockLangStore) SetUserLang(ctx context.Context, userID, lang string) error {
	if m.Langs == nil {
		m.Langs = make(map[string]string)
	}
	m.Langs[userID] = lang
	return nil
}

func (m *MockLangStore) GetUserLang(ctx context.Context, userID string) (string, error) {
	lang, ok := m.Langs[userID]
	if !ok {
		return "", ErrNotFound
	}
	return lang, nil
}

type MockStatStore struct {
	Stats map[string]int
}

func (m *MockStatStore) IncrementByLang(ctx context.Context, lang string, date string) error {
	key := lang + "_" + date
	if m.Stats == nil {
		m.Stats = make(map[string]int)
	}
	m.Stats[key]++
	return nil
}

func (m *MockStatStore) Get(ctx context.Context, lang string, date string) (int, error) {
	key := lang + "_" + date
	if m.Stats == nil {
		return 0, ErrNotFound
	}
	count, ok := m.Stats[key]
	if !ok {
		return 0, ErrNotFound
	}
	return count, nil
}
