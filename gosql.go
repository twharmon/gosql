package gosql

import (
	"database/sql"
	"errors"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ErrNotFound is returned when a query for one result returns no
// results.
var ErrNotFound = errors.New("no result found")

// New returns a reference to DB.
func New(db *sql.DB) *DB {
	return &DB{
		db:     db,
		models: make(map[string]*model),
	}
}
