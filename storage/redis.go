package storage

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

func (rs *redisStorage) AddSlackTeam(teamID, appID, accessToken, refreshToken string) (*SlackTeam, error) {
	team := &SlackTeam{
		ID:           teamID,
		AppID:        appID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err := rs.setSlackTeam(team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (rs *redisStorage) GetSlackTeam(teamID string) (*SlackTeam, error) {
	jsonValue, err := rs.redis.Get(teamID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get slack team: %w", err)
	}
	team := &SlackTeam{}
	if err = json.Unmarshal([]byte(jsonValue), team); err != nil {
		return nil, fmt.Errorf("error unmarshaling team value: %w", err)
	}
	return team, nil
}

func (rs *redisStorage) UpdateSlackTeamTokens(teamID, accessToken, refreshToken string) (*SlackTeam, error) {
	team, err := rs.GetSlackTeam(teamID)
	if err != nil {
		return nil, err
	}
	team.AccessToken = accessToken
	team.RefreshToken = refreshToken
	err = rs.setSlackTeam(team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (rs *redisStorage) setSlackTeam(team *SlackTeam) error {
	jsonValue, err := json.Marshal(team)
	if err != nil {
		return fmt.Errorf("failed to json-marshal team object: %w", err)
	}
	if err = rs.redis.Set(team.ID, jsonValue, 0).Err(); err != nil {
		return fmt.Errorf("error setting team data on redis: %w", err)
	}
	return nil
}
