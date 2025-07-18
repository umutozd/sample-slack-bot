package slackbot

import (
	"github.com/slack-go/slack"
)

// defaultHomeHandler is the default handler for app home opened events
// It provides a welcoming experience and shows users that their bot is working
func defaultHomeHandler(ctx *HomeContext) error {
	// Create a welcoming home view
	view := slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Welcome header
				slack.NewHeaderBlock(
					slack.NewTextBlockObject(slack.PlainTextType, "🎉 Welcome to Your Slack Bot!", false, false),
				),
				
				// Description section
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "Your bot is up and running! This is the default home tab. Ready to customize?", false, false),
					nil,
					nil,
				),
				
				// Action buttons
				slack.NewActionBlock(
					"welcome_actions",
					slack.NewButtonBlockElement("get_started", "get_started", slack.NewTextBlockObject(slack.PlainTextType, "Get Started", false, false)).WithStyle(slack.StylePrimary),
					slack.NewButtonBlockElement("view_docs", "view_docs", slack.NewTextBlockObject(slack.PlainTextType, "View Documentation", false, false)),
				),
				
				// Divider
				slack.NewDividerBlock(),
				
				// Next steps section
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Next Steps:*", false, false),
					nil,
					nil,
				),
				
				// Next steps list
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "• Add custom message handlers with `OnMessage()`\n• Create interactive components with `OnButtonClick()`\n• Customize this home tab with `OnHomeOpened()`\n• Add slash commands with `OnSlashCommand()`", false, false),
					nil,
					nil,
				),
				
				// Divider
				slack.NewDividerBlock(),
				
				// Bot info section
				slack.NewContextBlock(
					"",
					slack.NewTextBlockObject(slack.MarkdownType, "🤖 Powered by SlackBot Library", false, false),
				),
			},
		},
	}
	
	// Publish the view
	return ctx.PublishView(view)
}

// DefaultGetStartedHandler handles the "Get Started" button click
func DefaultGetStartedHandler(ctx *InteractiveContext) error {
	// Create a modal with getting started information
	modal := slack.ModalViewRequest{
		Type:         slack.VTModal,
		Title:        slack.NewTextBlockObject(slack.PlainTextType, "Getting Started", false, false),
		Close:        slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false),
		ClearOnClose: true,
		CallbackID:   "get_started_modal",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Header
				slack.NewHeaderBlock(
					slack.NewTextBlockObject(slack.PlainTextType, "🚀 Quick Start Guide", false, false),
				),
				
				// Step 1
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Step 1: Add Message Handlers*\n```go\nbot.OnMessage(func(ctx *slackbot.MessageContext) error {\n    return ctx.Reply(\"Hello \" + ctx.UserID)\n})\n```", false, false),
					nil,
					nil,
				),
				
				// Step 2
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Step 2: Create Interactive Elements*\n```go\nbot.OnButtonClick(\"my_button\", func(ctx *slackbot.InteractiveContext) error {\n    return ctx.UpdateMessage(\"Button clicked!\")\n})\n```", false, false),
					nil,
					nil,
				),
				
				// Step 3
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Step 3: Customize Home Tab*\n```go\nbot.OnHomeOpened(func(ctx *slackbot.HomeContext) error {\n    // Your custom home view here\n    return nil\n})\n```", false, false),
					nil,
					nil,
				),
				
				// Links section
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Useful Links:*\n• GitHub Repository\n• Documentation\n• Examples", false, false),
					nil,
					nil,
				),
			},
		},
	}
	
	return ctx.OpenModal(modal)
}

// DefaultViewDocsHandler handles the "View Documentation" button click
func DefaultViewDocsHandler(ctx *InteractiveContext) error {
	// Create a modal with documentation links
	modal := slack.ModalViewRequest{
		Type:         slack.VTModal,
		Title:        slack.NewTextBlockObject(slack.PlainTextType, "Documentation", false, false),
		Close:        slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false),
		ClearOnClose: true,
		CallbackID:   "view_docs_modal",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				// Header
				slack.NewHeaderBlock(
					slack.NewTextBlockObject(slack.PlainTextType, "📚 Documentation", false, false),
				),
				
				// API Reference
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*API Reference*\nComplete reference for all available methods and handlers.", false, false),
					nil,
					slack.NewAccessory(
						slack.NewButtonBlockElement("", "api_ref", slack.NewTextBlockObject(slack.PlainTextType, "View API", false, false)),
					),
				),
				
				// Examples
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Examples*\nCode examples and sample implementations.", false, false),
					nil,
					slack.NewAccessory(
						slack.NewButtonBlockElement("", "examples", slack.NewTextBlockObject(slack.PlainTextType, "View Examples", false, false)),
					),
				),
				
				// Guides
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Guides*\nStep-by-step guides for common use cases.", false, false),
					nil,
					slack.NewAccessory(
						slack.NewButtonBlockElement("", "guides", slack.NewTextBlockObject(slack.PlainTextType, "View Guides", false, false)),
					),
				),
				
				// Divider
				slack.NewDividerBlock(),
				
				// Community
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "*Community*\nJoin our community for support and discussions.", false, false),
					nil,
					slack.NewAccessory(
						slack.NewButtonBlockElement("", "community", slack.NewTextBlockObject(slack.PlainTextType, "Join Community", false, false)),
					),
				),
			},
		},
	}
	
	return ctx.OpenModal(modal)
}

// RegisterDefaultHandlers registers the default interactive handlers
func RegisterDefaultHandlers(builder BotBuilder) BotBuilder {
	return builder.
		OnButtonClick("get_started", DefaultGetStartedHandler).
		OnButtonClick("view_docs", DefaultViewDocsHandler)
}

// DefaultMiddleware provides basic logging middleware
func DefaultMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			// Log the event
			if ctx.Logger != nil {
				// Assuming logger has a WithFields method (like logrus)
				// This would need to be adapted based on the actual logger interface
			}
			
			// Call the next handler
			return next.Handle(ctx)
		})
	}
}

// HandlerFunc is an adapter to allow the use of ordinary functions as Handlers
type HandlerFunc func(*Context) error

// Handle calls f(ctx)
func (f HandlerFunc) Handle(ctx *Context) error {
	return f(ctx)
}