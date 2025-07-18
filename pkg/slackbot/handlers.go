package slackbot

import (
	"time"

	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
)

// Handler function types for different Slack events and interactions

// MessageHandler handles incoming messages
type MessageHandler func(ctx *MessageContext) error

// HomeHandler handles app home opened events
type HomeHandler func(ctx *HomeContext) error

// InteractiveHandler handles interactive component interactions (buttons, menus, modals)
type InteractiveHandler func(ctx *InteractiveContext) error

// SlashCommandHandler handles slash commands
type SlashCommandHandler func(ctx *SlashCommandContext) error

// EventHandler handles any Slack event (middleware-style)
type EventHandler func(ctx *EventContext) error

// Middleware is a function that wraps handlers with additional functionality
type Middleware func(next Handler) Handler

// Handler is the base handler interface
type Handler interface {
	Handle(ctx *Context) error
}

// Context represents the base context for all handlers
type Context struct {
	Bot        Bot
	TeamID     string
	UserID     string
	ChannelID  string
	TriggerID  string
	EventID    string
	Timestamp  string
	Client     *slack.Client
	Storage    storage.Storage
	Logger     interface{} // Using interface{} to avoid logrus dependency
	RequestID  string
	StartTime  time.Time
	
	// Internal data
	data map[string]interface{}
}

// MessageContext provides context for message events
type MessageContext struct {
	*Context
	Text        string
	ThreadTS    string
	IsThread    bool
	IsDirect    bool
	IsChannel   bool
	IsMention   bool
	Files       []slack.File
	Attachments []slack.Attachment
	Blocks      []slack.Block
	User        *slack.User
	Channel     *slack.Channel
	Team        *slack.Team
}

// HomeContext provides context for home tab events
type HomeContext struct {
	*Context
	Tab     string
	IsFirst bool
	User    *slack.User
}

// InteractiveContext provides context for interactive component events
type InteractiveContext struct {
	*Context
	ActionID    string
	ActionValue string
	ActionType  string
	BlockID     string
	ViewID      string
	CallbackID  string
	User        *slack.User
	Channel     *slack.Channel
	Message     *slack.Message
	View        *slack.View
	
	// For button clicks
	Button *slack.ButtonBlockElement
	
	// For select menus
	SelectedOption  *slack.OptionBlockObject
	SelectedOptions []*slack.OptionBlockObject
	
	// For modal submissions
	StateValues map[string]map[string]slack.BlockAction
}

// SlashCommandContext provides context for slash command events
type SlashCommandContext struct {
	*Context
	Command     string
	Text        string
	ResponseURL string
	User        *slack.User
	Channel     *slack.Channel
}

// EventContext provides context for any Slack event
type EventContext struct {
	*Context
	EventType string
	EventData interface{}
	RawEvent  map[string]interface{}
}

// Context helper methods

// Set stores a value in the context
func (c *Context) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

// Get retrieves a value from the context
func (c *Context) Get(key string) (interface{}, bool) {
	if c.data == nil {
		return nil, false
	}
	value, exists := c.data[key]
	return value, exists
}

// GetString retrieves a string value from the context
func (c *Context) GetString(key string) (string, bool) {
	if value, exists := c.Get(key); exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetUser retrieves the user information
func (c *Context) GetUser() (*slack.User, error) {
	if c.UserID == "" {
		return nil, nil
	}
	return c.Client.GetUserInfo(c.UserID)
}

// GetChannel retrieves the channel information
func (c *Context) GetChannel() (*slack.Channel, error) {
	if c.ChannelID == "" {
		return nil, nil
	}
	return c.Client.GetChannelInfo(c.ChannelID)
}

// GetTeam retrieves the team information from storage
func (c *Context) GetTeam() (*storage.Team, error) {
	if c.TeamID == "" {
		return nil, nil
	}
	return c.Storage.GetTeam(c.TeamID)
}

// MessageContext helper methods

// Reply sends a reply message to the same channel
func (mc *MessageContext) Reply(text string) error {
	_, _, err := mc.Client.PostMessage(mc.ChannelID, slack.MsgOptionText(text, false))
	return err
}

// ReplyWithBlocks sends a reply with Block Kit blocks
func (mc *MessageContext) ReplyWithBlocks(blocks []slack.Block) error {
	_, _, err := mc.Client.PostMessage(mc.ChannelID, slack.MsgOptionBlocks(blocks...))
	return err
}

// ReplyInThread sends a reply in the thread
func (mc *MessageContext) ReplyInThread(text string) error {
	threadTS := mc.ThreadTS
	if threadTS == "" {
		threadTS = mc.Timestamp
	}
	_, _, err := mc.Client.PostMessage(mc.ChannelID, 
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(threadTS),
	)
	return err
}

// ReplyEphemeral sends an ephemeral reply visible only to the user
func (mc *MessageContext) ReplyEphemeral(text string) error {
	_, err := mc.Client.PostEphemeral(mc.ChannelID, mc.UserID, slack.MsgOptionText(text, false))
	return err
}

// React adds a reaction to the message
func (mc *MessageContext) React(emoji string) error {
	return mc.Client.AddReaction(emoji, slack.ItemRef{
		Channel:   mc.ChannelID,
		Timestamp: mc.Timestamp,
	})
}

// DeleteMessage deletes the message
func (mc *MessageContext) DeleteMessage() error {
	_, _, err := mc.Client.DeleteMessage(mc.ChannelID, mc.Timestamp)
	return err
}

// UpdateMessage updates the message content
func (mc *MessageContext) UpdateMessage(text string) error {
	_, _, _, err := mc.Client.UpdateMessage(mc.ChannelID, mc.Timestamp, slack.MsgOptionText(text, false))
	return err
}

// StartThread starts a new thread from this message
func (mc *MessageContext) StartThread(text string) error {
	_, _, err := mc.Client.PostMessage(mc.ChannelID,
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(mc.Timestamp),
	)
	return err
}

// HomeContext helper methods

// PublishView publishes a view to the app home
func (hc *HomeContext) PublishView(view slack.HomeTabViewRequest) error {
	_, err := hc.Client.PublishView(hc.UserID, view, "")
	return err
}

// UpdateView updates the app home view
func (hc *HomeContext) UpdateView(view slack.HomeTabViewRequest) error {
	_, err := hc.Client.PublishView(hc.UserID, view, "")
	return err
}

// OpenModal opens a modal dialog
func (hc *HomeContext) OpenModal(modal slack.ModalViewRequest) error {
	_, err := hc.Client.OpenView(hc.TriggerID, modal)
	return err
}

// InteractiveContext helper methods

// UpdateMessage updates the message that triggered the interaction
func (ic *InteractiveContext) UpdateMessage(text string) error {
	if ic.Message == nil {
		return nil
	}
	_, _, _, err := ic.Client.UpdateMessage(ic.ChannelID, ic.Message.Timestamp, slack.MsgOptionText(text, false))
	return err
}

// UpdateMessageWithBlocks updates the message with blocks
func (ic *InteractiveContext) UpdateMessageWithBlocks(blocks []slack.Block) error {
	if ic.Message == nil {
		return nil
	}
	_, _, _, err := ic.Client.UpdateMessage(ic.ChannelID, ic.Message.Timestamp, slack.MsgOptionBlocks(blocks...))
	return err
}

// OpenModal opens a modal dialog
func (ic *InteractiveContext) OpenModal(modal slack.ModalViewRequest) error {
	_, err := ic.Client.OpenView(ic.TriggerID, modal)
	return err
}

// UpdateModal updates the current modal
func (ic *InteractiveContext) UpdateModal(modal slack.ModalViewRequest) error {
	_, err := ic.Client.UpdateView(modal, "", "", ic.ViewID)
	return err
}

// Acknowledge acknowledges the interaction
func (ic *InteractiveContext) Acknowledge() error {
	// This is handled automatically by the library
	return nil
}

// RespondEphemeral sends an ephemeral response
func (ic *InteractiveContext) RespondEphemeral(text string) error {
	_, err := ic.Client.PostEphemeral(ic.ChannelID, ic.UserID, slack.MsgOptionText(text, false))
	return err
}

// SlashCommandContext helper methods

// Respond responds to the slash command
func (sc *SlashCommandContext) Respond(text string) error {
	_, _, err := sc.Client.PostMessage(sc.ChannelID, slack.MsgOptionText(text, false))
	return err
}

// RespondEphemeral responds with an ephemeral message
func (sc *SlashCommandContext) RespondEphemeral(text string) error {
	_, err := sc.Client.PostEphemeral(sc.ChannelID, sc.UserID, slack.MsgOptionText(text, false))
	return err
}

// RespondWithBlocks responds with Block Kit blocks
func (sc *SlashCommandContext) RespondWithBlocks(blocks []slack.Block) error {
	_, _, err := sc.Client.PostMessage(sc.ChannelID, slack.MsgOptionBlocks(blocks...))
	return err
}

// OpenModal opens a modal dialog
func (sc *SlashCommandContext) OpenModal(modal slack.ModalViewRequest) error {
	_, err := sc.Client.OpenView(sc.TriggerID, modal)
	return err
}