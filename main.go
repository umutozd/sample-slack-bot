package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/umutozd/sample-slack-bot/server"
)

var config = &server.Config{
	Debug:             true,
	Port:              8080,
	RedisAddress:      "localhost:6379",
	RedisPassword:     "",
	RedisDB:           0,
	SlackClientID:     os.Getenv("SLACK_CLIENT_ID"),
	SlackClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
}

func main() {
	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	s, err := server.NewServer(config)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create new server instance")
	}
	s.ListenAndServe()
}
