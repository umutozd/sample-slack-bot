package main

import (
	"os"
	"log"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot"
)

func main() {
	// Get environment variables
	clientID := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		log.Fatal("SLACK_CLIENT_ID and SLACK_CLIENT_SECRET environment variables are required")
	}
	
	// Create a simple bot - this is all you need!
	bot := slackbot.New().
		WithClientCredentials(clientID, clientSecret).
		WithDebug(true).
		OnMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Hello " + ctx.UserID + "! You said: " + ctx.Text)
		}).
		OnDirectMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Thanks for the direct message!")
		}).
		OnButtonClick("hello", func(ctx *slackbot.InteractiveContext) error {
			return ctx.UpdateMessage("Button clicked! 🎉")
		})
	
	// Start the bot
	log.Println("Starting simple bot...")
	if err := bot.Start(":8080"); err != nil {
		log.Fatal("Failed to start bot:", err)
	}
}