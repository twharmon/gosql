package gosql

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// DB is a wrapper around sql.DB.
type DB struct {
	db     *sql.DB
	models map[string]*model
}

func (db *DB) register(typ reflect.Type) error {
	m := new(model)
	m.typ = typ
	m.name = m.typ.Name()
	m.table = toSnakeCase(m.name)
	m.primaryFieldIndecies = nil
	for i := 0; i < m.typ.NumField(); i++ {
		f := m.typ.Field(i)
		if tag, ok := f.Tag.Lookup("idx"); ok && tag == "primary" {
			m.primaryFieldIndecies = append(m.primaryFieldIndecies, i)
		}
		if tag, ok := f.Tag.Lookup("col"); ok {
			if tag == "-" {
				continue
			}
			m.fields = append(m.fields, tag)
		} else {
			m.fields = append(m.fields, toSnakeCase(f.Name))
		}
	}
	if err := db.mustBeValid(m); err != nil {
		return err
	}
	db.models[m.name] = m
	return nil
}

func (db *DB) getModelOf(t reflect.Type) (*model, error) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a pointer to your model struct")
	}
	m := db.models[t.Name()]
	if m == nil {
		if err := db.register(t); err != nil {
			return nil, err
		}
		m = db.models[t.Name()]
	}
	return m, nil
}

func (db *DB) mustBeValid(m *model) error {
	if db.models[m.name] != nil {
		return fmt.Errorf("model %s found more than once", m.name)
	}
	if len(m.primaryFieldIndecies) == 0 {
		return fmt.Errorf("model %s must have at least one field tagged `idx:\"primary\"`", m.name)
	}
	return nil
}

// Begin starts a transaction.
func (db *DB) Begin() (*Tx, error) {
	sqlTx, err := db.db.Begin()
	tx := Tx{
		tx: sqlTx,
		db: db,
	}
	return &tx, err
}

// Insert insterts a row in the database.
func (db *DB) Insert(obj interface{}) (sql.Result, error) {
	m, err := db.getModelOf(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	return db.db.Exec(m.getInsertQuery(v), m.getArgs(v)...)
}

// Update updates a row in the database.
func (db *DB) Update(obj interface{}) (sql.Result, error) {
	m, err := db.getModelOf(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	return db.db.Exec(m.getUpdateQuery(), m.getArgsPrimaryLast(v)...)
}

// Delete deletes a row from the database.
func (db *DB) Delete(obj interface{}) (sql.Result, error) {
	m, err := db.getModelOf(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(obj).Elem()
	var inserts []interface{}
	for _, i := range m.primaryFieldIndecies {
		inserts = append(inserts, v.Field(i).Interface())
	}
	return db.db.Exec(m.getDeleteQuery(), inserts...)
}

// Exec is a wrapper around sql.DB.Exec().
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

// Query is a wrapper around sql.DB.Query().
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

// QueryRow is a wrapper around sql.DB.QueryRow().
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

// Select selects columns of a table.
func (db *DB) Select(fields ...string) *SelectQuery {
	sq := new(SelectQuery)
	sq.db = db
	sq.querier = db.db
	sq.fields = fields
	return sq
}

// ManualUpdate starts a query for manually updating rows in a table.
func (db *DB) ManualUpdate(table string) *UpdateQuery {
	uq := new(UpdateQuery)
	uq.db = db
	uq.execer = db.db
	uq.table = table
	return uq
}

// Count starts a query for counting rows in a table.
func (db *DB) Count(table string, count string) *CountQuery {
	cq := new(CountQuery)
	cq.db = db
	cq.queryRower = db.db
	cq.table = table
	cq.count = count
	return cq
}

// ManualDelete starts a query for manually deleting rows in a table.
func (db *DB) ManualDelete(table string) *DeleteQuery {
	dq := new(DeleteQuery)
	dq.db = db
	dq.execer = db.db
	dq.table = table
	return dq
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
