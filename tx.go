package gosql

import (
	"database/sql"
	"reflect"
)

// Tx .
type Tx struct {
	tx *sql.Tx
	db *DB
}

// Commit .
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback .
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// Insert insterts a row in the database.
func (t *Tx) Insert(obj interface{}) (sql.Result, error) {
	m, err := t.db.getModelOf(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	return t.tx.Exec(m.getInsertQuery(v), m.getArgs(v)...)
}

// Update updates a row in the database.
func (t *Tx) Update(obj interface{}) (sql.Result, error) {
	m, err := t.db.getModelOf(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	return t.tx.Exec(m.getUpdateQuery(), m.getArgsPrimaryLast(v)...)
}

// Delete deletes a row from the database.
func (t *Tx) Delete(obj interface{}) (sql.Result, error) {
	m, err := t.db.getModelOf(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	var inserts []interface{}
	for _, i := range m.primaryFieldIndecies {
		inserts = append(inserts, v.Field(i).Interface())
	}
	return t.tx.Exec(m.getDeleteQuery(), inserts...)
}

// Exec is a wrapper around sql.DB.Exec().
func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

// Query is a wrapper around sql.DB.Query().
func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

// QueryRow is a wrapper around sql.DB.QueryRow().
func (t *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

// Select selects columns of a table.
func (t *Tx) Select(fields ...string) *SelectQuery {
	sq := new(SelectQuery)
	sq.db = t.db
	sq.querier = t.tx
	sq.fields = fields
	return sq
}

// ManualUpdate starts a query for manually updating rows in a table.
func (t *Tx) ManualUpdate(table string) *UpdateQuery {
	uq := new(UpdateQuery)
	uq.db = t.db
	uq.execer = t.tx
	uq.table = table
	return uq
}

// Count starts a query for counting rows in a table.
func (t *Tx) Count(table string, count string) *CountQuery {
	cq := new(CountQuery)
	cq.db = t.db
	cq.queryRower = t.tx
	cq.table = table
	cq.count = count
	return cq
}

// ManualDelete starts a query for manually deleting rows in a table.
func (t *Tx) ManualDelete(table string) *DeleteQuery {
	dq := new(DeleteQuery)
	dq.db = t.db
	dq.execer = t.tx
	dq.table = table
	return dq
}
