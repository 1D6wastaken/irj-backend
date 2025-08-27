package framework

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type (
	ctxKey int

	dbTracer struct{}

	withTx[T any] interface {
		WithTx(tx pgx.Tx) T
	}

	DB[T withTx[T]] struct {
		pool    *pgxpool.Pool
		Queries T
	}
)

const (
	ctxKeyStart ctxKey = iota
	ctxKeyStatement
	ctxKeyArgs
)

func (dbTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, ctxKeyStart, time.Now())
	ctx = context.WithValue(ctx, ctxKeyStatement, data.SQL)
	ctx = context.WithValue(ctx, ctxKeyArgs, data.Args)

	return ctx
}

func (dbTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	start, ok := ctx.Value(ctxKeyStart).(time.Time)
	if !ok {
		return
	}

	statement, ok := ctx.Value(ctxKeyStatement).(string)
	if !ok || statement == "SELECT VERSION()" {
		return
	}

	args, ok := ctx.Value(ctxKeyArgs).([]any)
	if !ok {
		return
	}

	zerolog.Ctx(ctx).Debug().Dict("sql", zerolog.Dict().
		Str("statement", statement).
		Interface("args", args).
		Int64("rowsAffected", data.CommandTag.RowsAffected()).
		Err(data.Err).
		Str("duration", time.Since(start).String()),
	).Msg("SQL result")
}

func NewDB[T withTx[T]](
	ctx context.Context,
	dsn string,
	queryBuilder func(db *pgxpool.Pool) T, afterConnect func(context.Context, *pgx.Conn) error) (*DB[T], error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.Tracer = dbTracer{}

	if afterConnect != nil {
		config.AfterConnect = afterConnect
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB[T]{
		pool:    pool,
		Queries: queryBuilder(pool),
	}, nil
}

func (s *DB[T]) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

//nolint:gocritic
func (s *DB[T]) begin(ctx context.Context, options pgx.TxOptions) (T, func(), func() error, error) {
	tx, err := s.pool.BeginTx(ctx, options)
	if err != nil {
		return s.Queries, nil, nil, err
	}

	return s.Queries.WithTx(tx),
		func() {
			_ = tx.Rollback(ctx)
		},
		func() error {
			return tx.Commit(ctx)
		},
		nil
}

func (s *DB[T]) Begin(ctx context.Context) (T, func(), func() error, error) {
	//nolint:exhaustruct
	return s.begin(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
}

func (s *DB[T]) BeginRead(ctx context.Context) (T, func(), func() error, error) {
	//nolint:exhaustruct
	return s.begin(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
}

func (s *DB[T]) Close() error {
	s.pool.Close()

	return nil
}
