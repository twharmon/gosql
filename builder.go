package gosql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	sqlActionSelect = iota
	sqlActionDelete
)

// QueryBuilder .
type QueryBuilder struct {
	db     *DB
	model  *model
	action int
	fields []string
	joins  []string
	where  string
	args   []interface{}
	limit  uint64
}

// Query .
func (db *DB) Query() *QueryBuilder {
	qb := new(QueryBuilder)
	qb.db = db
	return qb
}

// Select .
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.action = sqlActionSelect
	qb.fields = fields
	return qb
}

// Delete .
func (qb *QueryBuilder) Delete(from string) *QueryBuilder {
	qb.action = sqlActionDelete
	qb.model = &model{table: from}
	return qb
}

// Join .
func (qb *QueryBuilder) Join(join string) *QueryBuilder {
	qb.joins = append(qb.joins, join)
	return qb
}

// Where .
func (qb *QueryBuilder) Where(where string, args ...interface{}) *QueryBuilder {
	qb.where = where
	qb.args = args
	return qb
}

// Limit .
func (qb *QueryBuilder) Limit(limit uint64) *QueryBuilder {
	qb.limit = limit
	return qb
}

// To .
func (qb *QueryBuilder) To(out interface{}) error {
	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("out must be a pointer")
	}
	e := t.Elem()
	switch e.Kind() {
	case reflect.Struct:
		qb.model = models[e.Name()]
		if qb.model == nil {
			return fmt.Errorf("you must first register %s", e.Name())
		}
		return qb.toOne(out)
	case reflect.Slice:
		ptr := e.Elem()
		if ptr.Kind() != reflect.Ptr {
			return fmt.Errorf("out must be a slice of pointers")
		}
		strct := ptr.Elem()
		if strct.Kind() != reflect.Struct {
			return fmt.Errorf("out must be a slice of pointers to structs")
		}
		qb.model = models[strct.Name()]
		if qb.model == nil {
			return fmt.Errorf("you must first register %s", strct.Name())
		}
		return qb.toMany(e, out)
	default:
		return fmt.Errorf("models must be a struct or slice (%s found)", e.Kind().String())
	}
}

// Exec .
func (qb *QueryBuilder) Exec() error {
	_, err := qb.db.db.Exec(qb.string(), qb.args...)
	return err
}

func (qb *QueryBuilder) toOne(out interface{}) error {
	row := qb.db.db.QueryRow(qb.string(), qb.args...)
	err := row.Scan(qb.getDests(reflect.ValueOf(out).Elem())...)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (qb *QueryBuilder) toMany(sliceType reflect.Type, outs interface{}) error {
	rows, err := qb.db.db.Query(qb.string(), qb.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	newOuts := reflect.MakeSlice(sliceType, 0, int(qb.limit))
	for i := 0; rows.Next(); i++ {
		newOut := reflect.New(qb.model.typ)
		if err := rows.Scan(qb.getDests(newOut.Elem())...); err != nil {
			return err
		}
		newOuts = reflect.Append(newOuts, newOut)
	}
	reflect.Indirect(reflect.ValueOf(outs)).Set(newOuts)
	return nil
}

func (qb *QueryBuilder) string() string {
	var q strings.Builder
	switch qb.action {
	case sqlActionSelect:
		q.WriteString("select ")
		for i := 0; i < len(qb.fields)-1; i++ {
			if qb.joins != nil {
				q.WriteString(qb.model.table)
				q.WriteString(".")
			}
			q.WriteString(qb.fields[i])
			q.WriteString(", ")
		}
		if qb.joins != nil {
			q.WriteString(qb.model.table)
			q.WriteString(".")
		}
		q.WriteString(qb.fields[len(qb.fields)-1])

		q.WriteString(" from ")
		q.WriteString(qb.model.table)
	case sqlActionDelete:
		q.WriteString("delete from ")
		q.WriteString(qb.model.table)
	}
	for _, join := range qb.joins {
		q.WriteString(" join ")
		q.WriteString(join)
	}
	if qb.where != "" {
		q.WriteString(" where ")
		q.WriteString(qb.where)
	}
	if qb.limit > 0 {
		q.WriteString(" limit ")
		q.WriteString(strconv.FormatUint(qb.limit, 10))
	}
	return q.String()
}

func (qb *QueryBuilder) getDests(v reflect.Value) []interface{} {
	if qb.fields[0] == "*" {
		scans := make([]interface{}, qb.model.fieldCount)
		for i := 0; i < qb.model.fieldCount; i++ {
			scans[i] = v.Field(i).Addr().Interface()
		}
		return scans
	}
	fieldCount := len(qb.fields)
	scans := make([]interface{}, fieldCount)
	for i := 0; i < fieldCount; i++ {
		scans[i] = v.Field(qb.model.getFieldIndexByName(qb.fields[i])).Addr().Interface()
	}
	return scans
}
