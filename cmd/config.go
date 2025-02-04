package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vlkhvnn/TestON/internal/discord"
	"github.com/vlkhvnn/TestON/internal/store"
	"github.com/vlkhvnn/TestON/internal/wikimedia"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	bot    discord.Bot
}

type config struct {
	token string
	db    dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := wikimedia.StartStream(ctx, &app.store, app.logger); err != nil {
			app.logger.Errorw("Failed to start Wikimedia stream", "error", err)
			cancel()
		}
	}()

	if err := app.bot.Start(); err != nil {
		return err
	}

	defer app.bot.Stop()

	app.logger.Info("Application started. Press CTRL-C to exit.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	app.logger.Info("Shutting down...")
	return nil
}
