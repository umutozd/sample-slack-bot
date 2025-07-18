# SlackBot Library Examples

This directory contains examples showing how to use the SlackBot library to build Slack bots with varying levels of complexity.

## 🚀 Quick Start

### Simple Bot (5 lines of code!)

The `simple/` example shows how to create a working Slack bot with just 5 lines of code:

```go
slackbot.New().
    WithClientCredentials("client_id", "client_secret").
    OnMessage(func(ctx *slackbot.MessageContext) error {
        return ctx.Reply("Hello " + ctx.UserID + "! You said: " + ctx.Text)
    }).
    Start(":8080")
```

### Advanced Bot (Full-featured)

The `advanced/` example demonstrates a comprehensive bot with:

- Custom middleware (logging, rate limiting, error recovery)
- Multiple event handlers (messages, interactions, reactions)
- UI components (modals, buttons, home tabs)
- Storage integration (Redis with memory fallback)
- Slash commands
- Pattern matching and conditional responses

## 📋 Prerequisites

1. **Slack App**: Create a Slack app at https://api.slack.com/apps
2. **Bot Token**: Enable bot functionality and get your tokens
3. **Environment Variables**: Set up your credentials

## 🔧 Setup Instructions

### 1. Create a Slack App

1. Go to https://api.slack.com/apps and click "Create New App"
2. Choose "From scratch"
3. Give your app a name and select your workspace
4. Click "Create App"

### 2. Configure OAuth & Permissions

1. Go to "OAuth & Permissions" in your app settings
2. Add the following **Bot Token Scopes**:
   - `app_mentions:read` - View messages that directly mention your bot
   - `channels:read` - View basic information about public channels
   - `chat:write` - Send messages as your bot
   - `im:read` - View messages in direct messages
   - `im:write` - Start direct messages with people
   - `users:read` - View people in your workspace
   - `reactions:read` - View emoji reactions and their associated content
   - `files:read` - View files shared in channels and conversations

3. Add **Redirect URLs**:
   - `http://localhost:8080/slack/oauth` (for local development)
   - Your production URL when deploying

### 3. Set Environment Variables

```bash
export SLACK_CLIENT_ID="your_client_id_here"
export SLACK_CLIENT_SECRET="your_client_secret_here"
export REDIS_URL="redis://localhost:6379"  # Optional, for advanced example
```

### 4. Configure Event Subscriptions

1. Go to "Event Subscriptions" in your app settings
2. Enable Events
3. Set Request URL to: `http://localhost:8080/slack/events`
4. Subscribe to these **bot events**:
   - `app_home_opened`
   - `message.channels`
   - `message.im`
   - `reaction_added`
   - `team_join`
   - `member_joined_channel`
   - `file_shared`

### 5. Configure Interactive Components

1. Go to "Interactivity & Shortcuts" in your app settings
2. Enable Interactivity
3. Set Request URL to: `http://localhost:8080/slack/interactive`

### 6. Configure Slash Commands (Optional)

1. Go to "Slash Commands" in your app settings
2. Create new commands:
   - Command: `/status`
   - Request URL: `http://localhost:8080/slack/slash`
   - Description: "Check bot status"
   
   - Command: `/config`
   - Request URL: `http://localhost:8080/slack/slash`
   - Description: "Configure bot settings"

### 7. Enable App Home (Optional)

1. Go to "App Home" in your app settings
2. Enable "Home Tab"
3. Enable "Messages Tab"

## 🏃‍♂️ Running the Examples

### Simple Example

```bash
cd examples/simple
go run main.go
```

### Advanced Example

```bash
cd examples/advanced
go run main.go
```

### Install to Slack

1. Start your bot
2. Visit: `http://localhost:8080/slack/install`
3. Authorize the app in your Slack workspace
4. Start interacting with your bot!

## 🎯 What Each Example Demonstrates

### Simple Bot (`simple/main.go`)

- **Minimal Setup**: Just 5 lines of code
- **Message Handling**: Responds to all messages
- **Button Interactions**: Handles button clicks
- **Default Storage**: Uses in-memory storage automatically
- **Default Home Tab**: Shows welcome message automatically

**Features:**
- ✅ Responds to any message
- ✅ Handles button clicks
- ✅ Automatic OAuth flow
- ✅ Default home tab
- ✅ Zero configuration required

### Advanced Bot (`advanced/main.go`)

- **Full Middleware Stack**: Logging, rate limiting, error recovery, metrics
- **Multiple Event Types**: Messages, reactions, file sharing, user events
- **Rich UI Components**: Modals, buttons, select menus, home tabs
- **Storage Integration**: Redis with fallback to memory
- **Pattern Matching**: Regex and text-based message routing
- **Slash Commands**: Custom commands with modal responses
- **User Preferences**: Persistent settings and state management

**Features:**
- ✅ All simple bot features
- ✅ Redis storage with memory fallback
- ✅ Custom middleware pipeline
- ✅ Rich interactive components
- ✅ Pattern-based message handling
- ✅ Slash command support
- ✅ User state management
- ✅ Reaction tracking
- ✅ File sharing notifications
- ✅ Welcome flow for new users
- ✅ Comprehensive error handling

## 🧪 Testing Your Bot

### Basic Tests

1. **Message Test**: Send a message to your bot
2. **Mention Test**: Mention your bot in a channel
3. **Button Test**: Click buttons in bot responses
4. **Home Tab Test**: Open the app home tab
5. **Slash Command Test**: Use `/status` or `/config`

### Advanced Tests

1. **Modal Test**: Click "Show Modal" and submit form
2. **Pattern Test**: Say "weather in Seattle"
3. **Help Test**: Say "help" to see commands
4. **Reaction Test**: Add 🎉 reaction to a message
5. **File Test**: Share a file in a channel

## 📚 Key Concepts

### Builder Pattern

The library uses a fluent builder pattern for easy configuration:

```go
bot := slackbot.New().
    WithClientCredentials("id", "secret").
    WithStorage(storage).
    WithMiddleware(middleware...).
    OnMessage(handler).
    OnButtonClick("action", handler)
```

### Context-Rich Handlers

All handlers receive rich context objects with helper methods:

```go
func handleMessage(ctx *slackbot.MessageContext) error {
    user, _ := ctx.GetUser()
    channel, _ := ctx.GetChannel()
    
    return ctx.Reply("Hello " + user.Name + "!")
}
```

### UI Builder

Create rich UI components with a fluent API:

```go
modal := ui.NewModal("callback", "Title").
    Section("Description").
    Input("field", "Label:", ui.PlainText()).
    Submit("Save").
    Cancel("Cancel")
```

### Middleware

Add cross-cutting concerns like logging and rate limiting:

```go
bot.WithMiddleware(
    middleware.LoggingMiddleware(),
    middleware.RateLimitMiddleware(100, time.Minute),
    middleware.ErrorRecoveryMiddleware(),
)
```

### Storage Abstraction

Use different storage backends seamlessly:

```go
// Memory storage (default)
storage := storage.NewMemoryStorage()

// Redis storage
storage, _ := storage.NewRedisStorage("localhost:6379", "", 0)

// PostgreSQL storage (coming soon)
storage, _ := storage.NewPostgresStorage("postgres://...")
```

## 🚀 Next Steps

1. **Customize**: Modify the examples to fit your use case
2. **Deploy**: Use Docker, Heroku, or your preferred platform
3. **Scale**: Add Redis/PostgreSQL for production
4. **Extend**: Add custom middleware and handlers
5. **Share**: Distribute your bot to other workspaces

## 🆘 Troubleshooting

### Common Issues

1. **"Invalid signing secret"**: Make sure your signing secret is correct
2. **"Token not found"**: Complete the OAuth flow first
3. **"Events not received"**: Check your event subscription URL
4. **"Buttons not working"**: Verify interactive components URL
5. **"Storage errors"**: Check Redis connection or use memory storage

### Debug Mode

Enable debug logging to see what's happening:

```go
bot.WithDebug(true)
```

### Logs

Check the console output for detailed logs:

```
🚀 Bot server starting on port 8080
📊 Using Redis storage
🔄 Request started
✅ Request completed in 45ms
```

## 📖 API Reference

See the main README.md for complete API documentation and advanced usage patterns.

## 💡 Tips

1. **Start Simple**: Begin with the simple example and add features gradually
2. **Use Middleware**: Leverage built-in middleware for common patterns
3. **Handle Errors**: Always handle errors gracefully in your handlers
4. **Test Locally**: Use ngrok or similar tools for local testing
5. **Monitor Usage**: Use the metrics middleware to track performance
6. **Cache Data**: Use storage for frequently accessed data
7. **Be Responsive**: Keep handlers fast and use async patterns for long operations

Happy bot building! 🤖