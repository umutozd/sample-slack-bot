package main

import (
	"log"
	"os"

	"github.com/umutozd/sample-slack-bot/pkg/slackbot"
)

func main() {
	// Create a simple bot with just 5 lines of code!
	bot := slackbot.New().
		WithClientCredentials(os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET")).
		OnMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Hello " + ctx.UserID + "! You said: " + ctx.Text)
		}).
		OnButtonClick("hello", func(ctx *slackbot.InteractiveContext) error {
			return ctx.RespondEphemeral("Hello button was clicked!")
		})

	log.Println("🚀 Simple Slack Bot starting on port 8080...")
	log.Println("📝 Set environment variables:")
	log.Println("   SLACK_CLIENT_ID=your_client_id")
	log.Println("   SLACK_CLIENT_SECRET=your_client_secret")
	log.Println("🌐 Install URL: http://localhost:8080/slack/install")
	
	if err := bot.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}