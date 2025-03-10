package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Phone       string
	NickName    string
	Department  string
	Role        string
	Description string

	Ctime time.Time
	Utime time.Time
}
