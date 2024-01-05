package auth

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"golang.org/x/time/rate"
)

type RL struct {
	limiter *rate.Limiter
}

var (
	clients = make(map[string]*RL)
	mtx     sync.Mutex
)

// RateLimiter does per user rate limiting. use after auth handler
func RateLimiter() routing.Handler {
	return func(c *routing.Context) error {
		userID := c.Get("user_id").(string)
		mtx.Lock()
		defer mtx.Unlock()
		limiter, exists := clients[userID]
		if !exists {
			limiter = &RL{
				limiter: rate.NewLimiter(rate.Every(time.Minute), 10), // change this to per second
			}
			clients[userID] = limiter
		}
		if !limiter.limiter.Allow() {
			c.Abort()
			c.Response.WriteHeader(http.StatusTooManyRequests)
			c.Response.Write([]byte("Too many requests"))
			return nil
		}

		return nil
	}
}

// Handler returns a JWT-based authentication middleware.
func Handler(verificationKey string) routing.Handler {
	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: handleToken})
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func handleToken(c *routing.Context, token *jwt.Token) error {
	ctx := WithUser(
		c.Request.Context(),
		token.Claims.(jwt.MapClaims)["id"].(string),
		token.Claims.(jwt.MapClaims)["name"].(string),
	)
	c.Set("username", token.Claims.(jwt.MapClaims)["name"].(string))
	c.Set("user_id", token.Claims.(jwt.MapClaims)["id"].(string))
	c.Request = c.Request.WithContext(ctx)
	return nil
}

type contextKey int

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, id, name string) context.Context {
	return context.WithValue(ctx, userKey, entity.User{ID: id, Name: name})
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentUser(ctx context.Context) Identity {
	if user, ok := ctx.Value(userKey).(entity.User); ok {
		return user
	}
	return nil
}

// MockAuthHandler creates a mock authentication middleware for testing purpose.
// If the request contains an Authorization header whose value is "TEST", then
// it considers the user is authenticated as "Tester" whose ID is "100".
// It fails the authentication otherwise.
func MockAuthHandler(c *routing.Context) error {
	if c.Request.Header.Get("Authorization") != "TEST" {
		return errors.Unauthorized("")
	}
	ctx := WithUser(c.Request.Context(), "testuser", "Tester")
	c.Set("user_id", "testuser")

	c.Request = c.Request.WithContext(ctx)
	return nil
}

// MockAuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}
