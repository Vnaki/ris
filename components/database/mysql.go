package database

import (
	"time"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

type Database struct {
	Host                  string `yaml:"host"`
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	Database              string `yaml:"database"`
	MaxLifetime           int    `yaml:"max_lifetime"`
	MinIdleConnections    int    `yaml:"min_idle_connections"`
	MaxOpenConnections    int    `yaml:"max_open_connections"`
	MaxTransactionRetries int    `yaml:"max_transaction_retries"`
}

func New() *Database {
	return &Database{
		Host:                  "127.0.0.1:5432",
		User:                  "mysql",
		Database:              "mall",
		Password:              "mall@qazwsx",
		MinIdleConnections:    5,
		MaxOpenConnections:    200,
		MaxLifetime:           8,
		MaxTransactionRetries: 3,
	}
}

func (c *Database) Connect() (db.Session, error) {
	conn := mysql.ConnectionURL{
		Host:     c.Host,
		User:     c.User,
		Password: c.Password,
		Database: c.Database,
		Options:  map[string]string{},
	}

	sess, err := mysql.Open(conn)
	if err != nil {
		return nil, err
	}

	//db.LC().SetLevel()
	//db.LC().SetLogger(app.Logger())

	//sess.LC().SetLogger(&customLogger{})
	sess.SetMaxIdleConns(c.MinIdleConnections)
	sess.SetMaxOpenConns(c.MaxOpenConnections)
	sess.SetMaxTransactionRetries(c.MaxTransactionRetries)
	sess.SetConnMaxLifetime(time.Hour * time.Duration(c.MaxLifetime))
	sess.SetPreparedStatementCache(false)

	return sess, nil
}
