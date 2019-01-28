package gosql

import (
	"fmt"
	"reflect"
)

// Insert .
func (db *DB) Insert(obj interface{}) error {
	t := reflect.TypeOf(obj)
	if !isPointer(t) {
		return fmt.Errorf("obj must be a pointer to your model struct")
	}
	v := reflect.ValueOf(obj).Elem()
	m := models[t.Elem().Name()]
	res, err := db.db.Exec(m.insertQuery, getArgs(m, v, obj)...)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	v.Field(0).SetUint(uint64(id))
	return nil
}
