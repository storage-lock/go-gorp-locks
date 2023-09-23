package gorp_locks

import (
	"database/sql"
	sqldb_storage "github.com/storage-lock/go-sqldb-storage"
	"github.com/storage-lock/go-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"gopkg.in/gorp.v1"
)

// GorpLockFactory gorp的锁工厂，可以创建好多锁啥的
type GorpLockFactory struct {
	dbMap *gorp.DbMap
	*storage_lock_factory.StorageLockFactory[*sql.DB]
}

// NewGorpLockFactory 从*gorp.DbMap中创建
func NewGorpLockFactory(dbMap *gorp.DbMap) (*GorpLockFactory, error) {
	connectionManager := NewGorpConnectionManager(dbMap)

	gorpStorage, err := CreateStorageForGorp(dbMap)
	if err != nil {
		return nil, err
	}

	factory := storage_lock_factory.NewStorageLockFactory[*sql.DB](gorpStorage, connectionManager)

	return &GorpLockFactory{
		dbMap:              dbMap,
		StorageLockFactory: factory,
	}, nil
}

// CreateStorageForGorp 尝试从gorp创建Storage
func CreateStorageForGorp(dbMap *gorp.DbMap) (storage.Storage, error) {
	return sqldb_storage.NewStorage(dbMap.Db)
}
