package middleware

import (
	"log"
	"time"
)

// Middleware defines the middleware function signature
type Middleware func(next Handler) Handler

// Handler defines the handler function signature
type Handler func() error

// LoggingMiddleware logs all requests
func LoggingMiddleware() Middleware {
	return func(next Handler) Handler {
		return func() error {
			start := time.Now()
			
			log.Printf("🔄 Request started")
			
			err := next()
			
			duration := time.Since(start)
			
			if err != nil {
				log.Printf("❌ Request failed after %v: %v", duration, err)
			} else {
				log.Printf("✅ Request completed in %v", duration)
			}
			
			return err
		}
	}
}

// AuthMiddleware provides authentication
func AuthMiddleware(authFunc func() bool) Middleware {
	return func(next Handler) Handler {
		return func() error {
			if !authFunc() {
				log.Printf("🔒 Authentication failed")
				return nil // Don't proceed with the request
			}
			
			return next()
		}
	}
}

// RateLimitMiddleware provides rate limiting
func RateLimitMiddleware(limit int, window time.Duration) Middleware {
	// This is a simple implementation. In production, you'd use a more sophisticated rate limiter
	requests := make(map[string][]time.Time)
	
	return func(next Handler) Handler {
		return func() error {
			now := time.Now()
			key := "global" // In a real implementation, this would be per-user/team
			
			// Clean old requests
			userRequests := requests[key]
			var validRequests []time.Time
			for _, reqTime := range userRequests {
				if now.Sub(reqTime) < window {
					validRequests = append(validRequests, reqTime)
				}
			}
			
			// Check rate limit
			if len(validRequests) >= limit {
				log.Printf("⚠️ Rate limit exceeded for %s", key)
				return nil // Don't proceed with the request
			}
			
			// Add current request
			validRequests = append(validRequests, now)
			requests[key] = validRequests
			
			return next()
		}
	}
}

// RetryMiddleware provides retry logic
func RetryMiddleware(maxRetries int) Middleware {
	return func(next Handler) Handler {
		return func() error {
			var lastErr error
			
			for i := 0; i <= maxRetries; i++ {
				err := next()
				if err == nil {
					return nil
				}
				
				lastErr = err
				
				if i < maxRetries {
					backoff := time.Duration(i+1) * time.Second
					log.Printf("🔄 Retry %d/%d after %v: %v", i+1, maxRetries, backoff, err)
					time.Sleep(backoff)
				}
			}
			
			log.Printf("❌ All retries failed: %v", lastErr)
			return lastErr
		}
	}
}

// ErrorRecoveryMiddleware recovers from panics
func ErrorRecoveryMiddleware() Middleware {
	return func(next Handler) Handler {
		return func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("🚨 Panic recovered: %v", r)
					err = nil // Don't propagate panic as error
				}
			}()
			
			return next()
		}
	}
}

// MetricsMiddleware collects metrics
func MetricsMiddleware() Middleware {
	var (
		totalRequests int
		successCount  int
		errorCount    int
		totalTime     time.Duration
	)
	
	return func(next Handler) Handler {
		return func() error {
			start := time.Now()
			totalRequests++
			
			err := next()
			
			duration := time.Since(start)
			totalTime += duration
			
			if err != nil {
				errorCount++
			} else {
				successCount++
			}
			
			// Log metrics every 100 requests
			if totalRequests%100 == 0 {
				avgTime := totalTime / time.Duration(totalRequests)
				successRate := float64(successCount) / float64(totalRequests) * 100
				
				log.Printf("📊 Metrics: %d requests, %.1f%% success rate, avg %v", 
					totalRequests, successRate, avgTime)
			}
			
			return err
		}
	}
}

// CachingMiddleware provides caching capabilities
func CachingMiddleware(cache map[string]interface{}) Middleware {
	return func(next Handler) Handler {
		return func() error {
			// This is a simplified caching middleware
			// In practice, you'd need a way to generate cache keys and manage cache invalidation
			return next()
		}
	}
}

// ValidationMiddleware validates requests
func ValidationMiddleware(validator func() error) Middleware {
	return func(next Handler) Handler {
		return func() error {
			if err := validator(); err != nil {
				log.Printf("❌ Validation failed: %v", err)
				return err
			}
			
			return next()
		}
	}
}

// TimeoutMiddleware provides timeout functionality
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Handler) Handler {
		return func() error {
			done := make(chan error, 1)
			
			go func() {
				done <- next()
			}()
			
			select {
			case err := <-done:
				return err
			case <-time.After(timeout):
				log.Printf("⏰ Request timeout after %v", timeout)
				return nil // Don't propagate timeout as error
			}
		}
	}
}

// Chain combines multiple middleware
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply applies middleware to a handler
func Apply(handler Handler, middlewares ...Middleware) Handler {
	return Chain(middlewares...)(handler)
}