package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/storage"
)

const (
	mainActionsBlockID      = "main-actions"
	modalTopicSelectBlockID = "modal-topic-select"

	actionToggleText       = "action-toggle-text"
	actionOpenModal        = "action-open-modal"
	actionModalTopicSelect = "action-modal-topic-select"

	modalCallbackID = "modal-callback"
)

func (s *server) Interactive(w http.ResponseWriter, r *http.Request) {
	cb, ok := parseInteractivePayload(w, r)
	if !ok {
		return
	}

	team, err := s.st.GetSlackTeam(cb.Team.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeJsonResponse(w, http.StatusNotFound, newHttpError(fmt.Sprintf("team with ID '%s' not found", cb.Team.ID), err))
		} else {
			writeJsonResponse(w, http.StatusInternalServerError, newHttpError(fmt.Sprintf("failed to find team with ID '%s'", cb.Team.ID), err))
		}
		return
	}

	// check block actions first (e.g. button clicks)
	if len(cb.ActionCallback.BlockActions) > 0 {
		for _, act := range cb.ActionCallback.BlockActions {
			switch act.ActionID {
			case actionToggleText:
				if act.Value == "version-1" {
					s.publishHomeView(team, cb.User.ID, 2)
				} else {
					s.publishHomeView(team, cb.User.ID, 1)
				}
			case actionOpenModal:
				if _, err := slack.New(team.AccessToken).OpenView(cb.TriggerID, newModalView("")); err != nil {
					logrus.WithError(detailedSlackError(err)).Error("failed to open modal view")
				}
			case actionModalTopicSelect:
				topic := act.SelectedOption.Value

				if _, err := slack.New(team.AccessToken).UpdateView(newModalView(topic), "", "", cb.View.ID); err != nil {
					logrus.WithError(detailedSlackError(err)).Error("failed to open modal view")
				}
			default:
				logrus.WithField("action_id", act.ActionID).Warnf("Interactive: unhandled block action")
			}
		}
	}

}

func newModalView(topic string) slack.ModalViewRequest {
	selectElement := &slack.SelectBlockElement{
		Type:        "static_select",
		Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "Choose one!", false, false),
		ActionID:    actionModalTopicSelect,
		Options: []*slack.OptionBlockObject{
			{
				Text:  slack.NewTextBlockObject(slack.PlainTextType, "Option 1", false, false),
				Value: "option-1",
			},
			{
				Text:  slack.NewTextBlockObject(slack.PlainTextType, "Option 2", false, false),
				Value: "option-2",
			},
			{
				Text:  slack.NewTextBlockObject(slack.PlainTextType, "Option 3", false, false),
				Value: "option-3",
			},
		},
	}
	descriptionText := ""
	switch topic {
	case "option-1":
		selectElement.InitialOption = selectElement.Options[0]
		descriptionText = "You have selected Option 1!"
	case "option-2":
		selectElement.InitialOption = selectElement.Options[1]
		descriptionText = "You have selected Option 2!"
	case "option-3":
		selectElement.InitialOption = selectElement.Options[2]
		descriptionText = "You have selected Option 3!"
	}
	view := slack.ModalViewRequest{
		Type:         slack.VTModal,
		Title:        slack.NewTextBlockObject(slack.PlainTextType, "Sample Modal", false, false),
		Close:        slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false),
		Submit:       slack.NewTextBlockObject(slack.PlainTextType, "Close - but positive", false, false),
		ClearOnClose: true,
		CallbackID:   modalCallbackID,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					Type:           slack.MBTInput,
					BlockID:        modalTopicSelectBlockID,
					Label:          slack.NewTextBlockObject(slack.PlainTextType, "Select topic", false, false),
					DispatchAction: true, // changes to the UI elements (e.g. selections) are sent to the server as callbacks
					Element:        selectElement,
				},
			},
		},
	}
	if descriptionText != "" {
		view.Blocks.BlockSet = append(view.Blocks.BlockSet, slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.PlainTextType, descriptionText, false, false),
			nil,
			nil,
		))
	}
	return view
}

func parseInteractivePayload(w http.ResponseWriter, r *http.Request) (callback *slack.InteractionCallback, ok bool) {
	// interactive payload's body consists of a url-encoded form with a single
	// field "payload" whose value is a JSON document
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, newHttpError("failed to read request body", err))
		return
	}
	parsed, err := url.ParseQuery(string(body))
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, newHttpError("failed to parse request body as form-urlencoded", err))
		return
	}
	payload := parsed.Get("payload")
	if payload == "" {
		writeJsonResponse(w, http.StatusBadRequest, newHttpError("form must have non-empty payload field", err))
		return
	}
	logrus.Debugf("Interactive payload:\n%s", payload)

	cb := &slack.InteractionCallback{}
	if err = json.Unmarshal([]byte(payload), cb); err != nil {
		writeJsonResponse(w, http.StatusBadRequest, newHttpError("invalid interaction json payload", err))
		return
	}

	ok = true
	callback = cb
	return
}
