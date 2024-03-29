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

type having struct {
	conjunction string
	condition   string
}

// Querier .
type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// SelectQuery holds information for a select query.
type SelectQuery struct {
	db         *DB
	querier    Querier
	model      *model
	fields     []string
	joins      []string
	wheres     []*where
	whereArgs  []interface{}
	havings    []*having
	havingArgs []interface{}
	groupBy    string
	order      string
	many       bool
	limit      int64
	offset     int64
}

// Join joins another table to this query.
func (sq *SelectQuery) Join(join string) *SelectQuery {
	sq.joins = append(sq.joins, fmt.Sprintf(" join %s", join))
	return sq
}

// LeftJoin joins another table to this query.
func (sq *SelectQuery) LeftJoin(join string) *SelectQuery {
	sq.joins = append(sq.joins, fmt.Sprintf(" left join %s", join))
	return sq
}

// Where specifies which rows will be returned.
func (sq *SelectQuery) Where(condition string, args ...interface{}) *SelectQuery {
	w := &where{
		conjunction: " and ",
		condition:   condition,
	}
	sq.wheres = append(sq.wheres, w)
	sq.whereArgs = append(sq.whereArgs, args...)
	return sq
}

// Having specifies which rows will be returned.
func (sq *SelectQuery) Having(condition string, args ...interface{}) *SelectQuery {
	h := &having{
		conjunction: " and ",
		condition:   condition,
	}
	sq.havings = append(sq.havings, h)
	sq.havingArgs = append(sq.havingArgs, args...)
	return sq
}

// OrWhere specifies which rows will be returned.
func (sq *SelectQuery) OrWhere(condition string, args ...interface{}) *SelectQuery {
	w := &where{
		conjunction: " or ",
		condition:   condition,
	}
	sq.wheres = append(sq.wheres, w)
	sq.whereArgs = append(sq.whereArgs, args...)
	return sq
}

// OrHaving specifies which rows will be returned.
func (sq *SelectQuery) OrHaving(condition string, args ...interface{}) *SelectQuery {
	h := &having{
		conjunction: " or ",
		condition:   condition,
	}
	sq.havings = append(sq.havings, h)
	sq.havingArgs = append(sq.havingArgs, args...)
	return sq
}

// GroupBy specifies how to group the results.
func (sq *SelectQuery) GroupBy(bys ...string) *SelectQuery {
	sq.groupBy = strings.Join(bys, ", ")
	return sq
}

// OrderBy orders the results by the given criteria.
func (sq *SelectQuery) OrderBy(orderBy string) *SelectQuery {
	sq.order = orderBy
	return sq
}

// Limit limits the number of results returned by the query.
func (sq *SelectQuery) Limit(limit int64) *SelectQuery {
	sq.limit = limit
	return sq
}

// Offset specifies the offset value in the query.
func (sq *SelectQuery) Offset(offset int64) *SelectQuery {
	sq.offset = offset
	return sq
}

// Get sets the result of the query to out. Get() can only take a pointer
// to a struct, a pointer to a slice of structs, or a pointer to a slice
// of pointers to structs.
func (sq *SelectQuery) Get(out interface{}) error {
	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("out must be a pointer")
	}
	t = t.Elem()
	switch t.Kind() {
	case reflect.Struct:
		var err error
		sq.model, err = sq.db.getModelOf(t)
		if err != nil {
			return err
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
			var err error
			sq.model, err = sq.db.getModelOf(el)
			if err != nil {
				return err
			}
			return sq.toMany(t, out)
		case reflect.Struct:
			var err error
			sq.model, err = sq.db.getModelOf(el)
			if err != nil {
				return err
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
	args := sq.whereArgs
	args = append(args, sq.havingArgs...)
	rows, err := sq.querier.Query(sq.String(), args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	fieldCount := len(columns)
	found := false
	for rows.Next() {
		dests := make([]interface{}, fieldCount)
		for j := 0; j < fieldCount; j++ {
			fieldIdx := sq.model.getFieldIndexByName(columns[j])
			if fieldIdx < 0 {
				return fmt.Errorf("no field for column %s", columns[j])
			}
			dests[j] = e.Field(fieldIdx).Addr().Interface()
		}
		if err := rows.Scan(dests...); err != nil {
			return err
		}
		found = true
	}
	if !found {
		return ErrNotFound
	}
	return nil
}

func (sq *SelectQuery) toMany(sliceType reflect.Type, outs interface{}) error {
	sq.many = true
	args := sq.whereArgs
	args = append(args, sq.havingArgs...)
	rows, err := sq.querier.Query(sq.String(), args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	var newOuts reflect.Value
	if sq.limit == 0 {
		newOuts = reflect.MakeSlice(sliceType, 8, 8)
	} else {
		newOuts = reflect.MakeSlice(sliceType, int(sq.limit), int(sq.limit))
	}
	i := 0
	columns, _ := rows.Columns()
	fieldCount := len(columns)
	fieldIndecies := make([]int, fieldCount)
	for j := 0; j < fieldCount; j++ {
		fieldIndecies[j] = sq.model.getFieldIndexByName(columns[j])
		if fieldIndecies[j] < 0 {
			return fmt.Errorf("no field for column %s", columns[j])
		}
	}
	dests := make([]interface{}, fieldCount)
	for rows.Next() {
		if newOuts.Len() == i {
			newOuts.SetCap(newOuts.Len() * 2)
			newOuts.SetLen(newOuts.Len() * 2)
		}
		newOut := newOuts.Index(i)
		newOut.Set(reflect.New(sq.model.typ))
		for j := 0; j < fieldCount; j++ {
			dests[j] = newOut.Elem().Field(fieldIndecies[j]).Addr().Interface()
		}
		if err := rows.Scan(dests...); err != nil {
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
	args := sq.whereArgs
	args = append(args, sq.havingArgs...)
	rows, err := sq.querier.Query(sq.String(), args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	var newOuts reflect.Value
	if sq.limit == 0 {
		newOuts = reflect.MakeSlice(sliceType, 8, 8)
	} else {
		newOuts = reflect.MakeSlice(sliceType, int(sq.limit), int(sq.limit))
	}
	i := 0
	columns, _ := rows.Columns()
	fieldCount := len(columns)
	fieldIndecies := make([]int, fieldCount)
	for j := 0; j < fieldCount; j++ {
		fieldIndecies[j] = sq.model.getFieldIndexByName(columns[j])
		if fieldIndecies[j] < 0 {
			return fmt.Errorf("no field for column %s", columns[j])
		}
	}
	dests := make([]interface{}, fieldCount)
	newOut := newOuts.Index(0)
	for rows.Next() {
		if newOuts.Len() == i {
			newOuts.SetCap(newOuts.Len() * 2)
			newOuts.SetLen(newOuts.Len() * 2)
		}
		newOut = newOuts.Index(i)
		for j := 0; j < fieldCount; j++ {
			dests[j] = newOut.Field(fieldIndecies[j]).Addr().Interface()
		}
		if err := rows.Scan(dests...); err != nil {
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

	if sq.groupBy != "" {
		q.WriteString(" group by ")
		q.WriteString(sq.groupBy)
	}

	for i, having := range sq.havings {
		if i == 0 {
			q.WriteString(" having ")
		} else {
			q.WriteString(having.conjunction)
		}
		q.WriteString(having.condition)
	}

	if sq.order != "" {
		q.WriteString(" order by ")
		q.WriteString(sq.order)
	}
	if sq.many {
		if sq.limit > 0 {
			q.WriteString(" limit ")
			q.WriteString(strconv.FormatInt(sq.limit, 10))
		}
	} else {
		q.WriteString(" limit 1")
	}
	if sq.offset > 0 {
		q.WriteString(" offset ")
		q.WriteString(strconv.FormatInt(sq.offset, 10))
	}
	return q.String()
}
