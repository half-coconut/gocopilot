package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱或者手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateByEmail(ctx context.Context, user User) error
}

type GORMUserDAO struct {
	// 这里是 gorm的db 实现的结构体
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 新增一条 User 数据
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictCodeNo = 1062
		if mysqlErr.Number == uniqueConflictCodeNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	return user, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&user).Error
	return user, err
}

func (dao *GORMUserDAO) UpdateByEmail(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	err := dao.db.Model(&user).WithContext(ctx).Where("email=?", user.Email).Updates(map[string]interface{}{
		"FullName":    user.FullName,
		"Department":  user.Department,
		"Phone":       user.Phone,
		"Role":        user.Role,
		"Description": user.Description,
		"Utime":       now,
	}).Error
	return err
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Phone       sql.NullString `gorm:"unique"`
	FullName    sql.NullString
	Department  sql.NullString
	Role        sql.NullString
	Avatar      sql.NullString
	Description sql.NullString

	Ctime int64
	Utime int64
}
