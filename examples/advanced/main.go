package main

import (
	"log"
	"os"
	"regexp"
	"strings"
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
	
	// Create an advanced bot with custom storage and middleware
	bot := slackbot.New().
		WithClientCredentials(clientID, clientSecret).
		WithStorage(storage.NewMemoryStorage()).
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
			return ctx.Reply("Thanks for the DM!")
		}).
		OnChannelMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.React("👋")
		}).
		OnThreadMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.ReplyInThread("Thread reply!")
		}).
		OnMentionMessage(func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Thanks for mentioning me!")
		}).
		
		// Pattern-based handlers
		OnMessageContaining("hello", func(ctx *slackbot.MessageContext) error {
			return ctx.Reply("Hello there! 👋")
		}).
		OnMessageContaining("help", func(ctx *slackbot.MessageContext) error {
			helpMessage := ui.NewMessage().
				Section("*Available Commands:*").
				Section("• Say 'hello' for a greeting").
				Section("• Say 'help' for this message").
				Section("• Say 'weather' for weather info").
				ButtonRow(
					ui.Button("get_started", "Get Started").Primary(),
					ui.Button("contact_support", "Contact Support"),
				)
			
			options := helpMessage.Build()
			_, _, err := ctx.Client.PostMessage(ctx.ChannelID, options...)
			return err
		}).
		OnMessageMatching(regexp.MustCompile(`weather in (\w+)`), func(ctx *slackbot.MessageContext) error {
			// Extract city from message
			re := regexp.MustCompile(`weather in (\w+)`)
			matches := re.FindStringSubmatch(ctx.Text)
			if len(matches) > 1 {
				city := matches[1]
				return ctx.Reply("Weather in " + city + ": Sunny, 25°C 🌞")
			}
			return ctx.Reply("I couldn't understand which city you're asking about.")
		}).
		
		// Interactive handlers
		OnButtonClick("get_started", func(ctx *slackbot.InteractiveContext) error {
			modal := ui.NewModal("onboarding", "Welcome!").
				Section("Welcome to our bot! Let's get you set up.").
				Input("name", "Your Name:", ui.PlainText().WithPlaceholder("Enter your name")).
				Input("email", "Your Email:", ui.Email().WithPlaceholder("Enter your email")).
				SelectMenu("department", "Department:", []ui.Option{
					ui.NewOption("engineering", "Engineering"),
					ui.NewOption("design", "Design"),
					ui.NewOption("product", "Product"),
				}).
				Submit("Complete Setup").
				Cancel("Cancel")
			
			return ctx.OpenModal(modal.Build())
		}).
		OnButtonClick("contact_support", func(ctx *slackbot.InteractiveContext) error {
			return ctx.RespondEphemeral("Please email support@company.com for help!")
		}).
		OnModalSubmit("onboarding", func(ctx *slackbot.InteractiveContext) error {
			// Process the onboarding form
			name := "User" // In real implementation, extract from ctx.StateValues
			welcomeMsg := ui.WelcomeMessage(name)
			
			options := welcomeMsg.Build()
			_, _, err := ctx.Client.PostMessage(ctx.ChannelID, options...)
			return err
		}).
		
		// Custom home tab
		OnHomeOpened(func(ctx *slackbot.HomeContext) error {
			homeView := ui.NewHomeView().
				Header("🚀 Advanced Bot Dashboard").
				Section("Welcome to your personalized bot experience!").
				ButtonRow(
					ui.Button("quick_action", "Quick Action").Primary(),
					ui.Button("settings", "Settings"),
					ui.Button("help", "Help"),
				).
				Divider().
				Section("*Recent Activity*").
				Section("• Message processed at 2:30 PM").
				Section("• Button clicked at 2:25 PM").
				Section("• Home opened at 2:20 PM").
				Divider().
				Section("*Bot Statistics*").
				Section("📊 Messages processed: 42").
				Section("👥 Active users: 15").
				Section("⚡ Uptime: 99.9%")
			
			return ctx.PublishView(homeView.Build())
		}).
		
		// Slash commands
		OnSlashCommand("/status", func(ctx *slackbot.SlashCommandContext) error {
			statusMsg := ui.NewMessage().
				Section("*Bot Status* ✅").
				Section("• All systems operational").
				Section("• Response time: 42ms").
				Section("• Memory usage: 15MB").
				ButtonRow(
					ui.Button("refresh", "Refresh").Primary(),
					ui.Button("details", "View Details"),
				)
			
			options := statusMsg.Build()
			_, _, err := ctx.Client.PostMessage(ctx.ChannelID, options...)
			return err
		}).
		OnSlashCommand("/weather", func(ctx *slackbot.SlashCommandContext) error {
			location := strings.TrimSpace(ctx.Text)
			if location == "" {
				return ctx.RespondEphemeral("Please specify a location: /weather London")
			}
			
			weatherMsg := ui.NewMessage().
				Section("*Weather in " + location + "* 🌤️").
				Section("Temperature: 23°C").
				Section("Conditions: Partly cloudy").
				Section("Humidity: 65%").
				Section("Wind: 5 mph SW")
			
			options := weatherMsg.Build()
			_, _, err := ctx.Client.PostMessage(ctx.ChannelID, options...)
			return err
		}).
		
		// Global event handler for logging
		OnAnyEvent(func(ctx *slackbot.EventContext) error {
			log.Printf("Event received: %s from user %s", ctx.EventType, ctx.UserID)
			return nil
		})
	
	// Start the bot
	log.Println("Starting advanced bot with full features...")
	if err := bot.Start(":8080"); err != nil {
		log.Fatal("Failed to start bot:", err)
	}
}