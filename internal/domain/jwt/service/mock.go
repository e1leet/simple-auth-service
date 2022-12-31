package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func NewMock() *MockService {
	return &MockService{}
}

func (s *MockService) CreateAccessToken(userID int) (AccessToken, error) {
	args := s.Called(userID)
	return args.Get(0).(AccessToken), args.Error(1)
}

func (s *MockService) CreateRefreshToken(ctx context.Context, userID int) (RefreshToken, error) {
	args := s.Called(ctx, userID)
	return args.Get(0).(RefreshToken), args.Error(1)
}

func (s *MockService) RecreateRefreshToken(ctx context.Context, userID int, refresh RefreshToken) (RefreshToken, error) {
	args := s.Called(ctx, userID, refresh)
	return args.Get(0).(RefreshToken), args.Error(1)
}
