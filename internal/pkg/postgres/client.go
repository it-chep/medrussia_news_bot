package postgres

import (
	"context"
	"fmt"
	"medrussia_news_bot/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type PgxClientWrapper struct {
	pool *pgxpool.Pool
}

func (c *PgxClientWrapper) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, arguments...)
}

func (c *PgxClientWrapper) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

func (c *PgxClientWrapper) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

func (c *PgxClientWrapper) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

func (c *PgxClientWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, txOptions)
}

func NewClient(ctx context.Context, pg config.StorageConfig) (client Client, err error) {
	const op = "sorkin_bot.pkg.client.postgres.NewClient"
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?pool_max_conns=%d", pg.User, pg.Password, pg.Host, pg.Port, pg.Database, pg.MaxConnects)

	poolConfig, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &PgxClientWrapper{pool: pool}, nil
}

func ConnectWithRetry(fn func() error, maxRetry int, timeout time.Duration) (err error) {
	for maxRetry > 0 {
		if err = fn(); err != nil {
			time.Sleep(timeout)
			maxRetry--
			continue
		}
		return nil
	}
	return
}
