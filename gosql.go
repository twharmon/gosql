package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DB .
type DB struct {
	db *sql.DB
}

// ErrNotFound .
var ErrNotFound = errors.New("no result found")

// Conn .
func Conn(user string, pass string, host string, db string) (*DB, error) {
	d, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", user, pass, host, db))
	return &DB{d}, err
}

// Register .
func Register(ptrs ...interface{}) {
	for _, p := range ptrs {
		if reflect.TypeOf(p).Kind() != reflect.Ptr {
			panic("ptrs must be pointers to your model structs")
		}
		m := new(model)
		m.typ = reflect.TypeOf(p).Elem()
		m.name = m.typ.Name()
		m.table = strings.ToLower(m.name)

		m.mustBeValid()

		for i := 0; i < m.typ.NumField(); i++ {
			if !isField(m.typ.Field(i)) {
				continue
			}
			m.fields = append(m.fields, strings.ToLower(m.typ.Field(i).Name))
		}
		m.fieldCount = len(m.fields)
		models[m.name] = m
	}

	for _, m := range models {
		m.setInsertQuery()
		m.setUpdateQuery()
		m.setDeleteQuery()
	}
}
