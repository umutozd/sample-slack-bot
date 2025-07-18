package slackbot

import (
	"fmt"
	"time"
)

// Common middleware implementations

// LoggingMiddleware logs incoming events and their processing time
func LoggingMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			start := time.Now()
			
			// Log the incoming event
			fmt.Printf("[%s] Event: %s, Team: %s, User: %s, Channel: %s\n",
				start.Format(time.RFC3339),
				ctx.EventID,
				ctx.TeamID,
				ctx.UserID,
				ctx.ChannelID,
			)
			
			// Process the event
			err := next.Handle(ctx)
			
			// Log the completion
			duration := time.Since(start)
			status := "SUCCESS"
			if err != nil {
				status = "ERROR"
			}
			
			fmt.Printf("[%s] Event: %s completed in %v - Status: %s\n",
				time.Now().Format(time.RFC3339),
				ctx.EventID,
				duration,
				status,
			)
			
			if err != nil {
				fmt.Printf("[%s] Error: %v\n", time.Now().Format(time.RFC3339), err)
			}
			
			return err
		})
	}
}

// AuthMiddleware provides authentication/authorization for events
func AuthMiddleware(authFunc func(*Context) bool) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			if !authFunc(ctx) {
				return fmt.Errorf("unauthorized access denied for user %s", ctx.UserID)
			}
			return next.Handle(ctx)
		})
	}
}

// RateLimitMiddleware provides rate limiting based on user/team
func RateLimitMiddleware(limit int, window time.Duration) Middleware {
	// Simple in-memory rate limiter
	// In production, you'd want to use Redis or similar
	userCounts := make(map[string]int)
	userWindows := make(map[string]time.Time)
	
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			key := fmt.Sprintf("%s:%s", ctx.TeamID, ctx.UserID)
			now := time.Now()
			
			// Check if window has expired
			if windowStart, exists := userWindows[key]; exists {
				if now.Sub(windowStart) > window {
					// Reset window
					userCounts[key] = 0
					userWindows[key] = now
				}
			} else {
				// First request for this user
				userCounts[key] = 0
				userWindows[key] = now
			}
			
			// Check rate limit
			if userCounts[key] >= limit {
				return fmt.Errorf("rate limit exceeded for user %s", ctx.UserID)
			}
			
			// Increment count
			userCounts[key]++
			
			return next.Handle(ctx)
		})
	}
}

// RetryMiddleware provides automatic retries for failed handlers
func RetryMiddleware(maxRetries int) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			var lastErr error
			
			for attempt := 0; attempt <= maxRetries; attempt++ {
				err := next.Handle(ctx)
				if err == nil {
					return nil
				}
				
				lastErr = err
				
				// Don't retry on the last attempt
				if attempt < maxRetries {
					backoff := time.Duration(attempt+1) * time.Second
					fmt.Printf("[%s] Retry %d/%d failed, retrying in %v: %v\n",
						time.Now().Format(time.RFC3339),
						attempt+1,
						maxRetries,
						backoff,
						err,
					)
					time.Sleep(backoff)
				}
			}
			
			return fmt.Errorf("handler failed after %d retries: %w", maxRetries, lastErr)
		})
	}
}

// MetricsMiddleware provides basic metrics collection
func MetricsMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			start := time.Now()
			
			err := next.Handle(ctx)
			
			duration := time.Since(start)
			
			// In a real implementation, you'd send these metrics to a monitoring system
			fmt.Printf("[METRICS] Event: %s, Duration: %v, Success: %t\n",
				ctx.EventID,
				duration,
				err == nil,
			)
			
			return err
		})
	}
}

// ErrorRecoveryMiddleware recovers from panics and logs them
func ErrorRecoveryMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[%s] PANIC recovered: %v\n", time.Now().Format(time.RFC3339), r)
				}
			}()
			
			return next.Handle(ctx)
		})
	}
}

// CachingMiddleware provides response caching
func CachingMiddleware(ttl time.Duration) Middleware {
	cache := make(map[string]interface{})
	cacheTimes := make(map[string]time.Time)
	
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			// Create cache key
			key := fmt.Sprintf("%s:%s:%s", ctx.TeamID, ctx.UserID, ctx.EventID)
			
			// Check if we have a cached response
			if cachedTime, exists := cacheTimes[key]; exists {
				if time.Since(cachedTime) < ttl {
					fmt.Printf("[%s] Cache hit for key: %s\n", time.Now().Format(time.RFC3339), key)
					return nil // Return cached response
				}
			}
			
			// Process the request
			err := next.Handle(ctx)
			
			// Cache the response if successful
			if err == nil {
				cache[key] = true // In a real implementation, you'd cache the actual response
				cacheTimes[key] = time.Now()
			}
			
			return err
		})
	}
}

// ValidationMiddleware validates incoming events
func ValidationMiddleware(validator func(*Context) error) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			// Validate the context
			if err := validator(ctx); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			
			return next.Handle(ctx)
		})
	}
}

// RequestIDMiddleware adds a unique request ID to each event
func RequestIDMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			// Generate a simple request ID
			ctx.RequestID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), ctx.EventID)
			
			return next.Handle(ctx)
		})
	}
}

// TimeoutMiddleware adds a timeout to handler execution
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			done := make(chan error, 1)
			
			go func() {
				done <- next.Handle(ctx)
			}()
			
			select {
			case err := <-done:
				return err
			case <-time.After(timeout):
				return fmt.Errorf("handler timed out after %v", timeout)
			}
		})
	}
}

// Chain chains multiple middleware together
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Common middleware combinations

// DefaultMiddlewareStack provides a sensible default middleware stack
func DefaultMiddlewareStack() []Middleware {
	return []Middleware{
		ErrorRecoveryMiddleware(),
		RequestIDMiddleware(),
		LoggingMiddleware(),
		TimeoutMiddleware(30 * time.Second),
		MetricsMiddleware(),
	}
}

// ProductionMiddlewareStack provides a production-ready middleware stack
func ProductionMiddlewareStack() []Middleware {
	return []Middleware{
		ErrorRecoveryMiddleware(),
		RequestIDMiddleware(),
		LoggingMiddleware(),
		RateLimitMiddleware(60, time.Minute), // 60 requests per minute
		TimeoutMiddleware(30 * time.Second),
		RetryMiddleware(3),
		MetricsMiddleware(),
	}
}

// DevelopmentMiddlewareStack provides a development-friendly middleware stack
func DevelopmentMiddlewareStack() []Middleware {
	return []Middleware{
		ErrorRecoveryMiddleware(),
		RequestIDMiddleware(),
		LoggingMiddleware(),
		TimeoutMiddleware(60 * time.Second), // Longer timeout for debugging
		MetricsMiddleware(),
	}
}