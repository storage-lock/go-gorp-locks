package gorp_locks

import (
	"context"
	"database/sql"
	storage_lock_test_helper "github.com/storage-lock/go-storage-lock-test-helper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"
	"os"
	"testing"
)

func TestNewGorpLockBySqlDb(t *testing.T) {

	mysqlDsn := os.Getenv("STORAGE_LOCK_MYSQL_DSN")
	assert.NotEmpty(t, mysqlDsn)

	db, err := sql.Open("mysql", mysqlDsn)
	assert.Nil(t, err)

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	factory, err := GetGorpLockFactory(context.Background(), dbmap)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = 10
	storage_lock_test_helper.EveryOnePlayTimes = 100
	storage_lock_test_helper.TestStorageLock(t, factory)
}
