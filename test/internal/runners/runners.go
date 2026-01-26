package runners

import "testing"

type StorageKey interface {
	isKey()
}

// Storage handles type-safe storage of values.
type Storage interface {
	Store(StorageKey, any)
	Load(StorageKey) (any, bool)
}

// RunnerContext provides testing capabilities and access to shared storage.
type RunnerContext interface {
	testing.TB
	Storage() Storage
}

// runnerCtx implements RunnerContext.
type runnerCtx struct {
	testing.TB
	storage *typedStorage
}

func (c *runnerCtx) Storage() Storage {
	return c.storage
}

func NewRunnerContext(t testing.TB) RunnerContext {
	return &runnerCtx{
		TB: t,
		storage: &typedStorage{
			values: make(map[StorageKey]any),
		},
	}
}
