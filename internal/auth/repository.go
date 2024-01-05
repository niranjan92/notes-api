package auth

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

type UserRepo interface {
	Get(ctx context.Context, name string) (entity.User, error)
	GetByName(ctx context.Context, name string) (entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) UserRepo {
	return repository{db, logger}
}

// Get function actually gets by name
func (r repository) Get(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

func (r repository) GetByName(ctx context.Context, name string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"name": name}).One(&user)
	return user, err
}

func (r repository) Create(ctx context.Context, user entity.User) error {
	return r.db.With(ctx).Model(&user).Insert()
}

func (r repository) Update(ctx context.Context, user entity.User) error {
	return r.db.With(ctx).Model(&user).Update()
}

func (r repository) Delete(ctx context.Context, id string) error {
	item, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&item).Delete()
}
