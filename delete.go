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
	v := reflect.ValueOf(obj).Elem()
	m := models[t.Elem().Name()]
	_, err := db.db.Exec(m.deleteQuery, m.getIDArg(v))
	return err
}
