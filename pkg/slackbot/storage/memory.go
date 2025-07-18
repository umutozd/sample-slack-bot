package storage

import (
	"encoding/json"
	"sync"
	"time"
)

// memoryStorage implements the Storage interface using in-memory maps
type memoryStorage struct {
	teams    map[string]*Team
	users    map[string]*User
	channels map[string]*Channel
	data     map[string]interface{}
	ttls     map[string]time.Time
	mu       sync.RWMutex
}

// Team operations
func (m *memoryStorage) StoreTeam(team *Team) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.teams[team.ID]; exists {
		return ErrExists
	}
	
	team.InstalledAt = time.Now()
	team.UpdatedAt = time.Now()
	m.teams[team.ID] = team
	return nil
}

func (m *memoryStorage) GetTeam(teamID string) (*Team, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	team, exists := m.teams[teamID]
	if !exists {
		return nil, ErrNotFound
	}
	
	// Return a copy to prevent external modification
	teamCopy := *team
	return &teamCopy, nil
}

func (m *memoryStorage) UpdateTeam(team *Team) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.teams[team.ID]; !exists {
		return ErrNotFound
	}
	
	team.UpdatedAt = time.Now()
	m.teams[team.ID] = team
	return nil
}

func (m *memoryStorage) DeleteTeam(teamID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.teams[teamID]; !exists {
		return ErrNotFound
	}
	
	delete(m.teams, teamID)
	return nil
}

// User operations
func (m *memoryStorage) StoreUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.users[user.ID]; exists {
		return ErrExists
	}
	
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *memoryStorage) GetUser(userID string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	user, exists := m.users[userID]
	if !exists {
		return nil, ErrNotFound
	}
	
	// Return a copy to prevent external modification
	userCopy := *user
	return &userCopy, nil
}

func (m *memoryStorage) UpdateUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.users[user.ID]; !exists {
		return ErrNotFound
	}
	
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *memoryStorage) DeleteUser(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.users[userID]; !exists {
		return ErrNotFound
	}
	
	delete(m.users, userID)
	return nil
}

// Channel operations
func (m *memoryStorage) StoreChannel(channel *Channel) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.channels[channel.ID]; exists {
		return ErrExists
	}
	
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()
	m.channels[channel.ID] = channel
	return nil
}

func (m *memoryStorage) GetChannel(channelID string) (*Channel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	channel, exists := m.channels[channelID]
	if !exists {
		return nil, ErrNotFound
	}
	
	// Return a copy to prevent external modification
	channelCopy := *channel
	return &channelCopy, nil
}

func (m *memoryStorage) UpdateChannel(channel *Channel) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.channels[channel.ID]; !exists {
		return ErrNotFound
	}
	
	channel.UpdatedAt = time.Now()
	m.channels[channel.ID] = channel
	return nil
}

func (m *memoryStorage) DeleteChannel(channelID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.channels[channelID]; !exists {
		return ErrNotFound
	}
	
	delete(m.channels, channelID)
	return nil
}

// Custom data operations
func (m *memoryStorage) Set(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.data[key] = value
	return nil
}

func (m *memoryStorage) Get(key string, dest interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	value, exists := m.data[key]
	if !exists {
		return ErrNotFound
	}
	
	// Use JSON marshaling/unmarshaling to copy the value
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

func (m *memoryStorage) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.data[key]; !exists {
		return ErrNotFound
	}
	
	delete(m.data, key)
	delete(m.ttls, key) // Also remove TTL if exists
	return nil
}

func (m *memoryStorage) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Check if key exists and hasn't expired
	if _, exists := m.data[key]; !exists {
		return false
	}
	
	// Check TTL
	if expiry, hasTTL := m.ttls[key]; hasTTL {
		if time.Now().After(expiry) {
			// Key has expired, remove it
			delete(m.data, key)
			delete(m.ttls, key)
			return false
		}
	}
	
	return true
}

// Batch operations
func (m *memoryStorage) SetMany(data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for key, value := range data {
		m.data[key] = value
	}
	return nil
}

func (m *memoryStorage) GetMany(keys []string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make(map[string]interface{})
	for _, key := range keys {
		if value, exists := m.data[key]; exists {
			// Check TTL
			if expiry, hasTTL := m.ttls[key]; hasTTL {
				if time.Now().After(expiry) {
					continue // Skip expired keys
				}
			}
			result[key] = value
		}
	}
	return result, nil
}

func (m *memoryStorage) DeleteMany(keys []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, key := range keys {
		delete(m.data, key)
		delete(m.ttls, key)
	}
	return nil
}

// TTL operations
func (m *memoryStorage) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.data[key] = value
	m.ttls[key] = time.Now().Add(ttl)
	return nil
}

func (m *memoryStorage) GetTTL(key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	expiry, exists := m.ttls[key]
	if !exists {
		return 0, ErrNotFound
	}
	
	remaining := time.Until(expiry)
	if remaining <= 0 {
		return 0, ErrNotFound
	}
	
	return remaining, nil
}

// Transaction operations
func (m *memoryStorage) Transaction(fn func(Storage) error) error {
	// For in-memory storage, we implement a simple transaction by creating a copy
	// and rolling back on error
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Create backups
	teamsBackup := make(map[string]*Team)
	usersBackup := make(map[string]*User)
	channelsBackup := make(map[string]*Channel)
	dataBackup := make(map[string]interface{})
	ttlsBackup := make(map[string]time.Time)
	
	// Copy current state
	for k, v := range m.teams {
		teamCopy := *v
		teamsBackup[k] = &teamCopy
	}
	for k, v := range m.users {
		userCopy := *v
		usersBackup[k] = &userCopy
	}
	for k, v := range m.channels {
		channelCopy := *v
		channelsBackup[k] = &channelCopy
	}
	for k, v := range m.data {
		dataBackup[k] = v
	}
	for k, v := range m.ttls {
		ttlsBackup[k] = v
	}
	
	// Execute the transaction
	err := fn(m)
	if err != nil {
		// Rollback
		m.teams = teamsBackup
		m.users = usersBackup
		m.channels = channelsBackup
		m.data = dataBackup
		m.ttls = ttlsBackup
		return err
	}
	
	return nil
}

// Cleanup
func (m *memoryStorage) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Clear all data
	m.teams = make(map[string]*Team)
	m.users = make(map[string]*User)
	m.channels = make(map[string]*Channel)
	m.data = make(map[string]interface{})
	m.ttls = make(map[string]time.Time)
	
	return nil
}

// Cleanup expired entries (should be called periodically)
func (m *memoryStorage) cleanupExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := time.Now()
	for key, expiry := range m.ttls {
		if now.After(expiry) {
			delete(m.data, key)
			delete(m.ttls, key)
		}
	}
}

// StartCleanupRoutine starts a goroutine to periodically clean up expired entries
func (m *memoryStorage) StartCleanupRoutine(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for range ticker.C {
			m.cleanupExpired()
		}
	}()
}