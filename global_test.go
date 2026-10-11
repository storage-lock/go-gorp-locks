package gorp_locks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGorpStorageLockFactoryBeanFactory_InitiallyInitialized(t *testing.T) {
	assert.NotNil(t, gorpStorageLockFactoryBeanFactory)
}

func TestNewGorpLock_NilDbMap(t *testing.T) {
	assert.Panics(t, func() {
		NewGorpLock(context.Background(), nil, "test-lock")
	})
}

func TestNewGorpLockWithOptions_NilDbMap(t *testing.T) {
	assert.Panics(t, func() {
		NewGorpLockWithOptions(context.Background(), nil, nil)
	})
}

func TestGetGorpLockFactory_NilDbMap(t *testing.T) {
	assert.Panics(t, func() {
		GetGorpLockFactory(context.Background(), nil)
	})
}

