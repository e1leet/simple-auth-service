package dao

import "context"

type DAO interface {
	Create(ctx context.Context, session *SessionStorage) (int, error)
	GetByID(ctx context.Context, id int) (*SessionStorage, error)
	GetByToken(ctx context.Context, token string) (*SessionStorage, error)
	DeleteByID(ctx context.Context, id int) error
}
