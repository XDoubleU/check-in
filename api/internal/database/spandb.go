package database

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SpanDB struct {
	DB DB
}

func (spandb SpanDB) Exec(ctx context.Context, sql string,
	arguments ...any) (pgconn.CommandTag, error) {
	span := startSpan(ctx, sql)
	defer span.Finish()

	return spandb.DB.Exec(span.Context(), sql, arguments...)
}

func (spandb SpanDB) Query(ctx context.Context, sql string,
	args ...any) (pgx.Rows, error) {
	span := startSpan(ctx, sql)
	defer span.Finish()

	//nolint:sqlclosecheck // user is supposed to close query
	return spandb.DB.Query(span.Context(), sql, args...)
}

func (spandb SpanDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	span := startSpan(ctx, sql)
	defer span.Finish()

	return spandb.DB.QueryRow(span.Context(), sql, args...)
}

func (spandb SpanDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return spandb.DB.Begin(ctx)
}

func startSpan(ctx context.Context, sql string) *sentry.Span {
	span := sentry.StartSpan(ctx, "db.query", sentry.WithDescription(sql))
	span.SetData("db.system", "postgresql")

	return span
}
