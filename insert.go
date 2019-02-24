package gosql

import (
	"fmt"
	"reflect"
)

// Insert .
func (db *DB) Insert(obj interface{}) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer to your model struct")
	}
	m := models[t.Elem().Name()]
	if m == nil {
		return fmt.Errorf("you must first register %s", t.Elem().Name())
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
