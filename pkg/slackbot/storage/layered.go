package storage

import (
	"time"
)

// layeredStorage implements a two-layer storage system with cache and persistent layers
type layeredStorage struct {
	cache      Storage
	persistent Storage
}

// Team operations
func (l *layeredStorage) StoreTeam(team *Team) error {
	// Store in both cache and persistent storage
	if err := l.persistent.StoreTeam(team); err != nil {
		return err
	}
	
	// Try to store in cache, but don't fail if it doesn't work
	l.cache.StoreTeam(team)
	return nil
}

func (l *layeredStorage) GetTeam(teamID string) (*Team, error) {
	// Try cache first
	if team, err := l.cache.GetTeam(teamID); err == nil {
		return team, nil
	}
	
	// Fall back to persistent storage
	team, err := l.persistent.GetTeam(teamID)
	if err != nil {
		return nil, err
	}
	
	// Store in cache for future requests
	l.cache.StoreTeam(team)
	return team, nil
}

func (l *layeredStorage) UpdateTeam(team *Team) error {
	// Update in persistent storage first
	if err := l.persistent.UpdateTeam(team); err != nil {
		return err
	}
	
	// Update in cache
	l.cache.UpdateTeam(team)
	return nil
}

func (l *layeredStorage) DeleteTeam(teamID string) error {
	// Delete from persistent storage
	if err := l.persistent.DeleteTeam(teamID); err != nil {
		return err
	}
	
	// Delete from cache
	l.cache.DeleteTeam(teamID)
	return nil
}

// User operations
func (l *layeredStorage) StoreUser(user *User) error {
	// Store in both cache and persistent storage
	if err := l.persistent.StoreUser(user); err != nil {
		return err
	}
	
	// Try to store in cache, but don't fail if it doesn't work
	l.cache.StoreUser(user)
	return nil
}

func (l *layeredStorage) GetUser(userID string) (*User, error) {
	// Try cache first
	if user, err := l.cache.GetUser(userID); err == nil {
		return user, nil
	}
	
	// Fall back to persistent storage
	user, err := l.persistent.GetUser(userID)
	if err != nil {
		return nil, err
	}
	
	// Store in cache for future requests
	l.cache.StoreUser(user)
	return user, nil
}

func (l *layeredStorage) UpdateUser(user *User) error {
	// Update in persistent storage first
	if err := l.persistent.UpdateUser(user); err != nil {
		return err
	}
	
	// Update in cache
	l.cache.UpdateUser(user)
	return nil
}

func (l *layeredStorage) DeleteUser(userID string) error {
	// Delete from persistent storage
	if err := l.persistent.DeleteUser(userID); err != nil {
		return err
	}
	
	// Delete from cache
	l.cache.DeleteUser(userID)
	return nil
}

// Channel operations
func (l *layeredStorage) StoreChannel(channel *Channel) error {
	// Store in both cache and persistent storage
	if err := l.persistent.StoreChannel(channel); err != nil {
		return err
	}
	
	// Try to store in cache, but don't fail if it doesn't work
	l.cache.StoreChannel(channel)
	return nil
}

func (l *layeredStorage) GetChannel(channelID string) (*Channel, error) {
	// Try cache first
	if channel, err := l.cache.GetChannel(channelID); err == nil {
		return channel, nil
	}
	
	// Fall back to persistent storage
	channel, err := l.persistent.GetChannel(channelID)
	if err != nil {
		return nil, err
	}
	
	// Store in cache for future requests
	l.cache.StoreChannel(channel)
	return channel, nil
}

func (l *layeredStorage) UpdateChannel(channel *Channel) error {
	// Update in persistent storage first
	if err := l.persistent.UpdateChannel(channel); err != nil {
		return err
	}
	
	// Update in cache
	l.cache.UpdateChannel(channel)
	return nil
}

func (l *layeredStorage) DeleteChannel(channelID string) error {
	// Delete from persistent storage
	if err := l.persistent.DeleteChannel(channelID); err != nil {
		return err
	}
	
	// Delete from cache
	l.cache.DeleteChannel(channelID)
	return nil
}

// Custom data operations
func (l *layeredStorage) Set(key string, value interface{}) error {
	// Store in both cache and persistent storage
	if err := l.persistent.Set(key, value); err != nil {
		return err
	}
	
	// Try to store in cache, but don't fail if it doesn't work
	l.cache.Set(key, value)
	return nil
}

func (l *layeredStorage) Get(key string, dest interface{}) error {
	// Try cache first
	if err := l.cache.Get(key, dest); err == nil {
		return nil
	}
	
	// Fall back to persistent storage
	if err := l.persistent.Get(key, dest); err != nil {
		return err
	}
	
	// Store in cache for future requests
	l.cache.Set(key, dest)
	return nil
}

func (l *layeredStorage) Delete(key string) error {
	// Delete from persistent storage
	if err := l.persistent.Delete(key); err != nil {
		return err
	}
	
	// Delete from cache
	l.cache.Delete(key)
	return nil
}

func (l *layeredStorage) Exists(key string) bool {
	// Check cache first
	if l.cache.Exists(key) {
		return true
	}
	
	// Check persistent storage
	return l.persistent.Exists(key)
}

// Batch operations
func (l *layeredStorage) SetMany(data map[string]interface{}) error {
	// Store in both cache and persistent storage
	if err := l.persistent.SetMany(data); err != nil {
		return err
	}
	
	// Try to store in cache, but don't fail if it doesn't work
	l.cache.SetMany(data)
	return nil
}

func (l *layeredStorage) GetMany(keys []string) (map[string]interface{}, error) {
	// Try to get from cache first
	cacheResult, _ := l.cache.GetMany(keys)
	
	// Find missing keys
	missingKeys := []string{}
	for _, key := range keys {
		if _, exists := cacheResult[key]; !exists {
			missingKeys = append(missingKeys, key)
		}
	}
	
	// Get missing keys from persistent storage
	if len(missingKeys) > 0 {
		persistentResult, err := l.persistent.GetMany(missingKeys)
		if err != nil {
			return cacheResult, err
		}
		
		// Merge results
		for key, value := range persistentResult {
			cacheResult[key] = value
		}
		
		// Store missing keys in cache
		l.cache.SetMany(persistentResult)
	}
	
	return cacheResult, nil
}

func (l *layeredStorage) DeleteMany(keys []string) error {
	// Delete from persistent storage
	if err := l.persistent.DeleteMany(keys); err != nil {
		return err
	}
	
	// Delete from cache
	l.cache.DeleteMany(keys)
	return nil
}

// TTL operations
func (l *layeredStorage) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	// Store in both cache and persistent storage
	if err := l.persistent.SetWithTTL(key, value, ttl); err != nil {
		return err
	}
	
	// Try to store in cache, but don't fail if it doesn't work
	l.cache.SetWithTTL(key, value, ttl)
	return nil
}

func (l *layeredStorage) GetTTL(key string) (time.Duration, error) {
	// Try cache first
	if ttl, err := l.cache.GetTTL(key); err == nil {
		return ttl, nil
	}
	
	// Fall back to persistent storage
	return l.persistent.GetTTL(key)
}

// Transaction operations
func (l *layeredStorage) Transaction(fn func(Storage) error) error {
	// For layered storage, we execute the transaction on the persistent layer
	// and invalidate the cache after successful completion
	return l.persistent.Transaction(func(s Storage) error {
		err := fn(s)
		if err == nil {
			// Clear cache after successful transaction
			// This is a simple approach - in production, you might want to be more selective
			if closer, ok := l.cache.(interface{ Close() error }); ok {
				closer.Close()
			}
		}
		return err
	})
}

// Cleanup
func (l *layeredStorage) Close() error {
	// Close both layers
	var cacheErr, persistentErr error
	
	if l.cache != nil {
		cacheErr = l.cache.Close()
	}
	
	if l.persistent != nil {
		persistentErr = l.persistent.Close()
	}
	
	// Return the first error encountered
	if cacheErr != nil {
		return cacheErr
	}
	return persistentErr
}