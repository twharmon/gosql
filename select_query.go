package gosql

import (
	"database/sql"
	"errors"
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
	many   bool
	limit  int64
	offset int64
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

// To sets the result of the query to out. To() can only take a pointer
// to a struct, a pointer to a slice of structs, or a pointer to a slice
// of pointers to structs.
func (sq *SelectQuery) To(out interface{}) error {
	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("out must be a pointer")
	}
	t = t.Elem()
	switch t.Kind() {
	case reflect.Struct:
		sq.model = models[t.Name()]
		if sq.model == nil {
			return fmt.Errorf("you must first register %s", t.Name())
		}
		return sq.toOne(out)
	case reflect.Slice:
		el := t.Elem()
		switch el.Kind() {
		case reflect.Ptr:
			el = el.Elem()
			if el.Kind() != reflect.Struct {
				break
			}
			sq.model = models[el.Name()]
			if sq.model == nil {
				return fmt.Errorf("you must first register %s", el.Name())
			}
			return sq.toMany(t, out)
		case reflect.Struct:
			sq.model = models[el.Name()]
			if sq.model == nil {
				return fmt.Errorf("you must first register %s", el.Name())
			}
			return sq.toManyValues(t, out)
		}
	}
	return fmt.Errorf("out must be a struct, slice of structs, or slice of pointers to structs (%s found)", t.Kind().String())
}

func (sq *SelectQuery) toOne(out interface{}) error {
	e := reflect.ValueOf(out).Elem()
	if !e.IsValid() {
		return errors.New("out must not be a nil pointer")
	}
	row := sq.db.db.QueryRow(sq.String(), sq.args...)
	err := row.Scan(sq.getDests(e)...)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (sq *SelectQuery) toMany(sliceType reflect.Type, outs interface{}) error {
	sq.many = true
	rows, err := sq.db.db.Query(sq.String(), sq.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	newOuts := reflect.MakeSlice(sliceType, int(sq.limit), int(sq.limit))
	i := 0
	for rows.Next() {
		newOuts.Index(i).Set(reflect.New(sq.model.typ))
		if err := rows.Scan(sq.getDests(newOuts.Index(i).Elem())...); err != nil {
			return err
		}
		i++
	}
	v := reflect.Indirect(reflect.ValueOf(outs))
	v.Set(newOuts)
	v.SetLen(i)
	return nil
}

func (sq *SelectQuery) toManyValues(sliceType reflect.Type, outs interface{}) error {
	sq.many = true
	rows, err := sq.db.db.Query(sq.String(), sq.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	newOuts := reflect.MakeSlice(sliceType, int(sq.limit), int(sq.limit))
	i := 0
	for rows.Next() {
		if err := rows.Scan(sq.getDests(newOuts.Index(i))...); err != nil {
			return err
		}
		i++
	}
	v := reflect.Indirect(reflect.ValueOf(outs))
	v.Set(newOuts)
	v.SetLen(i)
	return nil
}

// String returns the string representation of SelectQuery.
func (sq *SelectQuery) String() string {
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
	if sq.many {
		q.WriteString(" limit ")
		q.WriteString(strconv.FormatInt(sq.limit, 10))
	} else {
		q.WriteString(" limit 1")
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
