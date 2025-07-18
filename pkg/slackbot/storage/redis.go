package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisStorage is a Redis-based implementation of Storage
type RedisStorage struct {
	client *redis.Client
	prefix string
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(addr, password string, db int) (Storage, error) {
	return NewRedisStorageWithOptions(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// NewRedisStorageWithOptions creates a new Redis storage instance with custom options
func NewRedisStorageWithOptions(options *redis.Options) (Storage, error) {
	client := redis.NewClient(options)
	
	// Test connection
	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return &RedisStorage{
		client: client,
		prefix: "slackbot:",
	}, nil
}

// WithPrefix sets a custom prefix for all Redis keys
func (r *RedisStorage) WithPrefix(prefix string) *RedisStorage {
	r.prefix = prefix
	return r
}

// key returns the prefixed key
func (r *RedisStorage) key(suffix string) string {
	return r.prefix + suffix
}

// StoreTeam stores a team in Redis
func (r *RedisStorage) StoreTeam(team *Team) error {
	team.UpdatedAt = time.Now()
	if team.InstalledAt.IsZero() {
		team.InstalledAt = time.Now()
	}
	
	data, err := json.Marshal(team)
	if err != nil {
		return fmt.Errorf("failed to marshal team: %w", err)
	}
	
	return r.client.Set(r.key("teams:"+team.ID), data, 0).Err()
}

// GetTeam retrieves a team from Redis
func (r *RedisStorage) GetTeam(teamID string) (*Team, error) {
	data, err := r.client.Get(r.key("teams:" + teamID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("team not found: %s", teamID)
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	
	var team Team
	if err := json.Unmarshal([]byte(data), &team); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team: %w", err)
	}
	
	return &team, nil
}

// UpdateTeam updates an existing team in Redis
func (r *RedisStorage) UpdateTeam(team *Team) error {
	// Check if team exists
	if _, err := r.GetTeam(team.ID); err != nil {
		return err
	}
	
	return r.StoreTeam(team)
}

// DeleteTeam deletes a team from Redis
func (r *RedisStorage) DeleteTeam(teamID string) error {
	result := r.client.Del(r.key("teams:" + teamID))
	if result.Err() != nil {
		return fmt.Errorf("failed to delete team: %w", result.Err())
	}
	
	if result.Val() == 0 {
		return fmt.Errorf("team not found: %s", teamID)
	}
	
	return nil
}

// StoreUser stores a user in Redis
func (r *RedisStorage) StoreUser(user *User) error {
	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	
	return r.client.Set(r.key("users:"+user.ID), data, 0).Err()
}

// GetUser retrieves a user from Redis
func (r *RedisStorage) GetUser(userID string) (*User, error) {
	data, err := r.client.Get(r.key("users:" + userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found: %s", userID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	var user User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	
	return &user, nil
}

// UpdateUser updates an existing user in Redis
func (r *RedisStorage) UpdateUser(user *User) error {
	// Check if user exists
	if _, err := r.GetUser(user.ID); err != nil {
		return err
	}
	
	return r.StoreUser(user)
}

// DeleteUser deletes a user from Redis
func (r *RedisStorage) DeleteUser(userID string) error {
	result := r.client.Del(r.key("users:" + userID))
	if result.Err() != nil {
		return fmt.Errorf("failed to delete user: %w", result.Err())
	}
	
	if result.Val() == 0 {
		return fmt.Errorf("user not found: %s", userID)
	}
	
	return nil
}

// StoreChannel stores a channel in Redis
func (r *RedisStorage) StoreChannel(channel *Channel) error {
	channel.UpdatedAt = time.Now()
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = time.Now()
	}
	
	data, err := json.Marshal(channel)
	if err != nil {
		return fmt.Errorf("failed to marshal channel: %w", err)
	}
	
	return r.client.Set(r.key("channels:"+channel.ID), data, 0).Err()
}

// GetChannel retrieves a channel from Redis
func (r *RedisStorage) GetChannel(channelID string) (*Channel, error) {
	data, err := r.client.Get(r.key("channels:" + channelID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("channel not found: %s", channelID)
		}
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	
	var channel Channel
	if err := json.Unmarshal([]byte(data), &channel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal channel: %w", err)
	}
	
	return &channel, nil
}

// UpdateChannel updates an existing channel in Redis
func (r *RedisStorage) UpdateChannel(channel *Channel) error {
	// Check if channel exists
	if _, err := r.GetChannel(channel.ID); err != nil {
		return err
	}
	
	return r.StoreChannel(channel)
}

// DeleteChannel deletes a channel from Redis
func (r *RedisStorage) DeleteChannel(channelID string) error {
	result := r.client.Del(r.key("channels:" + channelID))
	if result.Err() != nil {
		return fmt.Errorf("failed to delete channel: %w", result.Err())
	}
	
	if result.Val() == 0 {
		return fmt.Errorf("channel not found: %s", channelID)
	}
	
	return nil
}

// Set stores a value with the given key
func (r *RedisStorage) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	return r.client.Set(r.key("data:"+key), data, 0).Err()
}

// Get retrieves a value by key
func (r *RedisStorage) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.key("data:" + key)).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get key: %w", err)
	}
	
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	
	return nil
}

// Delete removes a key
func (r *RedisStorage) Delete(key string) error {
	return r.client.Del(r.key("data:" + key)).Err()
}

// Exists checks if a key exists
func (r *RedisStorage) Exists(key string) bool {
	result := r.client.Exists(r.key("data:" + key))
	return result.Val() > 0
}

// SetMany stores multiple key-value pairs
func (r *RedisStorage) SetMany(data map[string]interface{}) error {
	pipe := r.client.Pipeline()
	
	for key, value := range data {
		valueData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		pipe.Set(r.key("data:"+key), valueData, 0)
	}
	
	_, err := pipe.Exec()
	return err
}

// GetMany retrieves multiple values by keys
func (r *RedisStorage) GetMany(keys []string) (map[string]interface{}, error) {
	pipe := r.client.Pipeline()
	
	redisKeys := make([]string, len(keys))
	for i, key := range keys {
		redisKeys[i] = r.key("data:" + key)
		pipe.Get(redisKeys[i])
	}
	
	cmds, err := pipe.Exec()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}
	
	result := make(map[string]interface{})
	
	for i, cmd := range cmds {
		stringCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		
		data, err := stringCmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue // Key doesn't exist
			}
			return nil, fmt.Errorf("failed to get key %s: %w", keys[i], err)
		}
		
		var value interface{}
		if err := json.Unmarshal([]byte(data), &value); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", keys[i], err)
		}
		
		result[keys[i]] = value
	}
	
	return result, nil
}

// SetWithTTL stores a value with a time-to-live
func (r *RedisStorage) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	return r.client.Set(r.key("data:"+key), data, ttl).Err()
}

// Ping checks if Redis is available
func (r *RedisStorage) Ping() error {
	return r.client.Ping().Err()
}

// Close closes the Redis connection
func (r *RedisStorage) Close() error {
	return r.client.Close()
}