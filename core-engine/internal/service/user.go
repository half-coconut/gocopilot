package service

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("邮箱/用户或者密码不正确")
)

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	l    logger.LoggerV1
}

func NewUserService(repo repository.UserRepository, l logger.LoggerV1) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}

func (u *userService) Signup(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = u.repo.Create(ctx, user)
	if err != nil {
		return err
	}
	return err
}

func (u *userService) Login(ctx context.Context, email, password string) (domain.User, error) {

	var user domain.User
	user, err := u.repo.FindByEmail(ctx, email)
	// 校验 email 是否存在
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	// 返回系统错误
	if err != nil {
		return domain.User{}, err
	}
	// 校验 password 是否匹配
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, err
}

func (u *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	// 更新非敏感信息
	err := u.repo.UpdateByEmail(ctx, user)
	return err
}

func (u *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := u.repo.FindById(ctx, id)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, err
}
