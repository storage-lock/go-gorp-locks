package gorp_locks

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"
)

func testGorpDbMap(t *testing.T) *gorp.DbMap {
	t.Helper()
	db, err := sql.Open("sqlite3", "file:gorp_test.db?mode=memory")
	assert.Nil(t, err)
	t.Cleanup(func() { db.Close() })
	return &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
}

func TestNewGorpConnectionManager(t *testing.T) {
	dbMap := testGorpDbMap(t)
	m := NewGorpConnectionManager(dbMap)
	assert.NotNil(t, m)
	assert.Equal(t, GorpConnectionManagerName, m.Name())

	ctx := context.Background()
	got, err := m.Take(ctx)
	assert.Nil(t, err)
	assert.Equal(t, dbMap.Db, got)

	assert.Nil(t, m.Return(ctx, got))
	assert.Nil(t, m.Shutdown(ctx))
}

func TestGorpConnectionManager_ShutdownNil(t *testing.T) {
	m := NewGorpConnectionManager(nil)
	assert.Nil(t, m.Shutdown(context.Background()))
}

func TestNewGorpLockFactory(t *testing.T) {
	dbMap := testGorpDbMap(t)
	factory, err := NewGorpLockFactory(dbMap)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 端到端加解锁
	ctx := context.Background()
	lock, err := factory.CreateLock("test-gorp-lock")
	assert.Nil(t, err)
	assert.NotNil(t, lock)
	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestCreateStorageForGorp_Error(t *testing.T) {
	// Db 为 nil：NewStorage 返回 db can not nil
	_, err := CreateStorageForGorp(&gorp.DbMap{})
	assert.NotNil(t, err)

	// sqlmock 驱动不被识别：GetDriverNameForSqlDb 失败
	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()
	_, err = CreateStorageForGorp(&gorp.DbMap{Db: mockDB, Dialect: gorp.SqliteDialect{}})
	assert.NotNil(t, err)

	// NewGorpLockFactory 的错误分支
	_, err = NewGorpLockFactory(&gorp.DbMap{})
	assert.NotNil(t, err)
}

func TestGorpGlobal_Error(t *testing.T) {
	old := newGorpLockFactory
	newGorpLockFactory = func(dbMap *gorp.DbMap) (*GorpLockFactory, error) {
		return nil, errors.New("boom")
	}
	defer func() { newGorpLockFactory = old }()

	ctx := context.Background()
	boomDbMap := &gorp.DbMap{}
	_, err := GetGorpLockFactory(ctx, boomDbMap)
	assert.NotNil(t, err)
	_, err = NewGorpLock(ctx, boomDbMap, "x")
	assert.NotNil(t, err)
	_, err = NewGorpLockWithOptions(ctx, boomDbMap, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewGorpLock_Success(t *testing.T) {
	dbMap := testGorpDbMap(t)
	ctx := context.Background()

	factory, err := GetGorpLockFactory(ctx, dbMap)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetGorpLockFactory(ctx, dbMap)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewGorpLock(ctx, dbMap, "test-gorp-lock-global")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-gorp-lock-global-opt")
	lock2, err := NewGorpLockWithOptions(ctx, dbMap, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)
}
