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

// DSN for tests; adjust if needed.
var testDSN = "postgres://postgres:1234@localhost:5432/teston?sslmode=disable"

func initTestDB(t *testing.T, db *sql.DB) {
	eventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		event_id TEXT NOT NULL,
		lang TEXT NOT NULL,
		title TEXT NOT NULL,
		username TEXT NOT NULL,
		comment TEXT,
		timestamp BIGINT NOT NULL,
		wiki TEXT NOT NULL,
		server_name TEXT NOT NULL
	);
	`
	_, err := db.Exec(eventsTable)
	require.NoError(t, err, "failed to create events table")

	statsTable := `
	CREATE TABLE IF NOT EXISTS stats (
		id SERIAL PRIMARY KEY,
		lang TEXT NOT NULL,
		date DATE NOT NULL,
		count INT NOT NULL DEFAULT 0,
		UNIQUE(lang, date)
	);
	`
	_, err = db.Exec(statsTable)
	require.NoError(t, err, "failed to create stats table")

	userLangTable := `
	CREATE TABLE IF NOT EXISTS user_languages (
		user_id TEXT PRIMARY KEY,
		lang TEXT NOT NULL
	);
	`
	_, err = db.Exec(userLangTable)
	require.NoError(t, err, "failed to create user_languages table")
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", testDSN)
	require.NoError(t, err, "failed to connect to test database")

	initTestDB(t, db)

	cleanQueries := []string{
		"TRUNCATE TABLE events RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE stats RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE user_languages RESTART IDENTITY CASCADE;",
	}
	for _, q := range cleanQueries {
		_, err := db.Exec(q)
		require.NoError(t, err)
	}

	return db
}

func TestEventStore_AddAndGetRecent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	initTestDB(t, db)

	eventStore := &EventStore{db}
	ctx := context.Background()

	now := time.Now().Unix()
	event := &models.RecentChangeEvent{
		ID:         "1001",
		Type:       "edit",
		Title:      "Test Page",
		User:       "TestUser",
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
	initTestDB(t, db)

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
	initTestDB(t, db)

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
