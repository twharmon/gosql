package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ErrNotFound .
var ErrNotFound = errors.New("Record not found")

// DB .
type DB struct {
	conn *sql.DB
}

var queries struct {
	inserts map[string]string
	updates map[string]string
	deletes map[string]string
}

func init() {
	queries.inserts = make(map[string]string)
	queries.updates = make(map[string]string)
	queries.deletes = make(map[string]string)
}

// PrepareAndValidate .
func (db *DB) PrepareAndValidate(models ...interface{}) error {
	for _, model := range models {
		t := reflect.TypeOf(model).Elem()
		v := reflect.ValueOf(model).Elem()
		if queries.inserts[t.Name()] != "" {
			return fmt.Errorf("model %s found more than once", t.Name())
		}
		hasValidID := false
		for i := 0; i < t.NumField(); i++ {
			f := v.Field(i).Type().String()
			if t.Field(i).Name == "ID" && (f == "uint" || f == "uint64") {
				hasValidID = true
			}
		}
		if !hasValidID {
			return fmt.Errorf("model %s must have field 'ID' with type uint or uint64", t.Name())
		}
		queries.inserts[t.Name()] = getInsertPreparedQuery(t, model)
		queries.deletes[t.Name()] = getDeletePreparedQuery(t, model)
	}
	return nil
}

// Conn .
func Conn(user string, pass string, host string, db string) (*DB, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", user, pass, host, db))
	return &DB{conn}, err
}

// FindOne .
func (db *DB) FindOne(out interface{}, where string, args ...interface{}) error {
	v := reflect.ValueOf(out).Elem()
	q := fmt.Sprintf(
		"select * from %s where %s limit 1",
		strings.ToLower(reflect.TypeOf(out).Elem().Name()),
		where,
	)
	err := db.conn.QueryRow(q, args...).Scan(getDestinations(v)...)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

// Insert .
func (db *DB) Insert(obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	res, err := db.conn.Exec(queries.inserts[reflect.TypeOf(obj).Elem().Name()], getArgs(v, obj)...)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	v.FieldByName("ID").SetUint(uint64(id))
	return nil
}

// Delete .
func (db *DB) Delete(obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	res, err := db.conn.Exec(queries.deletes[reflect.TypeOf(obj).Elem().Name()], getIDArg(v, obj))
	a, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if a == 0 {
		return ErrNotFound
	}
	return err
}

func getInsertPreparedQuery(t reflect.Type, model interface{}) string {
	queryStart := fmt.Sprintf("insert into %s (", strings.ToLower(t.Name()))
	queryEnd := "values ("
	for i := 0; i < t.NumField(); i++ {
		queryStart = queryStart + strings.ToLower(t.Field(i).Name)
		queryEnd = queryEnd + "?"
		if i+1 != t.NumField() {
			queryEnd = queryEnd + ", "
			queryStart = queryStart + ", "
		} else {
			queryStart = queryStart + ") "
			queryEnd = queryEnd + ")"
		}
	}
	return queryStart + queryEnd
}

func getDeletePreparedQuery(t reflect.Type, model interface{}) string {
	return fmt.Sprintf("delete from %s where id = ?", strings.ToLower(t.Name()))
}

func getArgs(v reflect.Value, obj interface{}) []interface{} {
	args := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		args[i] = v.Field(i).Interface()
	}
	return args
}

func getIDArg(v reflect.Value, obj interface{}) interface{} {
	return v.FieldByName("ID").Interface()
}

func getDestinations(v reflect.Value) []interface{} {
	scans := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		scans[i] = v.Field(i).Addr().Interface()
	}
	return scans
}
