# SlackBot - Plug-and-Play Slack Bot Library

A Material UI-style library for building Slack bots in Go. Mind-numbingly easy to use, yet incredibly powerful and configurable.

## ⚡ Quick Start

Create a working Slack bot in just 5 lines of code:

```go
slackbot.New().
    WithClientCredentials("client_id", "client_secret").
    OnMessage(func(ctx *slackbot.MessageContext) error {
        return ctx.Reply("Hello " + ctx.UserID)
    }).
    Start(":8080")
```

That's it! Your bot is running with:
- ✅ Automatic in-memory storage
- ✅ Default welcome home tab
- ✅ Built-in error handling
- ✅ Comprehensive logging
- ✅ OAuth installation flow

## 🎯 Design Philosophy

Inspired by Material UI's developer experience:
- **Plug-and-play**: Works immediately with sensible defaults
- **5-year-old-can-do-it-easy**: Minimal learning curve
- **Highly configurable**: Customize everything when needed
- **Production-ready**: Built-in best practices and patterns

## 🚀 Features

### Core Features
- **Builder Pattern API**: Fluent configuration
- **Default Behaviors**: In-memory storage, welcome home tab, error handling
- **Rich Event System**: 25+ Slack events with advanced patterns
- **Interactive Components**: Buttons, modals, select menus
- **UI Builder**: Create rich Slack interfaces easily
- **Middleware System**: Logging, rate limiting, authentication, etc.
- **Multiple Storage Backends**: Memory, Redis, PostgreSQL, SQLite
- **Context-Rich Handlers**: Access to user, channel, team information

### Developer Experience
- **Zero Configuration**: Works out of the box
- **Type Safety**: Comprehensive type definitions
- **Extensive Documentation**: Examples and guides
- **Testing Utilities**: Built-in test helpers
- **CLI Tools**: Project scaffolding and generators

## 📦 Installation

```bash
go get github.com/umutozd/sample-slack-bot/pkg/slackbot
```

## 🏃 Running the Examples

1. Set your environment variables:
   ```bash
   export SLACK_CLIENT_ID="your-client-id"
   export SLACK_CLIENT_SECRET="your-client-secret"
   ```

2. Run the simple example:
   ```bash
   cd examples/simple
   go run main.go
   ```

3. Run the advanced example:
   ```bash
   cd examples/advanced
   go run main.go
   ```

## 🎨 UI Builder

Create rich Slack interfaces with ease:

```go
homeView := ui.NewHomeView().
    Header("Welcome to My Bot").
    Section("This is what my bot does!").
    ButtonRow(
        ui.Button("action1", "Primary Action").Primary(),
        ui.Button("action2", "Secondary Action"),
    ).
    Divider().
    Section("More content here")

return ctx.PublishView(homeView.Build())
```

## 🔧 Configuration

### Storage Options

```go
// Default: In-memory storage
bot := slackbot.New() // Uses memory storage automatically

// Redis storage
bot := slackbot.New().
    WithStorage(storage.NewRedisStorage("localhost:6379", "", 0))

// PostgreSQL storage
bot := slackbot.New().
    WithStorage(storage.NewPostgresStorage("postgres://..."))

// Layered storage (cache + persistent)
bot := slackbot.New().
    WithStorage(storage.NewLayeredStorage(
        storage.NewMemoryStorage(),
        storage.NewRedisStorage("localhost:6379", "", 0),
    ))
```

### Middleware

```go
bot := slackbot.New().
    WithMiddleware(
        slackbot.LoggingMiddleware(),
        slackbot.RateLimitMiddleware(60, time.Minute),
        slackbot.AuthMiddleware(func(ctx *slackbot.Context) bool {
            return ctx.UserID != "blocked_user"
        }),
        slackbot.ErrorRecoveryMiddleware(),
    )
```

## 📋 Event Handling

### Message Events
```go
bot.OnMessage(handler)           // Any message
bot.OnDirectMessage(handler)     // Direct messages only
bot.OnChannelMessage(handler)    // Channel messages only
bot.OnThreadMessage(handler)     // Thread replies only
bot.OnMentionMessage(handler)    // When bot is mentioned
```

### Pattern Matching
```go
bot.OnMessageContaining("hello", handler)
bot.OnMessageMatching(regexp.MustCompile(`weather in (\w+)`), handler)
```

### Interactive Components
```go
bot.OnButtonClick("action_id", handler)
bot.OnSelectMenu("menu_id", handler)
bot.OnModalSubmit("modal_id", handler)
bot.OnSlashCommand("/command", handler)
```

### Advanced Patterns
```go
bot.OnAnyEvent(func(ctx *slackbot.EventContext) error {
    // Global event handler
    return nil
})
```

## 🏠 Default Home Tab

When users first open your bot, they see a welcoming default home tab:

![Default Home Tab](docs/images/default-home.png)

- 🎉 Welcome message
- 📚 Getting started information
- 🔗 Links to documentation
- 💡 Next steps guidance

Override with your own custom home tab:

```go
bot.OnHomeOpened(func(ctx *slackbot.HomeContext) error {
    // Your custom home tab here
    return ctx.PublishView(myCustomHomeView)
})
```

## 🧪 Testing

Built-in testing utilities:

```go
func TestMyBot(t *testing.T) {
    bot := slackbot.NewTestBot()
    
    // Simulate events
    bot.SimulateMessage("user123", "channel456", "hello")
    bot.SimulateButtonClick("my_button", "value")
    
    // Assert responses
    bot.AssertMessageSent("Hello user123")
    bot.AssertModalOpened("welcome_modal")
}
```

## 🏗️ Architecture

```
pkg/slackbot/
├── bot.go              # Main bot interface and builder
├── handlers.go         # Handler types and context
├── defaults.go         # Default behaviors
├── middleware.go       # Middleware system
├── storage/           # Storage abstractions
│   ├── storage.go     # Interface definitions
│   ├── memory.go      # In-memory implementation
│   ├── redis.go       # Redis implementation
│   ├── postgres.go    # PostgreSQL implementation
│   └── layered.go     # Layered storage
└── ui/                # UI builder system
    └── ui.go          # Block Kit builders
```

## 🚦 Migration from Legacy Code

The library maintains backward compatibility while providing new features:

### Before (Legacy)
```go
// Old server-based approach
config := &server.Config{
    Debug:             true,
    Port:              8080,
    RedisAddress:      "localhost:6379",
    SlackClientID:     os.Getenv("SLACK_CLIENT_ID"),
    SlackClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
}
s, err := server.NewServer(config)
s.ListenAndServe()
```

### After (New Library)
```go
// New builder-based approach
slackbot.New().
    WithClientCredentials(
        os.Getenv("SLACK_CLIENT_ID"),
        os.Getenv("SLACK_CLIENT_SECRET"),
    ).
    WithDebug(true).
    OnMessage(myMessageHandler).
    Start(":8080")
```

## 🎯 Success Metrics

- **Ease of Use**: ✅ 5-line minimum for working bot
- **Plug-and-play**: ✅ Works out of the box with defaults
- **Immediate satisfaction**: ✅ Default home tab shows instant success
- **Zero dependencies**: ✅ In-memory storage by default
- **Flexibility**: ✅ Supports Redis, PostgreSQL, custom middleware
- **Production-ready**: ✅ Built-in error handling, logging, metrics

## 🛠️ Development

### Project Structure
```
/
├── main.go                    # Legacy server (kept for compatibility)
├── main-new.go               # New library demonstration
├── server/                   # Legacy server code
├── storage/                  # Legacy storage code
├── pkg/slackbot/            # New library code
├── examples/                # Usage examples
│   ├── simple/              # Basic example
│   ├── advanced/            # Full-featured example
│   └── README.md           # Example documentation
└── README.md               # This file
```

### Building and Testing
```bash
# Build the library
go build ./pkg/slackbot/...

# Run examples
go run examples/simple/main.go
go run examples/advanced/main.go

# Test the library
go test ./pkg/slackbot/...
```

## 📚 Documentation

- [Examples](examples/README.md) - Complete usage examples
- [API Reference](docs/api.md) - Detailed API documentation
- [Migration Guide](docs/migration.md) - Upgrading from legacy code
- [Best Practices](docs/best-practices.md) - Production deployment guide

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Material UI's developer experience
- Built on top of the excellent [slack-go/slack](https://github.com/slack-go/slack) library
- Thanks to all contributors who made this transformation possible

---

**From a simple sample to a professional library** - this is what modern Slack bot development should look like! 🚀