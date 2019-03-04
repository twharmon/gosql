package gosql

import (
	"reflect"
)

// Insert .
func (db *DB) Insert(obj interface{}) error {
	m, err := getModelOf(obj)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(obj).Elem()
	res, err := db.db.Exec(m.insertQuery, m.getArgs(v)...)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	v.Field(0).SetInt(id)
	return nil
}
