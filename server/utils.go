package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func parseRequestBody(r *http.Request, unmarshalTo any) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if err = json.Unmarshal(body, unmarshalTo); err != nil {
		return nil, fmt.Errorf("failed to json-unmarshal request body: %w", err)
	}

	return body, nil
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, body any) {
	if statusCode >= 400 {
		logrus.Warnf("writing response with status=%d and body: %#v", statusCode, body)
	}
	w.Header().Set("Content-Type", "application/json")

	var writeError error
	response, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeError = w.Write([]byte(fmt.Sprintf(`{"error": "%v", "message": "server failed to write response"}`, err)))
	} else {
		w.WriteHeader(statusCode)
		_, writeError = w.Write(response)
	}

	if writeError != nil {
		logrus.WithError(err).Errorf("failed to write json response with status=%d and body: %s", statusCode, string(response))
	}
}

func detailedSlackError(err error) error {
	if e, ok := err.(slack.SlackErrorResponse); ok {
		metadata, _ := json.Marshal(e.ResponseMetadata)
		return fmt.Errorf("%v - %s", e.Err, string(metadata))
	}

	return err
}
