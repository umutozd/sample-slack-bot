package slackbot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// handleInstall handles the OAuth installation flow
func (bot *Bot) handleInstall(w http.ResponseWriter, r *http.Request) {
	// Create OAuth URL
	oauthURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=app_mentions:read,channels:read,chat:write,im:read,im:write,users:read&user_scope=",
		bot.clientID,
	)
	
	http.Redirect(w, r, oauthURL, http.StatusTemporaryRedirect)
}

// handleOAuth handles the OAuth callback
func (bot *Bot) handleOAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}
	
	// Exchange code for token
	resp, err := slack.GetOAuthV2Response(http.DefaultClient, bot.clientID, bot.clientSecret, code, "")
	if err != nil {
		log.Printf("Error getting OAuth token: %v", err)
		http.Error(w, "Failed to get OAuth token", http.StatusInternalServerError)
		return
	}
	
	// Store team information
	team := &storage.Team{
		ID:           resp.Team.ID,
		Name:         resp.Team.Name,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		BotUserID:    resp.BotUserID,
		Scope:        resp.Scope,
	}
	
	if err := bot.storage.StoreTeam(team); err != nil {
		log.Printf("Error storing team: %v", err)
		http.Error(w, "Failed to store team", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<html>
		<head><title>Installation Complete</title></head>
		<body>
			<h1>🎉 Installation Complete!</h1>
			<p>Your Slack bot has been successfully installed to <strong>%s</strong>.</p>
			<p>You can now close this window and return to Slack.</p>
		</body>
		</html>
	`, resp.Team.Name)
}

// handleEvents handles Slack events
func (bot *Bot) handleEvents(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	// Parse the event
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Printf("Error parsing event: %v", err)
		http.Error(w, "Failed to parse event", http.StatusBadRequest)
		return
	}
	
	// Handle URL verification
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &r); err != nil {
			log.Printf("Error unmarshaling challenge: %v", err)
			http.Error(w, "Failed to unmarshal challenge", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
		return
	}
	
	// Handle callback events
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		
		// Create base context
		baseCtx := bot.newContext(
			eventsAPIEvent.TeamID,
			"", // Will be set based on event type
			"", // Will be set based on event type
			"",
			eventsAPIEvent.EventID,
			"",
		)
		
		if baseCtx == nil {
			log.Printf("Failed to create base context for team %s", eventsAPIEvent.TeamID)
			http.Error(w, "Failed to create context", http.StatusInternalServerError)
			return
		}
		
		// Handle different event types
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			bot.handleMessageEvent(baseCtx, ev)
		case *slackevents.AppHomeOpenedEvent:
			bot.handleAppHomeOpenedEvent(baseCtx, ev)
		case *slackevents.ReactionAddedEvent:
			bot.handleReactionAddedEvent(baseCtx, ev)
		case *slackevents.TeamJoinEvent:
			bot.handleTeamJoinEvent(baseCtx, ev)
		case *slackevents.ChannelJoinedEvent:
			bot.handleChannelJoinedEvent(baseCtx, ev)
		default:
			if bot.debug {
				log.Printf("Unhandled event type: %s", innerEvent.Type)
			}
		}
	}
	
	w.WriteHeader(http.StatusOK)
}

// handleMessageEvent handles message events
func (bot *Bot) handleMessageEvent(baseCtx *Context, event *slackevents.MessageEvent) {
	// Skip bot messages
	if event.BotID != "" {
		return
	}
	
	// Update context with event-specific data
	baseCtx.UserID = event.User
	baseCtx.ChannelID = event.Channel
	baseCtx.Timestamp = event.TimeStamp
	
	// Create message context
	msgCtx := bot.newMessageContext(baseCtx, &slack.MessageEvent{
		Msg: slack.Msg{
			Type:            event.Type,
			Channel:         event.Channel,
			User:            event.User,
			Text:            event.Text,
			TimeStamp:       event.TimeStamp,
			ThreadTimeStamp: event.ThreadTimeStamp,
			Files:           event.Files,
		},
	})
	
	// Check for mentions
	if strings.Contains(event.Text, "<@"+baseCtx.Bot.storage.GetTeam(baseCtx.TeamID).BotUserID+">") {
		msgCtx.IsMention = true
	}
	
	// Apply middleware
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchMessageEvent(msgCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchMessageEvent dispatches message events to appropriate handlers
func (bot *Bot) dispatchMessageEvent(ctx *MessageContext) error {
	var errors []error
	
	// Handle pattern-based handlers first
	for text, handlers := range bot.messageContainingHandlers {
		if strings.Contains(strings.ToLower(ctx.Text), strings.ToLower(text)) {
			for _, handler := range handlers {
				if err := handler(ctx); err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	
	for pattern, handlers := range bot.messageMatchingHandlers {
		if pattern.MatchString(ctx.Text) {
			for _, handler := range handlers {
				if err := handler(ctx); err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	
	// Handle specific event types
	if ctx.IsMention {
		for _, handler := range bot.mentionMessageHandlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	if ctx.IsThread {
		for _, handler := range bot.threadMessageHandlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	if ctx.IsDirect {
		for _, handler := range bot.directMessageHandlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	if ctx.IsChannel {
		for _, handler := range bot.channelMessageHandlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	// Handle general message handlers
	for _, handler := range bot.messageHandlers {
		if err := handler(ctx); err != nil {
			errors = append(errors, err)
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// handleAppHomeOpenedEvent handles app home opened events
func (bot *Bot) handleAppHomeOpenedEvent(baseCtx *Context, event *slackevents.AppHomeOpenedEvent) {
	// Update context with event-specific data
	baseCtx.UserID = event.User
	
	// Create home context
	homeCtx := bot.newHomeContext(baseCtx, event)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchHomeEvent(homeCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchHomeEvent dispatches home events to appropriate handlers
func (bot *Bot) dispatchHomeEvent(ctx *HomeContext) error {
	var errors []error
	
	for _, handler := range bot.homeOpenedHandlers {
		if err := handler(ctx); err != nil {
			errors = append(errors, err)
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// handleReactionAddedEvent handles reaction added events
func (bot *Bot) handleReactionAddedEvent(baseCtx *Context, event *slackevents.ReactionAddedEvent) {
	// Update context with event-specific data
	baseCtx.UserID = event.User
	baseCtx.ChannelID = event.Item.Channel
	baseCtx.Timestamp = event.Item.Timestamp
	
	// Create reaction context
	reactionCtx := bot.newReactionContext(baseCtx, event)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchReactionEvent(reactionCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchReactionEvent dispatches reaction events to appropriate handlers
func (bot *Bot) dispatchReactionEvent(ctx *ReactionContext) error {
	var errors []error
	
	for _, handler := range bot.reactionAddedHandlers {
		if err := handler(ctx); err != nil {
			errors = append(errors, err)
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// handleTeamJoinEvent handles team join events
func (bot *Bot) handleTeamJoinEvent(baseCtx *Context, event *slackevents.TeamJoinEvent) {
	// Update context with event-specific data
	baseCtx.UserID = event.User.ID
	
	// Create user context
	userCtx := bot.newUserContext(baseCtx, &event.User, "joined")
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchUserEvent(userCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchUserEvent dispatches user events to appropriate handlers
func (bot *Bot) dispatchUserEvent(ctx *UserContext) error {
	var errors []error
	
	if ctx.Action == "joined" {
		for _, handler := range bot.userJoinedHandlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// handleChannelJoinedEvent handles channel joined events
func (bot *Bot) handleChannelJoinedEvent(baseCtx *Context, event *slackevents.ChannelJoinedEvent) {
	// Update context with event-specific data
	baseCtx.ChannelID = event.Channel.ID
	
	// Create channel context
	channelCtx := bot.newChannelContext(baseCtx, &event.Channel, "joined")
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchChannelEvent(channelCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchChannelEvent dispatches channel events to appropriate handlers
func (bot *Bot) dispatchChannelEvent(ctx *ChannelContext) error {
	var errors []error
	
	if ctx.Action == "joined" {
		for _, handler := range bot.channelJoinedHandlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// Helper functions to create contexts
func (bot *Bot) newUserContext(baseCtx *Context, user *slackevents.User, action string) *UserContext {
	slackUser := &slack.User{
		ID:       user.ID,
		Name:     user.Name,
		RealName: user.RealName,
		Profile: slack.UserProfile{
			Email: user.Profile.Email,
		},
	}
	
	return &UserContext{
		Context: baseCtx,
		User:    slackUser,
		Action:  action,
	}
}

func (bot *Bot) newChannelContext(baseCtx *Context, channel *slackevents.Channel, action string) *ChannelContext {
	slackChannel := &slack.Channel{
		GroupConversation: slack.GroupConversation{
			Conversation: slack.Conversation{
				ID:   channel.ID,
				Name: channel.Name,
			},
		},
	}
	
	return &ChannelContext{
		Context: baseCtx,
		Channel: slackChannel,
		Action:  action,
	}
}