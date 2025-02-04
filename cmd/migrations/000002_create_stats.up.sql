CREATE TABLE IF NOT EXISTS stats (
    id SERIAL PRIMARY KEY,
    lang TEXT NOT NULL,
    date DATE NOT NULL,
    count INT NOT NULL DEFAULT 0,
    UNIQUE(lang, date)
);