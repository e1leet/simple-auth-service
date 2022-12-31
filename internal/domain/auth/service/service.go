package service

import (
	"context"
	"fmt"

	jwtService "github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	sessionDAO "github.com/e1leet/simple-auth-service/internal/domain/session/dao"
	"github.com/e1leet/simple-auth-service/internal/domain/user"
	userDAO "github.com/e1leet/simple-auth-service/internal/domain/user/dao"
	"github.com/e1leet/simple-auth-service/internal/utils/password/manager"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Register(ctx context.Context, u *user.User) error
	Login(ctx context.Context, u *user.User) (jwtService.AccessToken, jwtService.RefreshToken, error)
	UpdateAccessToken(
		ctx context.Context,
		userID int,
		refresh jwtService.RefreshToken,
	) (jwtService.AccessToken, jwtService.RefreshToken, error)
	Logout(ctx context.Context, refresh jwtService.RefreshToken) error
}

type service struct {
	jwtService      jwtService.Service
	userDAO         userDAO.DAO
	sessionDAO      sessionDAO.DAO
	passwordManager manager.Manager
	logger          zerolog.Logger
}

func New(jwt jwtService.Service, usr userDAO.DAO, session sessionDAO.DAO, password manager.Manager) Service {
	return &service{
		userDAO:         usr,
		sessionDAO:      session,
		passwordManager: password,
		jwtService:      jwt,
		logger:          log.With().Str("component", "authService").Logger(),
	}
}

func (s *service) Register(ctx context.Context, u *user.User) error {
	s.logger.Info().Str("username", u.Username).Msg("register user")

	storage := user.ToStorage(u)
	storage.Password = s.passwordManager.HashPassword(storage.Password)

	if _, err := s.userDAO.Create(ctx, storage); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	return nil
}

func (s *service) Login(ctx context.Context, u *user.User) (jwtService.AccessToken, jwtService.RefreshToken, error) {
	got, err := s.userDAO.GetByUsername(ctx, u.Username)
	if err != nil {
		switch err {
		case userDAO.ErrUserNotFound:
			s.logger.Warn().Str("username", u.Username).Err(userDAO.ErrUserNotFound).Send()
			return "", "", fmt.Errorf("failed to login: %w", err)
		default:
			s.logger.Err(err).Send()
			return "", "", err
		}
	}

	if !s.passwordManager.CheckPassword(u.Password, got.Password) {
		s.logger.Warn().Str("username", u.Username).Err(ErrIncorrectPassword).Send()
		return "", "", fmt.Errorf("failed to login: %w", ErrIncorrectPassword)
	}

	refresh, err := s.jwtService.CreateRefreshToken(ctx, got.ID)
	if err != nil {
		s.logger.Err(err).Send()
		return "", "", fmt.Errorf("failed to login: %w", err)
	}

	access, err := s.jwtService.CreateAccessToken(got.ID)
	if err != nil {
		s.logger.Err(err).Send()
		return "", "", fmt.Errorf("failed to login: %w", err)
	}

	return access, refresh, nil
}

func (s *service) UpdateAccessToken(
	ctx context.Context,
	userID int,
	refresh jwtService.RefreshToken,
) (jwtService.AccessToken, jwtService.RefreshToken, error) {
	newRefresh, err := s.jwtService.RecreateRefreshToken(ctx, userID, refresh)
	if err != nil {
		s.logger.Err(err).Send()
		return "", "", fmt.Errorf("failed to update access token: %w", err)
	}

	newAccess, err := s.jwtService.CreateAccessToken(userID)
	if err != nil {
		s.logger.Err(err).Send()
		return "", "", fmt.Errorf("failed to update access token: %w", err)
	}

	return newAccess, newRefresh, nil
}

func (s *service) Logout(ctx context.Context, refresh jwtService.RefreshToken) error {
	if err := s.sessionDAO.DeleteByToken(ctx, string(refresh)); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}
