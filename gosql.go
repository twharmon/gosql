package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sort"
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
			m.fields = append(m.fields, strings.ToLower(m.typ.Field(i).Name))
		}
		m.fieldCount = len(m.fields)
		models[m.name] = m
	}

	for _, m := range models {
		for i := 0; i < m.typ.NumField(); i++ {
			f := m.typ.Field(i)
			if m.isOneToMany(f) {
				otm := new(oneToMany)
				otm.typStr = f.Type.String()
				otm.field = f
				m.oneToManys = append(m.oneToManys, otm)
			} else if m.isManyToOne(f) {
				mto := new(manyToOne)
				mto.typStr = f.Type.String()
				mto.field = f
				mto.column = strings.ToLower(f.Name) + "_id"
				m.manyToOnes = append(m.manyToOnes, mto)
			} else if m.isManyToMany(f) {
				mtm := new(manyToMany)
				mtm.typStr = f.Type.String()
				mtm.field = f
				fTableName := strings.ToLower(strings.Split(f.Type.String(), ".")[1])
				tables := []string{fTableName, strings.ToLower(m.table)}
				sort.Strings(tables)
				mtm.table = tables[0] + "_" + tables[1]
				mtm.column = fTableName + "_id"
				m.manyToManys = append(m.manyToManys, mtm)
			}
		}
	}

	for _, m := range models {
		m.setInsertQuery()
		m.setUpdateQuery()
		m.setDeleteQuery()
	}

	return nil
}
