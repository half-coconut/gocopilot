package repository

import (
	"context"
	"database/sql"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

//go:generate mockgen -source=user.go -package=mocks -destination=mocks/user.mock.go UserRepository
type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateByEmail(ctx context.Context, user domain.User) error
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache, l logger.LoggerV1) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
	l     logger.LoggerV1
}

func (u *CacheUserRepository) Create(ctx context.Context, user domain.User) error {
	defer func() {
		err := u.cache.Delete(ctx, user.Id)
		if err != nil {
			return
		}
	}()
	err := u.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Password: user.Password,
	})
	if err != nil {
		return err
	}
	return err
}
func (u *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 使用缓存后，先从缓存里取
	// 缓存里的 user 就是 domain.User
	// 如果缓存没报错，返回缓存的
	val, err := u.cache.Get(ctx, id)
	if err == nil {
		return val, err
	}
	// 如果缓存报错，再查库
	var user dao.User
	user, err = u.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	// 查库之后，set 进缓存
	us := u.entityToDomain(user)
	// 异步调用 userCache
	//go func() {
	//
	//}()
	err = u.cache.Set(ctx, us)
	if err != nil {
		return domain.User{}, err
	}

	return u.entityToDomain(user), err
}
func (u *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user dao.User
	user, err := u.dao.FindByEmail(ctx, email)
	return u.entityToDomain(user), err
}

func (u *CacheUserRepository) UpdateByEmail(ctx context.Context, user domain.User) error {
	defer func() {
		err := u.cache.Delete(ctx, user.Id)
		if err != nil {
			u.l.Info("User 缓存删除失败", logger.Error(err))
			return
		}
		u.l.Info("User 缓存删除成功")
	}()
	err := u.dao.UpdateByEmail(ctx, u.domainToEntity(user))
	return err
}

func (u *CacheUserRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,

		Phone:       user.Phone.String,
		FullName:    user.FullName.String,
		Department:  user.Department.String,
		Role:        user.Role.String,
		Avatar:      user.Avatar.String,
		Description: user.Description.String,

		Ctime: time.UnixMilli(user.Ctime),
		Utime: time.UnixMilli(user.Utime),
	}
}

func (u *CacheUserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		FullName: sql.NullString{
			String: user.FullName,
			Valid:  user.FullName != "",
		},
		Department: sql.NullString{
			String: user.Department,
			Valid:  user.Department != "",
		},
		Role: sql.NullString{
			String: user.Role,
			Valid:  user.Role != "",
		},
		Avatar: sql.NullString{
			String: user.Avatar,
			Valid:  user.Avatar != "",
		},

		Description: sql.NullString{
			String: user.Description,
			Valid:  user.Description != "",
		},
	}
}
