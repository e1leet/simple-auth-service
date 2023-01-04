package dao

import (
	"context"
	"errors"
	"strings"

	"github.com/e1leet/simple-auth-service/pkg/client/postgresql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type postgresqlDAO struct {
	client postgresql.Client
	logger zerolog.Logger
}

func NewPostgresql(client postgresql.Client) DAO {
	return &postgresqlDAO{
		client: client,
		logger: log.With().Str("component", "userDAO").Logger(),
	}
}

func (s *postgresqlDAO) Create(ctx context.Context, u *UserStorage) (int, error) {
	q := `
		INSERT INTO public.usr 
		    (username, password) 
		VALUES 
		    ($1, $2)
		RETURNING id, created_at
	`
	s.logger.Trace().Str("query", postgresql.FormatQuery(q)).Send()

	if err := s.client.QueryRow(ctx, q, u.Username, u.Password).Scan(&u.ID, &u.CreatedAt); err != nil {
		s.logger.Error().Err(err).Send()

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			strings.Contains(pgErr.Message, "usr_username_key") &&
			pgErr.SQLState() == "23505" {
			return 0, ErrUsernameAlreadyUsed
		}

		return 0, err
	}

	return u.ID, nil
}

func (s *postgresqlDAO) GetById(ctx context.Context, id int) (*UserStorage, error) {
	q := `
		SELECT id, username, password, created_at 
		FROM public.usr 
		WHERE usr.id = $1
	`
	s.logger.Trace().Str("query", postgresql.FormatQuery(q)).Send()

	usr := &UserStorage{}

	if err := s.client.QueryRow(ctx, q, id).Scan(
		&usr.ID,
		&usr.Username,
		&usr.Password,
		&usr.CreatedAt,
	); err != nil {
		s.logger.Error().Err(err).Send()
		return nil, err
	}

	return usr, nil
}

func (s *postgresqlDAO) GetByUsername(ctx context.Context, username string) (*UserStorage, error) {
	q := `
		SELECT id, username, password, created_at 
		FROM public.usr 
		WHERE usr.username = $1
	`
	s.logger.Trace().Str("query", postgresql.FormatQuery(q)).Send()

	usr := &UserStorage{}

	if err := s.client.QueryRow(ctx, q, username).Scan(
		&usr.ID,
		&usr.Username,
		&usr.Password, &usr.CreatedAt,
	); err != nil {
		s.logger.Error().Err(err).Send()
		return nil, err
	}

	return usr, nil
}

func (s *postgresqlDAO) DeleteById(ctx context.Context, id int) error {
	q := `
		DELETE FROM public.usr
		WHERE id = $1
		RETURNING id
	`
	s.logger.Trace().Str("query", postgresql.FormatQuery(q)).Send()

	err := s.client.QueryRow(ctx, q).Scan()
	if err != nil {
		s.logger.Error().Err(err).Send()
		return err
	}

	return nil
}
