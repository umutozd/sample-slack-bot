package ui

import (
	"github.com/slack-go/slack"
)

// HomeView builder for creating home tab views
type HomeView struct {
	blocks []slack.Block
}

// NewHomeView creates a new home view builder
func NewHomeView() *HomeView {
	return &HomeView{
		blocks: []slack.Block{},
	}
}

// Header adds a header section
func (h *HomeView) Header(text string) *HomeView {
	h.blocks = append(h.blocks, slack.NewHeaderBlock(
		slack.NewTextBlockObject("plain_text", text, false, false),
	))
	return h
}

// Section adds a section block
func (h *HomeView) Section(text string) *HomeView {
	h.blocks = append(h.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil, nil,
	))
	return h
}

// SectionWithFields adds a section with fields
func (h *HomeView) SectionWithFields(text string, fields []string) *HomeView {
	textFields := make([]*slack.TextBlockObject, len(fields))
	for i, field := range fields {
		textFields[i] = slack.NewTextBlockObject("mrkdwn", field, false, false)
	}
	
	h.blocks = append(h.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		textFields, nil,
	))
	return h
}

// ButtonRow adds a row of buttons
func (h *HomeView) ButtonRow(buttons ...Button) *HomeView {
	if len(buttons) == 0 {
		return h
	}
	
	elements := make([]slack.BlockElement, len(buttons))
	for i, button := range buttons {
		elements[i] = button.build()
	}
	
	h.blocks = append(h.blocks, slack.NewActionBlock("", elements...))
	return h
}

// Button creates a button
func (h *HomeView) Button(actionID, text string) *HomeView {
	button := Button{
		actionID: actionID,
		text:     text,
		style:    "default",
	}
	
	h.blocks = append(h.blocks, slack.NewActionBlock("", button.build()))
	return h
}

// Divider adds a divider
func (h *HomeView) Divider() *HomeView {
	h.blocks = append(h.blocks, slack.NewDividerBlock())
	return h
}

// Input adds an input field
func (h *HomeView) Input(actionID, label string, inputType InputType) *HomeView {
	var element slack.BlockElement
	
	switch inputType.Type {
	case "plain_text":
		element = slack.NewPlainTextInputBlockElement(
			slack.NewTextBlockObject("plain_text", inputType.Placeholder, false, false),
			actionID,
		)
	case "email":
		element = slack.NewEmailInputBlockElement(
			slack.NewTextBlockObject("plain_text", inputType.Placeholder, false, false),
			actionID,
		)
	default:
		element = slack.NewPlainTextInputBlockElement(
			slack.NewTextBlockObject("plain_text", inputType.Placeholder, false, false),
			actionID,
		)
	}
	
	h.blocks = append(h.blocks, slack.NewInputBlock(
		actionID,
		slack.NewTextBlockObject("plain_text", label, false, false),
		element,
	))
	return h
}

// SelectMenu adds a select menu
func (h *HomeView) SelectMenu(actionID, label string, options []Option) *HomeView {
	slackOptions := make([]*slack.OptionBlockObject, len(options))
	for i, option := range options {
		slackOptions[i] = slack.NewOptionBlockObject(
			option.Value,
			slack.NewTextBlockObject("plain_text", option.Text, false, false),
			slack.NewTextBlockObject("plain_text", option.Description, false, false),
		)
	}
	
	selectMenu := slack.NewOptionsSelectBlockElement(
		"static_select",
		slack.NewTextBlockObject("plain_text", "Select an option", false, false),
		actionID,
		slackOptions...,
	)
	
	h.blocks = append(h.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("plain_text", label, false, false),
		nil, selectMenu,
	))
	return h
}

// If conditionally adds blocks
func (h *HomeView) If(condition bool, fn func(*HomeView) *HomeView) *HomeView {
	if condition {
		return fn(h)
	}
	return h
}

// Blocks returns the built blocks
func (h *HomeView) Blocks() []slack.Block {
	return h.blocks
}

// Button represents a button component
type Button struct {
	actionID string
	text     string
	style    string
	value    string
	url      string
}

// Button creates a new button
func Button(actionID, text string) Button {
	return Button{
		actionID: actionID,
		text:     text,
		style:    "default",
	}
}

// Primary sets the button style to primary
func (b Button) Primary() Button {
	b.style = "primary"
	return b
}

// Secondary sets the button style to secondary
func (b Button) Secondary() Button {
	b.style = "default"
	return b
}

// Danger sets the button style to danger
func (b Button) Danger() Button {
	b.style = "danger"
	return b
}

// WithValue sets the button value
func (b Button) WithValue(value string) Button {
	b.value = value
	return b
}

// WithURL sets the button URL (for link buttons)
func (b Button) WithURL(url string) Button {
	b.url = url
	return b
}

// build creates the Slack button block element
func (b Button) build() slack.BlockElement {
	if b.url != "" {
		return slack.NewButtonBlockElement(
			b.actionID,
			b.value,
			slack.NewTextBlockObject("plain_text", b.text, false, false),
		).WithURL(b.url)
	}
	
	button := slack.NewButtonBlockElement(
		b.actionID,
		b.value,
		slack.NewTextBlockObject("plain_text", b.text, false, false),
	)
	
	if b.style != "default" {
		button = button.WithStyle(slack.StyleMap[b.style])
	}
	
	return button
}

// Option represents a select menu option
type Option struct {
	Value       string
	Text        string
	Description string
}

// Option creates a new option
func Option(value, text string) Option {
	return Option{
		Value: value,
		Text:  text,
	}
}

// WithDescription adds a description to the option
func (o Option) WithDescription(description string) Option {
	o.Description = description
	return o
}

// InputType represents different input types
type InputType struct {
	Type        string
	Placeholder string
	Required    bool
	Multiline   bool
}

// PlainText creates a plain text input
func PlainText() InputType {
	return InputType{
		Type:        "plain_text",
		Placeholder: "Enter text...",
	}
}

// Email creates an email input
func Email() InputType {
	return InputType{
		Type:        "email",
		Placeholder: "Enter email address...",
	}
}

// Required marks the input as required
func (i InputType) Required() InputType {
	i.Required = true
	return i
}

// WithPlaceholder sets the placeholder text
func (i InputType) WithPlaceholder(placeholder string) InputType {
	i.Placeholder = placeholder
	return i
}

// Multiline enables multiline input
func (i InputType) Multiline() InputType {
	i.Multiline = true
	return i
}

// Message builder for creating messages
type Message struct {
	text        string
	blocks      []slack.Block
	attachments []slack.Attachment
	ephemeral   bool
}

// NewMessage creates a new message builder
func NewMessage() *Message {
	return &Message{
		blocks:      []slack.Block{},
		attachments: []slack.Attachment{},
	}
}

// Text sets the message text
func (m *Message) Text(text string) *Message {
	m.text = text
	return m
}

// Section adds a section block
func (m *Message) Section(text string) *Message {
	m.blocks = append(m.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil, nil,
	))
	return m
}

// ButtonRow adds a row of buttons
func (m *Message) ButtonRow(buttons ...Button) *Message {
	if len(buttons) == 0 {
		return m
	}
	
	elements := make([]slack.BlockElement, len(buttons))
	for i, button := range buttons {
		elements[i] = button.build()
	}
	
	m.blocks = append(m.blocks, slack.NewActionBlock("", elements...))
	return m
}

// Ephemeral marks the message as ephemeral
func (m *Message) Ephemeral() *Message {
	m.ephemeral = true
	return m
}

// Blocks returns the built blocks
func (m *Message) Blocks() []slack.Block {
	return m.blocks
}

// IsEphemeral returns whether the message is ephemeral
func (m *Message) IsEphemeral() bool {
	return m.ephemeral
}

// GetText returns the message text
func (m *Message) GetText() string {
	return m.text
}

// Modal builder for creating modals
type Modal struct {
	callbackID string
	title      string
	blocks     []slack.Block
	submit     string
	cancel     string
}

// NewModal creates a new modal builder
func NewModal(callbackID, title string) *Modal {
	return &Modal{
		callbackID: callbackID,
		title:      title,
		blocks:     []slack.Block{},
		submit:     "Submit",
		cancel:     "Cancel",
	}
}

// Section adds a section block
func (m *Modal) Section(text string) *Modal {
	m.blocks = append(m.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil, nil,
	))
	return m
}

// Input adds an input field
func (m *Modal) Input(actionID, label string, inputType InputType) *Modal {
	var element slack.BlockElement
	
	switch inputType.Type {
	case "plain_text":
		element = slack.NewPlainTextInputBlockElement(
			slack.NewTextBlockObject("plain_text", inputType.Placeholder, false, false),
			actionID,
		)
	case "email":
		element = slack.NewEmailInputBlockElement(
			slack.NewTextBlockObject("plain_text", inputType.Placeholder, false, false),
			actionID,
		)
	default:
		element = slack.NewPlainTextInputBlockElement(
			slack.NewTextBlockObject("plain_text", inputType.Placeholder, false, false),
			actionID,
		)
	}
	
	m.blocks = append(m.blocks, slack.NewInputBlock(
		actionID,
		slack.NewTextBlockObject("plain_text", label, false, false),
		element,
	))
	return m
}

// SelectMenu adds a select menu
func (m *Modal) SelectMenu(actionID, label string, options []Option) *Modal {
	slackOptions := make([]*slack.OptionBlockObject, len(options))
	for i, option := range options {
		slackOptions[i] = slack.NewOptionBlockObject(
			option.Value,
			slack.NewTextBlockObject("plain_text", option.Text, false, false),
			slack.NewTextBlockObject("plain_text", option.Description, false, false),
		)
	}
	
	selectMenu := slack.NewOptionsSelectBlockElement(
		"static_select",
		slack.NewTextBlockObject("plain_text", "Select an option", false, false),
		actionID,
		slackOptions...,
	)
	
	m.blocks = append(m.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("plain_text", label, false, false),
		nil, selectMenu,
	))
	return m
}

// Submit sets the submit button text
func (m *Modal) Submit(text string) *Modal {
	m.submit = text
	return m
}

// Cancel sets the cancel button text
func (m *Modal) Cancel(text string) *Modal {
	m.cancel = text
	return m
}

// Build creates the Slack modal view request
func (m *Modal) Build() slack.ModalViewRequest {
	return slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: m.callbackID,
		Title:      slack.NewTextBlockObject("plain_text", m.title, false, false),
		Blocks:     slack.Blocks{BlockSet: m.blocks},
		Submit:     slack.NewTextBlockObject("plain_text", m.submit, false, false),
		Close:      slack.NewTextBlockObject("plain_text", m.cancel, false, false),
	}
}

// Pre-built components for common use cases
func WelcomeMessage(userName string) *Message {
	return NewMessage().
		Text("Welcome to the team!").
		Section("👋 Welcome to the team, " + userName + "! We're excited to have you here.").
		ButtonRow(
			Button("getting_started", "Getting Started").Primary(),
			Button("team_info", "Team Info"),
		)
}

func ErrorMessage(error string) *Message {
	return NewMessage().
		Text("An error occurred").
		Section("❌ " + error).
		Ephemeral()
}

func SuccessMessage(message string) *Message {
	return NewMessage().
		Text("Success!").
		Section("✅ " + message)
}

func ConfirmationDialog(title, message string) *Modal {
	return NewModal("confirmation", title).
		Section(message).
		Submit("Confirm").
		Cancel("Cancel")
}

func LoadingMessage() *Message {
	return NewMessage().
		Text("Loading...").
		Section("🔄 Please wait while we process your request...")
}