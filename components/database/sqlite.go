package database

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
)

type SQLiteDatabase struct {
	Database string `yaml:"database"`
}

func NewSQLite() *SQLiteDatabase {
	return &SQLiteDatabase{
		Database: "./data/sqlite.db",
	}
}

func (c *SQLiteDatabase) Connect() (db.Session, error) {
	// options @see https://github.com/mattn/go-sqlite3#connection-string
	conn := sqlite.ConnectionURL{
		Database: c.Database,
		Options:  map[string]string{
			"mode": "rw", // 数据库访问模式, ro/rw/rwc/memory
			"cache": "private", // 缓存模式, shared or private
			"_busy_timeout": "5000", // 处理超时, 单位毫秒
			"_txlock": "deferred", // 事务锁, deferred/immediate/exclusive
		},
	}

	sess, err := sqlite.Open(conn)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
