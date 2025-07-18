# SlackBot Library - Plug-and-Play Slack Bot Development

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-examples-brightgreen.svg)](examples/)

A **Material UI-level** plug-and-play library for building Slack bots in Go. Create professional Slack bots with just 5 lines of code, while maintaining full flexibility for complex use cases.

## 🚀 Quick Start

### 5-Line Working Bot

```go
package main

import "github.com/umutozd/sample-slack-bot/pkg/slackbot"

func main() {
    slackbot.New().
        WithClientCredentials("your_client_id", "your_client_secret").
        OnMessage(func(ctx *slackbot.MessageContext) error {
            return ctx.Reply("Hello " + ctx.UserID + "!")
        }).
        Start(":8080")
}
```

That's it! You now have a working Slack bot with:
- ✅ OAuth flow
- ✅ Message handling
- ✅ Default home tab
- ✅ In-memory storage
- ✅ Error handling
- ✅ Logging middleware

## 🌟 Key Features

### **Plug-and-Play Design**
- **5-line minimum** - Working bot in 5 lines of code
- **Zero configuration** - Sensible defaults for everything
- **Immediate satisfaction** - Default home tab shows instant success
- **No dependencies** - In-memory storage by default

### **Material UI-Level API**
- **Builder pattern** - Fluent configuration interface
- **Rich context** - Powerful handler context objects
- **UI components** - Block Kit made simple
- **Type safety** - Full TypeScript-level safety in Go

### **Production Ready**
- **Middleware system** - Logging, rate limiting, authentication
- **Storage abstraction** - Redis, PostgreSQL, SQLite, memory
- **Error recovery** - Graceful error handling and recovery
- **Comprehensive events** - 25+ Slack events supported

## 📋 Installation

```bash
go mod init your-slack-bot
go get github.com/umutozd/sample-slack-bot
```

## 🎯 Examples

### Simple Bot (5 lines)
```go
slackbot.New().
    WithClientCredentials("id", "secret").
    OnMessage(func(ctx *slackbot.MessageContext) error {
        return ctx.Reply("Hello!")
    }).
    Start(":8080")
```

### Full-Featured Bot
```go
bot := slackbot.New().
    WithClientCredentials("id", "secret").
    WithRedisStorage("localhost:6379", "", 0).
    WithMiddleware(
        middleware.LoggingMiddleware(),
        middleware.RateLimitMiddleware(100, time.Minute),
        middleware.ErrorRecoveryMiddleware(),
    ).
    OnMessage(handleMessage).
    OnDirectMessage(handleDirectMessage).
    OnHomeOpened(handleHomeOpened).
    OnButtonClick("approve", handleApprove).
    OnModalSubmit("form", handleForm).
    OnSlashCommand("/status", handleStatus).
    OnReactionAdded(handleReaction).
    OnUserJoined(handleUserJoined)

bot.Start(":8080")
```

## 🏗️ Architecture

### Builder Pattern API
```go
bot := slackbot.New().
    WithClientCredentials("id", "secret").
    WithStorage(storage).
    WithMiddleware(middleware...).
    OnMessage(handler).
    OnButtonClick("action", handler).
    Start(":8080")
```

### Rich Context System
```go
func handleMessage(ctx *slackbot.MessageContext) error {
    user, _ := ctx.GetUser()
    channel, _ := ctx.GetChannel()
    
    return ctx.Reply("Hello " + user.Name + "!")
}
```

### UI Builder System
```go
modal := ui.NewModal("callback", "Title").
    Section("Description").
    Input("field", "Label:", ui.PlainText()).
    ButtonRow(
        ui.Button("save", "Save").Primary(),
        ui.Button("cancel", "Cancel"),
    ).
    Submit("Submit").
    Cancel("Cancel")
```

### Storage Abstraction
```go
// Memory storage (default)
storage := storage.NewMemoryStorage()

// Redis storage
storage, _ := storage.NewRedisStorage("localhost:6379", "", 0)

// PostgreSQL storage
storage, _ := storage.NewPostgresStorage("postgres://...")
```

### Middleware System
```go
bot.WithMiddleware(
    middleware.LoggingMiddleware(),
    middleware.RateLimitMiddleware(100, time.Minute),
    middleware.AuthMiddleware(authFunc),
    middleware.ErrorRecoveryMiddleware(),
)
```

## 📚 Comprehensive Event Support

### Message Events
- `OnMessage` - Any message
- `OnDirectMessage` - Direct messages only
- `OnChannelMessage` - Channel messages only
- `OnThreadMessage` - Thread replies
- `OnMentionMessage` - When bot is mentioned
- `OnMessageContaining` - Messages containing text
- `OnMessageMatching` - Messages matching regex

### Interactive Components
- `OnButtonClick` - Button interactions
- `OnSelectMenu` - Select menu interactions
- `OnModalSubmit` - Modal form submissions
- `OnSlashCommand` - Slash commands

### App Events
- `OnHomeOpened` - App home tab opened
- `OnReactionAdded` - Reaction events
- `OnUserJoined` - User joined workspace
- `OnChannelJoined` - Bot joined channel
- `OnFileShared` - File sharing events

## 🔧 Middleware

### Built-in Middleware
- **LoggingMiddleware** - Request/response logging
- **RateLimitMiddleware** - Rate limiting protection
- **AuthMiddleware** - Authentication and authorization
- **ErrorRecoveryMiddleware** - Panic recovery
- **MetricsMiddleware** - Performance metrics
- **RetryMiddleware** - Automatic retries

### Custom Middleware
```go
func CustomMiddleware() middleware.Middleware {
    return func(next middleware.Handler) middleware.Handler {
        return func() error {
            // Pre-processing
            err := next()
            // Post-processing
            return err
        }
    }
}
```

## 💾 Storage Options

### Memory Storage (Default)
```go
storage := storage.NewMemoryStorage()
```

### Redis Storage
```go
storage, err := storage.NewRedisStorage("localhost:6379", "", 0)
```

### PostgreSQL Storage
```go
storage, err := storage.NewPostgresStorage("postgres://user:pass@localhost/db")
```

### SQLite Storage
```go
storage, err := storage.NewSQLiteStorage("./bot.db")
```

## 🎨 UI Components

### Home Views
```go
home := ui.NewHomeView().
    Header("Welcome!").
    Section("Description").
    ButtonRow(
        ui.Button("action1", "Button 1").Primary(),
        ui.Button("action2", "Button 2"),
    ).
    Divider().
    Section("More content")
```

### Messages
```go
msg := ui.NewMessage().
    Text("Hello!").
    Section("Message content").
    ButtonRow(
        ui.Button("yes", "Yes").Primary(),
        ui.Button("no", "No").Danger(),
    )
```

### Modals
```go
modal := ui.NewModal("callback", "Title").
    Section("Instructions").
    Input("name", "Name:", ui.PlainText().Required()).
    SelectMenu("option", "Choose:", options).
    Submit("Submit").
    Cancel("Cancel")
```

## 🔄 Default Behaviors

### Automatic Defaults
- **In-memory storage** when no storage specified
- **Welcome home tab** with getting started guide
- **Error handlers** for graceful error responses
- **Logging middleware** for request tracking
- **OAuth flow** with automatic token management

### Default Home Tab
```go
// Automatic welcome home tab shows:
// - 🎉 Welcome message
// - Getting started guide
// - Quick start code example
// - Links to documentation
```

## 🚀 Production Deployment

### Docker
```dockerfile
FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o bot
EXPOSE 8080
CMD ["./bot"]
```

### Environment Variables
```bash
SLACK_CLIENT_ID=your_client_id
SLACK_CLIENT_SECRET=your_client_secret
REDIS_URL=redis://localhost:6379
DATABASE_URL=postgres://user:pass@localhost/db
```

### Health Checks
```go
bot.WithMiddleware(middleware.HealthCheckMiddleware("/health"))
```

## 📖 Documentation

### Examples
- [`examples/simple/`](examples/simple/) - 5-line working bot
- [`examples/advanced/`](examples/advanced/) - Full-featured bot
- [`examples/README.md`](examples/README.md) - Comprehensive setup guide

### API Reference
- [`pkg/slackbot/`](pkg/slackbot/) - Core library
- [`pkg/slackbot/storage/`](pkg/slackbot/storage/) - Storage interfaces
- [`pkg/slackbot/ui/`](pkg/slackbot/ui/) - UI components
- [`pkg/slackbot/middleware/`](pkg/slackbot/middleware/) - Middleware system

## 🎯 Success Metrics

### Ease of Use
- ✅ **5-line minimum** - Working bot in 5 lines
- ✅ **Zero configuration** - Works out of the box
- ✅ **Immediate satisfaction** - Default home tab shows success
- ✅ **No dependencies** - In-memory storage by default

### Flexibility
- ✅ **Pluggable storage** - Redis, PostgreSQL, SQLite, memory
- ✅ **Middleware system** - Authentication, rate limiting, logging
- ✅ **25+ events** - Comprehensive Slack event coverage
- ✅ **UI components** - Rich Block Kit component library

### Developer Experience
- ✅ **Builder pattern** - Material UI-level fluent API
- ✅ **Rich context** - Powerful handler context objects
- ✅ **Type safety** - Full compile-time type checking
- ✅ **Comprehensive docs** - Examples, guides, API reference

## 🔗 Slack App Setup

### 1. Create Slack App
1. Go to https://api.slack.com/apps
2. Click "Create New App" → "From scratch"
3. Enter app name and select workspace

### 2. Configure Permissions
Add these **Bot Token Scopes**:
- `app_mentions:read`
- `channels:read`
- `chat:write`
- `im:read`
- `im:write`
- `users:read`
- `reactions:read`
- `files:read`

### 3. Set URLs
- **OAuth Redirect URL**: `https://yourapp.com/slack/oauth`
- **Event Request URL**: `https://yourapp.com/slack/events`
- **Interactive Request URL**: `https://yourapp.com/slack/interactive`

### 4. Subscribe to Events
- `app_home_opened`
- `message.channels`
- `message.im`
- `reaction_added`
- `team_join`
- `member_joined_channel`
- `file_shared`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and examples
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Material UI's ease of use
- Built on the excellent [slack-go/slack](https://github.com/slack-go/slack) library
- Designed for the modern Go ecosystem

---

**Transform your Slack bot development experience with plug-and-play simplicity!** 🚀

Get started with 5 lines of code, scale to enterprise-level complexity. The choice is yours.