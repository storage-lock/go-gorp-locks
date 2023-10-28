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

func TestMysql(t *testing.T) {
	mysqlDsn := os.Getenv("STORAGE_LOCK_MYSQL_DSN")
	assert.NotEmpty(t, mysqlDsn)

	db, err := sql.Open("mysql", mysqlDsn)
	assert.Nil(t, err)

	testSqlDb(t, db, 30)
}

func TestPostgreSQL(t *testing.T) {
	postgresqlDsn := os.Getenv("STORAGE_LOCK_POSTGRESQL_DSN")
	assert.NotEmpty(t, postgresqlDsn)

	db, err := sql.Open("postgres", postgresqlDsn)
	assert.Nil(t, err)

	testSqlDb(t, db, 30)
}

func TestSqlite3(t *testing.T) {
	dbPath := os.Getenv("STORAGE_LOCK_SQLITE3_DB_PATH")
	assert.NotEmpty(t, dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	assert.Nil(t, err)

	testSqlDb(t, db, 3)
}

// 单元测试的公共逻辑提取
func testSqlDb(t *testing.T, db *sql.DB, playNum int) {
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	factory, err := GetGorpLockFactory(context.Background(), dbMap)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = playNum
	storage_lock_test_helper.EveryOnePlayTimes = 100
	storage_lock_test_helper.TestStorageLock(t, factory)
}
