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

// Conn returns a reference to DB.
func Conn(db *sql.DB) *DB {
	return &DB{db}
}

// Insert .
func (db *DB) Insert(obj interface{}) error {
	m, err := getModelOf(obj)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(obj).Elem()
	_, err = db.db.Exec(m.insertQuery, m.getArgs(v)...)
	return err
}

// Update .
func (db *DB) Update(obj interface{}) error {
	m, err := getModelOf(obj)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(obj).Elem()
	_, err = db.db.Exec(m.updateQuery, m.getArgsIDLast(v)...)
	return err
}

// Delete .
func (db *DB) Delete(obj interface{}) error {
	m, err := getModelOf(obj)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(obj).Elem()
	_, err = db.db.Exec(m.deleteQuery, m.getIDArg(v))
	return err
}

// Exec .
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

// Select .
func (db *DB) Select(fields ...string) *SelectQuery {
	sq := new(SelectQuery)
	sq.db = db
	sq.fields = fields
	sq.limit = 1000
	return sq
}

// Table .
func (db *DB) Table(table string) *TableQuery {
	tq := new(TableQuery)
	tq.db = db
	tq.table = table
	return tq
}
