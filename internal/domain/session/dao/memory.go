package dao

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type memoryDAO struct {
	sessions map[int]*SessionStorage
	mu       sync.Mutex
	logger   zerolog.Logger
}

func NewMemory() DAO {
	return &memoryDAO{
		sessions: make(map[int]*SessionStorage),
		logger:   log.With().Str("component", "sessionDAO").Logger(),
	}
}

func (s *memoryDAO) Create(ctx context.Context, session *SessionStorage) (int, error) {
	s.logger.Info().Int("userID", session.UserID).Msg("create session")

	s.mu.Lock()
	defer s.mu.Unlock()

	session.ID = len(s.sessions) + 1
	session.Token = uuid.New().String()
	session.CreatedAt = time.Now().UTC()

	s.sessions[session.ID] = session

	return session.ID, nil
}

func (s *memoryDAO) GetByID(ctx context.Context, id int) (*SessionStorage, error) {
	s.logger.Info().Int("id", id).Msg("get by id")

	s.mu.Lock()
	defer s.mu.Unlock()

	if session, ok := s.sessions[id]; ok {
		return session, nil
	}

	return nil, ErrSessionNotFound
}

func (s *memoryDAO) GetByToken(ctx context.Context, token string) (*SessionStorage, error) {
	s.logger.Info().Msg("get by token")

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.Token == token {
			return session, nil
		}
	}

	return nil, ErrSessionNotFound
}

func (s *memoryDAO) DeleteByID(ctx context.Context, id int) error {
	s.logger.Info().Int("id", id).Msg("delete by id")

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)

	return nil
}

func (s *memoryDAO) DeleteByToken(ctx context.Context, token string) error {
	s.logger.Info().Msg("delete by token")

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.Token == token {
			delete(s.sessions, session.ID)
			return nil
		}
	}

	return nil
}
