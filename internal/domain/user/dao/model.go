package dao

import "time"

type UserStorage struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
}
