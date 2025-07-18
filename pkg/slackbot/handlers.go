package slackbot

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/pkg/slackbot/storage"
)

// Handler function types
type MessageHandler func(ctx *MessageContext) error
type HomeHandler func(ctx *HomeContext) error
type InteractiveHandler func(ctx *InteractiveContext) error
type SlashCommandHandler func(ctx *SlashCommandContext) error
type ReactionHandler func(ctx *ReactionContext) error
type FileHandler func(ctx *FileContext) error
type UserHandler func(ctx *UserContext) error
type ChannelHandler func(ctx *ChannelContext) error
type MiddlewareHandler func(ctx *EventContext, next func()) error

// Base context for all handlers
type Context struct {
	Bot         *Bot
	TeamID      string
	UserID      string
	ChannelID   string
	TriggerID   string
	EventID     string
	Timestamp   string
	Client      *slack.Client
	Storage     storage.Storage
	
	// Cached data
	user    *slack.User
	channel *slack.Channel
	team    *storage.Team
}

// GetUser returns the user associated with this context
func (c *Context) GetUser() (*slack.User, error) {
	if c.user != nil {
		return c.user, nil
	}
	
	user, err := c.Client.GetUserInfo(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	
	c.user = user
	return user, nil
}

// GetChannel returns the channel associated with this context
func (c *Context) GetChannel() (*slack.Channel, error) {
	if c.channel != nil {
		return c.channel, nil
	}
	
	channel, err := c.Client.GetConversationInfo(c.ChannelID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel info: %w", err)
	}
	
	c.channel = channel
	return channel, nil
}

// GetTeam returns the team associated with this context
func (c *Context) GetTeam() (*storage.Team, error) {
	if c.team != nil {
		return c.team, nil
	}
	
	team, err := c.Storage.GetTeam(c.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	
	c.team = team
	return team, nil
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
	RawMessage  *slack.MessageEvent
}

// Reply sends a reply message
func (ctx *MessageContext) Reply(text string) error {
	_, _, err := ctx.Client.PostMessage(ctx.ChannelID, slack.MsgOptionText(text, false))
	return err
}

// ReplyInThread sends a reply in a thread
func (ctx *MessageContext) ReplyInThread(text string) error {
	threadTS := ctx.ThreadTS
	if threadTS == "" {
		threadTS = ctx.Timestamp
	}
	
	_, _, err := ctx.Client.PostMessage(ctx.ChannelID, 
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(threadTS))
	return err
}

// ReplyEphemeral sends an ephemeral reply
func (ctx *MessageContext) ReplyEphemeral(text string) error {
	_, err := ctx.Client.PostEphemeral(ctx.ChannelID, ctx.UserID, slack.MsgOptionText(text, false))
	return err
}

// React adds a reaction to the message
func (ctx *MessageContext) React(emoji string) error {
	return ctx.Client.AddReaction(emoji, slack.ItemRef{
		Channel:   ctx.ChannelID,
		Timestamp: ctx.Timestamp,
	})
}

// DeleteMessage deletes the message
func (ctx *MessageContext) DeleteMessage() error {
	_, _, err := ctx.Client.DeleteMessage(ctx.ChannelID, ctx.Timestamp)
	return err
}

// UpdateMessage updates the message
func (ctx *MessageContext) UpdateMessage(text string) error {
	_, _, _, err := ctx.Client.UpdateMessage(ctx.ChannelID, ctx.Timestamp, slack.MsgOptionText(text, false))
	return err
}

// ReplyWithBlocks sends a reply with blocks
func (ctx *MessageContext) ReplyWithBlocks(blocks []slack.Block) error {
	_, _, err := ctx.Client.PostMessage(ctx.ChannelID, slack.MsgOptionBlocks(blocks...))
	return err
}

// StartThread starts a new thread
func (ctx *MessageContext) StartThread(text string) error {
	_, _, err := ctx.Client.PostMessage(ctx.ChannelID, 
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(ctx.Timestamp))
	return err
}

// HomeContext provides context for home events
type HomeContext struct {
	*Context
	Tab     string
	IsFirst bool
	View    *slack.HomeTabViewRequest
}

// PublishView publishes a view to the home tab
func (ctx *HomeContext) PublishView(view *slack.HomeTabViewRequest) error {
	_, err := ctx.Client.PublishView(ctx.UserID, *view, "")
	return err
}

// UpdateView updates the home view
func (ctx *HomeContext) UpdateView(view *slack.HomeTabViewRequest) error {
	_, err := ctx.Client.PublishView(ctx.UserID, *view, "")
	return err
}

// OpenModal opens a modal
func (ctx *HomeContext) OpenModal(modal slack.ModalViewRequest) error {
	_, err := ctx.Client.OpenView(ctx.TriggerID, modal)
	return err
}

// InteractiveContext provides context for interactive events
type InteractiveContext struct {
	*Context
	ActionID     string
	ActionValue  string
	ActionType   string
	BlockID      string
	ViewID       string
	ResponseURL  string
	Interaction  slack.InteractionCallback
}

// UpdateMessage updates the message that triggered the interaction
func (ctx *InteractiveContext) UpdateMessage(text string) error {
	_, _, _, err := ctx.Client.UpdateMessage(ctx.ChannelID, ctx.Timestamp, slack.MsgOptionText(text, false))
	return err
}

// OpenModal opens a modal
func (ctx *InteractiveContext) OpenModal(modal slack.ModalViewRequest) error {
	_, err := ctx.Client.OpenView(ctx.TriggerID, modal)
	return err
}

// UpdateModal updates an existing modal
func (ctx *InteractiveContext) UpdateModal(modal slack.ModalViewRequest) error {
	_, err := ctx.Client.UpdateView(modal, "", "", ctx.ViewID)
	return err
}

// Acknowledge sends an acknowledgment
func (ctx *InteractiveContext) Acknowledge() error {
	return nil // Slack expects empty 200 response
}

// RespondEphemeral sends an ephemeral response
func (ctx *InteractiveContext) RespondEphemeral(text string) error {
	_, err := ctx.Client.PostEphemeral(ctx.ChannelID, ctx.UserID, slack.MsgOptionText(text, false))
	return err
}

// SlashCommandContext provides context for slash command events
type SlashCommandContext struct {
	*Context
	Command     string
	Text        string
	ResponseURL string
	SlashCommand slack.SlashCommand
}

// Respond sends a response to the slash command
func (ctx *SlashCommandContext) Respond(text string) error {
	_, _, err := ctx.Client.PostMessage(ctx.ChannelID, slack.MsgOptionText(text, false))
	return err
}

// RespondEphemeral sends an ephemeral response
func (ctx *SlashCommandContext) RespondEphemeral(text string) error {
	_, err := ctx.Client.PostEphemeral(ctx.ChannelID, ctx.UserID, slack.MsgOptionText(text, false))
	return err
}

// RespondDelayed sends a delayed response
func (ctx *SlashCommandContext) RespondDelayed(text string) error {
	// This would typically use the response URL for delayed responses
	return ctx.Respond(text)
}

// ReactionContext provides context for reaction events
type ReactionContext struct {
	*Context
	Reaction string
	ItemUser string
	Item     slack.Item
}

// FileContext provides context for file events
type FileContext struct {
	*Context
	File   slack.File
	Action string
}

// UserContext provides context for user events
type UserContext struct {
	*Context
	User   *slack.User
	Action string
}

// ChannelContext provides context for channel events
type ChannelContext struct {
	*Context
	Channel *slack.Channel
	Action  string
}

// EventContext provides context for any event (used in middleware)
type EventContext struct {
	*Context
	EventType string
	RawEvent  interface{}
}

// Utility functions for context creation
func (bot *Bot) newContext(teamID, userID, channelID, triggerID, eventID, timestamp string) *Context {
	// Get team from storage to create client
	team, err := bot.storage.GetTeam(teamID)
	if err != nil {
		// Handle error appropriately
		return nil
	}
	
	client := slack.New(team.AccessToken)
	
	return &Context{
		Bot:       bot,
		TeamID:    teamID,
		UserID:    userID,
		ChannelID: channelID,
		TriggerID: triggerID,
		EventID:   eventID,
		Timestamp: timestamp,
		Client:    client,
		Storage:   bot.storage,
	}
}

func (bot *Bot) newMessageContext(baseCtx *Context, event *slack.MessageEvent) *MessageContext {
	return &MessageContext{
		Context:     baseCtx,
		Text:        event.Text,
		ThreadTS:    event.ThreadTimeStamp,
		IsThread:    event.ThreadTimeStamp != "",
		IsDirect:    event.Channel[0] == 'D',
		IsChannel:   event.Channel[0] == 'C',
		IsMention:   false, // This would need to be determined from the text
		Files:       event.Files,
		Attachments: event.Attachments,
		RawMessage:  event,
	}
}

func (bot *Bot) newHomeContext(baseCtx *Context, event *slack.AppHomeOpenedEvent) *HomeContext {
	return &HomeContext{
		Context: baseCtx,
		Tab:     event.Tab,
		IsFirst: false, // This would need to be determined
	}
}

func (bot *Bot) newInteractiveContext(baseCtx *Context, interaction slack.InteractionCallback) *InteractiveContext {
	var actionID, actionValue, actionType, blockID string
	
	if len(interaction.ActionCallback.BlockActions) > 0 {
		action := interaction.ActionCallback.BlockActions[0]
		actionID = action.ActionID
		actionValue = action.Value
		actionType = action.Type
		blockID = action.BlockID
	}
	
	return &InteractiveContext{
		Context:     baseCtx,
		ActionID:    actionID,
		ActionValue: actionValue,
		ActionType:  actionType,
		BlockID:     blockID,
		ViewID:      interaction.View.ID,
		ResponseURL: interaction.ResponseURL,
		Interaction: interaction,
	}
}

func (bot *Bot) newSlashCommandContext(baseCtx *Context, command slack.SlashCommand) *SlashCommandContext {
	return &SlashCommandContext{
		Context:      baseCtx,
		Command:      command.Command,
		Text:         command.Text,
		ResponseURL:  command.ResponseURL,
		SlashCommand: command,
	}
}

func (bot *Bot) newReactionContext(baseCtx *Context, event *slack.ReactionAddedEvent) *ReactionContext {
	return &ReactionContext{
		Context:  baseCtx,
		Reaction: event.Reaction,
		ItemUser: event.ItemUser,
		Item:     event.Item,
	}
}

func (bot *Bot) newEventContext(baseCtx *Context, eventType string, rawEvent interface{}) *EventContext {
	return &EventContext{
		Context:   baseCtx,
		EventType: eventType,
		RawEvent:  rawEvent,
	}
}