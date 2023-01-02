package service

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	sessionDAO "github.com/e1leet/simple-auth-service/internal/domain/session/dao"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	log.Logger = log.Output(io.Discard)
}

func TestService_CreateAccessToken(t *testing.T) {
	for _, userID := range []int{1, 2, 3, 4} {
		t.Run(fmt.Sprintf("create_access_token__user_id:%d", userID), func(t *testing.T) {
			const jwtSecret = "something"
			s := New(nil, jwtSecret, time.Hour, time.Hour)

			access, err := s.CreateAccessToken(userID)
			assert.NoError(t, err)

			token, err := jwt.Parse(string(access), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error in parsing")
				}
				return []byte(jwtSecret), nil
			})
			assert.NoError(t, err)

			claims, ok := token.Claims.(jwt.MapClaims)
			assert.True(t, ok && token.Valid)

			assert.Equal(t, float64(userID), claims["user_id"].(float64))
			assert.True(t, claims["exp"].(float64) > float64(time.Now().Unix()))
		})
	}
}

func TestService_CreateRefreshToken(t *testing.T) {
	t.Run("create_refresh_token", func(t *testing.T) {
		sessionMock := sessionDAO.NewMock()
		sessionMock.
			On("Create", mock.Anything, mock.Anything).
			Return(1, nil)

		s := New(sessionMock, "something", time.Hour, time.Hour)
		_, err := s.CreateRefreshToken(context.Background(), 1)
		assert.NoError(t, err)
	})
}

func TestService_RecreateRefreshToken(t *testing.T) {
	t.Run("session_not_found", func(t *testing.T) {
		sessionMock := sessionDAO.NewMock()
		sessionMock.
			On("GetByToken", mock.Anything, mock.Anything).
			Return(&sessionDAO.SessionStorage{}, sessionDAO.ErrSessionNotFound)

		s := New(sessionMock, "something", time.Hour, time.Hour)

		_, err := s.RecreateRefreshToken(context.Background(), 1, "wow")
		assert.ErrorIs(t, err, sessionDAO.ErrSessionNotFound)
	})

	t.Run("token_expired", func(t *testing.T) {
		sessionMock := sessionDAO.NewMock()
		sessionMock.
			On("GetByToken", mock.Anything, mock.Anything).
			Return(&sessionDAO.SessionStorage{ExpiresIn: time.Now().Unix() - 1000}, nil)

		s := New(sessionMock, "something", time.Hour, time.Hour)

		_, err := s.RecreateRefreshToken(context.Background(), 1, "wow")
		assert.ErrorIs(t, err, ErrRefreshTokenExpired)
	})

	t.Run("recreate_refresh_token", func(t *testing.T) {
		sessionMock := sessionDAO.NewMock()
		sessionMock.
			On("GetByToken", mock.Anything, mock.Anything).
			Return(&sessionDAO.SessionStorage{ExpiresIn: time.Now().Add(time.Hour).Unix()}, nil)
		sessionMock.
			On("Create", mock.Anything, mock.Anything).
			Return(1, nil)

		s := New(sessionMock, "something", time.Hour, time.Hour)

		_, err := s.RecreateRefreshToken(context.Background(), 1, "wow")
		assert.NoError(t, err)
	})
}
