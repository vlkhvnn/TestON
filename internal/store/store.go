package store

import (
	"sync"
	"time"

	"github.com/vlkhvnn/TestON/internal/models"
)

// Store holds recent Wikimedia events and daily stats.
// - recentChanges: language → list of events (up to 10)
// - dailyStats: language → (date string "yyyy-mm-dd" → count)
type Store struct {
	mu            sync.Mutex
	recentChanges map[string][]models.RecentChangeEvent
	dailyStats    map[string]map[string]int
}

// NewStore creates and returns a new store instance.
func NewStore() *Store {
	return &Store{
		recentChanges: make(map[string][]models.RecentChangeEvent),
		dailyStats:    make(map[string]map[string]int),
	}
}

// AddEvent adds a new event for the given language.
func (s *Store) AddEvent(lang string, event models.RecentChangeEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prepend the event to the list.
	events := s.recentChanges[lang]
	events = append([]models.RecentChangeEvent{event}, events...)
	if len(events) > 10 {
		events = events[:10]
	}
	s.recentChanges[lang] = events

	// Update daily stats.
	t := time.Unix(event.Timestamp, 0)
	dateStr := t.Format("2006-01-02")
	if s.dailyStats[lang] == nil {
		s.dailyStats[lang] = make(map[string]int)
	}
	s.dailyStats[lang][dateStr]++
}

// GetEvents returns the recent events for the given language.
func (s *Store) GetEvents(lang string) []models.RecentChangeEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.recentChanges[lang]
}

// GetStats returns the number of events for the given language on the given date (format: yyyy-mm-dd).
func (s *Store) GetStats(lang, date string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dailyStats[lang] == nil {
		return 0
	}
	return s.dailyStats[lang][date]
}
