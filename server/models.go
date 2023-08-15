package server

type slackEvent struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	TeamID    string `json:"team_id"`
	Event     struct {
		Type       string `json:"type"`
		Text       string `json:"text"`
		User       string `json:"user"`
		Channel    string `json:"channel"`
		BotProfile struct {
			Name string `json:"name"`
		} `json:"bot_profile"`
		Tab string `json:"tab"`
		// When event is a reply to a thread, this field contains "<parent-message-timestamp>.<reply-message-id>"
		ThreadTS string `json:"thread_ts"`
	} `json:"event"`
	EventID string `json:"event_id"`
}

type httpError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func newHttpError(message string, err error) *httpError {
	he := &httpError{
		Message: message,
	}
	if err != nil {
		he.Error = err.Error()
	}
	return he
}
