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

// ErrNotFound .
var ErrNotFound = errors.New("no result found")

// Register .
func Register(structs ...interface{}) {
	for _, s := range structs {
		register(s)
	}
}

func register(s interface{}) {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("you can only register structs, %s found", reflect.TypeOf(s).Kind()))
	}
	m := new(model)
	m.typ = typ
	m.name = m.typ.Name()
	m.table = toSnakeCase(m.name)
	m.mustBeValid()
	for i := 0; i < m.typ.NumField(); i++ {
		if !isField(m.typ.Field(i)) {
			continue
		}
		tag, ok := m.typ.Field(i).Tag.Lookup("gosql")
		if ok && tag == "primary" {
			m.primaryFieldIndex = i
		}
		m.fields = append(m.fields, toSnakeCase(m.typ.Field(i).Name))
	}
	m.fieldCount = len(m.fields)
	models[m.name] = m
	m.setInsertQuery()
	m.setUpdateQuery()
	m.setDeleteQuery()
}

// Conn returns a reference to DB.
func Conn(db *sql.DB) *DB {
	return &DB{db}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
