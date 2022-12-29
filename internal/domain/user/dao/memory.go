package dao

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type memoryDAO struct {
	users  map[int]*UserStorage
	mu     sync.Mutex
	logger zerolog.Logger
}

func NewMemory() DAO {
	return &memoryDAO{
		users:  make(map[int]*UserStorage),
		logger: log.With().Str("component", "userDAO").Logger(),
	}
}

func (s *memoryDAO) Create(ctx context.Context, u *UserStorage) (int, error) {
	s.logger.Info().Str("username", u.Username).Msg("create user")

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, usr := range s.users {
		if usr.Username == u.Username {
			s.logger.Warn().Str("username", u.Username).Err(ErrUsernameAlreadyUsed).Send()
			return 0, ErrUsernameAlreadyUsed
		}
	}

	u.CreatedAt = time.Now().UTC()
	u.ID = len(s.users) + 1
	s.users[u.ID] = u

	return u.ID, nil
}

func (s *memoryDAO) GetById(ctx context.Context, id int) (*UserStorage, error) {
	s.logger.Info().Int("id", id).Msg("get by id")

	s.mu.Lock()
	defer s.mu.Unlock()

	if usr, ok := s.users[id]; ok {
		return usr, nil
	}

	s.logger.Warn().Int("id", id).Err(ErrUserNotFound).Send()

	return nil, ErrUserNotFound
}

func (s *memoryDAO) GetByUsername(ctx context.Context, username string) (*UserStorage, error) {
	s.logger.Info().Str("username", username).Msg("get by username")

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, usr := range s.users {
		if usr.Username == username {
			return usr, nil
		}
	}

	s.logger.Warn().Str("username", username).Err(ErrUserNotFound).Send()

	return nil, ErrUserNotFound
}

func (s *memoryDAO) DeleteById(ctx context.Context, id int) error {
	s.logger.Info().Int("id", id).Msg("delete by id")

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, id)

	return nil
}
