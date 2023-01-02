package dao

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDAO struct {
	mock.Mock
}

func NewMock() *MockDAO {
	return &MockDAO{}
}

func (s *MockDAO) Create(ctx context.Context, u *UserStorage) (int, error) {
	args := s.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (s *MockDAO) GetById(ctx context.Context, id int) (*UserStorage, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*UserStorage), args.Error(1)
}

func (s *MockDAO) GetByUsername(ctx context.Context, username string) (*UserStorage, error) {
	args := s.Called(ctx, username)
	return args.Get(0).(*UserStorage), args.Error(1)
}

func (s *MockDAO) DeleteById(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
