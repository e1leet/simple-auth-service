package service

import (
	"context"
	"testing"

	jwtService "github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	sessionDAO "github.com/e1leet/simple-auth-service/internal/domain/session/dao"
	"github.com/e1leet/simple-auth-service/internal/domain/user"
	userDAO "github.com/e1leet/simple-auth-service/internal/domain/user/dao"
	"github.com/e1leet/simple-auth-service/internal/utils/password/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		usrMock := userDAO.NewMock()
		usrMock.On("Create", mock.Anything, mock.Anything).Return(1, nil)

		passwordMock := manager.NewMock()
		passwordMock.
			On("HashPassword", mock.Anything).
			Return("123")

		s := New(nil, usrMock, nil, passwordMock)

		err := s.Register(context.Background(), &user.User{
			Username: "test",
			Password: "test",
		})
		assert.NoError(t, err)
	})

	t.Run("user_already_exists", func(t *testing.T) {
		usrMock := userDAO.NewMock()
		usrMock.
			On("Create", mock.Anything, mock.Anything).
			Return(0, userDAO.ErrUsernameAlreadyUsed)

		passwordMock := manager.NewMock()
		passwordMock.
			On("HashPassword", mock.Anything).
			Return("123")

		s := New(nil, usrMock, nil, passwordMock)

		err := s.Register(context.Background(), &user.User{
			Username: "test",
			Password: "test",
		})
		assert.ErrorIs(t, err, userDAO.ErrUsernameAlreadyUsed)
	})
}

func TestService_Login(t *testing.T) {
	t.Run("user_doesnt_exist", func(t *testing.T) {
		usrMock := userDAO.NewMock()
		usrMock.
			On("GetByUsername", mock.Anything, mock.Anything).
			Return(&userDAO.UserStorage{}, userDAO.ErrUserNotFound)

		s := New(nil, usrMock, nil, nil)

		_, _, err := s.Login(context.Background(), &user.User{Username: "test", Password: "test"})
		assert.Error(t, err, userDAO.ErrUserNotFound)
	})

	t.Run("incorrect_password", func(t *testing.T) {
		usrMock := userDAO.NewMock()
		usrMock.
			On("GetByUsername", mock.Anything, mock.Anything).
			Return(&userDAO.UserStorage{}, userDAO.ErrUserNotFound)

		passwordMock := manager.NewMock()
		passwordMock.
			On("CheckPassword", mock.Anything, mock.Anything).
			Return(false)

		s := New(nil, usrMock, nil, passwordMock)

		_, _, err := s.Login(context.Background(), &user.User{Username: "test", Password: "test"})
		assert.Error(t, err, userDAO.ErrUserNotFound)
	})

	t.Run("login", func(t *testing.T) {
		usrMock := userDAO.NewMock()
		usrMock.
			On("GetByUsername", mock.Anything, mock.Anything).
			Return(&userDAO.UserStorage{}, nil)

		passwordMock := manager.NewMock()
		passwordMock.On("CheckPassword", mock.Anything, mock.Anything).Return(true)

		jwtMock := jwtService.NewMock()
		expectedRefresh := jwtService.RefreshToken("wow")
		jwtMock.
			On("CreateRefreshToken", mock.Anything, mock.Anything).
			Return(expectedRefresh, nil)

		expectedAccess := jwtService.AccessToken("wow")
		jwtMock.
			On("CreateAccessToken", mock.Anything).
			Return(expectedAccess, nil)

		s := New(jwtMock, usrMock, nil, passwordMock)

		access, refresh, err := s.Login(context.Background(), &user.User{Username: "test", Password: "test"})
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, expectedAccess, access)
		assert.Equal(t, expectedRefresh, refresh)

	})
}

func TestService_UpdateAccessToken(t *testing.T) {
	t.Run("refresh_token_not_found", func(t *testing.T) {
		jwtMock := jwtService.NewMock()
		jwtMock.
			On("RecreateRefreshToken", mock.Anything, mock.Anything, mock.Anything).
			Return(jwtService.RefreshToken(""), sessionDAO.ErrSessionNotFound)
		s := New(jwtMock, nil, nil, nil)

		_, _, err := s.UpdateAccessToken(context.Background(), 123, "123")
		assert.ErrorIs(t, err, sessionDAO.ErrSessionNotFound)
	})

	t.Run("refresh_token_expired", func(t *testing.T) {
		jwtMock := jwtService.NewMock()
		jwtMock.
			On("RecreateRefreshToken", mock.Anything, mock.Anything, mock.Anything).
			Return(jwtService.RefreshToken(""), jwtService.ErrRefreshTokenExpired)

		s := New(jwtMock, nil, nil, nil)

		_, _, err := s.UpdateAccessToken(context.Background(), 123, "123")
		assert.ErrorIs(t, err, jwtService.ErrRefreshTokenExpired)
	})

	t.Run("update_access_token", func(t *testing.T) {
		jwtMock := jwtService.NewMock()
		expectedRefresh := jwtService.RefreshToken("test")
		jwtMock.
			On("RecreateRefreshToken", mock.Anything, mock.Anything, mock.Anything).
			Return(expectedRefresh, nil)

		expectedAccess := jwtService.AccessToken("wow")
		jwtMock.
			On("CreateAccessToken", mock.Anything).
			Return(expectedAccess, nil)

		s := New(jwtMock, nil, nil, nil)

		access, refresh, err := s.UpdateAccessToken(context.Background(), 123, "123")
		assert.NoError(t, err)
		assert.Equal(t, expectedAccess, access)
		assert.Equal(t, expectedRefresh, refresh)
	})
}

func TestService_Logout(t *testing.T) {
	t.Run("logout", func(t *testing.T) {
		sessionMock := sessionDAO.NewMock()
		sessionMock.On("DeleteByToken", mock.Anything, mock.Anything).Return(nil)

		s := New(nil, nil, sessionMock, nil)

		err := s.Logout(context.Background(), "test")
		assert.NoError(t, err)
	})
}
