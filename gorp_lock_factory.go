package gorp_locks

import (
	"database/sql"
	sqldb_storage "github.com/storage-lock/go-sqldb-storage"
	"github.com/storage-lock/go-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"gopkg.in/gorp.v1"
)

type GorpLockFactory struct {
	dbMap *gorp.DbMap
	*storage_lock_factory.StorageLockFactory[*sql.DB]
}

func NewGorpLockFactory(dbMap *gorp.DbMap) (*GorpLockFactory, error) {
	connectionManager := NewGorpConnectionManager(dbMap)

	storage, err := CreateStorageForGorp(dbMap, connectionManager)
	if err != nil {
		return nil, err
	}

	factory := storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager)

	return &GorpLockFactory{
		dbMap:              dbMap,
		StorageLockFactory: factory,
	}, nil
}

// CreateStorageForGorp 尝试从gorp创建Storage
func CreateStorageForGorp(dbMap *gorp.DbMap, connectionManager storage.ConnectionManager[*sql.DB]) (storage.Storage, error) {
	return sqldb_storage.NewStorageBySqlDb(dbMap.Db, connectionManager)
}
