package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, uri string, maxAttempts int) (pool *pgxpool.Pool, err error) {
	logger := log.With().Logger()
	logger.Info().Msg("connect to postgres")

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, uri)
		if err != nil {
			logger.Error().Err(err).Send()
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		logger.Fatal().Err(err).Msg("error do with tries")
	}

	return pool, nil
}
