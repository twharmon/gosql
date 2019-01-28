package gosql

import (
	"database/sql"
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

// Conn .
func Conn(user string, pass string, host string, db string) (*DB, error) {
	d, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", user, pass, host, db))
	return &DB{d}, err
}

// MustPrepare .
func MustPrepare(ptrs ...interface{}) error {
	for _, p := range ptrs {
		if !isPointer(reflect.TypeOf(p)) {
			return fmt.Errorf("ptrs must be pointers to your model structs")
		}
		m := new(model)
		m.typ = reflect.TypeOf(p).Elem()
		m.name = m.typ.Name()
		m.table = strings.ToLower(m.name)

		if err := m.mustBeValid(); err != nil {
			return err
		}

		for i := 0; i < m.typ.NumField(); i++ {
			if !isField(m.typ.Field(i)) {
				continue
			}
			m.fields = append(m.fields, m.typ.Field(i))
		}
		m.fieldCount = len(m.fields)
		models[m.name] = m
	}

	for _, m := range models {
		for i := 0; i < m.typ.NumField(); i++ {
			f := m.typ.Field(i)
			if isOneToMany(m, f) {
				m.oneToManys = append(m.oneToManys, f)
			} else if isManyToOne(m, f) {
				m.manyToOnes = append(m.manyToOnes, f)
			} else if isOneToOne(m, f) {
				m.oneToOnes = append(m.oneToOnes, f)
			} else if isManyToMany(m, f) {
				m.manyToManys = append(m.manyToManys, f)
			}
		}
	}

	for _, m := range models {
		m.setInsertQuery()
		m.setUpdateQuery()
		m.setDeleteQuery()
		fmt.Println(m.insertQuery)
		fmt.Println(m.updateQuery)
		fmt.Println(m.deleteQuery)
	}

	return nil
}
