package tests

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"check-in/api/internal/database"
)

type MainTestEnv struct {
	TestTx  pgx.Tx
	TestCtx context.Context
	TestDB  *pgxpool.Pool
}

func SetupGlobal(dbDsn string, dbMaxConns int,
	dbMaxIdletime string) (*MainTestEnv, error) {
	logger.SetLogger(logger.NullLogger)

	testDB, err := database.Connect(
		dbDsn,
		dbMaxConns,
		dbMaxIdletime,
	)
	if err != nil {
		return nil, err
	}

	testCtx := context.Background()
	testTx, err := testDB.Begin(testCtx)
	if err != nil {
		return nil, err
	}

	mainTestEnv := MainTestEnv{
		TestTx:  testTx,
		TestCtx: testCtx,
		TestDB:  testDB,
	}

	return &mainTestEnv, nil
}

func TeardownGlobal(mainTestEnv *MainTestEnv) error {
	err := mainTestEnv.TestTx.Rollback(mainTestEnv.TestCtx)
	if err != nil {
		return err
	}

	mainTestEnv.TestDB.Close()
	return nil
}
