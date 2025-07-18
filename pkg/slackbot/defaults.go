package slackbot

import (
	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/ui"
)

// defaultHomeHandler provides a default welcome home tab
func defaultHomeHandler(ctx *HomeContext) error {
	view := ui.NewHomeView().
		Header("🎉 Welcome to Your Slack Bot!").
		Section("Your bot is up and running. Ready to customize?").
		ButtonRow(
			ui.Button("get_started", "Get Started"),
			ui.Button("view_docs", "View Documentation"),
		).
		Divider().
		Section("*Next Steps:*\n• Add custom message handlers\n• Create interactive components\n• Customize this home tab").
		Section("*Quick Start:*\n```go\nbot := slackbot.New().\n    WithClientCredentials(\"id\", \"secret\").\n    OnMessage(func(ctx *slackbot.MessageContext) error {\n        return ctx.Reply(\"Hello!\")\n    }).\n    Start(\":8080\")\n```").
		Divider().
		Section("Need help? Check out the examples and documentation to get started building your custom Slack bot!")
	
	homeView := &slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: view.Blocks(),
	}
	
	return ctx.PublishView(homeView)
}

// defaultErrorHandler provides a default error handler
func defaultErrorHandler(ctx *MessageContext, err error) error {
	if ctx.Bot.debug {
		return ctx.ReplyEphemeral("❌ An error occurred: " + err.Error())
	}
	return ctx.ReplyEphemeral("❌ Sorry, something went wrong. Please try again.")
}

// defaultWelcomeMessage provides a default welcome message for new users
func defaultWelcomeMessage(ctx *UserContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return err
	}
	
	// Send a DM to the new user
	channel, _, _, err := ctx.Client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user.ID},
	})
	if err != nil {
		return err
	}
	
	_, _, err = ctx.Client.PostMessage(channel.ID, slack.MsgOptionText(
		"👋 Welcome to the team, "+user.Name+"! I'm your friendly Slack bot. Feel free to message me anytime!", false))
	
	return err
}

// defaultButtonHandler provides a default handler for unhandled buttons
func defaultButtonHandler(ctx *InteractiveContext) error {
	switch ctx.ActionID {
	case "get_started":
		return ctx.RespondEphemeral("🚀 Great! Start by customizing your bot's message handlers. Check out the examples for inspiration!")
	case "view_docs":
		return ctx.RespondEphemeral("📚 Documentation and examples are available in the repository. Look for the `examples/` directory!")
	default:
		return ctx.RespondEphemeral("🤔 This button doesn't have a handler yet. Add one using `OnButtonClick(\"" + ctx.ActionID + "\", handler)`")
	}
}

// Default configuration values
const (
	DefaultPort         = 8080
	DefaultDebug        = false
	DefaultLogLevel     = "info"
	DefaultTimeout      = 30 // seconds
	DefaultMaxRetries   = 3
	DefaultRateLimit    = 100 // requests per minute
	DefaultCacheSize    = 1000
	DefaultCacheTTL     = 300 // seconds
)

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Port:         DefaultPort,
		Debug:        DefaultDebug,
		LogLevel:     DefaultLogLevel,
		Timeout:      DefaultTimeout,
		MaxRetries:   DefaultMaxRetries,
		RateLimit:    DefaultRateLimit,
		CacheSize:    DefaultCacheSize,
		CacheTTL:     DefaultCacheTTL,
	}
}

// Config holds configuration options for the bot
type Config struct {
	Port         int
	Debug        bool
	LogLevel     string
	Timeout      int
	MaxRetries   int
	RateLimit    int
	CacheSize    int
	CacheTTL     int
	
	// Custom handlers
	ErrorHandler func(*MessageContext, error) error
}

// GetErrorHandler returns the error handler or default if none set
func (c *Config) GetErrorHandler() func(*MessageContext, error) error {
	if c.ErrorHandler != nil {
		return c.ErrorHandler
	}
	return defaultErrorHandler
}

// Predefined messages and responses
var (
	WelcomeMessages = []string{
		"👋 Hello there! I'm your friendly Slack bot.",
		"🎉 Welcome! I'm here to help you get things done.",
		"🚀 Hi! Ready to supercharge your Slack experience?",
		"✨ Hey there! I'm your new Slack assistant.",
		"🤖 Hello! I'm a bot, but I promise I'm friendly!",
	}
	
	ErrorMessages = []string{
		"❌ Oops! Something went wrong. Please try again.",
		"🔧 I'm having trouble with that. Mind giving it another shot?",
		"⚠️ Something's not quite right. Let's try that again.",
		"🤔 Hmm, I hit a snag. Please retry your request.",
		"🚨 Sorry, I encountered an error. Please try again.",
	}
	
	SuccessMessages = []string{
		"✅ Done! That worked perfectly.",
		"🎉 Success! Your request has been completed.",
		"👍 All set! Everything went smoothly.",
		"🌟 Perfect! Your task is complete.",
		"✨ Fantastic! That's been taken care of.",
	}
)

// Helper functions for random responses
func RandomWelcomeMessage() string {
	return WelcomeMessages[0] // For simplicity, return first. In production, you'd randomize.
}

func RandomErrorMessage() string {
	return ErrorMessages[0] // For simplicity, return first. In production, you'd randomize.
}

func RandomSuccessMessage() string {
	return SuccessMessages[0] // For simplicity, return first. In production, you'd randomize.
}