package main

import (
	"github.com/joho/godotenv"
	"github.com/vlkhvnn/TestON/internal/db"
	"github.com/vlkhvnn/TestON/internal/discord"
	"github.com/vlkhvnn/TestON/internal/env"
	"github.com/vlkhvnn/TestON/internal/store"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	cfg := config{
		token: env.GetString("DISCORD_TOKEN", ""),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:1234@localhost:5432/teston?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatalf("DB connection pool error: %v", err)
	}

	defer db.Close()
	logger.Info("DB connection pool established")

	store := store.NewStorage(db)

	bot, err := discord.NewBot(cfg.token, store)
	if err != nil {
		logger.Fatalf("Error starting discord bot: %v", err)
	}

	app := application{
		config: cfg,
		store:  store,
		logger: logger,
		bot:    *bot,
	}

	logger.Fatal(app.run())
}
