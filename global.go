package gorp_locks

import (
	"context"
	"database/sql"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"gopkg.in/gorp.v1"
)

// 用于管理工厂的工厂，主要是为了能够同时支持多个实例
var gorpStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*gorp.DbMap, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*gorp.DbMap, *sql.DB]()

// NewGorpLock
// @Description: 基于gorp创建一把分布式锁
// @param ctx
// @param dbMap gorp的map，用于进行各种gorp的操作，锁操作需要用到这些api
// @param lockId 要创建的锁的ID，ID相同即被认为是同一把锁，同一把锁会互相排斥
// @return *storage_lock.StorageLock
// @return error
func NewGorpLock(ctx context.Context, dbMap *gorp.DbMap, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetGorpLockFactory(ctx, dbMap)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

// NewGorpLockWithOptions
// @Description:
// @param ctx
// @param dbMap
// @param options
// @return *storage_lock.StorageLock
// @return error
func NewGorpLockWithOptions(ctx context.Context, dbMap *gorp.DbMap, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetGorpLockFactory(ctx, dbMap)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

// GetGorpLockFactory
// @Description: 从gorp创建一个锁工厂，后续可以拿着这个锁工厂再创建
// @param ctx 超时控制之类的
// @param dbMap gorp对应的dbMap
// @return *storage_lock_factory.StorageLockFactory[*sql.DB]
// @return error
func GetGorpLockFactory(ctx context.Context, dbMap *gorp.DbMap) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return gorpStorageLockFactoryBeanFactory.GetOrInit(ctx, dbMap, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		lockFactory, err := NewGorpLockFactory(dbMap)
		if err != nil {
			return nil, err
		}
		return lockFactory.StorageLockFactory, nil
	})
}
