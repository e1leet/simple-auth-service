package service

import (
	"context"
	"fmt"
	"time"

	sessionDAO "github.com/e1leet/simple-auth-service/internal/domain/session/dao"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateAccessToken(userID int) (AccessToken, error)
	CreateRefreshToken(ctx context.Context, userID int) (RefreshToken, error)
	RecreateRefreshToken(ctx context.Context, userID int, refresh RefreshToken) (RefreshToken, error)
}

type service struct {
	sessionDAO       sessionDAO.DAO
	jwtSecret        string
	accessExpiresIn  time.Duration
	refreshExpiresIn time.Duration
	logger           zerolog.Logger
}

func New(dao sessionDAO.DAO, jwtSecret string, accessExpiresIn time.Duration, refreshExpiresIn time.Duration) Service {
	return &service{
		sessionDAO:       dao,
		jwtSecret:        jwtSecret,
		accessExpiresIn:  accessExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
		logger:           log.With().Str("component", "jwtService").Logger(),
	}
}

func (s *service) CreateAccessToken(userID int) (AccessToken, error) {
	s.logger.Info().Int("userID", userID).Msg("create access token")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.accessExpiresIn).Unix(),
	})

	access, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Err(err).Send()
		return "", fmt.Errorf("failed to create access token: %w", err)
	}

	return AccessToken(access), nil
}

func (s *service) CreateRefreshToken(ctx context.Context, userID int) (RefreshToken, error) {
	s.logger.Info().Int("userID", userID).Msg("create refresh token")

	session := &sessionDAO.SessionStorage{
		UserID:    userID,
		ExpiresIn: time.Now().Add(s.refreshExpiresIn).Unix(),
	}

	if _, err := s.sessionDAO.Create(ctx, session); err != nil {
		s.logger.Err(err).Send()
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return RefreshToken(session.Token), nil
}

func (s *service) RecreateRefreshToken(ctx context.Context, userID int, refresh RefreshToken) (RefreshToken, error) {
	s.logger.Info().Int("userID", userID).Msg("recreate refresh token")

	session, err := s.sessionDAO.GetByToken(ctx, string(refresh))
	if err != nil {
		s.logger.Err(err).Send()
		return "", fmt.Errorf("failed to recreate refresh token: %w", err)
	}

	if session.ExpiresIn < time.Now().Unix() {
		return "", fmt.Errorf("failed to recreate refresh token: %w", ErrRefreshTokenExpired)
	}

	newRefresh, err := s.CreateRefreshToken(ctx, userID)
	if err != nil {
		s.logger.Err(err).Send()
		return "", fmt.Errorf("failed to recreate access token: %w", err)
	}

	return newRefresh, nil
}
