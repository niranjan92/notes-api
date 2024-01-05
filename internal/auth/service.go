package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
	Signup(ctx context.Context, username, password string) (string, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
}

type service struct {
	signingKey      string
	tokenExpiration int
	logger          log.Logger
	uRepo           UserRepo
}

// NewService creates a new authentication service.
func NewService(userRepo UserRepo, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{signingKey, tokenExpiration, logger, userRepo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if identity := s.authenticate(ctx, username, password); identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized("")
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) Identity {
	logger := s.logger.With(ctx, "user", username)

	dbUser, err := s.uRepo.GetByName(ctx, username)
	if err != nil {
		logger.Infof("user not found: %v", err)
		return nil
	}

	if dbUser.Name == username && dbUser.Password == password {
		// TODO: salt, hash then compare password
		logger.Debugf("authentication successful")
		return dbUser
	}
	// if username == "demo" && password == "pass" {
	// 	logger.Infof("authentication successful")
	// 	return entity.User{ID: "100", Name: "demo"}
	// }

	logger.Infof("authentication failed")
	return nil
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   identity.GetID(),
		"name": identity.GetName(),
		"exp":  time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}

func (s service) Signup(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("username and password cannot be empty")
	}

	dbUser, err := s.uRepo.GetByName(ctx, username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", fmt.Errorf("err getting user %w", err)
	}
	if dbUser.Name == username {
		return "", fmt.Errorf("user already exists")
	}
	id := entity.GenerateID()

	newUser := entity.User{
		ID:        id,
		Name:      username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.uRepo.Create(ctx, newUser)
	if err != nil {
		return "", fmt.Errorf("error creating user")
	}
	token, err := s.generateJWT(newUser)
	if err != nil {
		s.logger.Errorf("error generating token: %v", err)
		return "", fmt.Errorf("error generating token")
	}

	return token, nil
}
