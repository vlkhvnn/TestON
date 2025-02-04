package main

import (
	"log"
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
	// Start the Wikimedia stream consumer in the background.
	go wikimedia.StartStream(&app.store)

	// Create and start the Discord bot.
	bot, err := discord.NewBot(app.config.token, app.store)
	if err != nil {
		return err
	}
	if err := bot.Start(); err != nil {
		return err
	}
	defer bot.Stop()
	log.Println("Bot is running. Press CTRL-C to exit.")

	// Wait for a termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down.")
	return nil
}
