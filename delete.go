package gosql

import (
	"fmt"
	"reflect"
)

// Delete .
func (db *DB) Delete(obj interface{}) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer to your model struct")
	}
	m := models[t.Elem().Name()]
	if m == nil {
		return fmt.Errorf("you must first register %s", t.Elem().Name())
	}
	v := reflect.ValueOf(obj).Elem()
	_, err := db.db.Exec(m.deleteQuery, m.getIDArg(v))
	return err
}
