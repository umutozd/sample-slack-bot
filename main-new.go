package main

import (
	"os"
	"log"
	"time"

	"github.com/umutozd/sample-slack-bot/pkg/slackbot"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/ui"
)

func main() {
	// Get environment variables
	clientID := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		log.Fatal("SLACK_CLIENT_ID and SLACK_CLIENT_SECRET environment variables are required")
	}
	
	// Create a new bot using the plug-and-play library
	bot := slackbot.New().
		WithClientCredentials(clientID, clientSecret).
		WithStorage(storage.NewMemoryStorage()). // Uses in-memory storage by default
		WithMiddleware(
			slackbot.LoggingMiddleware(),
			slackbot.RateLimitMiddleware(60, time.Minute),
			slackbot.ErrorRecoveryMiddleware(),
		).
		WithDebug(true).
		
		// Message handlers
		OnMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Got your message: " + ctx.Text)
		}).
		OnDirectMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Thanks for the direct message!")
		}).
		OnMessageContaining("hello", func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Hello there! 👋")
		}).
		OnMessageContaining("help", func(ctx *slackbot.MessageContext) error {
			helpMessage := ui.NewMessage().
				Section("*Available Commands:*").
				Section("• Say 'hello' for a greeting").
				Section("• Say 'help' for this message").
				ButtonRow(
					ui.Button("get_started", "Get Started").Primary(),
					ui.Button("contact_support", "Contact Support"),
				)
			
			options := helpMessage.Build()
			_, _, err := ctx.Client.PostMessage(ctx.ChannelID, options...)
			return err
		}).
		
		// Interactive handlers
		OnButtonClick("get_started", func(ctx *slackbot.InteractiveContext) error {
			return ctx.RespondEphemeral("Welcome! This is the new SlackBot library in action! 🎉")
		}).
		OnButtonClick("contact_support", func(ctx *slackbot.InteractiveContext) error {
			return ctx.RespondEphemeral("Please email support@company.com for help!")
		}).
		
		// Custom home tab (replaces the default)
		OnHomeOpened(func(ctx *slackbot.HomeContext) error {
			homeView := ui.NewHomeView().
				Header("🎉 Welcome to the New SlackBot Library!").
				Section("This bot is now powered by the new plug-and-play library.").
				Section("*What's New:*").
				Section("• 🚀 Material UI-style developer experience").
				Section("• 💾 Automatic in-memory storage").
				Section("• 🔧 Built-in middleware system").
				Section("• 🎨 Rich UI builder components").
				Section("• 📝 Fluent configuration API").
				ButtonRow(
					ui.Button("explore", "Explore Features").Primary(),
					ui.Button("docs", "View Documentation"),
				).
				Divider().
				Section("*Quick Start:*").
				Section("```go\nbot := slackbot.New().\n    WithClientCredentials(id, secret).\n    OnMessage(handler).\n    Start(\":8080\")\n```")
			
			return ctx.PublishView(homeView.Build())
		}).
		OnButtonClick("explore", func(ctx *slackbot.InteractiveContext) error {
			modal := ui.NewModal("features", "Library Features").
				Section("🎯 **Core Features:**").
				Section("• Builder pattern API").
				Section("• Default behaviors & fallbacks").
				Section("• Comprehensive event handling").
				Section("• Rich UI components").
				Section("• Multiple storage backends").
				Section("• Middleware system").
				Section("• Interactive workflows").
				Divider().
				Section("🛠️ **Developer Experience:**").
				Section("• 5-line minimum for working bot").
				Section("• Sensible defaults").
				Section("• Type-safe handlers").
				Section("• Extensive documentation").
				Section("• Production-ready patterns").
				Cancel("Close")
			
			return ctx.OpenModal(modal.Build())
		}).
		OnButtonClick("docs", func(ctx *slackbot.InteractiveContext) error {
			return ctx.RespondEphemeral("📚 Documentation: Check out the examples/ directory and README.md for comprehensive guides!")
		})
	
	// Start the bot
	log.Println("🚀 Starting SlackBot with the new plug-and-play library...")
	log.Println("🏠 Default home tab will show unless customized")
	log.Println("💾 Using in-memory storage by default")
	log.Println("🔧 Middleware stack includes logging and error recovery")
	log.Println("🎨 UI builder available for rich interactions")
	
	if err := bot.Start(":8080"); err != nil {
		log.Fatal("Failed to start bot:", err)
	}
}