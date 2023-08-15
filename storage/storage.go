package storage

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

var ErrNotFound = errors.New("entry not found")

type Storage interface {
	AddSlackTeam(teamID, appID, accessToken, refreshToken string) (*SlackTeam, error)
	GetSlackTeam(teamID string) (*SlackTeam, error)
	UpdateSlackTeamTokens(teamID, accessToken, refreshToken string) (*SlackTeam, error)
}

type redisStorage struct {
	redis *redis.Client
}

func NewRedisStorage(address, password string, db int) (Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis server: %w", err)
	}
	return &redisStorage{redis: client}, nil
}
