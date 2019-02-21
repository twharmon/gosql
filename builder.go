package gosql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type join struct {
	relative interface{}
	fields   []string
}

// QueryBuilder .
type QueryBuilder struct {
	db     *DB
	fields []string
	where  string
	args   []interface{}
	limit  uint64
	query  strings.Builder
}

// Select .
func (db *DB) Select(fields ...string) *QueryBuilder {
	qb := new(QueryBuilder)
	qb.db = db
	qb.fields = fields
	qb.limit = 0
	qb.query.WriteString("select ")
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
	if !isPointer(t) {
		return fmt.Errorf("out must be a pointer")
	}
	e := t.Elem()
	switch e.Kind().String() {
	case "struct":
		return qb.toOne(e, out)
	case "slice":
		return qb.toMany(out)
	default:
		return fmt.Errorf("models must be a struct or slice (%s found)", e.Kind().String())
	}
}

func (qb *QueryBuilder) toOne(t reflect.Type, out interface{}) error {
	qb.makeQuery(t)
	row := qb.db.db.QueryRow(qb.query.String(), qb.args...)
	m := models[t.Name()]
	return row.Scan(qb.getDests(m, reflect.ValueOf(out).Elem())...)
}

func (qb *QueryBuilder) toMany(outs interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(outs))
	t := reflect.TypeOf(outs).Elem()
	if !isPointer(t.Elem()) {
		return fmt.Errorf("outs must be a slice of pointers")
	}
	e := t.Elem().Elem()
	m := models[e.Name()]
	qb.makeQuery(e)
	rows, err := qb.db.db.Query(qb.query.String(), qb.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	newOuts := reflect.MakeSlice(t, 0, int(qb.limit))
	for i := 0; rows.Next(); i++ {
		newOut := reflect.New(e)
		if err := rows.Scan(qb.getDests(m, newOut.Elem())...); err != nil {
			return err
		}
		newOuts = reflect.Append(newOuts, newOut)
	}

	v.Set(newOuts)
	return nil
}

func (qb *QueryBuilder) makeQuery(t reflect.Type) {
	m := models[t.Name()]
	if qb.fields[0] == "*" {
		for i := 0; i < m.fieldCount-1; i++ {
			qb.query.WriteString(m.fields[i])
			qb.query.WriteString(", ")
		}
		qb.query.WriteString(m.fields[m.fieldCount-1])
	} else {
		for i := 0; i < len(qb.fields)-1; i++ {
			qb.query.WriteString(qb.fields[i])
			qb.query.WriteString(", ")
		}
		qb.query.WriteString(qb.fields[len(qb.fields)-1])
	}
	qb.query.WriteString(" from ")
	qb.query.WriteString(m.table)

	if qb.where != "" {
		qb.query.WriteString(" where ")
		qb.query.WriteString(qb.where)
	}
	if qb.limit > 0 {
		qb.query.WriteString(" limit ")
		qb.query.WriteString(strconv.FormatUint(qb.limit, 10))
	}
}

func (qb *QueryBuilder) getDests(m *model, v reflect.Value) []interface{} {
	if qb.fields[0] == "*" {
		scans := make([]interface{}, m.fieldCount)
		for i := 0; i < m.fieldCount; i++ {
			scans[i] = v.Field(i).Addr().Interface()
		}
		return scans
	}
	fieldCount := len(qb.fields)
	scans := make([]interface{}, fieldCount)
	for i := 0; i < fieldCount; i++ {
		scans[i] = v.Field(m.getFieldIndexByName(qb.fields[i])).Addr().Interface()
	}
	return scans
}
