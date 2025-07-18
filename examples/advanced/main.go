package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/umutozd/sample-slack-bot/pkg/slackbot"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/middleware"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/ui"
)

func main() {
	// Create Redis storage (falls back to memory if Redis unavailable)
	var stor storage.Storage
	if redisAddr := os.Getenv("REDIS_URL"); redisAddr != "" {
		if redisStorage, err := storage.NewRedisStorage(redisAddr, "", 0); err == nil {
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

	// Create a full-featured bot with custom middleware and handlers
	bot := slackbot.New().
		WithClientCredentials(os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET")).
		WithStorage(stor).
		WithDebug(true).
		WithMiddleware(
			middleware.LoggingMiddleware(),
			middleware.ErrorRecoveryMiddleware(),
			middleware.RateLimitMiddleware(100, 60), // 100 requests per minute
			middleware.MetricsMiddleware(),
		).
		// Message handlers
		OnMessage(handleGeneralMessage).
		OnDirectMessage(handleDirectMessage).
		OnMentionMessage(handleMention).
		OnMessageContaining("help", handleHelpMessage).
		OnMessageMatching(regexp.MustCompile(`(?i)weather\s+in\s+(\w+)`), handleWeatherMessage).
		// Interactive handlers
		OnButtonClick("get_started", handleGetStarted).
		OnButtonClick("show_modal", handleShowModal).
		OnButtonClick("toggle_feature", handleToggleFeature).
		OnSelectMenu("user_role", handleRoleSelection).
		OnModalSubmit("user_preferences", handleUserPreferences).
		// Event handlers
		OnHomeOpened(handleCustomHome).
		OnUserJoined(handleUserJoined).
		OnReactionAdded(handleReactionAdded).
		OnFileShared(handleFileShared).
		// Slash commands
		OnSlashCommand("status", handleStatusCommand).
		OnSlashCommand("config", handleConfigCommand)

	log.Println("🚀 Advanced Slack Bot starting on port 8080...")
	log.Println("📝 Set environment variables:")
	log.Println("   SLACK_CLIENT_ID=your_client_id")
	log.Println("   SLACK_CLIENT_SECRET=your_client_secret")
	log.Println("   REDIS_URL=redis://localhost:6379 (optional)")
	log.Println("🌐 Install URL: http://localhost:8080/slack/install")
	
	if err := bot.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Message handlers
func handleGeneralMessage(ctx *slackbot.MessageContext) error {
	// Store message count for user
	var count int
	ctx.Storage.Get("msg_count_"+ctx.UserID, &count)
	count++
	ctx.Storage.Set("msg_count_"+ctx.UserID, count)
	
	// Don't respond to every message, just track it
	return nil
}

func handleDirectMessage(ctx *slackbot.MessageContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	return ctx.Reply("👋 Hi " + user.Name + "! I received your direct message: \"" + ctx.Text + "\"")
}

func handleMention(ctx *slackbot.MessageContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	msg := ui.NewMessage().
		Text("You mentioned me!").
		Section("Hey " + user.Name + "! 👋 You mentioned me in: \"" + ctx.Text + "\"").
		ButtonRow(
			ui.Button("get_started", "Get Started").Primary(),
			ui.Button("show_modal", "Show Modal"),
			ui.Button("toggle_feature", "Toggle Feature"),
		)
	
	return ctx.ReplyWithBlocks(msg.Blocks())
}

func handleHelpMessage(ctx *slackbot.MessageContext) error {
	helpText := `
*🤖 Bot Commands & Features:*

*Message Commands:*
• Mention me (@bot) for interactive options
• Say "help" to see this message
• Say "weather in [city]" to get weather info
• Direct message me for personal assistance

*Slash Commands:*
• \`/status\` - Check bot status
• \`/config\` - Configure bot settings

*Features:*
• Automatic welcome for new team members
• Reaction tracking and responses
• File sharing notifications
• Custom home tab experience
• User preference management
`
	
	return ctx.Reply(helpText)
}

func handleWeatherMessage(ctx *slackbot.MessageContext) error {
	// Extract city from regex match
	re := regexp.MustCompile(`(?i)weather\s+in\s+(\w+)`)
	matches := re.FindStringSubmatch(ctx.Text)
	
	if len(matches) > 1 {
		city := strings.Title(matches[1])
		return ctx.Reply("🌤️ Weather in " + city + ": 72°F, Partly Cloudy (This is a demo - integrate with a real weather API!)")
	}
	
	return ctx.Reply("🤔 I couldn't find the city name. Try: \"weather in Seattle\"")
}

// Interactive handlers
func handleGetStarted(ctx *slackbot.InteractiveContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	msg := ui.NewMessage().
		Text("Getting started guide").
		Section("🎉 Welcome " + user.Name + "! Here's how to get started:").
		Section("1. 💬 Send me a direct message\n2. 🏠 Check out the home tab\n3. ⚙️ Use `/config` to customize settings\n4. 🆘 Say \"help\" for more commands").
		ButtonRow(
			ui.Button("show_modal", "Open Settings").Primary(),
			ui.Button("toggle_feature", "Toggle Feature"),
		)
	
	return ctx.UpdateMessage(msg.GetText())
}

func handleShowModal(ctx *slackbot.InteractiveContext) error {
	modal := ui.NewModal("user_preferences", "User Preferences").
		Section("Customize your bot experience:").
		Input("name", "Display Name:", ui.PlainText().WithPlaceholder("Enter your preferred name")).
		Input("email", "Email:", ui.Email().WithPlaceholder("your@email.com")).
		SelectMenu("timezone", "Timezone:", []ui.Option{
			ui.Option("UTC", "UTC"),
			ui.Option("EST", "Eastern Time"),
			ui.Option("PST", "Pacific Time"),
			ui.Option("CST", "Central Time"),
		}).
		SelectMenu("notifications", "Notifications:", []ui.Option{
			ui.Option("all", "All notifications"),
			ui.Option("mentions", "Mentions only"),
			ui.Option("none", "No notifications"),
		}).
		Submit("Save Settings").
		Cancel("Cancel")
	
	return ctx.OpenModal(modal.Build())
}

func handleToggleFeature(ctx *slackbot.InteractiveContext) error {
	// Get current feature state
	var featureEnabled bool
	ctx.Storage.Get("feature_enabled_"+ctx.UserID, &featureEnabled)
	
	// Toggle it
	featureEnabled = !featureEnabled
	ctx.Storage.Set("feature_enabled_"+ctx.UserID, featureEnabled)
	
	status := "OFF"
	if featureEnabled {
		status = "ON"
	}
	
	return ctx.RespondEphemeral("🔧 Feature toggled " + status + " for your account!")
}

func handleRoleSelection(ctx *slackbot.InteractiveContext) error {
	selectedOption := ctx.GetSelectedOption("user_role")
	if selectedOption == nil {
		return ctx.RespondEphemeral("❌ No role selected")
	}
	
	// Store user role
	ctx.Storage.Set("user_role_"+ctx.UserID, selectedOption.Value)
	
	return ctx.RespondEphemeral("✅ Your role has been set to: " + selectedOption.Text.Text)
}

func handleUserPreferences(ctx *slackbot.InteractiveContext) error {
	values := ctx.GetSubmissionValues()
	
	// Save user preferences
	preferences := map[string]string{
		"name":          values["name"],
		"email":         values["email"],
		"timezone":      values["timezone"],
		"notifications": values["notifications"],
	}
	
	ctx.Storage.Set("user_preferences_"+ctx.UserID, preferences)
	
	return ctx.RespondEphemeral("✅ Your preferences have been saved!")
}

// Event handlers
func handleCustomHome(ctx *slackbot.HomeContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	// Get user stats
	var msgCount int
	ctx.Storage.Get("msg_count_"+ctx.UserID, &msgCount)
	
	var featureEnabled bool
	ctx.Storage.Get("feature_enabled_"+ctx.UserID, &featureEnabled)
	
	featureStatus := "OFF"
	if featureEnabled {
		featureStatus = "ON"
	}
	
	homeView := ui.NewHomeView().
		Header("🏠 Welcome " + user.Name + "!").
		Section("This is your personalized home tab. Here's what I know about you:").
		SectionWithFields("*Your Stats:*", []string{
			"*Messages sent:* " + string(rune(msgCount)),
			"*Feature status:* " + featureStatus,
			"*Member since:* " + user.Profile.Email,
		}).
		Divider().
		Section("*Quick Actions:*").
		ButtonRow(
			ui.Button("show_modal", "⚙️ Settings").Primary(),
			ui.Button("toggle_feature", "🔧 Toggle Feature"),
		).
		Divider().
		Section("*Getting Started:*").
		Section("• 💬 Send me a message to start chatting\n• 🔍 Try saying \"help\" for all commands\n• 🌤️ Ask for weather: \"weather in [city]\"\n• ⚡ Use `/status` slash command").
		Divider().
		Section("_This home tab updates based on your activity and preferences!_")
	
	return ctx.PublishView(&slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: homeView.Blocks(),
	})
}

func handleUserJoined(ctx *slackbot.UserContext) error {
	user := ctx.User
	
	// Send welcome message to general channel (you'd need to get the channel ID)
	// For now, just log it
	log.Printf("👋 New user joined: %s (%s)", user.Name, user.ID)
	
	return nil
}

func handleReactionAdded(ctx *slackbot.ReactionContext) error {
	// Track reaction statistics
	var reactionCount int
	ctx.Storage.Get("reaction_count_"+ctx.Reaction, &reactionCount)
	reactionCount++
	ctx.Storage.Set("reaction_count_"+ctx.Reaction, reactionCount)
	
	// Special handling for specific reactions
	if ctx.Reaction == "tada" {
		// React back with celebration
		return ctx.Client.AddReaction("sparkles", slack.ItemRef{
			Channel:   ctx.ChannelID,
			Timestamp: ctx.Item.Timestamp,
		})
	}
	
	return nil
}

func handleFileShared(ctx *slackbot.FileContext) error {
	file := ctx.File
	
	// Log file sharing
	log.Printf("📎 File shared: %s (%s) by %s", file.Name, file.Filetype, ctx.UserID)
	
	// Could implement file processing, virus scanning, etc.
	return nil
}

// Slash command handlers
func handleStatusCommand(ctx *slackbot.SlashCommandContext) error {
	// Get some stats
	var totalMessages int
	ctx.Storage.Get("total_messages", &totalMessages)
	
	status := `
*🤖 Bot Status:*
• Status: ✅ Online and running
• Total messages processed: ` + string(rune(totalMessages)) + `
• Uptime: Since last restart
• Storage: Connected and healthy
• Features: All systems operational

*Your Stats:*
• Messages sent: Loading...
• Feature status: Use settings to check
	`
	
	return ctx.RespondEphemeral(status)
}

func handleConfigCommand(ctx *slackbot.SlashCommandContext) error {
	modal := ui.NewModal("bot_config", "Bot Configuration").
		Section("Configure bot settings and preferences:").
		SelectMenu("notification_level", "Notification Level:", []ui.Option{
			ui.Option("all", "All notifications"),
			ui.Option("mentions", "Mentions only"),
			ui.Option("important", "Important only"),
			ui.Option("none", "No notifications"),
		}).
		SelectMenu("response_style", "Response Style:", []ui.Option{
			ui.Option("friendly", "Friendly and casual"),
			ui.Option("professional", "Professional"),
			ui.Option("concise", "Brief and concise"),
			ui.Option("detailed", "Detailed responses"),
		}).
		Input("custom_greeting", "Custom Greeting:", ui.PlainText().WithPlaceholder("Enter a custom greeting message")).
		Submit("Save Configuration").
		Cancel("Cancel")
	
	return ctx.OpenModal(modal.Build())
}