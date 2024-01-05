package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/qiangxue/go-rest-api/internal/entity"
	errs "github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func Test_service_Authenticate(t *testing.T) {
	logger, _ := log.NewForTest()

	s := NewService(&mockRepository{}, "test", 100, logger)
	_, err := s.Login(context.Background(), "unknown", "bad")
	assert.Equal(t, errs.Unauthorized(""), err)

	token, err := s.Signup(context.Background(), "demo", "pass")
	_ = token
	assert.Nil(t, err)
	token, err = s.Login(context.Background(), "demo", "pass")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func Test_service_authenticate(t *testing.T) {
	logger, _ := log.NewForTest()
	s := service{"test", 100, logger, &mockRepository{}}
	assert.Nil(t, s.authenticate(context.Background(), "unknown", "bad"))

	_, err := s.Signup(context.Background(), "demo", "pass")
	assert.Nil(t, err)
	assert.NotNil(t, s.authenticate(context.Background(), "demo", "pass"))
}

func Test_service_GenerateJWT(t *testing.T) {
	logger, _ := log.NewForTest()
	s := service{"test", 100, logger, &mockRepository{}}
	token, err := s.generateJWT(entity.User{
		ID:   "100",
		Name: "demo",
	})
	if assert.Nil(t, err) {
		assert.NotEmpty(t, token)
	}
}

type mockRepository struct {
	items []entity.User
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.User, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.User{}, sql.ErrNoRows
}

func (m mockRepository) GetByName(ctx context.Context, name string) (entity.User, error) {
	for _, item := range m.items {
		if item.Name == name {
			return item, nil
		}
	}
	return entity.User{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.User, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, user entity.User) error {
	if user.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, user)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, user entity.User) error {
	if user.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == user.ID {
			m.items[i] = user
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
