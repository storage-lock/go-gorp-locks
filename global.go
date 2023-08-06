package gorp_locks

import (
	storage_lock "github.com/storage-lock/go-storage-lock"
	"gopkg.in/gorp.v1"
	"sync"
)

var GlobalGorpLockFactory *GorpLockFactory
var globalGorpLockFactoryOnce sync.Once
var globalGorpLockFactoryErr error

func InitGlobalGorpLockFactory(dbMap *gorp.DbMap) error {
	factory, err := NewGorpLockFactory(dbMap)
	if err != nil {
		return err
	}
	GlobalGorpLockFactory = factory
	return nil
}

func NewGorpLock(dbMap *gorp.DbMap, lockId string) (*storage_lock.StorageLock, error) {
	globalGorpLockFactoryOnce.Do(func() {
		globalGorpLockFactoryErr = InitGlobalGorpLockFactory(dbMap)
	})
	if globalGorpLockFactoryErr != nil {
		return nil, globalGorpLockFactoryErr
	}
	return GlobalGorpLockFactory.CreateLock(lockId)
}
