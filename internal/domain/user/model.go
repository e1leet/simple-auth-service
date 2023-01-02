package user

import (
	"time"

	"github.com/e1leet/simple-auth-service/internal/domain/user/dao"
)

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
}

func ToStorage(user *User) *dao.UserStorage {
	return &dao.UserStorage{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}
