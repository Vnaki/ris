package types

import "github.com/upper/db/v4"

type Database interface {
	Connect() (db.Session, error)
}
