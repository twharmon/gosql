package gosql

import (
	"fmt"
	"reflect"
)

// Update .
func (db *DB) Update(obj interface{}) error {
	t := reflect.TypeOf(obj)
	if !isPointer(t) {
		return fmt.Errorf("obj must be a pointer to your model struct")
	}
	v := reflect.ValueOf(obj).Elem()
	m := models[t.Elem().Name()]
	_, err := db.db.Exec(m.updateQuery, getArgsIDLast(m, v, obj)...)
	return err
}
