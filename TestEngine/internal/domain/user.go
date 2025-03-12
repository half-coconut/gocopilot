package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Phone       string
	FullName    string
	Department  string
	Role        string
	Avatar      string
	Description string

	Ctime time.Time
	Utime time.Time
}
