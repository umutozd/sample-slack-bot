package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Storage defines the interface for storing bot data
type Storage interface {
	// Teams
	StoreTeam(team *Team) error
	GetTeam(teamID string) (*Team, error)
	UpdateTeam(team *Team) error
	DeleteTeam(teamID string) error
	
	// Users
	StoreUser(user *User) error
	GetUser(userID string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(userID string) error
	
	// Channels
	StoreChannel(channel *Channel) error
	GetChannel(channelID string) (*Channel, error)
	UpdateChannel(channel *Channel) error
	DeleteChannel(channelID string) error
	
	// Custom data
	Set(key string, value interface{}) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) bool
	
	// Batch operations
	SetMany(data map[string]interface{}) error
	GetMany(keys []string) (map[string]interface{}, error)
	
	// Expiry
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	
	// Health check
	Ping() error
}

// Team represents a Slack team/workspace
type Team struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	BotUserID    string    `json:"bot_user_id"`
	BotToken     string    `json:"bot_token"`
	Scope        string    `json:"scope"`
	InstalledAt  time.Time `json:"installed_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// User represents a Slack user
type User struct {
	ID        string            `json:"id"`
	TeamID    string            `json:"team_id"`
	Name      string            `json:"name"`
	RealName  string            `json:"real_name"`
	Email     string            `json:"email"`
	IsBot     bool              `json:"is_bot"`
	IsAdmin   bool              `json:"is_admin"`
	IsOwner   bool              `json:"is_owner"`
	Timezone  string            `json:"timezone"`
	Locale    string            `json:"locale"`
	Profile   map[string]string `json:"profile"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Channel represents a Slack channel
type Channel struct {
	ID          string    `json:"id"`
	TeamID      string    `json:"team_id"`
	Name        string    `json:"name"`
	IsChannel   bool      `json:"is_channel"`
	IsGroup     bool      `json:"is_group"`
	IsIM        bool      `json:"is_im"`
	IsPrivate   bool      `json:"is_private"`
	IsArchived  bool      `json:"is_archived"`
	IsGeneral   bool      `json:"is_general"`
	CreatorID   string    `json:"creator_id"`
	Topic       string    `json:"topic"`
	Purpose     string    `json:"purpose"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MemoryStorage is an in-memory implementation of Storage
type MemoryStorage struct {
	teams    map[string]*Team
	users    map[string]*User
	channels map[string]*Channel
	data     map[string]interface{}
	ttl      map[string]time.Time
	mutex    sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() Storage {
	return &MemoryStorage{
		teams:    make(map[string]*Team),
		users:    make(map[string]*User),
		channels: make(map[string]*Channel),
		data:     make(map[string]interface{}),
		ttl:      make(map[string]time.Time),
	}
}

// StoreTeam stores a team
func (m *MemoryStorage) StoreTeam(team *Team) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	team.UpdatedAt = time.Now()
	if team.InstalledAt.IsZero() {
		team.InstalledAt = time.Now()
	}
	
	m.teams[team.ID] = team
	return nil
}

// GetTeam retrieves a team by ID
func (m *MemoryStorage) GetTeam(teamID string) (*Team, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	team, exists := m.teams[teamID]
	if !exists {
		return nil, fmt.Errorf("team not found: %s", teamID)
	}
	
	return team, nil
}

// UpdateTeam updates an existing team
func (m *MemoryStorage) UpdateTeam(team *Team) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.teams[team.ID]; !exists {
		return fmt.Errorf("team not found: %s", team.ID)
	}
	
	team.UpdatedAt = time.Now()
	m.teams[team.ID] = team
	return nil
}

// DeleteTeam deletes a team
func (m *MemoryStorage) DeleteTeam(teamID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.teams[teamID]; !exists {
		return fmt.Errorf("team not found: %s", teamID)
	}
	
	delete(m.teams, teamID)
	return nil
}

// StoreUser stores a user
func (m *MemoryStorage) StoreUser(user *User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	
	m.users[user.ID] = user
	return nil
}

// GetUser retrieves a user by ID
func (m *MemoryStorage) GetUser(userID string) (*User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	user, exists := m.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	
	return user, nil
}

// UpdateUser updates an existing user
func (m *MemoryStorage) UpdateUser(user *User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.users[user.ID]; !exists {
		return fmt.Errorf("user not found: %s", user.ID)
	}
	
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

// DeleteUser deletes a user
func (m *MemoryStorage) DeleteUser(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.users[userID]; !exists {
		return fmt.Errorf("user not found: %s", userID)
	}
	
	delete(m.users, userID)
	return nil
}

// StoreChannel stores a channel
func (m *MemoryStorage) StoreChannel(channel *Channel) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	channel.UpdatedAt = time.Now()
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = time.Now()
	}
	
	m.channels[channel.ID] = channel
	return nil
}

// GetChannel retrieves a channel by ID
func (m *MemoryStorage) GetChannel(channelID string) (*Channel, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	channel, exists := m.channels[channelID]
	if !exists {
		return nil, fmt.Errorf("channel not found: %s", channelID)
	}
	
	return channel, nil
}

// UpdateChannel updates an existing channel
func (m *MemoryStorage) UpdateChannel(channel *Channel) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.channels[channel.ID]; !exists {
		return fmt.Errorf("channel not found: %s", channel.ID)
	}
	
	channel.UpdatedAt = time.Now()
	m.channels[channel.ID] = channel
	return nil
}

// DeleteChannel deletes a channel
func (m *MemoryStorage) DeleteChannel(channelID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.channels[channelID]; !exists {
		return fmt.Errorf("channel not found: %s", channelID)
	}
	
	delete(m.channels, channelID)
	return nil
}

// Set stores a value with the given key
func (m *MemoryStorage) Set(key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.data[key] = value
	return nil
}

// Get retrieves a value by key
func (m *MemoryStorage) Get(key string, dest interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	value, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}
	
	// Check TTL
	if ttl, hasTTL := m.ttl[key]; hasTTL && time.Now().After(ttl) {
		delete(m.data, key)
		delete(m.ttl, key)
		return fmt.Errorf("key expired: %s", key)
	}
	
	// Simple type assertion for common types
	switch dest := dest.(type) {
	case *string:
		if str, ok := value.(string); ok {
			*dest = str
		} else {
			return fmt.Errorf("value is not a string")
		}
	case *int:
		if i, ok := value.(int); ok {
			*dest = i
		} else {
			return fmt.Errorf("value is not an int")
		}
	case *bool:
		if b, ok := value.(bool); ok {
			*dest = b
		} else {
			return fmt.Errorf("value is not a bool")
		}
	default:
		// For complex types, use JSON marshaling/unmarshaling
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		
		if err := json.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}
	}
	
	return nil
}

// Delete removes a key
func (m *MemoryStorage) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	delete(m.data, key)
	delete(m.ttl, key)
	return nil
}

// Exists checks if a key exists
func (m *MemoryStorage) Exists(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	_, exists := m.data[key]
	return exists
}

// SetMany stores multiple key-value pairs
func (m *MemoryStorage) SetMany(data map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	for key, value := range data {
		m.data[key] = value
	}
	
	return nil
}

// GetMany retrieves multiple values by keys
func (m *MemoryStorage) GetMany(keys []string) (map[string]interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := make(map[string]interface{})
	
	for _, key := range keys {
		if value, exists := m.data[key]; exists {
			// Check TTL
			if ttl, hasTTL := m.ttl[key]; hasTTL && time.Now().After(ttl) {
				delete(m.data, key)
				delete(m.ttl, key)
				continue
			}
			result[key] = value
		}
	}
	
	return result, nil
}

// SetWithTTL stores a value with a time-to-live
func (m *MemoryStorage) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.data[key] = value
	m.ttl[key] = time.Now().Add(ttl)
	return nil
}

// Ping checks if the storage is available
func (m *MemoryStorage) Ping() error {
	return nil // Memory storage is always available
}

// Cleanup removes expired keys (should be called periodically)
func (m *MemoryStorage) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	now := time.Now()
	for key, expiry := range m.ttl {
		if now.After(expiry) {
			delete(m.data, key)
			delete(m.ttl, key)
		}
	}
}