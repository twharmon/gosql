package gosql

import "reflect"

func getArgs(m *model, v reflect.Value, obj interface{}) []interface{} {
	args := make([]interface{}, m.fieldCount-1)
	for i := 1; i < m.fieldCount; i++ {
		args[i-1] = v.Field(i).Interface()
	}
	return args
}

func getArgsIDLast(m *model, v reflect.Value, obj interface{}) []interface{} {
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
