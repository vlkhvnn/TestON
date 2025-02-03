package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/vlkhvnn/TestON/internal/discord"
	"github.com/vlkhvnn/TestON/internal/env"
	"github.com/vlkhvnn/TestON/internal/store"
	"github.com/vlkhvnn/TestON/internal/wikimedia"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	err := godotenv.Load("../../.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	token := env.GetString("DISCORD_TOKEN", "")

	// Initialize the in-memory event store.
	eventStore := store.NewStore()

	// Start the Wikimedia stream consumer in the background.
	go wikimedia.StartStream(eventStore)

	// Create and start the Discord bot.
	bot, err := discord.NewBot(token, eventStore)
	if err != nil {
		logger.Fatalf("Error creating Discord bot: %v", err)
	}
	if err := bot.Start(); err != nil {
		logger.Fatalf("Error starting Discord bot: %v", err)
	}
	defer bot.Stop()

	if err != nil {
		logger.Fatalf("Error creating a bot session %s", err)
	}
	log.Println("Bot is running. Press CTRL-C to exit.")

	// Wait for a termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down.")

}
