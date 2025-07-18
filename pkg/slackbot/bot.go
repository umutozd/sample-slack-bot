package slackbot

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
)

// Bot represents the main Slack bot interface
type Bot interface {
	Start(addr string) error
	Stop() error
	SendMessage(channelID, text string) error
	SendEphemeral(channelID, userID, text string) error
}

// BotBuilder provides a fluent interface for building and configuring a Slack bot
type BotBuilder interface {
	// Configuration
	WithClientCredentials(clientID, clientSecret string) BotBuilder
	WithSigningSecret(secret string) BotBuilder
	WithStorage(storage storage.Storage) BotBuilder
	WithConfig(config *Config) BotBuilder
	WithMiddleware(middleware ...Middleware) BotBuilder
	WithPort(port int) BotBuilder
	WithDebug(debug bool) BotBuilder

	// Core Event Handlers
	OnMessage(handler MessageHandler) BotBuilder
	OnDirectMessage(handler MessageHandler) BotBuilder
	OnChannelMessage(handler MessageHandler) BotBuilder
	OnThreadMessage(handler MessageHandler) BotBuilder
	OnMentionMessage(handler MessageHandler) BotBuilder
	OnHomeOpened(handler HomeHandler) BotBuilder

	// Interactive Component Handlers
	OnButtonClick(actionID string, handler InteractiveHandler) BotBuilder
	OnAnyButtonClick(handler InteractiveHandler) BotBuilder
	OnSelectMenu(actionID string, handler InteractiveHandler) BotBuilder
	OnModalSubmit(callbackID string, handler InteractiveHandler) BotBuilder
	OnSlashCommand(command string, handler SlashCommandHandler) BotBuilder

	// Advanced Event Patterns
	OnMessageContaining(text string, handler MessageHandler) BotBuilder
	OnMessageMatching(pattern *regexp.Regexp, handler MessageHandler) BotBuilder
	OnAnyEvent(handler EventHandler) BotBuilder

	// Build and start
	Build() Bot
	Start(addr string) error
}

// Config holds the configuration for the Slack bot
type Config struct {
	ClientID      string
	ClientSecret  string
	SigningSecret string
	Port          int
	Debug         bool
	LogLevel      string
	Storage       storage.Storage
	Middleware    []Middleware
	HomeHandler   HomeHandler
}

// bot is the internal implementation of the Bot interface
type bot struct {
	config *Config
	server *http.Server
	client *slack.Client
	
	// Event handlers
	messageHandlers        []messageHandlerEntry
	directMessageHandlers  []MessageHandler
	channelMessageHandlers []MessageHandler
	threadMessageHandlers  []MessageHandler
	mentionMessageHandlers []MessageHandler
	homeHandlers           []HomeHandler
	
	// Interactive handlers
	buttonHandlers       map[string]InteractiveHandler
	anyButtonHandlers    []InteractiveHandler
	selectMenuHandlers   map[string]InteractiveHandler
	modalSubmitHandlers  map[string]InteractiveHandler
	slashCommandHandlers map[string]SlashCommandHandler
	
	// Advanced handlers
	messageContainingHandlers []messageContainingEntry
	messageMatchingHandlers   []messageMatchingEntry
	eventHandlers            []EventHandler
	
	// Middleware
	middleware []Middleware
	
	// Internal state
	mu       sync.RWMutex
	running  bool
	stopChan chan struct{}
}

// messageHandlerEntry represents a message handler with its metadata
type messageHandlerEntry struct {
	handler MessageHandler
	filter  func(*MessageContext) bool
}

// messageContainingEntry represents a message handler that triggers on specific text
type messageContainingEntry struct {
	text    string
	handler MessageHandler
}

// messageMatchingEntry represents a message handler that triggers on regex patterns
type messageMatchingEntry struct {
	pattern *regexp.Regexp
	handler MessageHandler
}

// botBuilder is the internal implementation of the BotBuilder interface
type botBuilder struct {
	config *Config
	bot    *bot
}

// New creates a new BotBuilder with default configuration
func New() BotBuilder {
	config := &Config{
		Port:        8080,
		Debug:       false,
		LogLevel:    "info",
		Storage:     storage.NewMemoryStorage(),
		Middleware:  []Middleware{},
		HomeHandler: defaultHomeHandler,
	}
	
	return &botBuilder{
		config: config,
		bot: &bot{
			config:                    config,
			buttonHandlers:           make(map[string]InteractiveHandler),
			anyButtonHandlers:        []InteractiveHandler{},
			selectMenuHandlers:       make(map[string]InteractiveHandler),
			modalSubmitHandlers:      make(map[string]InteractiveHandler),
			slashCommandHandlers:     make(map[string]SlashCommandHandler),
			messageHandlers:          []messageHandlerEntry{},
			directMessageHandlers:    []MessageHandler{},
			channelMessageHandlers:   []MessageHandler{},
			threadMessageHandlers:    []MessageHandler{},
			mentionMessageHandlers:   []MessageHandler{},
			homeHandlers:             []HomeHandler{},
			messageContainingHandlers: []messageContainingEntry{},
			messageMatchingHandlers:   []messageMatchingEntry{},
			eventHandlers:            []EventHandler{},
			middleware:               []Middleware{},
			stopChan:                 make(chan struct{}),
		},
	}
}

// Configuration methods
func (b *botBuilder) WithClientCredentials(clientID, clientSecret string) BotBuilder {
	b.config.ClientID = clientID
	b.config.ClientSecret = clientSecret
	return b
}

func (b *botBuilder) WithSigningSecret(secret string) BotBuilder {
	b.config.SigningSecret = secret
	return b
}

func (b *botBuilder) WithStorage(storage storage.Storage) BotBuilder {
	b.config.Storage = storage
	return b
}

func (b *botBuilder) WithConfig(config *Config) BotBuilder {
	b.config = config
	return b
}

func (b *botBuilder) WithMiddleware(middleware ...Middleware) BotBuilder {
	b.config.Middleware = append(b.config.Middleware, middleware...)
	b.bot.middleware = append(b.bot.middleware, middleware...)
	return b
}

func (b *botBuilder) WithPort(port int) BotBuilder {
	b.config.Port = port
	return b
}

func (b *botBuilder) WithDebug(debug bool) BotBuilder {
	b.config.Debug = debug
	return b
}

// Event handler methods
func (b *botBuilder) OnMessage(handler MessageHandler) BotBuilder {
	b.bot.messageHandlers = append(b.bot.messageHandlers, messageHandlerEntry{
		handler: handler,
		filter:  func(ctx *MessageContext) bool { return true },
	})
	return b
}

func (b *botBuilder) OnDirectMessage(handler MessageHandler) BotBuilder {
	b.bot.directMessageHandlers = append(b.bot.directMessageHandlers, handler)
	return b
}

func (b *botBuilder) OnChannelMessage(handler MessageHandler) BotBuilder {
	b.bot.channelMessageHandlers = append(b.bot.channelMessageHandlers, handler)
	return b
}

func (b *botBuilder) OnThreadMessage(handler MessageHandler) BotBuilder {
	b.bot.threadMessageHandlers = append(b.bot.threadMessageHandlers, handler)
	return b
}

func (b *botBuilder) OnMentionMessage(handler MessageHandler) BotBuilder {
	b.bot.mentionMessageHandlers = append(b.bot.mentionMessageHandlers, handler)
	return b
}

func (b *botBuilder) OnHomeOpened(handler HomeHandler) BotBuilder {
	b.bot.homeHandlers = append(b.bot.homeHandlers, handler)
	return b
}

// Interactive handler methods
func (b *botBuilder) OnButtonClick(actionID string, handler InteractiveHandler) BotBuilder {
	b.bot.buttonHandlers[actionID] = handler
	return b
}

func (b *botBuilder) OnAnyButtonClick(handler InteractiveHandler) BotBuilder {
	b.bot.anyButtonHandlers = append(b.bot.anyButtonHandlers, handler)
	return b
}

func (b *botBuilder) OnSelectMenu(actionID string, handler InteractiveHandler) BotBuilder {
	b.bot.selectMenuHandlers[actionID] = handler
	return b
}

func (b *botBuilder) OnModalSubmit(callbackID string, handler InteractiveHandler) BotBuilder {
	b.bot.modalSubmitHandlers[callbackID] = handler
	return b
}

func (b *botBuilder) OnSlashCommand(command string, handler SlashCommandHandler) BotBuilder {
	b.bot.slashCommandHandlers[command] = handler
	return b
}

// Advanced pattern methods
func (b *botBuilder) OnMessageContaining(text string, handler MessageHandler) BotBuilder {
	b.bot.messageContainingHandlers = append(b.bot.messageContainingHandlers, messageContainingEntry{
		text:    text,
		handler: handler,
	})
	return b
}

func (b *botBuilder) OnMessageMatching(pattern *regexp.Regexp, handler MessageHandler) BotBuilder {
	b.bot.messageMatchingHandlers = append(b.bot.messageMatchingHandlers, messageMatchingEntry{
		pattern: pattern,
		handler: handler,
	})
	return b
}

func (b *botBuilder) OnAnyEvent(handler EventHandler) BotBuilder {
	b.bot.eventHandlers = append(b.bot.eventHandlers, handler)
	return b
}

// Build creates and returns the configured Bot
func (b *botBuilder) Build() Bot {
	b.bot.config = b.config
	
	// Set up logging
	if b.config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	
	// Use default home handler if none provided
	if len(b.bot.homeHandlers) == 0 {
		b.bot.homeHandlers = append(b.bot.homeHandlers, b.config.HomeHandler)
	}
	
	return b.bot
}

// Start builds and starts the bot
func (b *botBuilder) Start(addr string) error {
	bot := b.Build()
	return bot.Start(addr)
}

// Bot implementation methods
func (b *bot) Start(addr string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if b.running {
		return fmt.Errorf("bot is already running")
	}
	
	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/slack/install", b.handleInstall)
	mux.HandleFunc("/slack/events", b.handleEvents)
	mux.HandleFunc("/slack/interactive", b.handleInteractive)
	
	b.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	
	b.running = true
	
	logrus.Infof("Starting Slack bot on %s", addr)
	return b.server.ListenAndServe()
}

func (b *bot) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if !b.running {
		return fmt.Errorf("bot is not running")
	}
	
	close(b.stopChan)
	b.running = false
	
	if b.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return b.server.Shutdown(ctx)
	}
	
	return nil
}

func (b *bot) SendMessage(channelID, text string) error {
	if b.client == nil {
		return fmt.Errorf("slack client not initialized")
	}
	
	_, _, err := b.client.PostMessage(channelID, slack.MsgOptionText(text, false))
	return err
}

func (b *bot) SendEphemeral(channelID, userID, text string) error {
	if b.client == nil {
		return fmt.Errorf("slack client not initialized")
	}
	
	_, err := b.client.PostEphemeral(channelID, userID, slack.MsgOptionText(text, false))
	return err
}

// Placeholder HTTP handlers - these will be implemented in separate files
func (b *bot) handleInstall(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement OAuth installation flow
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("OAuth installation not yet implemented"))
}

func (b *bot) handleEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement event handling
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Event handling not yet implemented"))
}

func (b *bot) handleInteractive(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement interactive component handling
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Interactive handling not yet implemented"))
}