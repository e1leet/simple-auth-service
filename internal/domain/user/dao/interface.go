package dao

import "context"

type DAO interface {
	Create(ctx context.Context, storage *UserStorage) (int, error)
	GetById(ctx context.Context, id int) (*UserStorage, error)
	GetByUsername(ctx context.Context, username string) (*UserStorage, error)
	DeleteById(ctx context.Context, id int) error
}
