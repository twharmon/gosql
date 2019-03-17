package gosql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type where struct {
	conjunction string
	condition   string
}

// SelectQuery .
type SelectQuery struct {
	db     *DB
	model  *model
	fields []string
	joins  []string
	wheres []*where
	args   []interface{}
	order  string
	limit  int64
	offset int64
}

// Select .
func (db *DB) Select(fields ...string) *SelectQuery {
	sq := new(SelectQuery)
	sq.db = db
	sq.fields = fields
	return sq
}

// Join .
func (sq *SelectQuery) Join(join string) *SelectQuery {
	sq.joins = append(sq.joins, join)
	return sq
}

// Where .
func (sq *SelectQuery) Where(condition string, args ...interface{}) *SelectQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	sq.wheres = append(sq.wheres, w)
	sq.args = append(sq.args, args...)
	return sq
}

// OrWhere .
func (sq *SelectQuery) OrWhere(condition string, args ...interface{}) *SelectQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	sq.wheres = append(sq.wheres, w)
	sq.args = append(sq.args, args...)
	return sq
}

// OrderBy .
func (sq *SelectQuery) OrderBy(orderBy string) *SelectQuery {
	sq.order = orderBy
	return sq
}

// Limit .
func (sq *SelectQuery) Limit(limit int64) *SelectQuery {
	sq.limit = limit
	return sq
}

// Offset .
func (sq *SelectQuery) Offset(offset int64) *SelectQuery {
	sq.offset = offset
	return sq
}

// To .
func (sq *SelectQuery) To(out interface{}) error {
	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("out must be a pointer")
	}
	e := t.Elem()
	switch e.Kind() {
	case reflect.Struct:
		sq.model = models[e.Name()]
		if sq.model == nil {
			return fmt.Errorf("you must first register %s", e.Name())
		}
		return sq.toOne(out)
	case reflect.Slice:
		ptr := e.Elem()
		if ptr.Kind() != reflect.Ptr {
			return fmt.Errorf("out must be a slice of pointers")
		}
		strct := ptr.Elem()
		if strct.Kind() != reflect.Struct {
			return fmt.Errorf("out must be a slice of pointers to structs")
		}
		sq.model = models[strct.Name()]
		if sq.model == nil {
			return fmt.Errorf("you must first register %s", strct.Name())
		}
		return sq.toMany(e, out)
	default:
		return fmt.Errorf("models must be a struct or slice (%s found)", e.Kind().String())
	}
}

func (sq *SelectQuery) toOne(out interface{}) error {
	row := sq.db.db.QueryRow(sq.string(), sq.args...)
	err := row.Scan(sq.getDests(reflect.ValueOf(out).Elem())...)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (sq *SelectQuery) toMany(sliceType reflect.Type, outs interface{}) error {
	rows, err := sq.db.db.Query(sq.string(), sq.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	newOuts := reflect.MakeSlice(sliceType, 0, int(sq.limit))
	for i := 0; rows.Next(); i++ {
		newOut := reflect.New(sq.model.typ)
		if err := rows.Scan(sq.getDests(newOut.Elem())...); err != nil {
			return err
		}
		newOuts = reflect.Append(newOuts, newOut)
	}
	reflect.Indirect(reflect.ValueOf(outs)).Set(newOuts)
	return nil
}

func (sq *SelectQuery) string() string {
	var q strings.Builder
	q.WriteString("select ")
	for i := 0; i < len(sq.fields)-1; i++ {
		q.WriteString(sq.fields[i])
		q.WriteString(", ")
	}
	q.WriteString(sq.fields[len(sq.fields)-1])
	q.WriteString(" from ")
	q.WriteString(sq.model.table)
	for _, join := range sq.joins {
		q.WriteString(" join ")
		q.WriteString(join)
	}
	for i, where := range sq.wheres {
		if i == 0 {
			q.WriteString(" where ")
		} else {
			q.WriteString(where.conjunction)
		}
		q.WriteString(where.condition)
	}
	if sq.order != "" {
		q.WriteString(" order by ")
		q.WriteString(sq.order)
	}
	if sq.limit > 0 {
		q.WriteString(" limit ")
		q.WriteString(strconv.FormatInt(sq.limit, 10))
	}
	if sq.offset > 0 {
		q.WriteString(" offset ")
		q.WriteString(strconv.FormatInt(sq.offset, 10))
	}
	return q.String()
}

func (sq *SelectQuery) getDests(v reflect.Value) []interface{} {
	if strings.HasSuffix(sq.fields[0], "*") {
		scans := make([]interface{}, sq.model.fieldCount)
		for i := 0; i < sq.model.fieldCount; i++ {
			scans[i] = v.Field(i).Addr().Interface()
		}
		return scans
	}
	fieldCount := len(sq.fields)
	scans := make([]interface{}, fieldCount)
	for i := 0; i < fieldCount; i++ {
		scans[i] = v.Field(sq.model.getFieldIndexByName(sq.fields[i])).Addr().Interface()
	}
	return scans
}
