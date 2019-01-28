package gosql

import (
	"fmt"
	"reflect"
	"strings"
)

type model struct {
	name        string
	table       string
	typ         reflect.Type
	fields      []reflect.StructField
	fieldCount  int
	oneToManys  []reflect.StructField
	manyToOnes  []reflect.StructField
	manyToManys []reflect.StructField
	oneToOnes   []reflect.StructField
	insertQuery string
	updateQuery string
	deleteQuery string
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
	for i := 1; i < m.fieldCount; i++ {
		query.WriteString(strings.ToLower(m.fields[i].Name))
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
		query.WriteString(strings.ToLower(m.fields[i].Name))
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

func (m *model) mustBeValid() error {
	if models[m.name] != nil {
		return fmt.Errorf("model %s found more than once", m.name)
	}
	idField := m.typ.Field(0)
	if idField.Name != "ID" {
		return fmt.Errorf("first field of %s must be ID", m.name)
	}
	if idField.Type.Kind() != reflect.Uint && idField.Type.Kind() != reflect.Uint64 {
		return fmt.Errorf("%s.ID must have type uint or uint64", m.name)
	}
	return nil
}
