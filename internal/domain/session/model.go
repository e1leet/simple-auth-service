package session

import "time"

type Session struct {
	ID        int
	Token     string
	CreatedAt time.Time
	ExpiresIn int
	UserID    int
}
