package tests

import (
	"context"
	"sync"
)

type TestEnv struct {
	TestTx  TestTx
	TestCtx context.Context
}

func SetupSingle(mainTestEnv *MainTestEnv) TestEnv {
	testCtx := context.Background()

	var testTx TestTx
	var mu sync.Mutex
	for {
		tx, err := mainTestEnv.TestTx.Begin(testCtx)
		if err == nil {
			testTx = TestTx{
				tx: tx,
				mu: &mu,
			}
			break
		}
	}

	testEnv := TestEnv{
		TestTx:  testTx,
		TestCtx: testCtx,
	}

	return testEnv
}

func TeardownSingle(testEnv TestEnv) {
	for {
		if testEnv.TestTx.mu.TryLock() {
			break
		}
	}
	defer testEnv.TestTx.mu.Unlock()

	err := testEnv.TestTx.tx.Rollback(testEnv.TestCtx)
	if err != nil {
		panic(err)
	}
}
