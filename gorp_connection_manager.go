package gorp_locks

import (
	"context"
	"database/sql"
	"github.com/storage-lock/go-storage"
	"gopkg.in/gorp.v1"
)

// GorpConnectionManager 复用gorp的数据库连接（https://github.com/go-gorp/gorp）
type GorpConnectionManager struct {
	dbMap *gorp.DbMap
}

var _ storage.ConnectionManager[*sql.DB] = &GorpConnectionManager{}

func NewGorpConnectionManager(dbMap *gorp.DbMap) *GorpConnectionManager {
	return &GorpConnectionManager{
		dbMap: dbMap,
	}
}

const GorpConnectionManagerName = "gorp-connection-manager"

func (x *GorpConnectionManager) Name() string {
	return GorpConnectionManagerName
}

func (x *GorpConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	return x.dbMap.Db, nil
}

func (x *GorpConnectionManager) Return(ctx context.Context, db *sql.DB) error {
	return nil
}

func (x *GorpConnectionManager) Shutdown(ctx context.Context) error {
	if x.dbMap != nil {
		return x.dbMap.Db.Close()
	}
	return nil
}
