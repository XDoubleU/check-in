package tests

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TestTx struct {
	tx pgx.Tx
	mu *sync.Mutex
}

func waitOnLock(lock *sync.Mutex) {
	for {
		if lock.TryLock() {
			break
		}
	}
}

func (tx TestTx) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Exec(ctx, sql, arguments...)
}

func (tx TestTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	//nolint:sqlclosecheck // user is supposed to close query
	return tx.tx.Query(ctx, sql, args...)
}

func (tx TestTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.QueryRow(ctx, sql, args...)
}

func (tx TestTx) Begin(ctx context.Context) (pgx.Tx, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Begin(ctx)
}
