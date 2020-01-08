package gosql

import (
	"database/sql"
	"reflect"
)

// DB is a wrapper around sql.DB.
type DB struct {
	db *sql.DB
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (db *DB) SetMaxOpenConns(max int) {
	db.db.SetMaxOpenConns(max)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
//
// If n <= 0, no idle connections are retained.
//
// The default max idle connections is currently 2. This may change in
// a future release.
func (db *DB) SetMaxIdleConns(max int) {
	db.db.SetMaxIdleConns(max)
}

// Insert .
func (db *DB) Insert(obj interface{}) (sql.Result, error) {
	m, err := getModelOf(obj)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	return db.db.Exec(m.getInsertQuery(v), m.getArgs(v)...)
}

// Update .
func (db *DB) Update(obj interface{}) (sql.Result, error) {
	m, err := getModelOf(obj)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	return db.db.Exec(m.getUpdateQuery(v), m.getArgs(v)...)
}

// Exec .
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

// Query .
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

// QueryRow .
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

// Select .
func (db *DB) Select(fields ...string) *SelectQuery {
	sq := new(SelectQuery)
	sq.db = db
	sq.fields = fields
	sq.limit = 1000
	return sq
}

// ManualUpdate .
func (db *DB) ManualUpdate(table string) *UpdateQuery {
	uq := new(UpdateQuery)
	uq.db = db
	uq.table = table
	return uq
}

// Count .
func (db *DB) Count(table string, count string) *CountQuery {
	cq := new(CountQuery)
	cq.db = db
	cq.table = table
	cq.count = count
	return cq
}

// ManualDelete .
func (db *DB) ManualDelete(table string) *DeleteQuery {
	dq := new(DeleteQuery)
	dq.db = db
	dq.table = table
	return dq
}
