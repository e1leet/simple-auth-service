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

func (s *MockDAO) Create(ctx context.Context, session *SessionStorage) (int, error) {
	args := s.Called(ctx, session)
	return args.Int(0), args.Error(1)
}

func (s *MockDAO) GetByID(ctx context.Context, id int) (*SessionStorage, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*SessionStorage), args.Error(1)
}

func (s *MockDAO) GetByToken(ctx context.Context, token string) (*SessionStorage, error) {
	args := s.Called(ctx, token)
	return args.Get(0).(*SessionStorage), args.Error(1)
}

func (s *MockDAO) DeleteByID(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func (s *MockDAO) DeleteByToken(ctx context.Context, token string) error {
	args := s.Called(ctx, token)
	return args.Error(0)
}
