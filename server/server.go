package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/umutozd/sample-slack-bot/storage"
)

type server struct {
	cfg *Config

	st storage.Storage
}

type Server interface {
	ListenAndServe()
}

// NewServer initializes a new Server instance from the given Config.
func NewServer(config *Config) (Server, error) {
	st, err := storage.NewRedisStorage(config.RedisAddress, config.RedisPassword, config.RedisDB)
	if err != nil {
		return nil, fmt.Errorf("error creating redis storage: %w", err)
	}
	return &server{
		st:  st,
		cfg: config,
	}, nil
}

func (s *server) ListenAndServe() {
	router := http.NewServeMux()
	router.HandleFunc("/slack/install", s.AppInstallation)
	router.HandleFunc("/slack/events", s.Events)
	router.HandleFunc("/slack/interactive", s.Interactive)

	logrus.Infof("listening on port %d", s.cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", s.cfg.Port), router); err != nil {
		logrus.WithError(err).Fatal("server stopped with error")
	} else {
		logrus.Info("server stopped")
	}
}

func (s *server) refreshTeamTokensIfRequired(team *storage.SlackTeam) error {
	client := slack.New(team.AccessToken)
	_, err := client.GetTeamInfo()
	if err != nil {
		if !strings.Contains(err.Error(), "token_expired") {
			return fmt.Errorf("slack client unable to retrieve team info: %w", err)
		}

		resp, err := slack.RefreshOAuthV2Token(http.DefaultClient, s.cfg.SlackClientID, s.cfg.SlackClientSecret, team.RefreshToken)
		if err != nil {
			return fmt.Errorf("slack client unable to refresh tokens: %w", err)
		}

		updatedTeam, err := s.st.UpdateSlackTeamTokens(team.ID, resp.AccessToken, resp.RefreshToken)
		if err != nil {
			return fmt.Errorf("error updating team tokens in storage: %w", err)
		}

		team.AccessToken = updatedTeam.AccessToken
		team.RefreshToken = updatedTeam.RefreshToken
	}

	return nil
}
