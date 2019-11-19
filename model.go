package gosql

import (
	"fmt"
	"reflect"
	"strings"
)

type model struct {
	name              string
	table             string
	typ               reflect.Type
	fields            []string
	fieldCount        int
	primaryFieldIndex int
	insertQuery       string
	updateQuery       string
	deleteQuery       string
}

type modelMap map[string]*model

var models modelMap

func init() {
	models = make(modelMap)
}

func (m *model) setInsertQuery() {
	var query strings.Builder
	var values strings.Builder
	query.WriteString("insert into ")
	query.WriteString(m.table)
	query.WriteString(" (")
	for i := 0; i < m.fieldCount; i++ {
		if m.primaryFieldIndex == i {
			continue
		}
		query.WriteString(m.fields[i])
		values.WriteString("?")
		if i == m.fieldCount-1 {
			query.WriteString(") ")
			values.WriteString(")")
		} else {
			query.WriteString(", ")
			values.WriteString(", ")
		}
	}
	query.WriteString("values (")
	query.WriteString(values.String())
	m.insertQuery = query.String()
}

func (m *model) setUpdateQuery() {
	var query strings.Builder
	query.WriteString("update ")
	query.WriteString(m.table)
	query.WriteString(" set ")
	for i := 1; i < m.fieldCount; i++ {
		query.WriteString(m.fields[i])
		query.WriteString(" = ?")
		if i < m.fieldCount-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(" where id = ?")
	m.updateQuery = query.String()
}

func (m *model) setDeleteQuery() {
	var query strings.Builder
	query.WriteString("delete from ")
	query.WriteString(m.table)
	query.WriteString(" where id = ?")
	m.deleteQuery = query.String()
}

func (m *model) mustBeValid() {
	if models[m.name] != nil {
		panic(fmt.Sprintf("model %s found more than once", m.name))
	}
	// idField := m.typ.Field(0)
	// if idField.Name != "ID" {
	// 	panic(fmt.Sprintf("first field of %s must be ID", m.name))
	// }
	// if idField.Type.Kind() != reflect.Int64 {
	// 	panic(fmt.Sprintf("%s.ID must have type int64", m.name))
	// }
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
	for i := 0; i < m.fieldCount; i++ {
		if m.primaryFieldIndex == i {
			continue
		}
		args = append(args, v.Field(i).Interface())
	}
	return args
}

func getModelOf(obj interface{}) (*model, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("obj must be a pointer to your model struct")
	}
	e := t.Elem()
	if e.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a pointer to your model struct")
	}
	m := models[e.Name()]
	if m == nil {
		return nil, fmt.Errorf("you must first register %s", e.Name())
	}
	return m, nil
}
