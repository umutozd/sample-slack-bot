package slackbot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/middleware"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/ui"
)

// Bot represents the main Slack bot instance
type Bot struct {
	clientID      string
	clientSecret  string
	signingSecret string
	storage       storage.Storage
	port          int
	debug         bool
	middlewares   []middleware.Middleware

	// Event handlers
	messageHandlers         []MessageHandler
	directMessageHandlers   []MessageHandler
	channelMessageHandlers  []MessageHandler
	threadMessageHandlers   []MessageHandler
	mentionMessageHandlers  []MessageHandler
	homeOpenedHandlers      []HomeHandler
	buttonClickHandlers     map[string][]InteractiveHandler
	selectMenuHandlers      map[string][]InteractiveHandler
	modalSubmitHandlers     map[string][]InteractiveHandler
	slashCommandHandlers    map[string][]SlashCommandHandler
	reactionAddedHandlers   []ReactionHandler
	fileSharedHandlers      []FileHandler
	userJoinedHandlers      []UserHandler
	channelJoinedHandlers   []ChannelHandler

	// Pattern-based handlers
	messageContainingHandlers map[string][]MessageHandler
	messageMatchingHandlers   map[*regexp.Regexp][]MessageHandler

	// Middleware handlers
	anyEventHandlers []MiddlewareHandler

	server *http.Server
}

// BotBuilder provides a fluent interface for building bots
type BotBuilder struct {
	bot *Bot
}

// New creates a new bot builder
func New() *BotBuilder {
	bot := &Bot{
		port:                      8080,
		debug:                     false,
		storage:                   storage.NewMemoryStorage(),
		middlewares:               []middleware.Middleware{middleware.LoggingMiddleware()},
		buttonClickHandlers:       make(map[string][]InteractiveHandler),
		selectMenuHandlers:        make(map[string][]InteractiveHandler),
		modalSubmitHandlers:       make(map[string][]InteractiveHandler),
		slashCommandHandlers:      make(map[string][]SlashCommandHandler),
		messageContainingHandlers: make(map[string][]MessageHandler),
		messageMatchingHandlers:   make(map[*regexp.Regexp][]MessageHandler),
	}
	
	// Add default home handler
	bot.homeOpenedHandlers = []HomeHandler{defaultHomeHandler}
	
	return &BotBuilder{bot: bot}
}

// WithClientCredentials sets the Slack client credentials
func (b *BotBuilder) WithClientCredentials(clientID, clientSecret string) *BotBuilder {
	b.bot.clientID = clientID
	b.bot.clientSecret = clientSecret
	return b
}

// WithSigningSecret sets the Slack signing secret
func (b *BotBuilder) WithSigningSecret(signingSecret string) *BotBuilder {
	b.bot.signingSecret = signingSecret
	return b
}

// WithStorage sets the storage implementation
func (b *BotBuilder) WithStorage(storage storage.Storage) *BotBuilder {
	b.bot.storage = storage
	return b
}

// WithPort sets the server port
func (b *BotBuilder) WithPort(port int) *BotBuilder {
	b.bot.port = port
	return b
}

// WithDebug enables debug mode
func (b *BotBuilder) WithDebug(debug bool) *BotBuilder {
	b.bot.debug = debug
	return b
}

// WithMiddleware adds middleware to the bot
func (b *BotBuilder) WithMiddleware(middlewares ...middleware.Middleware) *BotBuilder {
	b.bot.middlewares = append(b.bot.middlewares, middlewares...)
	return b
}

// OnMessage adds a message handler
func (b *BotBuilder) OnMessage(handler MessageHandler) *BotBuilder {
	b.bot.messageHandlers = append(b.bot.messageHandlers, handler)
	return b
}

// OnDirectMessage adds a direct message handler
func (b *BotBuilder) OnDirectMessage(handler MessageHandler) *BotBuilder {
	b.bot.directMessageHandlers = append(b.bot.directMessageHandlers, handler)
	return b
}

// OnChannelMessage adds a channel message handler
func (b *BotBuilder) OnChannelMessage(handler MessageHandler) *BotBuilder {
	b.bot.channelMessageHandlers = append(b.bot.channelMessageHandlers, handler)
	return b
}

// OnThreadMessage adds a thread message handler
func (b *BotBuilder) OnThreadMessage(handler MessageHandler) *BotBuilder {
	b.bot.threadMessageHandlers = append(b.bot.threadMessageHandlers, handler)
	return b
}

// OnMentionMessage adds a mention message handler
func (b *BotBuilder) OnMentionMessage(handler MessageHandler) *BotBuilder {
	b.bot.mentionMessageHandlers = append(b.bot.mentionMessageHandlers, handler)
	return b
}

// OnHomeOpened adds a home opened handler
func (b *BotBuilder) OnHomeOpened(handler HomeHandler) *BotBuilder {
	b.bot.homeOpenedHandlers = append(b.bot.homeOpenedHandlers, handler)
	return b
}

// OnButtonClick adds a button click handler
func (b *BotBuilder) OnButtonClick(actionID string, handler InteractiveHandler) *BotBuilder {
	b.bot.buttonClickHandlers[actionID] = append(b.bot.buttonClickHandlers[actionID], handler)
	return b
}

// OnSelectMenu adds a select menu handler
func (b *BotBuilder) OnSelectMenu(actionID string, handler InteractiveHandler) *BotBuilder {
	b.bot.selectMenuHandlers[actionID] = append(b.bot.selectMenuHandlers[actionID], handler)
	return b
}

// OnModalSubmit adds a modal submit handler
func (b *BotBuilder) OnModalSubmit(callbackID string, handler InteractiveHandler) *BotBuilder {
	b.bot.modalSubmitHandlers[callbackID] = append(b.bot.modalSubmitHandlers[callbackID], handler)
	return b
}

// OnSlashCommand adds a slash command handler
func (b *BotBuilder) OnSlashCommand(command string, handler SlashCommandHandler) *BotBuilder {
	b.bot.slashCommandHandlers[command] = append(b.bot.slashCommandHandlers[command], handler)
	return b
}

// OnReactionAdded adds a reaction added handler
func (b *BotBuilder) OnReactionAdded(handler ReactionHandler) *BotBuilder {
	b.bot.reactionAddedHandlers = append(b.bot.reactionAddedHandlers, handler)
	return b
}

// OnFileShared adds a file shared handler
func (b *BotBuilder) OnFileShared(handler FileHandler) *BotBuilder {
	b.bot.fileSharedHandlers = append(b.bot.fileSharedHandlers, handler)
	return b
}

// OnUserJoined adds a user joined handler
func (b *BotBuilder) OnUserJoined(handler UserHandler) *BotBuilder {
	b.bot.userJoinedHandlers = append(b.bot.userJoinedHandlers, handler)
	return b
}

// OnChannelJoined adds a channel joined handler
func (b *BotBuilder) OnChannelJoined(handler ChannelHandler) *BotBuilder {
	b.bot.channelJoinedHandlers = append(b.bot.channelJoinedHandlers, handler)
	return b
}

// OnMessageContaining adds a handler for messages containing specific text
func (b *BotBuilder) OnMessageContaining(text string, handler MessageHandler) *BotBuilder {
	b.bot.messageContainingHandlers[text] = append(b.bot.messageContainingHandlers[text], handler)
	return b
}

// OnMessageMatching adds a handler for messages matching a regex pattern
func (b *BotBuilder) OnMessageMatching(pattern *regexp.Regexp, handler MessageHandler) *BotBuilder {
	b.bot.messageMatchingHandlers[pattern] = append(b.bot.messageMatchingHandlers[pattern], handler)
	return b
}

// OnAnyEvent adds a middleware-style handler for any event
func (b *BotBuilder) OnAnyEvent(handler MiddlewareHandler) *BotBuilder {
	b.bot.anyEventHandlers = append(b.bot.anyEventHandlers, handler)
	return b
}

// Start starts the bot server
func (b *BotBuilder) Start(addr string) error {
	if addr != "" {
		if addr[0] == ':' {
			port, err := strconv.Atoi(addr[1:])
			if err != nil {
				return fmt.Errorf("invalid port in address: %s", addr)
			}
			b.bot.port = port
		}
	}
	
	return b.bot.start()
}

// Build builds the bot (alternative to Start for advanced use cases)
func (b *BotBuilder) Build() *Bot {
	return b.bot
}

// start starts the bot server
func (bot *Bot) start() error {
	mux := http.NewServeMux()
	
	// Add OAuth endpoint
	mux.HandleFunc("/slack/install", bot.handleInstall)
	mux.HandleFunc("/slack/oauth", bot.handleOAuth)
	
	// Add event endpoints
	mux.HandleFunc("/slack/events", bot.handleEvents)
	mux.HandleFunc("/slack/interactive", bot.handleInteractive)
	
	bot.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", bot.port),
		Handler: mux,
	}
	
	if bot.debug {
		log.Printf("🚀 Bot server starting on port %d", bot.port)
	}
	
	return bot.server.ListenAndServe()
}

// Stop stops the bot server
func (bot *Bot) Stop() error {
	if bot.server == nil {
		return nil
	}
	return bot.server.Shutdown(context.Background())
}