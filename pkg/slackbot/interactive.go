package slackbot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/slack-go/slack"
)

// handleInteractive handles Slack interactive events (buttons, modals, etc.)
func (bot *Bot) handleInteractive(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	// Parse the form data
	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	
	// Get the payload
	payload := values.Get("payload")
	if payload == "" {
		log.Printf("Missing payload in request")
		http.Error(w, "Missing payload", http.StatusBadRequest)
		return
	}
	
	// Parse the interaction callback
	var interaction slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &interaction); err != nil {
		log.Printf("Error unmarshaling interaction: %v", err)
		http.Error(w, "Failed to unmarshal interaction", http.StatusBadRequest)
		return
	}
	
	// Create base context
	baseCtx := bot.newContext(
		interaction.Team.ID,
		interaction.User.ID,
		interaction.Channel.ID,
		interaction.TriggerID,
		"", // No event ID for interactive events
		"",
	)
	
	if baseCtx == nil {
		log.Printf("Failed to create base context for team %s", interaction.Team.ID)
		http.Error(w, "Failed to create context", http.StatusInternalServerError)
		return
	}
	
	// Handle different interaction types
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		bot.handleBlockActions(baseCtx, interaction)
	case slack.InteractionTypeViewSubmission:
		bot.handleViewSubmission(baseCtx, interaction)
	case slack.InteractionTypeViewClosed:
		bot.handleViewClosed(baseCtx, interaction)
	case slack.InteractionTypeShortcut:
		bot.handleShortcut(baseCtx, interaction)
	case slack.InteractionTypeSlashCommand:
		bot.handleSlashCommand(baseCtx, interaction)
	default:
		if bot.debug {
			log.Printf("Unhandled interaction type: %s", interaction.Type)
		}
	}
	
	w.WriteHeader(http.StatusOK)
}

// handleBlockActions handles block action interactions (buttons, select menus, etc.)
func (bot *Bot) handleBlockActions(baseCtx *Context, interaction slack.InteractionCallback) {
	// Create interactive context
	interactiveCtx := bot.newInteractiveContext(baseCtx, interaction)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchBlockActions(interactiveCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchBlockActions dispatches block actions to appropriate handlers
func (bot *Bot) dispatchBlockActions(ctx *InteractiveContext) error {
	var errors []error
	
	// Handle button clicks
	if ctx.ActionType == "button" {
		if handlers, exists := bot.buttonClickHandlers[ctx.ActionID]; exists {
			for _, handler := range handlers {
				if err := handler(ctx); err != nil {
					errors = append(errors, err)
				}
			}
		} else {
			// Use default button handler
			if err := defaultButtonHandler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
	
	// Handle select menus
	if ctx.ActionType == "static_select" || ctx.ActionType == "external_select" {
		if handlers, exists := bot.selectMenuHandlers[ctx.ActionID]; exists {
			for _, handler := range handlers {
				if err := handler(ctx); err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// handleViewSubmission handles modal submissions
func (bot *Bot) handleViewSubmission(baseCtx *Context, interaction slack.InteractionCallback) {
	// Create interactive context
	interactiveCtx := bot.newInteractiveContext(baseCtx, interaction)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchViewSubmission(interactiveCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchViewSubmission dispatches view submissions to appropriate handlers
func (bot *Bot) dispatchViewSubmission(ctx *InteractiveContext) error {
	var errors []error
	
	callbackID := ctx.Interaction.View.CallbackID
	
	if handlers, exists := bot.modalSubmitHandlers[callbackID]; exists {
		for _, handler := range handlers {
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

// handleViewClosed handles modal closures
func (bot *Bot) handleViewClosed(baseCtx *Context, interaction slack.InteractionCallback) {
	// Create interactive context
	interactiveCtx := bot.newInteractiveContext(baseCtx, interaction)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchViewClosed(interactiveCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchViewClosed dispatches view closures to appropriate handlers
func (bot *Bot) dispatchViewClosed(ctx *InteractiveContext) error {
	// Handle modal closures if needed
	if bot.debug {
		log.Printf("Modal closed: %s", ctx.Interaction.View.CallbackID)
	}
	
	return nil
}

// handleShortcut handles shortcuts
func (bot *Bot) handleShortcut(baseCtx *Context, interaction slack.InteractionCallback) {
	// Create interactive context
	interactiveCtx := bot.newInteractiveContext(baseCtx, interaction)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchShortcut(interactiveCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// dispatchShortcut dispatches shortcuts to appropriate handlers
func (bot *Bot) dispatchShortcut(ctx *InteractiveContext) error {
	// Handle shortcuts if needed
	if bot.debug {
		log.Printf("Shortcut triggered: %s", ctx.Interaction.CallbackID)
	}
	
	return nil
}

// handleSlashCommand handles slash commands
func (bot *Bot) handleSlashCommand(baseCtx *Context, interaction slack.InteractionCallback) {
	// This is actually handled in a different endpoint, but included here for completeness
	// Real slash commands come through a different webhook
	if bot.debug {
		log.Printf("Slash command received via interaction: %s", interaction.CallbackID)
	}
}

// handleSlashCommandWebhook handles slash commands from the dedicated webhook
func (bot *Bot) handleSlashCommandWebhook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	
	// Parse slash command
	slashCommand := slack.SlashCommand{
		Token:       r.FormValue("token"),
		TeamID:      r.FormValue("team_id"),
		TeamDomain:  r.FormValue("team_domain"),
		ChannelID:   r.FormValue("channel_id"),
		ChannelName: r.FormValue("channel_name"),
		UserID:      r.FormValue("user_id"),
		UserName:    r.FormValue("user_name"),
		Command:     r.FormValue("command"),
		Text:        r.FormValue("text"),
		ResponseURL: r.FormValue("response_url"),
		TriggerID:   r.FormValue("trigger_id"),
	}
	
	// Create base context
	baseCtx := bot.newContext(
		slashCommand.TeamID,
		slashCommand.UserID,
		slashCommand.ChannelID,
		slashCommand.TriggerID,
		"",
		"",
	)
	
	if baseCtx == nil {
		log.Printf("Failed to create base context for team %s", slashCommand.TeamID)
		http.Error(w, "Failed to create context", http.StatusInternalServerError)
		return
	}
	
	// Create slash command context
	slashCtx := bot.newSlashCommandContext(baseCtx, slashCommand)
	
	// Apply middleware and handle
	for _, middleware := range bot.middlewares {
		handler := middleware(func() error {
			return bot.dispatchSlashCommand(slashCtx)
		})
		if err := handler(); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
	
	w.WriteHeader(http.StatusOK)
}

// dispatchSlashCommand dispatches slash commands to appropriate handlers
func (bot *Bot) dispatchSlashCommand(ctx *SlashCommandContext) error {
	var errors []error
	
	// Remove leading slash from command
	command := strings.TrimPrefix(ctx.Command, "/")
	
	if handlers, exists := bot.slashCommandHandlers[command]; exists {
		for _, handler := range handlers {
			if err := handler(ctx); err != nil {
				errors = append(errors, err)
			}
		}
	} else {
		// Default slash command handler
		if err := ctx.RespondEphemeral(fmt.Sprintf("Unknown command: %s", ctx.Command)); err != nil {
			errors = append(errors, err)
		}
	}
	
	// If there are errors, return the first one
	if len(errors) > 0 {
		return errors[0]
	}
	
	return nil
}

// Helper functions for extracting values from interactions
func (ctx *InteractiveContext) GetSubmissionValues() map[string]string {
	values := make(map[string]string)
	
	if ctx.Interaction.View.State == nil {
		return values
	}
	
	for blockID, block := range ctx.Interaction.View.State.Values {
		for actionID, action := range block {
			if action.Type == "plain_text_input" {
				values[actionID] = action.Value
			} else if action.Type == "static_select" {
				if action.SelectedOption != nil {
					values[actionID] = action.SelectedOption.Value
				}
			}
		}
	}
	
	return values
}

// GetFormValue gets a specific form value from modal submission
func (ctx *InteractiveContext) GetFormValue(actionID string) string {
	values := ctx.GetSubmissionValues()
	return values[actionID]
}

// GetSelectedOption gets the selected option from a select menu
func (ctx *InteractiveContext) GetSelectedOption(actionID string) *slack.OptionBlockObject {
	if ctx.Interaction.View.State == nil {
		return nil
	}
	
	for _, block := range ctx.Interaction.View.State.Values {
		if action, exists := block[actionID]; exists {
			if action.Type == "static_select" {
				return action.SelectedOption
			}
		}
	}
	
	return nil
}

// IsFromModal checks if the interaction came from a modal
func (ctx *InteractiveContext) IsFromModal() bool {
	return ctx.Interaction.View.Type == slack.VTModal
}

// IsFromHomeTab checks if the interaction came from the home tab
func (ctx *InteractiveContext) IsFromHomeTab() bool {
	return ctx.Interaction.View.Type == slack.VTHomeTab
}

// IsFromMessage checks if the interaction came from a message
func (ctx *InteractiveContext) IsFromMessage() bool {
	return ctx.Interaction.Container.Type == slack.CTMessage
}