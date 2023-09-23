package gorp_locks

import (
	"context"
	"database/sql"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"gopkg.in/gorp.v1"
)

var gorpStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*gorp.DbMap, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*gorp.DbMap, *sql.DB]()

func NewGorpLock(ctx context.Context, dbMap *gorp.DbMap, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetGorpLockFactory(ctx, dbMap)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewGorpLockWithOptions(ctx context.Context, dbMap *gorp.DbMap, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetGorpLockFactory(ctx, dbMap)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetGorpLockFactory(ctx context.Context, dbMap *gorp.DbMap) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return gorpStorageLockFactoryBeanFactory.GetOrInit(ctx, dbMap, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		lockFactory, err := NewGorpLockFactory(dbMap)
		if err != nil {
			return nil, err
		}
		return lockFactory.StorageLockFactory, nil
	})
}
