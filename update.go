package gosql

import (
	"reflect"
)

// Update .
func (db *DB) Update(obj interface{}) error {
	m, err := getModelOf(obj)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(obj).Elem()
	_, err = db.db.Exec(m.updateQuery, m.getArgsIDLast(v)...)
	return err
}
