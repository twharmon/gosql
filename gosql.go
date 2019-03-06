package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
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

// ConnDB .
func ConnDB(db *sql.DB) *DB {
	return &DB{db}
}

// Register .
func Register(structs ...interface{}) {
	for _, s := range structs {
		if reflect.TypeOf(s).Kind() != reflect.Struct {
			panic("structs must be pointers to your model structs")
		}
		m := new(model)
		m.typ = reflect.TypeOf(s)
		m.name = m.typ.Name()
		m.table = strings.ToLower(m.name)
		m.mustBeValid()
		for i := 0; i < m.typ.NumField(); i++ {
			if !isField(m.typ.Field(i)) {
				continue
			}
			m.fields = append(m.fields, toSnakeCase(m.typ.Field(i).Name))
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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
