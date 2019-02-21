package gosql

import (
	"reflect"
)

var fieldKinds = []reflect.Kind{
	reflect.String,
	reflect.Bool,
	reflect.Uint,
	reflect.Uint64,
	reflect.Int,
	reflect.Int64,
	reflect.Float32,
	reflect.Float64,
}

func isField(s reflect.StructField) bool {
	k := s.Type.Kind()
	for i := 0; i < len(fieldKinds); i++ {
		if fieldKinds[i] == k {
			return true
		}
	}
	return false
}
