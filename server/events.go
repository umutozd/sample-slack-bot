package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/storage"
)

// Events handles incoming Slack events regarding the Slack bot.
//
// Event payload may contain a "challenge" field and in that case, the challenge data
// is sent back to the client as is. This is a mechanism to verify that the events endpoint
// of the bot is operational.
//
// Currently handled Slack events are "app_home_opened" and "message".
//
// Bot events are also received by this endpoint, for example the bot's message in a channel.
// Because we do not want to handle them, we ignore events whose payloads contain a non-empty
// "bot_profile" field.
func (s *server) Events(w http.ResponseWriter, r *http.Request) {
	event := &slackEvent{}
	if body, err := parseRequestBody(r, event); err != nil {
		writeJsonResponse(w, http.StatusBadRequest, newHttpError("failed to parse request body", err))
		return
	} else {
		logrus.Infof("Events payload:\n%s", string(body))
	}

	// slack api may send us challenges to check the health status of our server
	if event.Challenge != "" {
		logrus.WithField("challenge", event.Challenge).Infof("got challenge event")
		if _, err := w.Write([]byte(event.Challenge)); err != nil {
			writeJsonResponse(w, http.StatusInternalServerError, newHttpError("failed to write challenge response", err))
		}
		return
	}
	if event.Event.BotProfile.Name != "" {
		logrus.Debug("got bot event, skipping")
		return
	}
	logger := logrus.WithFields(logrus.Fields{
		"team_id":    event.TeamID,
		"user_id":    event.Event.User,
		"event_type": event.Type,
	})

	team, err := s.st.GetSlackTeam(event.TeamID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeJsonResponse(w, http.StatusNotFound, newHttpError(fmt.Sprintf("team with ID '%s' not found", event.TeamID), err))
		} else {
			writeJsonResponse(w, http.StatusInternalServerError, newHttpError(fmt.Sprintf("failed to find team with ID '%s'", event.TeamID), err))
		}
		return
	}

	logger.Info("refreshing team tokens")
	if err = s.refreshTeamTokensIfRequired(team); err != nil {
		writeJsonResponse(w, http.StatusInternalServerError, newHttpError("failed to refresh team tokens", err))
		return
	}

	switch event.Event.Type {
	case "app_home_opened":
		logger.Info("got app_home_opened event")
		s.publishHomeView(team, event.Event.User, 1)
	case "message":
		if event.Event.ThreadTS == "" {
			logger.Info("got direct message event")
			s.handleDirectMessage(w, event, team)
		} else {
			logger.Infof("got thread message event with thread_ts=%q", event.Event.ThreadTS)
			s.handleThreadMessage(w, event, team)
		}
	default:
		logger.Infof("got unhandled event %s", event.Event.Type)
	}
}

// publishHomeView sends a corresponding home view to the given user.
//
// toggleTextVersion parameter controls which version of the toggled text will be sent.
func (s *server) publishHomeView(team *storage.SlackTeam, userID string, toggleTextVersion int) {
	client := slack.New(team.AccessToken)

	var toggledText, toggleButtonValue string
	if toggleTextVersion == 1 {
		toggledText = "Version 1 of toggled text!"
		toggleButtonValue = "version-1"
	} else {
		toggledText = "Version 2 of toggled text!"
		toggleButtonValue = "version-2"
	}

	view := slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			// block kit UI
			BlockSet: []slack.Block{
				slack.NewHeaderBlock(
					slack.NewTextBlockObject(slack.PlainTextType, "Welcome to the Test App :tada:", false, false),
				),
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "Here's a very short and simple description of this bot! :grin:", false, false),
					nil,
					nil,
				),
				slack.NewActionBlock(
					mainActionsBlockID,
					slack.NewButtonBlockElement(actionToggleText, toggleButtonValue, slack.NewTextBlockObject(slack.PlainTextType, "Toggle Below Text", false, false)),
					slack.NewButtonBlockElement(actionOpenModal, "open", slack.NewTextBlockObject(slack.PlainTextType, "Open Modal", false, false)),
				),
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, toggledText, false, false),
					nil,
					nil,
				),
			},
		},
	}

	_, err := client.PublishView(userID, view, "")
	if err != nil {
		logrus.WithField("user_id", userID).WithError(detailedSlackError(err)).Error("failed to publish home tab view for user")
	}
}

func (s *server) handleDirectMessage(w http.ResponseWriter, event *slackEvent, team *storage.SlackTeam) {
	client := slack.New(team.AccessToken)

	channel, timestamp, _, err := client.SendMessage(
		event.Event.Channel,
		slack.MsgOptionText(fmt.Sprintf("Got your message:\n```%s```", event.Event.Text), false),
	)
	if err != nil {
		logrus.WithField("channel", event.Event.Channel).WithError(err).Error("error replying to direct message")
	} else {
		logrus.WithField("channel", channel).WithField("timestamp", timestamp).Info("replied to user's direct message")
	}
}

func (s *server) handleThreadMessage(w http.ResponseWriter, event *slackEvent, team *storage.SlackTeam) {
	client := slack.New(team.AccessToken)

	channel, timestamp, _, err := client.SendMessage(
		event.Event.Channel,
		slack.MsgOptionText(fmt.Sprintf("Got thread your message:\n```%s```", event.Event.Text), false),
		slack.MsgOptionTS(event.Event.ThreadTS),
	)
	if err != nil {
		logrus.WithField("channel", event.Event.Channel).WithError(err).Error("error replying to thread message")
	} else {
		logrus.WithField("channel", channel).WithField("timestamp", timestamp).Info("replied to user's thread message")
	}
}
