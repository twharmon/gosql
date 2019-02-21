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
	fields      []string
	fieldCount  int
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
	// todo: make sure order after id is fields... then relatives...
	return nil
}

func (m *model) getFieldIndexByName(name string) int {
	for i, f := range m.fields {
		if f == name {
			return i
		}
	}
	return -1
}

func (m *model) getManyToOneColumnByType(typ string) string {
	for _, mto := range m.manyToOnes {
		if strings.HasSuffix(mto.Type.String(), typ) {
			return strings.ToLower(mto.Name) + "_id"
		}
	}
	return ""
}

func getArgs(m *model, v reflect.Value) []interface{} {
	args := make([]interface{}, m.fieldCount-1)
	for i := 1; i < m.fieldCount; i++ {
		args[i-1] = v.Field(i).Interface()
	}
	return args
}

func getArgsIDLast(m *model, v reflect.Value) []interface{} {
	args := make([]interface{}, m.fieldCount)
	for i := 1; i < m.fieldCount; i++ {
		args[i-1] = v.Field(i).Interface()
	}
	args[m.fieldCount-1] = v.Field(0).Interface()
	return args
}

func getIDArg(v reflect.Value) interface{} {
	return v.Field(0).Interface()
}
