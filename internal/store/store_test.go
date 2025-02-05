package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlkhvnn/TestON/internal/models"
)

var testDSN = "postgres://postgres:1234@localhost:5432/teston?sslmode=disable"

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", testDSN)
	require.NoError(t, err)

	queries := []string{
		"TRUNCATE TABLE events RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE stats RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE user_languages RESTART IDENTITY CASCADE;",
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		require.NoError(t, err)
	}

	return db
}

func TestEventStore_AddAndGetRecent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	eventStore := &EventStore{db}
	ctx := context.Background()

	now := time.Now().Unix()
	event := &models.RecentChangeEvent{
		ID:         "1001",
		Type:       "edit",
		Title:      "Test Page",
		User:       "TestUser",
		Bot:        false,
		Minor:      false,
		Comment:    "Test comment",
		Timestamp:  now,
		Wiki:       "enwiki",
		ServerName: "en.wikipedia.org",
	}

	err := eventStore.Add(ctx, "en", event)
	require.NoError(t, err)

	events, err := eventStore.GetRecent(ctx, "en", 10)
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "Test Page", events[0].Title)
}

func TestStatStore_IncrementAndGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	statStore := &StatStore{db: db}
	ctx := context.Background()
	lang := "en"
	dateStr := time.Now().Format("2006-01-02")

	err := statStore.IncrementByLang(ctx, lang, dateStr)
	require.NoError(t, err)
	err = statStore.IncrementByLang(ctx, lang, dateStr)
	require.NoError(t, err)

	count, err := statStore.Get(ctx, lang, dateStr)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestLangStore_SetAndGetUserLang(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	langStore := &LangStore{db: db}
	ctx := context.Background()
	userID := "user123"
	expectedLang := "fr"

	err := langStore.SetUserLang(ctx, userID, expectedLang)
	require.NoError(t, err)

	lang, err := langStore.GetUserLang(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, expectedLang, lang)
}
