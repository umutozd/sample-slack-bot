package storage

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrNotFound = errors.New("entry not found")
	ErrExists   = errors.New("entry already exists")
)

// Storage defines the interface for bot data persistence
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
	
	// Custom data storage
	Set(key string, value interface{}) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) bool
	
	// Batch operations
	SetMany(data map[string]interface{}) error
	GetMany(keys []string) (map[string]interface{}, error)
	DeleteMany(keys []string) error
	
	// TTL operations
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
	
	// Transactions
	Transaction(fn func(Storage) error) error
	
	// Cleanup
	Close() error
}

// Team represents a Slack team/workspace
type Team struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Domain       string            `json:"domain"`
	AppID        string            `json:"app_id"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	BotUserID    string            `json:"bot_user_id"`
	Scopes       []string          `json:"scopes"`
	InstalledAt  time.Time         `json:"installed_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Metadata     map[string]string `json:"metadata"`
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
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Metadata  map[string]string `json:"metadata"`
}

// Channel represents a Slack channel
type Channel struct {
	ID          string            `json:"id"`
	TeamID      string            `json:"team_id"`
	Name        string            `json:"name"`
	Topic       string            `json:"topic"`
	Purpose     string            `json:"purpose"`
	IsPrivate   bool              `json:"is_private"`
	IsArchived  bool              `json:"is_archived"`
	MemberCount int               `json:"member_count"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// NewMemoryStorage creates a new in-memory storage implementation
func NewMemoryStorage() Storage {
	return &memoryStorage{
		teams:    make(map[string]*Team),
		users:    make(map[string]*User),
		channels: make(map[string]*Channel),
		data:     make(map[string]interface{}),
		ttls:     make(map[string]time.Time),
	}
}

// NewRedisStorage creates a new Redis-backed storage implementation
func NewRedisStorage(addr, password string, db int) (Storage, error) {
	return newRedisStorage(addr, password, db)
}

// NewPostgresStorage creates a new PostgreSQL-backed storage implementation
func NewPostgresStorage(dsn string) (Storage, error) {
	return newPostgresStorage(dsn)
}

// NewSQLiteStorage creates a new SQLite-backed storage implementation
func NewSQLiteStorage(filepath string) (Storage, error) {
	return newSQLiteStorage(filepath)
}

// NewLayeredStorage creates a layered storage with cache and persistent layers
func NewLayeredStorage(cache, persistent Storage) Storage {
	return &layeredStorage{
		cache:      cache,
		persistent: persistent,
	}
}