package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// Config holds the configuration for the Redis client.
// It uses time.Duration directly for timeouts and frequencies as is idiomatic in Go.
type Config struct {
	Host               string
	Port               string
	Password           string
	DB                 int // Use int for Redis DB
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleCheckFrequency time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
}

// Cache provides a wrapper around the Redis client for generic data caching.
type Cache struct {
	redisClient *redis.Client
}

// NewCache initializes and returns a new Cache instance.
// It accepts a context for potential future use (e.g., in ping) and adheres to Go conventions.
func NewCache(ctx context.Context, cfg *Config) (*Cache, error) {
	// Best Practice: Combine Host and Port into Addr in the config if possible,
	// but using fmt.Sprintf here is acceptable if the config is external.
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// NOTE: If using github.com/go-redis/redis/v7 as in the original code,
	// the `redis.NewClient` and `redis.Options` would be slightly different.
	// This example uses v8, which is generally preferred if not tied to v7.
	redisClient := redis.NewClient(&redis.Options{
		Addr:               addr,
		Password:           cfg.Password,
		DB:                 cfg.DB, // Use int DB
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		PoolSize:           cfg.PoolSize,
		PoolTimeout:        cfg.PoolTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		// Explicitly setting IdleTimeout to a default if not provided is good practice
		IdleTimeout: 5 * time.Minute, // A reasonable default
	})

	// Check the connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis at %s: %w", addr, err)
	}

	return &Cache{redisClient: redisClient}, nil
}

// Close closes the underlying Redis client connection.
// It's a best practice to provide a way to clean up resources.
func (c *Cache) Close() error {
	return c.redisClient.Close()
}

// Set stores a generic value 'T' associated with 'key' for a 'duration'.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	// Best Practice: Use interface{} instead of generics `[T any]` for Set/Get,
	// as JSON marshaling/unmarshaling works with interface{} and avoids complex type constraints.
	// NOTE: If you must use generics, the signature would be:
	// func (c *Cache) Set[T any](ctx context.Context, key string, value T, duration time.Duration) error { ... }

	v, err := json.Marshal(value)
	if err != nil {
		// Wrap error for better debugging
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	// Fix: Call the correct method on the redisClient, not a recursive call to c.Set
	// Best Practice: Use context.
	if err := c.redisClient.Set(ctx, key, v, duration).Err(); err != nil {
		return fmt.Errorf("redis SET command failed for key %s: %w", key, err)
	}
	return nil
}

// Get retrieves a value associated with 'key' and unmarshals it into 'dest'.
// 'dest' must be a pointer to the type you expect.
// Returns redis.Nil error if the key doesn't exist.
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	// NOTE: If you want a generic Get that returns the value, see the alternate below.
	// However, passing a pointer for unmarshaling is often simpler.

	// Best Practice: Use context.
	v, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		// Return the error, including redis.Nil for "key not found"
		if err == redis.Nil {
			return err // Key not found
		}
		return fmt.Errorf("redis GET command failed for key %s: %w", key, err)
	}

	err = json.Unmarshal([]byte(v), dest) // Unmarshal into the passed pointer 'dest'
	if err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Del deletes one or more keys from the cache.
// It returns an error if the underlying Redis operation fails, but it's not
// considered an error if the keys simply do not exist (Redis DEL semantics).
func (c *Cache) Del(ctx context.Context, keys ...string) error {
	// The DEL command returns the number of keys removed. It does not return an error
	// if a key does not exist. We only check for an underlying client/network error.
	if err := c.redisClient.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis DEL command failed: %w", err)
	}
	return nil
}
