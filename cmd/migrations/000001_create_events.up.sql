CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_id TEXT NOT NULL,
    lang TEXT NOT NULL,
    title TEXT NOT NULL,
    username TEXT NOT NULL,
    comment TEXT,
    timestamp BIGINT NOT NULL
);