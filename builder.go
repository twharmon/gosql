package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	sqlActionSelect = iota
	sqlActionDelete
	sqlActionCountFrom
)

// ErrBadQuery .
var ErrBadQuery = errors.New("bad query")

// QueryBuilder .
type QueryBuilder struct {
	db     *DB
	model  *model
	action int
	fields []string
	joins  []string
	where  string
	args   []interface{}
	order  string
	limit  uint64
}

// QueryBuilder2 .
type QueryBuilder2 struct {
	db     *DB
	action int
	from   string
	joins  []string
	where  string
	args   []interface{}
}

// Query .
func (db *DB) Query() *QueryBuilder {
	qb := new(QueryBuilder)
	qb.db = db
	return qb
}

// Query2 .
func (db *DB) Query2() *QueryBuilder2 {
	qb := new(QueryBuilder2)
	qb.db = db
	return qb
}

// From .
func (qb *QueryBuilder2) From(table string) *QueryBuilder2 {
	qb.from = table
	return qb
}

// Select .
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.action = sqlActionSelect
	qb.fields = fields
	return qb
}

// From .
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.model = &model{table: table}
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

// OrderBy .
func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.order = order
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

// Delete .
func (qb *QueryBuilder) Delete() error {
	qb.action = sqlActionDelete
	query, err := qb.string()
	if err != nil {
		return err
	}

	_, err = qb.db.db.Exec(query, qb.args...)
	return err
}

// Count .
func (qb *QueryBuilder) Count() (int64, error) {
	qb.action = sqlActionCountFrom
	query, err := qb.string()
	if err != nil {
		return 0, err
	}

	var count int64
	row := qb.db.db.QueryRow(query, qb.args...)
	err = row.Scan(&count)
	return count, err
}

func (qb *QueryBuilder) toOne(out interface{}) error {
	query, err := qb.string()
	if err != nil {
		return err
	}
	row := qb.db.db.QueryRow(query, qb.args...)
	err = row.Scan(qb.getDests(reflect.ValueOf(out).Elem())...)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (qb *QueryBuilder) toMany(sliceType reflect.Type, outs interface{}) error {
	query, err := qb.string()
	if err != nil {
		return err
	}
	rows, err := qb.db.db.Query(query, qb.args...)
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

func (qb *QueryBuilder) string() (string, error) {
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
	case sqlActionCountFrom:
		q.WriteString("select count(*) from ")
		q.WriteString(qb.model.table)
	default:
		return "", ErrBadQuery
	}
	for _, join := range qb.joins {
		q.WriteString(" join ")
		q.WriteString(join)
	}
	if qb.where != "" {
		q.WriteString(" where ")
		q.WriteString(qb.where)
	}
	if qb.order != "" {
		q.WriteString(" order by ")
		q.WriteString(qb.order)
	}
	if qb.limit > 0 {
		q.WriteString(" limit ")
		q.WriteString(strconv.FormatUint(qb.limit, 10))
	}
	return q.String(), nil
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
