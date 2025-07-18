# SlackBot Library Examples

This directory contains examples demonstrating how to use the SlackBot library.

## Simple Example

The `simple/main.go` example shows the most basic usage of the library:

```go
bot := slackbot.New().
    WithClientCredentials(clientID, clientSecret).
    WithDebug(true).
    OnMessage(func(ctx *slackbot.MessageContext) error {
        return ctx.Reply("Hello " + ctx.UserID + "! You said: " + ctx.Text)
    }).
    OnDirectMessage(func(ctx *slackbot.MessageContext) error {
        return ctx.Reply("Thanks for the direct message!")
    }).
    OnButtonClick("hello", func(ctx *slackbot.InteractiveContext) error {
        return ctx.UpdateMessage("Button clicked! 🎉")
    })

// Start the bot
bot.Start(":8080")
```

## Advanced Example

The `advanced/main.go` example demonstrates more sophisticated features:

- Custom storage backends
- Middleware for logging and rate limiting
- Pattern-based message handlers
- Interactive components (buttons, modals)
- Custom home tab
- Slash commands
- UI builder for creating rich messages

## Running the Examples

1. Set your environment variables:
   ```bash
   export SLACK_CLIENT_ID="your-client-id"
   export SLACK_CLIENT_SECRET="your-client-secret"
   ```

2. Run an example:
   ```bash
   cd examples/simple
   go run main.go
   ```

3. The bot will start on port 8080. Make sure to configure your Slack app's event subscriptions to point to your server.

## Key Features Demonstrated

### Builder Pattern
The library uses a fluent builder pattern for easy configuration:
```go
bot := slackbot.New().
    WithClientCredentials(clientID, clientSecret).
    WithStorage(storage.NewMemoryStorage()).
    WithMiddleware(middleware...).
    OnMessage(handler).
    OnHomeOpened(handler)
```

### Default Behaviors
- **Default Storage**: Uses in-memory storage when no storage is specified
- **Default Home Tab**: Shows a welcome message with getting started information
- **Default Middleware**: Includes basic logging and error recovery

### Event Handling
- Message events: `OnMessage`, `OnDirectMessage`, `OnChannelMessage`, `OnThreadMessage`, `OnMentionMessage`
- Interactive events: `OnButtonClick`, `OnSelectMenu`, `OnModalSubmit`
- Pattern matching: `OnMessageContaining`, `OnMessageMatching`
- Global events: `OnAnyEvent`

### UI Builder
Create rich Slack interfaces easily:
```go
homeView := ui.NewHomeView().
    Header("Welcome!").
    Section("Description text").
    ButtonRow(
        ui.Button("action1", "Click Me").Primary(),
        ui.Button("action2", "Or Me"),
    ).
    Divider().
    Section("More content")
```

### Middleware
Add cross-cutting concerns:
```go
bot.WithMiddleware(
    slackbot.LoggingMiddleware(),
    slackbot.RateLimitMiddleware(60, time.Minute),
    slackbot.ErrorRecoveryMiddleware(),
)
```

### Storage Options
- **Memory Storage**: `storage.NewMemoryStorage()` (default)
- **Redis Storage**: `storage.NewRedisStorage(addr, password, db)` (placeholder)
- **PostgreSQL Storage**: `storage.NewPostgresStorage(dsn)` (placeholder)
- **SQLite Storage**: `storage.NewSQLiteStorage(filepath)` (placeholder)
- **Layered Storage**: `storage.NewLayeredStorage(cache, persistent)` (cache + persistent)

## Next Steps

1. **Customize the home tab**: Replace the default home handler with your own
2. **Add more event handlers**: Handle additional Slack events
3. **Create interactive workflows**: Use modals and buttons for complex interactions
4. **Add persistence**: Configure Redis or PostgreSQL for production use
5. **Implement business logic**: Add your application-specific functionality
6. **Add monitoring**: Use middleware for metrics and logging
7. **Scale up**: Add load balancing and horizontal scaling as needed

## Material UI-Style Experience

The library aims to provide a Material UI-level experience for Slack bot development:

- **Easy to start**: 5 lines of code for a working bot
- **Plug-and-play**: Works out of the box with sensible defaults
- **Highly configurable**: Customize every aspect when needed
- **Rich components**: Pre-built UI components for common patterns
- **Powerful patterns**: Support for complex interactive workflows
- **Production-ready**: Built-in middleware, error handling, and scaling support