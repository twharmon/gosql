package gosql

import (
	"reflect"
)

// Delete .
func (db *DB) Delete(obj interface{}) error {
	m, err := getModelOf(obj)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(obj).Elem()
	_, err = db.db.Exec(m.deleteQuery, m.getIDArg(v))
	return err
}
