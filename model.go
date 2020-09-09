package gosql

import (
	"reflect"
	"strings"
)

type model struct {
	name                 string
	table                string
	typ                  reflect.Type
	fields               []string
	primaryFieldIndecies []int
}

func isIntIn(i int, arr []int) bool {
	for _, arrInt := range arr {
		if arrInt == i {
			return true
		}
	}
	return false
}

func (m *model) getInsertQuery(v reflect.Value) string {
	var query strings.Builder
	var values strings.Builder
	query.WriteString("insert into ")
	query.WriteString(m.table)
	query.WriteString(" (")
	for i := 0; i < len(m.fields); i++ {
		if isIntIn(i, m.primaryFieldIndecies) && v.Field(i).IsZero() {
			if i == len(m.fields)-1 {
				query.WriteString(") ")
				values.WriteString(")")
			}
			continue
		}
		query.WriteString(m.fields[i])
		values.WriteString("?")
		if i == len(m.fields)-1 {
			query.WriteString(") ")
			values.WriteString(")")
			continue
		}
		if !isIntIn(i+1, m.primaryFieldIndecies) || !v.Field(i+1).IsZero() {
			query.WriteString(", ")
			values.WriteString(", ")
		}
	}
	query.WriteString("values (")
	query.WriteString(values.String())
	return query.String()
}

func (m *model) getDeleteQuery() string {
	var query strings.Builder
	query.WriteString("delete from ")
	query.WriteString(m.table)
	query.WriteString(" where ")
	for i, index := range m.primaryFieldIndecies {
		query.WriteString(m.fields[index])
		query.WriteString(" = ?")
		if i < len(m.primaryFieldIndecies)-1 {
			query.WriteString(" and ")
		}
	}
	return query.String()
}

func (m *model) getUpdateQuery() string {
	var query strings.Builder
	query.WriteString("update ")
	query.WriteString(m.table)
	query.WriteString(" set ")
	for i := 0; i < len(m.fields); i++ {
		if isIntIn(i, m.primaryFieldIndecies) {
			continue
		}
		query.WriteString(m.fields[i])
		query.WriteString(" = ?")
		if i < len(m.fields)-1 && !isIntIn(i+1, m.primaryFieldIndecies) {
			query.WriteString(", ")
		}
	}
	query.WriteString(" where ")
	for i, index := range m.primaryFieldIndecies {
		query.WriteString(m.fields[index])
		query.WriteString(" = ?")
		if i < len(m.primaryFieldIndecies)-1 {
			query.WriteString(" and ")
		}
	}
	return query.String()
}

func (m *model) getFieldIndexByName(name string) int {
	for i, f := range m.fields {
		if name == f || strings.HasSuffix(name, "."+f) {
			return i
		}
	}
	return -1
}

func (m *model) getArgs(v reflect.Value) []interface{} {
	var args []interface{}
	for i := 0; i < len(m.fields); i++ {
		f := v.Field(i)
		if isIntIn(i, m.primaryFieldIndecies) && f.IsZero() {
			continue
		}
		args = append(args, f.Interface())
	}
	return args
}

func (m *model) getArgsPrimaryLast(v reflect.Value) []interface{} {
	var args []interface{}
	var primaryArgs []interface{}
	i := 0
	for {
		if i == len(m.fields) {
			break
		}
		arg := v.Field(i).Interface()
		if isIntIn(i, m.primaryFieldIndecies) {
			primaryArgs = append(primaryArgs, arg)
		} else {
			args = append(args, arg)
		}
		i++
	}
	args = append(args, primaryArgs...)
	return args
}
