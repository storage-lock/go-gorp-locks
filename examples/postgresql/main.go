package main

import (
	"context"
	"database/sql"
	gorp_locks "github.com/storage-lock/go-gorp-locks"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"gopkg.in/gorp.v1"
)

func main() {

	// 连接数据库
	postgresqlDsn := "host=127.0.0.1 user=postgres password=postgres port=5432 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", postgresqlDsn)
	if err != nil {
		panic(err)
	}

	// 创建gorp的客户端
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	// 从gorp的客户端创建分布式锁工厂
	factory, err := gorp_locks.GetGorpLockFactory(context.Background(), dbMap)
	if err != nil {
		panic(err)
	}

	// 使用工厂创建一把锁
	lockId := "task-id-10086"
	lock, err := factory.CreateLock(lockId)
	if err != nil {
		panic(err)
	}
	// 开始使用锁来竞争资源
	ownerId := storage_lock.NewOwnerIdGenerator().GenOwnerId()
	err = lock.Lock(context.Background(), ownerId)
	if err != nil {
		panic(err)
	}
	// 锁使用完要记得释放
	defer func() {
		err := lock.UnLock(context.Background(), ownerId)
		if err != nil {
			panic(err)
		}
	}()

	// 下面编写的代码都是全局互斥的，同一时间只会有一个owner的代码被执行
	// ...

}
