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
}

type modelMap map[string]*model

var models modelMap

func init() {
	models = make(modelMap)
}

func (m *model) getInsertQuery(v reflect.Value) string {
	var query strings.Builder
	var values strings.Builder
	query.WriteString("insert into ")
	query.WriteString(m.table)
	query.WriteString(" (")
	for i := 0; i < m.fieldCount; i++ {
		if m.primaryFieldIndex == i && v.Field(i).IsZero() {
			if i == m.fieldCount-1 {
				query.WriteString(") ")
				values.WriteString(")")
			}
			continue
		}
		query.WriteString(m.fields[i])
		values.WriteString("?")
		if i == m.fieldCount-1 {
			query.WriteString(") ")
			values.WriteString(")")
		} else if m.primaryFieldIndex != i+1 || m.primaryFieldIndex != m.fieldCount-1 {
			query.WriteString(", ")
			values.WriteString(", ")
		}
	}
	query.WriteString("values (")
	query.WriteString(values.String())
	return query.String()
}

func (m *model) getUpdateQuery(v reflect.Value) string {
	var query strings.Builder
	query.WriteString("update ")
	query.WriteString(m.table)
	query.WriteString(" set ")
	for i := 0; i < m.fieldCount; i++ {
		if m.primaryFieldIndex == i {
			continue
		}
		query.WriteString(m.fields[i])
		query.WriteString(" = ?")
		if i < m.fieldCount-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(" where ")
	query.WriteString(m.fields[m.primaryFieldIndex])
	query.WriteString(" = ?")
	return query.String()
}

func (m *model) mustBeValid() error {
	if models[m.name] != nil {
		return fmt.Errorf("model %s found more than once", m.name)
	}
	if m.primaryFieldIndex < 0 {
		return fmt.Errorf("model %s must have one and only one field tagged `gosql:\"primary\"`", m.name)
	}
	return nil
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
		f := v.Field(i)
		if m.primaryFieldIndex == i && f.IsZero() {
			continue
		}
		args = append(args, f.Interface())
	}
	return args
}

func (m *model) getArgsPrimaryLast(v reflect.Value) []interface{} {
	args := make([]interface{}, m.fieldCount)
	var primArg interface{}
	i := 0
	for {
		if i == m.fieldCount {
			break
		}
		arg := v.Field(i).Interface()
		if m.primaryFieldIndex == i {
			primArg = arg
		} else {
			if primArg == nil {
				args[i] = arg
			} else {
				args[i-1] = arg
			}
		}
		i++
	}
	args[m.fieldCount-1] = primArg
	return args
}

func getModelOf(obj interface{}) (*model, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("obj must be a pointer to your model struct")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a pointer to your model struct")
	}
	m := models[t.Name()]
	if m == nil {
		return nil, fmt.Errorf("you must first register %s", t.Name())
	}
	return m, nil
}
