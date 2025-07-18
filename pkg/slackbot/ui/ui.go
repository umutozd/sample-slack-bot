package ui

import (
	"github.com/slack-go/slack"
)

// HomeView represents a home tab view builder
type HomeView struct {
	blocks []slack.Block
}

// Message represents a message builder
type Message struct {
	blocks []slack.Block
	text   string
}

// Modal represents a modal builder
type Modal struct {
	title        *slack.TextBlockObject
	submit       *slack.TextBlockObject
	close        *slack.TextBlockObject
	callbackID   string
	blocks       []slack.Block
	clearOnClose bool
}

// NewHomeView creates a new home view builder
func NewHomeView() *HomeView {
	return &HomeView{
		blocks: []slack.Block{},
	}
}

// NewMessage creates a new message builder
func NewMessage() *Message {
	return &Message{
		blocks: []slack.Block{},
	}
}

// NewModal creates a new modal builder
func NewModal(callbackID, title string) *Modal {
	return &Modal{
		title:        slack.NewTextBlockObject(slack.PlainTextType, title, false, false),
		close:        slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false),
		callbackID:   callbackID,
		blocks:       []slack.Block{},
		clearOnClose: true,
	}
}

// HomeView builder methods

// Header adds a header block to the home view
func (hv *HomeView) Header(text string) *HomeView {
	hv.blocks = append(hv.blocks, slack.NewHeaderBlock(
		slack.NewTextBlockObject(slack.PlainTextType, text, false, false),
	))
	return hv
}

// Section adds a section block to the home view
func (hv *HomeView) Section(text string) *HomeView {
	hv.blocks = append(hv.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, text, false, false),
		nil,
		nil,
	))
	return hv
}

// Button adds a button to the home view
func (hv *HomeView) Button(actionID, text string) *HomeView {
	hv.blocks = append(hv.blocks, slack.NewActionBlock(
		"",
		slack.NewButtonBlockElement(actionID, actionID, slack.NewTextBlockObject(slack.PlainTextType, text, false, false)),
	))
	return hv
}

// ButtonRow adds a row of buttons to the home view
func (hv *HomeView) ButtonRow(buttons ...ButtonElement) *HomeView {
	elements := make([]slack.BlockElement, len(buttons))
	for i, btn := range buttons {
		elements[i] = btn.ToSlackButton()
	}
	
	hv.blocks = append(hv.blocks, slack.NewActionBlock("", elements...))
	return hv
}

// Divider adds a divider block to the home view
func (hv *HomeView) Divider() *HomeView {
	hv.blocks = append(hv.blocks, slack.NewDividerBlock())
	return hv
}

// Input adds an input block to the home view
func (hv *HomeView) Input(actionID, label string, element InputElement) *HomeView {
	hv.blocks = append(hv.blocks, slack.NewInputBlock(
		"",
		slack.NewTextBlockObject(slack.PlainTextType, label, false, false),
		element.ToSlackElement(),
	))
	return hv
}

// SelectMenu adds a select menu to the home view
func (hv *HomeView) SelectMenu(actionID, placeholder string, options []Option) *HomeView {
	slackOptions := make([]*slack.OptionBlockObject, len(options))
	for i, opt := range options {
		slackOptions[i] = slack.NewOptionBlockObject(opt.Value, slack.NewTextBlockObject(slack.PlainTextType, opt.Text, false, false))
	}
	
	selectElement := slack.NewOptionsSelectBlockElement(
		slack.OptTypeStatic,
		slack.NewTextBlockObject(slack.PlainTextType, placeholder, false, false),
		actionID,
		slackOptions...,
	)
	
	hv.blocks = append(hv.blocks, slack.NewActionBlock("", selectElement))
	return hv
}

// If adds conditional blocks based on a condition
func (hv *HomeView) If(condition bool, fn func(*HomeView) *HomeView) *HomeView {
	if condition {
		return fn(hv)
	}
	return hv
}

// Build returns the final HomeTabViewRequest
func (hv *HomeView) Build() slack.HomeTabViewRequest {
	return slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: hv.blocks,
		},
	}
}

// Message builder methods

// Text sets the text content of the message
func (m *Message) Text(text string) *Message {
	m.text = text
	return m
}

// Section adds a section block to the message
func (m *Message) Section(text string) *Message {
	m.blocks = append(m.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, text, false, false),
		nil,
		nil,
	))
	return m
}

// Button adds a button to the message
func (m *Message) Button(actionID, text string) *Message {
	m.blocks = append(m.blocks, slack.NewActionBlock(
		"",
		slack.NewButtonBlockElement(actionID, actionID, slack.NewTextBlockObject(slack.PlainTextType, text, false, false)),
	))
	return m
}

// ButtonRow adds a row of buttons to the message
func (m *Message) ButtonRow(buttons ...ButtonElement) *Message {
	elements := make([]slack.BlockElement, len(buttons))
	for i, btn := range buttons {
		elements[i] = btn.ToSlackButton()
	}
	
	m.blocks = append(m.blocks, slack.NewActionBlock("", elements...))
	return m
}

// Divider adds a divider block to the message
func (m *Message) Divider() *Message {
	m.blocks = append(m.blocks, slack.NewDividerBlock())
	return m
}

// Build returns the message options
func (m *Message) Build() []slack.MsgOption {
	options := []slack.MsgOption{}
	
	if m.text != "" {
		options = append(options, slack.MsgOptionText(m.text, false))
	}
	
	if len(m.blocks) > 0 {
		options = append(options, slack.MsgOptionBlocks(m.blocks...))
	}
	
	return options
}

// Modal builder methods

// Submit sets the submit button text
func (m *Modal) Submit(text string) *Modal {
	m.submit = slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
	return m
}

// Cancel sets the cancel button text
func (m *Modal) Cancel(text string) *Modal {
	m.close = slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
	return m
}

// Section adds a section block to the modal
func (m *Modal) Section(text string) *Modal {
	m.blocks = append(m.blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, text, false, false),
		nil,
		nil,
	))
	return m
}

// Input adds an input block to the modal
func (m *Modal) Input(actionID, label string, element InputElement) *Modal {
	m.blocks = append(m.blocks, slack.NewInputBlock(
		"",
		slack.NewTextBlockObject(slack.PlainTextType, label, false, false),
		element.ToSlackElement(),
	))
	return m
}

// SelectMenu adds a select menu to the modal
func (m *Modal) SelectMenu(actionID, placeholder string, options []Option) *Modal {
	slackOptions := make([]*slack.OptionBlockObject, len(options))
	for i, opt := range options {
		slackOptions[i] = slack.NewOptionBlockObject(opt.Value, slack.NewTextBlockObject(slack.PlainTextType, opt.Text, false, false))
	}
	
	selectElement := slack.NewOptionsSelectBlockElement(
		slack.OptTypeStatic,
		slack.NewTextBlockObject(slack.PlainTextType, placeholder, false, false),
		actionID,
		slackOptions...,
	)
	
	m.blocks = append(m.blocks, slack.NewInputBlock("", slack.NewTextBlockObject(slack.PlainTextType, placeholder, false, false), selectElement))
	return m
}

// Divider adds a divider block to the modal
func (m *Modal) Divider() *Modal {
	m.blocks = append(m.blocks, slack.NewDividerBlock())
	return m
}

// Build returns the final ModalViewRequest
func (m *Modal) Build() slack.ModalViewRequest {
	return slack.ModalViewRequest{
		Type:         slack.VTModal,
		Title:        m.title,
		Submit:       m.submit,
		Close:        m.close,
		CallbackID:   m.callbackID,
		ClearOnClose: m.clearOnClose,
		Blocks: slack.Blocks{
			BlockSet: m.blocks,
		},
	}
}

// Helper types and functions

// ButtonElement represents a button element
type ButtonElement struct {
	ActionID string
	Text     string
	Value    string
	Style    slack.Style
}

// Button creates a new button element
func Button(actionID, text string) ButtonElement {
	return ButtonElement{
		ActionID: actionID,
		Text:     text,
		Value:    actionID,
		Style:    slack.StyleDefault,
	}
}

// Primary sets the button style to primary
func (b ButtonElement) Primary() ButtonElement {
	b.Style = slack.StylePrimary
	return b
}

// Secondary sets the button style to secondary
func (b ButtonElement) Secondary() ButtonElement {
	b.Style = slack.StyleDefault
	return b
}

// Danger sets the button style to danger
func (b ButtonElement) Danger() ButtonElement {
	b.Style = slack.StyleDanger
	return b
}

// WithValue sets the button value
func (b ButtonElement) WithValue(value string) ButtonElement {
	b.Value = value
	return b
}

// ToSlackButton converts to slack.ButtonBlockElement
func (b ButtonElement) ToSlackButton() *slack.ButtonBlockElement {
	button := slack.NewButtonBlockElement(b.ActionID, b.Value, slack.NewTextBlockObject(slack.PlainTextType, b.Text, false, false))
	button.Style = b.Style
	return button
}

// InputElement represents an input element
type InputElement interface {
	ToSlackElement() slack.BlockElement
}

// PlainTextInput represents a plain text input
type PlainTextInput struct {
	Placeholder string
	Multiline   bool
	Required    bool
}

// PlainText creates a new plain text input
func PlainText() PlainTextInput {
	return PlainTextInput{
		Placeholder: "",
		Multiline:   false,
		Required:    false,
	}
}

// WithPlaceholder sets the placeholder text
func (p PlainTextInput) WithPlaceholder(placeholder string) PlainTextInput {
	p.Placeholder = placeholder
	return p
}

// Multiline sets the input to multiline
func (p PlainTextInput) Multiline() PlainTextInput {
	p.Multiline = true
	return p
}

// Required sets the input as required
func (p PlainTextInput) Required() PlainTextInput {
	p.Required = true
	return p
}

// ToSlackElement converts to slack element
func (p PlainTextInput) ToSlackElement() slack.BlockElement {
	element := slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject(slack.PlainTextType, p.Placeholder, false, false), "")
	element.Multiline = p.Multiline
	return element
}

// EmailInput represents an email input
type EmailInput struct {
	Placeholder string
}

// Email creates a new email input
func Email() EmailInput {
	return EmailInput{
		Placeholder: "Enter your email",
	}
}

// WithPlaceholder sets the placeholder text
func (e EmailInput) WithPlaceholder(placeholder string) EmailInput {
	e.Placeholder = placeholder
	return e
}

// ToSlackElement converts to slack element
func (e EmailInput) ToSlackElement() slack.BlockElement {
	return slack.NewEmailInputBlockElement(slack.NewTextBlockObject(slack.PlainTextType, e.Placeholder, false, false), "")
}

// Option represents a select option
type Option struct {
	Value string
	Text  string
}

// NewOption creates a new option
func NewOption(value, text string) Option {
	return Option{
		Value: value,
		Text:  text,
	}
}

// Pre-built components

// WelcomeMessage creates a welcome message
func WelcomeMessage(userName string) *Message {
	return NewMessage().
		Text("Welcome to our bot!").
		Section("Hello " + userName + "! 👋").
		ButtonRow(
			Button("get_started", "Get Started").Primary(),
			Button("help", "Help"),
		)
}

// ErrorMessage creates an error message
func ErrorMessage(errorText string) *Message {
	return NewMessage().
		Text("❌ Error").
		Section("*Error:* " + errorText)
}

// SuccessMessage creates a success message
func SuccessMessage(message string) *Message {
	return NewMessage().
		Text("✅ Success").
		Section("*Success:* " + message)
}

// ConfirmationDialog creates a confirmation modal
func ConfirmationDialog(title, message string) *Modal {
	return NewModal("confirm", title).
		Section(message).
		Submit("Confirm").
		Cancel("Cancel")
}

// LoadingIndicator creates a loading message
func LoadingIndicator() *Message {
	return NewMessage().
		Text("⏳ Loading...").
		Section("Please wait while we process your request.")
}