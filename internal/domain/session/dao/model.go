package dao

import "time"

type SessionStorage struct {
	ID        int
	Token     string
	CreatedAt time.Time
	ExpiresIn int
	UserID    int
}
