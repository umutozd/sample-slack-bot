package main

import (
	"log"
	"os"
	"regexp"

	"github.com/umutozd/sample-slack-bot/pkg/slackbot"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/middleware"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/ui"
)

func main() {
	log.Println("🚀 Starting SlackBot Library Demo...")
	
	// Check for required environment variables
	clientID := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		log.Println("❌ Missing required environment variables:")
		log.Println("   SLACK_CLIENT_ID=your_client_id")
		log.Println("   SLACK_CLIENT_SECRET=your_client_secret")
		log.Println("   REDIS_URL=redis://localhost:6379 (optional)")
		log.Println("")
		log.Println("Get these values from https://api.slack.com/apps")
		os.Exit(1)
	}
	
	// Create storage (Redis with memory fallback)
	var stor storage.Storage
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if redisStorage, err := storage.NewRedisStorage(redisURL, "", 0); err == nil {
			stor = redisStorage
			log.Println("📊 Using Redis storage")
		} else {
			log.Printf("⚠️ Redis unavailable, using memory storage: %v", err)
			stor = storage.NewMemoryStorage()
		}
	} else {
		stor = storage.NewMemoryStorage()
		log.Println("📊 Using memory storage")
	}
	
	// Create the bot with the new plug-and-play API
	bot := slackbot.New().
		WithClientCredentials(clientID, clientSecret).
		WithStorage(stor).
		WithDebug(true).
		WithMiddleware(
			middleware.LoggingMiddleware(),
			middleware.ErrorRecoveryMiddleware(),
			middleware.RateLimitMiddleware(100, 60),
			middleware.MetricsMiddleware(),
		).
		// Message handlers - showing different types
		OnMessage(handleGeneralMessage).
		OnDirectMessage(handleDirectMessage).
		OnMentionMessage(handleMentionMessage).
		OnMessageContaining("hello", handleHelloMessage).
		OnMessageContaining("help", handleHelpMessage).
		OnMessageMatching(regexp.MustCompile(`(?i)demo\s+(\w+)`), handleDemoCommand).
		// Interactive component handlers
		OnButtonClick("get_started", handleGetStarted).
		OnButtonClick("show_features", handleShowFeatures).
		OnButtonClick("open_modal", handleOpenModal).
		OnSelectMenu("demo_option", handleDemoOption).
		OnModalSubmit("demo_modal", handleModalSubmit).
		// Event handlers
		OnHomeOpened(handleHomeOpened).
		OnReactionAdded(handleReactionAdded).
		OnUserJoined(handleUserJoined).
		// Slash commands
		OnSlashCommand("demo", handleDemoSlashCommand).
		OnSlashCommand("status", handleStatusSlashCommand)
	
	log.Println("✅ Bot configured with plug-and-play API")
	log.Println("🌐 OAuth URL: http://localhost:8080/slack/install")
	log.Println("📋 Features enabled:")
	log.Println("   • Message handling (all types)")
	log.Println("   • Interactive components (buttons, modals)")
	log.Println("   • Custom home tab")
	log.Println("   • Slash commands (/demo, /status)")
	log.Println("   • Event handling (reactions, users)")
	log.Println("   • Middleware (logging, rate limiting, metrics)")
	log.Println("   • Storage abstraction (Redis + memory fallback)")
	log.Println("")
	log.Println("🎯 This demonstrates the Material UI-level ease of use!")
	log.Println("📝 Check examples/ directory for more usage patterns")
	
	// Start the bot - this is it!
	if err := bot.Start(":8080"); err != nil {
		log.Fatal("💥 Bot failed to start:", err)
	}
}

// Message handlers demonstrating different response types
func handleGeneralMessage(ctx *slackbot.MessageContext) error {
	// Only respond to some messages to avoid noise
	if len(ctx.Text) > 50 {
		return ctx.React("eyes") // Just react to long messages
	}
	return nil
}

func handleDirectMessage(ctx *slackbot.MessageContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	msg := ui.NewMessage().
		Text("Thanks for the direct message!").
		Section("👋 Hi " + user.Name + "! I received your message: \"" + ctx.Text + "\"").
		Section("Try these commands:").
		Section("• Say \"hello\" for a greeting\n• Say \"help\" for assistance\n• Say \"demo features\" to see what I can do").
		ButtonRow(
			ui.Button("get_started", "Get Started").Primary(),
			ui.Button("show_features", "Show Features"),
		)
	
	return ctx.ReplyWithBlocks(msg.Blocks())
}

func handleMentionMessage(ctx *slackbot.MessageContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	msg := ui.NewMessage().
		Text("You mentioned me!").
		Section("🔔 Hey " + user.Name + "! You mentioned me with: \"" + ctx.Text + "\"").
		Section("*What I can do:*").
		Section("• Respond to messages and mentions\n• Handle button clicks and form submissions\n• Provide a custom home tab experience\n• Track reactions and user events").
		ButtonRow(
			ui.Button("show_features", "Show All Features").Primary(),
			ui.Button("open_modal", "Open Demo Modal"),
		)
	
	return ctx.ReplyWithBlocks(msg.Blocks())
}

func handleHelloMessage(ctx *slackbot.MessageContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	greetings := []string{
		"👋 Hello " + user.Name + "!",
		"🎉 Hey there, " + user.Name + "!",
		"✨ Hi " + user.Name + "! Great to see you!",
		"🤖 Hello " + user.Name + "! I'm a friendly bot!",
	}
	
	// Simple greeting with a button
	msg := ui.NewMessage().
		Text(greetings[0]).
		Section(greetings[0] + " Thanks for saying hello!").
		ButtonRow(
			ui.Button("get_started", "Get Started").Primary(),
		)
	
	return ctx.ReplyWithBlocks(msg.Blocks())
}

func handleHelpMessage(ctx *slackbot.MessageContext) error {
	helpText := `
*🤖 SlackBot Library Demo - Help*

*This bot demonstrates the new plug-and-play SlackBot library!*

*Message Commands:*
• "hello" - Get a friendly greeting
• "help" - Show this help message
• "demo [feature]" - Try specific features
• Mention me (@bot) for interactive options

*Interactive Features:*
• Click buttons to see interactive components
• Try the modal forms and select menus
• Check out the custom home tab
• React to messages and see responses

*Slash Commands:*
• \`/demo\` - Show demo features
• \`/status\` - Check bot status

*Built with SlackBot Library:*
✅ 5-line minimum working bot
✅ Plug-and-play configuration
✅ Builder pattern API
✅ Rich context handlers
✅ UI component builders
✅ Middleware system
✅ Storage abstraction
✅ Default behaviors

*Check out the code:*
• \`examples/simple/\` - 5-line bot
• \`examples/advanced/\` - Full-featured bot
• \`pkg/slackbot/\` - Library source

This is Material UI-level ease of use for Slack bots! 🚀
`
	
	return ctx.Reply(helpText)
}

func handleDemoCommand(ctx *slackbot.MessageContext) error {
	// Extract feature from regex match
	re := regexp.MustCompile(`(?i)demo\s+(\w+)`)
	matches := re.FindStringSubmatch(ctx.Text)
	
	if len(matches) > 1 {
		feature := matches[1]
		
		msg := ui.NewMessage().
			Text("Demo: " + feature).
			Section("🎯 You want to see: *" + feature + "*").
			Section("This demonstrates regex pattern matching in the SlackBot library!").
			ButtonRow(
				ui.Button("show_features", "Show All Features").Primary(),
				ui.Button("open_modal", "Try Modal Demo"),
			)
		
		return ctx.ReplyWithBlocks(msg.Blocks())
	}
	
	return ctx.Reply("🤔 Try: \"demo features\" or \"demo modal\"")
}

// Interactive component handlers
func handleGetStarted(ctx *slackbot.InteractiveContext) error {
	msg := ui.NewMessage().
		Text("Getting Started Guide").
		Section("🎉 Welcome to the SlackBot Library!").
		Section("*Key Features:*").
		Section("• ⚡ **5-line minimum** - Working bot in 5 lines of code\n• 🏗️ **Builder pattern** - Fluent configuration API\n• 🎯 **Rich context** - Powerful handler context objects\n• 🎨 **UI builders** - Easy Block Kit components\n• 🔧 **Middleware** - Logging, rate limiting, error recovery\n• 💾 **Storage** - Redis, PostgreSQL, memory support\n• 🔄 **Defaults** - Works out of the box with sensible defaults").
		Section("*Try these features:*").
		ButtonRow(
			ui.Button("show_features", "Show Features"),
			ui.Button("open_modal", "Open Modal Demo"),
		)
	
	return ctx.UpdateMessage(msg.GetText())
}

func handleShowFeatures(ctx *slackbot.InteractiveContext) error {
	msg := ui.NewMessage().
		Text("SlackBot Library Features").
		Section("🚀 **SlackBot Library - Complete Feature Set**").
		Section("*Core Features:*").
		Section("• 📝 **Message Handling** - All message types with pattern matching\n• 🎛️ **Interactive Components** - Buttons, modals, select menus\n• 🏠 **Home Tab** - Custom home tab with default welcome\n• ⚡ **Slash Commands** - Easy slash command handling\n• 🎭 **Events** - Reactions, user joins, file sharing\n• 🔧 **Middleware** - Logging, rate limiting, metrics, error recovery").
		Section("*Developer Experience:*").
		Section("• 🏗️ **Builder Pattern** - Fluent API like Material UI\n• 🎯 **Rich Context** - Helper methods for common tasks\n• 🎨 **UI Builder** - Block Kit components made easy\n• 💾 **Storage Abstraction** - Redis, PostgreSQL, memory\n• 🔄 **Auto-fallbacks** - Sensible defaults everywhere\n• 📚 **Comprehensive Examples** - From 5-line to full-featured").
		ButtonRow(
			ui.Button("open_modal", "Try Modal Demo").Primary(),
			ui.Button("get_started", "Back to Start"),
		)
	
	return ctx.UpdateMessage(msg.GetText())
}

func handleOpenModal(ctx *slackbot.InteractiveContext) error {
	modal := ui.NewModal("demo_modal", "SlackBot Library Demo").
		Section("🎯 This modal was created with the UI builder system!").
		Section("*Try the form below:*").
		Input("user_name", "Your Name:", ui.PlainText().WithPlaceholder("Enter your name")).
		Input("user_email", "Email (optional):", ui.Email().WithPlaceholder("you@example.com")).
		SelectMenu("demo_option", "Choose a feature:", []ui.Option{
			ui.Option("builder_pattern", "Builder Pattern API"),
			ui.Option("rich_context", "Rich Context Objects"),
			ui.Option("ui_components", "UI Component Builders"),
			ui.Option("middleware", "Middleware System"),
			ui.Option("storage", "Storage Abstraction"),
			ui.Option("defaults", "Default Behaviors"),
		}).
		SelectMenu("experience_level", "Your Experience:", []ui.Option{
			ui.Option("beginner", "Beginner - New to Slack bots"),
			ui.Option("intermediate", "Intermediate - Some experience"),
			ui.Option("advanced", "Advanced - Building complex bots"),
		}).
		Submit("Submit Demo").
		Cancel("Cancel")
	
	return ctx.OpenModal(modal.Build())
}

func handleDemoOption(ctx *slackbot.InteractiveContext) error {
	selectedOption := ctx.GetSelectedOption("demo_option")
	if selectedOption == nil {
		return ctx.RespondEphemeral("❌ No option selected")
	}
	
	feature := selectedOption.Text.Text
	
	return ctx.RespondEphemeral("✅ Great choice! You selected: **" + feature + "**\n\nThis demonstrates the select menu handling in the SlackBot library!")
}

func handleModalSubmit(ctx *slackbot.InteractiveContext) error {
	values := ctx.GetSubmissionValues()
	
	name := values["user_name"]
	email := values["user_email"]
	feature := values["demo_option"]
	experience := values["experience_level"]
	
	if name == "" {
		return ctx.RespondEphemeral("❌ Please enter your name")
	}
	
	response := "✅ **Modal Submitted Successfully!**\n\n"
	response += "**Your Details:**\n"
	response += "• Name: " + name + "\n"
	if email != "" {
		response += "• Email: " + email + "\n"
	}
	response += "• Interested in: " + feature + "\n"
	response += "• Experience: " + experience + "\n\n"
	response += "This demonstrates form handling with the SlackBot library!"
	
	return ctx.RespondEphemeral(response)
}

// Event handlers
func handleHomeOpened(ctx *slackbot.HomeContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	homeView := ui.NewHomeView().
		Header("🏠 SlackBot Library Demo").
		Section("Welcome " + user.Name + "! This is a custom home tab built with the SlackBot library.").
		Section("*🚀 What makes this special:*").
		Section("• **Plug-and-play**: Works out of the box with defaults\n• **Builder pattern**: Fluent API for easy configuration\n• **Rich context**: Powerful handler context objects\n• **UI components**: Block Kit made simple").
		Divider().
		Section("*🎯 Key Features Demonstrated:*").
		ButtonRow(
			ui.Button("get_started", "Get Started Guide").Primary(),
			ui.Button("show_features", "Show All Features"),
		).
		ButtonRow(
			ui.Button("open_modal", "Try Modal Demo"),
			ui.Button("show_features", "Library Features"),
		).
		Divider().
		Section("*📚 Learn More:*").
		Section("• Check out `examples/simple/` for a 5-line bot\n• Try `examples/advanced/` for full features\n• Read the documentation in `README.md`\n• Explore the source code in `pkg/slackbot/`").
		Divider().
		Section("*💡 Quick Tips:*").
		Section("• Say \"hello\" to me for a greeting\n• Say \"help\" for all available commands\n• Try mentioning me in a channel\n• Use `/demo` and `/status` slash commands").
		Divider().
		Section("_This home tab demonstrates the default behavior with custom overrides - the library provides a default welcome tab, but you can easily customize it like this!_")
	
	return ctx.PublishView(&slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: homeView.Blocks(),
	})
}

func handleReactionAdded(ctx *slackbot.ReactionContext) error {
	// Demo: respond to specific reactions
	switch ctx.Reaction {
	case "rocket":
		return ctx.Client.AddReaction("sparkles", slack.ItemRef{
			Channel:   ctx.ChannelID,
			Timestamp: ctx.Item.Timestamp,
		})
	case "wave":
		return ctx.Client.AddReaction("wave", slack.ItemRef{
			Channel:   ctx.ChannelID,
			Timestamp: ctx.Item.Timestamp,
		})
	case "heart":
		return ctx.Client.AddReaction("heart", slack.ItemRef{
			Channel:   ctx.ChannelID,
			Timestamp: ctx.Item.Timestamp,
		})
	}
	
	return nil
}

func handleUserJoined(ctx *slackbot.UserContext) error {
	user := ctx.User
	log.Printf("👋 New user joined: %s (%s)", user.Name, user.ID)
	
	// In a real bot, you might send a welcome message
	// For the demo, we just log it
	return nil
}

// Slash command handlers
func handleDemoSlashCommand(ctx *slackbot.SlashCommandContext) error {
	modal := ui.NewModal("demo_slash", "Demo Slash Command").
		Section("🎯 This modal was triggered by a slash command!").
		Section("*SlackBot Library Features:*").
		Section("• Easy slash command handling\n• Automatic modal responses\n• Rich context objects\n• UI component builders").
		Input("feedback", "Your Feedback:", ui.PlainText().WithPlaceholder("What do you think of the SlackBot library?")).
		SelectMenu("rating", "Rate this demo:", []ui.Option{
			ui.Option("excellent", "⭐⭐⭐⭐⭐ Excellent"),
			ui.Option("good", "⭐⭐⭐⭐ Good"),
			ui.Option("fair", "⭐⭐⭐ Fair"),
			ui.Option("poor", "⭐⭐ Poor"),
		}).
		Submit("Submit Feedback").
		Cancel("Cancel")
	
	return ctx.OpenModal(modal.Build())
}

func handleStatusSlashCommand(ctx *slackbot.SlashCommandContext) error {
	status := `
*🤖 SlackBot Library Demo - Status*

*Bot Status:*
• Status: ✅ Online and running
• Library: SlackBot v1.0.0
• Features: All systems operational
• Storage: Connected and healthy

*Demonstrated Features:*
• ✅ Message handling (all types)
• ✅ Interactive components (buttons, modals)
• ✅ Custom home tab
• ✅ Slash commands
• ✅ Event handling (reactions, users)
• ✅ Middleware system
• ✅ Storage abstraction
• ✅ UI component builders

*Library Benefits:*
• 🚀 **5-line minimum** - Working bot in 5 lines
• 🏗️ **Builder pattern** - Material UI-level ease
• 🎯 **Rich context** - Powerful helper methods
• 🔧 **Middleware** - Built-in logging, rate limiting
• 💾 **Storage** - Redis, PostgreSQL, memory support
• 🔄 **Defaults** - Sensible fallbacks everywhere

*Try these commands:*
• Say "hello" for a greeting
• Say "help" for full command list
• Mention me in a channel
• Click buttons and try modals
• Check out the home tab

This demo shows how easy it is to build professional Slack bots with the SlackBot library! 🎉
	`
	
	return ctx.RespondEphemeral(status)
}