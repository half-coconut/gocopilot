package gorm

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGorm(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/testcopilot"))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Phone       sql.NullString `gorm:"unique"`
	NickName    sql.NullString
	Department  sql.NullString
	Role        sql.NullString
	Description sql.NullString

	Ctime int64
	Utime int64
}

func TestCrypto(t *testing.T) {
	pwd := []byte("Cc12345!")
	// 加密
	encrypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	// 比较
	err = bcrypt.CompareHashAndPassword(encrypted, pwd)
	require.NoError(t, err)
}
